package main

import (
	"html/template"

	"net/http"
	"os"

	"quillpen/editor"
	"quillpen/login"
	"quillpen/posts"
	"quillpen/signup"
	"quillpen/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/csrf"
	_ "github.com/gorilla/sessions"
)




var parsed_template *template.Template

func init() {

	parsed_template = template.Must(template.ParseFiles("templates/index.html", "templates/login.html","templates/signup.html","templates/posts.html"))
}








func main() {

	csrf_middleware := csrf.Protect([]byte("MyOwnSecret"),
	csrf.RequestHeader("Authenticity-Token"),
	csrf.FieldName("authenticity_token"),
	csrf.SameSite(csrf.SameSiteLaxMode),
	csrf.Secure(false),
)


	router := mux.NewRouter()
	router.Schemes("https")




	router.HandleFunc("/",IndexHandler).Methods("GET")
	router.Handle("/signup",handlers.MethodHandler{"GET":http.HandlerFunc(signup.SignUpForm),
	"POST":http.HandlerFunc(signup.SignUpHandler)})
	router.Handle("/login",handlers.MethodHandler{"GET":http.HandlerFunc(login.LoginForm),
	"POST":http.HandlerFunc(login.LoginHandler)})
	router.HandleFunc("/posts",posts.Read_Posts).Methods("GET")
	router.HandleFunc("/post",posts.Write_post).Methods("POST")
	router.HandleFunc("/post/{postid}",posts.Read_A_Post).Methods("GET")

	router.HandleFunc("/editor",editor.EditorSpace).Methods("GET")


    
	logged_handlers := handlers.LoggingHandler(os.Stdout,router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers,"application/json","application/x-www-form-urlencoded")
    compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080",csrf_middleware(compressedHandlers))


}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
    
	
	result_set := storage.FindMany(bson.D{},storage.POSTS_COLLECTIONS,10)
	if result_set.Posts == nil {
		panic("Unable to get any posts")
	}

	parsed_template.ExecuteTemplate(resp,"index",map[string]interface{}{
        csrf.TemplateTag: csrf.TemplateField(req),
		"posts": result_set.Posts,
    })


}

