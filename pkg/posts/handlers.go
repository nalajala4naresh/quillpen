package posts

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreatePost(resp http.ResponseWriter, req *http.Request) {
	var post Post

	defer req.Body.Close()
	data, _ := io.ReadAll(req.Body)
	json.Unmarshal(data, &post)

	// fill th post details
	post.Timestamp = time.Now()
	// now save the html bytes to Storage
	cerr := post.CreatePost()
	if cerr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
	resp.WriteHeader(http.StatusOK)
}

func ListPosts(resp http.ResponseWriter, req *http.Request) {
	var posts Post
	result_set, err := posts.ListPosts()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return

	}

	data, err := json.Marshal(result_set)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return

	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(data)
}

func GetPost(resp http.ResponseWriter, req *http.Request) {
	uri_params := mux.Vars(req)

	var post Post
	post.PostId = uri_params["postid"]

	result, err := post.GetPost()
	if err != nil {

		http.NotFound(resp, req)
		return

	}

	data, merr := json.Marshal(result)
	if merr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return

	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(data)
}
