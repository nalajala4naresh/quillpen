package main

import (
	
	"embed"
	"html/template"
	
	"net/http"
	"os"

	

	"quillpen/login"
	"quillpen/signup"
	"quillpen/posts"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/gorilla/sessions"

)




var cookie_store = sessions.NewCookieStore([]byte("NareshPass"))
//go:embed templates/index.html templates/login.html templates/signup.html
var fs embed.FS

var parsed_template *template.Template

func init() {

	parsed_template = template.Must(template.ParseFS(fs,"templates/index.html", "templates/login.html","templates/signup.html"))
}








func main() {

	router := mux.NewRouter()
	router.Schemes("https")


	router.HandleFunc("/",IndexHandler).Methods("GET")
	router.HandleFunc("/signup",signup.SignUpHandler).Methods("POST")
	router.HandleFunc("/login",login.LoginHandler)

    
	logged_handlers := handlers.LoggingHandler(os.Stdout,router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers,"application/json","application/x-www-form-urlencoded")
    compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080",compressedHandlers)

	


}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := cookie_store.Get(req, "_first")

	session.Values["name"] = "naresh"
	session.Values["age"]= 28
	err := session.Save(req, resp)
	if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
 	
	parsed_template.ExecuteTemplate(resp,"base",nil)
	posts.Write_posts(nil)
	
	
	


}

