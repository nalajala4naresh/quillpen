package storage

import (
	"errors"
	"fmt"
	"log"
	"github.com/quillpen/models"
	"time"

	"github.com/gocql/gocql"
)

// publishing new post
func CreatePost(post models.Post) error {

	q := "INSERT INTO POSTS (id, content,author,timestamp) VALUES (?, ?, ?,?)"
    query := Session.Query(q, post.PostId, post.Content, post.Author,post.Timestamp)
	err := query.Consistency(gocql.Quorum).Exec()
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

func GetPost(postid string) (*models.Post, error) {

	q := "SELECT * FROM POSTS WHERE ID = ? LIMIT 1"

	m := map[string]interface{}{}
	itr := Session.Query(q, postid).Consistency(gocql.EachQuorum).Iter()

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
