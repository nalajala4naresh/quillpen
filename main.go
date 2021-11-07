package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed index.html
var Templates embed.FS

type IndexCard struct {
	Title string
	Content string
}
func ArticleHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    
	fmt.Fprintf(resp, "Category: %v\n , Id: %v\n", vars["category"],vars["id"])
	
}


func main() {

	router := mux.NewRouter()
	router.Schemes("https")
	router.Use()
	


	router.HandleFunc("/",IndexHandler).Methods("GET")
	router.HandleFunc("/login",LoginHandler).Methods("POST")

	router.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
    
	logged_handlers := handlers.LoggingHandler(os.Stdout,router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers,"application/json")
    compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080",compressedHandlers)

	


}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
 mess := string("Hi form naresh")
 resp.Write([]byte(mess))


}

func LoginHandler(resp http.ResponseWriter, req *http.Request) {


}
