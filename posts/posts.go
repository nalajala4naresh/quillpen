package posts

import (
	"quillpen/storage"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"

)



var templates *template.Template
func init(){

	templates = template.Must(template.ParseFiles("templates/posts.html"))



}



func Write_post(resp http.ResponseWriter, req *http.Request) {
	




}

func Read_Posts(resp http.ResponseWriter, req *http.Request)  {


	posts := storage.ListPosts()

	templates.ExecuteTemplate(resp,"posts/list",posts)



}

func Read_A_Post(resp http.ResponseWriter, req *http.Request) {

	uri_params  := mux.Vars(req)
	post := storage.FindAPost(uri_params["postid"])

	if post == nil {
		panic("post not found")
	}


	templates.ExecuteTemplate(resp,"posts/one", post)



}
