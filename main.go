package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/quillpen/accounts"
	"github.com/quillpen/chat"
	"github.com/quillpen/posts"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/gorilla/sessions"
)

var parsed_template *template.Template

func init() {
	parsed_template = template.Must(template.ParseFiles("templates/index.html", "templates/login.html", "templates/signup.html", "templates/posts.html"))
}

func main() {
	router := mux.NewRouter()

	// router.Schemes("https")

	router.HandleFunc("/signin", accounts.SignInHandler).Methods("POST")
	router.HandleFunc("/users/{id}", accounts.UserHandler).Methods("GET", "PATCH")
	router.HandleFunc("/signup", accounts.SignUpHandler).Methods("POST")
	router.HandleFunc("/posts", posts.ListPosts).Methods("GET")
	router.HandleFunc("/post", posts.CreatePost).Methods("POST")
	router.HandleFunc("/post/{postid}", posts.GetPost).Methods("GET")
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/accounts/{email}", accounts.AccountLookUpHandler).Methods("GET")

	router.HandleFunc("/conversations", chat.ConversationsHandler).Methods("POST")
	router.HandleFunc("/conversation/{id}/{userid}", chat.ChatHandler)
	logged_handlers := handlers.LoggingHandler(os.Stdout, router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers, "application/json")
	compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080", compressedHandlers)
}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "/posts", http.StatusSeeOther)
}
