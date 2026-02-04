package server

import (
	"html/template"
	"net/http"
	"simple_twitter/internal/models"

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
			models.WithPassword(password),
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

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		data := map[string]string{"Error": "Username atau password salah!"}
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, data)
	}
}
func LogoutHandler(w http.ResponseWriter, r *http.Request)  {}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
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
	tokenValue := cookie.Value
	post := &models.Post{
		Content: content,
		User: models.User{
			ID: uuid.MustParse(tokenValue),
		},
	}
	if err := post.Save(); err != nil {
		log.Err(err)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {}

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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
