package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"simple_twitter/internal/models"
	"simple_twitter/internal/nats_st"
	"simple_twitter/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			log.Err(err)
			w.Write([]byte(err.Error()))
			return
		}

		tmpl.Execute(w, nil)
	case http.MethodPost:
		username := r.FormValue("username")
		displayname := r.FormValue("display_name")
		pass := r.FormValue("password")
		email := r.FormValue("email")
		if _, err := mail.ParseAddress(email); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		user := models.NewUser(
			models.WithPassword(pass),
			models.WithUsername(username),
			models.WithDisplayName(displayname),
		)
		if err := user.Save(); err != nil {
			log.Err(err)
			w.Write([]byte(err.Error()))
			return
		}
		go func() {
			msg := models.Mail{
				To:      email,
				Subj:    "Simple Twitter Account Creation",
				Content: "You're account has been created",
			}
			data, err := json.Marshal(msg)
			if err != nil {
				log.Err(err)
			}
			nats_st.NATSServer.Publish("mail_sender", data)
		}()
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		tmpl.Execute(w, nil)
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		user := models.NewUser(
			models.WithUsername(username),
		)

		if err := user.Login(); err != nil {
			log.Err(err)
			w.Write([]byte(err.Error()))
			return
		}

		if password == user.Password {
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    user.ID.String(),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			})
			log.Debug().Msg("redirecting to dashboard page")

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Err(err)
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 10MB in memory, rest goes to temp files)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Err(err).Msg("failed to parse multipart form")
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	tokenValue := cookie.Value

	post := &models.Post{
		Content: content,
		User: models.User{
			ID: uuid.MustParse(tokenValue),
		},
	}

	// Handle optional image upload
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		ct := header.Header.Get("Content-Type")
		if !storage.IsImageContentType(ct) {
			http.Error(w, "Only image files (JPEG, PNG, GIF, WebP) are allowed", http.StatusBadRequest)
			return
		}

		if storage.S3Client == nil {
			http.Error(w, "Image upload is not configured", http.StatusServiceUnavailable)
			return
		}

		// Compress image to max 2MB
		compressed, finalCT, err := storage.CompressImage(file, ct)
		if err != nil {
			log.Err(err).Msg("failed to compress image")
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		// Upload to S3
		key := fmt.Sprintf("tweets/%s/%d_%s", tokenValue, time.Now().UnixMilli(), header.Filename)
		imageURL, err := storage.S3Client.Upload(compressed, key, finalCT)
		if err != nil {
			log.Err(err).Msg("failed to upload image to S3")
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}

		post.ImageURL = imageURL
	}

	if err := post.Save(); err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func ListPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		log.Err(err)
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}
	posts, err := models.ListPost()
	if err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Posts        []models.Post
		LoggedUserID string
	}{
		Posts:        posts,
		LoggedUserID: cookie.Value,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	tmpl.Execute(w, data)
}

func DetailPostHandler(w http.ResponseWriter, r *http.Request) {}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("id", r.FormValue("id")).Msg("checking for ID")
	uid, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	p := models.Post{
		ID:      uid,
		Content: r.FormValue("content"),
	}
	if err := p.Update(); err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	p := models.Post{
		ID: uuid.MustParse(postID),
	}
	if err := p.Delete(); err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
