package storage

import (
	"errors"
	"fmt"
	"log"
	"quillpen/models"
	"time"

	"github.com/gocql/gocql"
)

// publishing new post
func CreatePost(post models.Post) error {

	q := "INSERT INTO POSTS (id, title, content, timestamp, tags) VALUES (?, ?, ?, ?, ?)"

	err := Session.Query(q, post.PostId, post.Title, post.Content, post.Timestamp, post.Tags).Exec()
	if err != nil {
		log.Printf("ERROR: fail create post, %s", err.Error())
	}

	return err

}

// listing top posts per category
func ListPosts() []*models.Post {

	q := "SELECT * FROM POSTS LIMIT 20"

	m := map[string]interface{}{}
	itr := Session.Query(q).Consistency(gocql.Quorum).Iter()
	var posts []*models.Post
	for itr.MapScan(m) {
		post := &models.Post{}
		post.Title = m["title"].(string)
		post.Content = m["content"].(string)
		post.PostId = m["id"].(gocql.UUID)
		post.Tags = m["tags"].([]string)
		post.Timestamp = m["timestamp"].(time.Time)
		fmt.Println(post.Content)

		posts = append(posts, post)
		m = map[string]interface{}{}
	}
	// handle for empty database page data
	return posts
}

func GetPost(postid string) (*models.Post, error) {

	q := "SELECT * FROM POSTS WHERE ID = ? LIMIT 1"

	m := map[string]interface{}{}
	itr := Session.Query(q, postid).Consistency(gocql.EachQuorum).Iter()

	for itr.MapScan(m) {
		post := &models.Post{}
		post.Title = m["title"].(string)
		post.Content = m["content"].(string)
		post.PostId = m["id"].(gocql.UUID)
		post.Tags = m["tags"].([]string)
		post.Timestamp = m["timestamp"].(time.Time)

		return post, nil
	}

	return nil, errors.New("document not found")

}
