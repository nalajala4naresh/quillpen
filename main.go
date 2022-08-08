package main

import (
	"html/template"

	"net/http"
	"os"

	"quillpen/accounts"
	"quillpen/editor"
	"quillpen/posts"
	"quillpen/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/gorilla/csrf"
	_ "github.com/gorilla/sessions"
)

var parsed_template *template.Template

func init() {

	parsed_template = template.Must(template.ParseFiles("templates/index.html", "templates/login.html", "templates/signup.html", "templates/posts.html"))
}

func main() {

	csrf_middleware := csrf.Protect([]byte("MyOwnSecret"),
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(false),
	)

	router := mux.NewRouter()

	// router.Schemes("https")

	router.HandleFunc("/", IndexHandler).Methods("GET")
	router.Handle("/signup", handlers.MethodHandler{"GET": http.HandlerFunc(accounts.SignUpForm),
		"POST": http.HandlerFunc(accounts.SignUpHandler)})
	router.Handle("/signin", handlers.MethodHandler{"GET": http.HandlerFunc(accounts.SignInForm),
		"POST": http.HandlerFunc(accounts.SignInHandler)})
	router.HandleFunc("/posts", posts.ListPosts).Methods("GET")
	router.HandleFunc("/post", posts.CreatePost).Methods("POST")
	router.HandleFunc("/post/{postid}", posts.GetPost).Methods("GET")

	router.HandleFunc("/editor", editor.EditorSpace).Methods("GET")

	logged_handlers := handlers.LoggingHandler(os.Stdout, router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers, "application/json", "application/x-www-form-urlencoded")
	compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080", csrf_middleware(compressedHandlers))

}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {

	result_set := storage.ListPosts()
	if result_set == nil {
		panic("Unable to get any posts")
	}

	for _, post := range result_set {

		if post.Content == "" {
			continue
		}


	}

	parsed_template.ExecuteTemplate(resp, "index", result_set)

}
