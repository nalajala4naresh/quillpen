package posts

import (
	"encoding/json"
	
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/gocql/gocql"
	"github.com/gorilla/mux"

	"quillpen/models"
	"quillpen/storage"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/posts.html"))

}

func CreatePost(resp http.ResponseWriter, req *http.Request) {


	var post models.Post

	defer req.Body.Close()
	data, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(data,&post)

	// fill th post details
	post.Timestamp = time.Now()
	// now save the html bytes to Storage
	cerr := storage.CreatePost(post)
	if cerr != nil {
		resp.WriteHeader(http.StatusInternalServerError)

	}
	resp.WriteHeader(http.StatusOK)


}

func ListPosts(resp http.ResponseWriter, req *http.Request) {

	result_set := storage.ListPosts()
	if result_set == nil {
		return

	}
	var cleanedPosts []models.Post
	for _, post := range result_set {

		if (*post).Content == "" {
			continue
		}

		cleanedPosts = append(cleanedPosts, *post)

	}
	data , err := json.Marshal(cleanedPosts)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return


	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(data)



}

func GetPost(resp http.ResponseWriter, req *http.Request) {

	uri_params := mux.Vars(req)

	result, err := storage.GetPost(uri_params["postid"])
	if err != nil {

		http.NotFound(resp, req)
		return

	}

	
	data , merr := json.Marshal(result)
	if merr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return


	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(data)

}
