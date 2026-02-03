package server

import "net/http"

func StartRouter() *http.ServeMux {
	mx := http.NewServeMux()
	mx.HandleFunc("/", ListPost)
	mx.HandleFunc("/user/login", LoginHandler)
	mx.HandleFunc("/user/logout", LogoutHandler)
	mx.HandleFunc("/user/register", CreateUserHandler)

	mx.HandleFunc("/posts/create", CreatePostHandler)
	mx.HandleFunc("/post", DetailPostHandler)
	mx.HandleFunc("/posts/update", UpdatePostHandler)
	mx.HandleFunc("/posts/delete/{id}", DeletePostHandler)
	return mx
}
