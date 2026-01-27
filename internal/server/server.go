package server

import (
	"html/template"
	"log"
	"net/http"
	"simple_twitter/internal/models"

	"github.com/google/uuid"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
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
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ListPost(w http.ResponseWriter, r *http.Request) {
	posts, err := models.ListPost()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewHTTPServer() *http.Server {
	mx := http.NewServeMux()
	mx.HandleFunc("/", ListPost)
	mx.HandleFunc("/user/login", LoginHandler)
	mx.HandleFunc("/user/logout", LogoutHandler)
	mx.HandleFunc("/user/register", CreateUserHandler)

	mx.HandleFunc("/posts/create", CreatePostHandler)
	mx.HandleFunc("/post", DetailPostHandler)
	mx.HandleFunc("/posts/update", UpdatePostHandler)
	mx.HandleFunc("/posts/delete/{id}", DeletePostHandler)

	addr := &http.Server{
		Addr:    ":8080",
		Handler: mx,
	}
	log.Println("http server is listening on port 8080..")

	return addr
}
