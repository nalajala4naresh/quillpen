package posts

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"

	"quillpen/models"
	"quillpen/storage"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/posts.html"))

}

func CreatePost(resp http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	var post models.Post
	fmt.Println(req.Form)

	// fill th post details
	post.Content = strings.TrimSpace(req.Form["content"][0])
	if post.Content == "" {
		ListPosts(resp, req)
	}
	if len(post.Content) <= 50 {
		post.Title = post.Content[:len(post.Content)]
		
	} else if len(post.Content) > 50  {
		post.Title = post.Content[:50]
	}
	
	post.Timestamp = time.Now()
	post.PostId = gocql.UUIDFromTime(post.Timestamp)
	// now save the html bytes to Storage
	err := storage.CreatePost(post)
	if err != nil {
		panic("Unable to write the post")

	}

	// redirect user to the same post he wrote, to make any modifications
	ListPosts(resp, req)

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
	templates.ExecuteTemplate(resp, "posts/list", cleanedPosts)

}

func GetPost(resp http.ResponseWriter, req *http.Request) {

	uri_params := mux.Vars(req)

	result, err := storage.GetPost(uri_params["postid"])
	if err != nil {

		http.NotFound(resp, req)
		return

	}

	templates.ExecuteTemplate(resp, "posts/one", result)

}
