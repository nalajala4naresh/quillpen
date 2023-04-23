package posts

import (
	"fmt"
	"log"
	"time"

	"github.com/quillpen/models"
	"github.com/quillpen/storage"
)

// publishing new post
func createPost(post models.Post) error {
	q := "INSERT INTO POSTS (id, content,author,timestamp) VALUES (?, ?, ?,?)"
	_, err := storage.Cassandra.Create(q, post.PostId, post.Content, post.Author, post.Timestamp)
	if err != nil {
		log.Printf("ERROR: fail create post, %s", err.Error())
	}

	return err
}

// listing top posts per category
func listPosts() ([]*models.Post, error) {
	q := "SELECT * FROM POSTS LIMIT 20"

	rawposts, err := storage.Cassandra.List(q)
	if err != nil {
		return nil, err
	}
	var posts []*models.Post
	for _, rawpost := range rawposts {
		post := new(models.Post)
		// post.Title = m["title"].(string)
		post.Content = rawpost["content"].(string)
		post.PostId = rawpost["id"].(string)
		// post.Tags = m["tags"].([]string)
		// post.Timestamp = m["timestamp"].(time.Time)

		posts = append(posts, post)
	}
	// handle for empty database page data
	fmt.Println(len(posts))
	return posts, nil
}

func getPost(postid string) (*models.Post, error) {
	q := "SELECT * FROM POSTS WHERE ID = ? LIMIT 1"

	rpost, err := storage.Cassandra.Get(q, postid)
	if err != nil {
		return nil, err
	}

	post := &models.Post{}
	// post.Title = m["title"].(string)
	post.Content = rpost["content"].(string)
	post.PostId = rpost["id"].(string)
	// post.Tags = m["tags"].([]string)
	post.Timestamp = rpost["timestamp"].(time.Time)

	return post, nil
}
