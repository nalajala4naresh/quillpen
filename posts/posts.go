package posts

import (
	"errors"
	"fmt"
	"log"
	"github.com/quillpen/models"
	"github.com/quillpen/storage"
	"time"

	"github.com/gocql/gocql"
)

// publishing new post
func createPost(post models.Post) error {

	q := "INSERT INTO POSTS (id, content,author,timestamp) VALUES (?, ?, ?,?)"
    query := storage.Session.Query(q, post.PostId, post.Content, post.Author,post.Timestamp)
	// Explictly providing consistency 
	err := query.Consistency(gocql.Quorum).Exec()
	if err != nil {
		log.Printf("ERROR: fail create post, %s", err.Error())
	}

	return err

}

// listing top posts per category
func listPosts() []*models.Post {

	q := "SELECT * FROM POSTS LIMIT 20"

	m := map[string]interface{}{}
	itr := storage.Session.Query(q).Iter()
	var posts []*models.Post
	for itr.MapScan(m) {
		post := new(models.Post)
		// post.Title = m["title"].(string)
		post.Content = m["content"].(string)
		post.PostId = m["id"].(string)
		// post.Tags = m["tags"].([]string)
		// post.Timestamp = m["timestamp"].(time.Time)

		posts = append(posts, post)
	}
	// handle for empty database page data
	fmt.Println(len(posts))
	return posts
}

func getPost(postid string) (*models.Post, error) {

	q := "SELECT * FROM POSTS WHERE ID = ? LIMIT 1"

	m := map[string]interface{}{}
	itr := storage.Session.Query(q, postid).Consistency(gocql.EachQuorum).Iter()

	for itr.MapScan(m) {
		post := &models.Post{}
		// post.Title = m["title"].(string)
		post.Content = m["content"].(string)
		post.PostId = m["id"].(string)
		// post.Tags = m["tags"].([]string)
		post.Timestamp = m["timestamp"].(time.Time)

		return post, nil
	}

	return nil, errors.New("document not found")

}
