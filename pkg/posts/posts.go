package posts

import (
	"log"
	"time"

	"github.com/quillpen/pkg/storage"
)

type Post struct {
	PostId    string    `json:"id" cql:"id"`
	Content   string    `json:"content" cql:"content"`
	Author    string    `json:"author" cql:"-"`
	Timestamp time.Time `json:"-" cql:"timestamp"`
	Tags      []string  `json:"-" cql:"tags"`
}

func (s *Post) ModelType() string {
	return "Post"
}

// publishing new post
func (s *Post) CreatePost() error {
	q := "INSERT INTO POSTS (id, content,author,timestamp) VALUES (?, ?, ?,?)"
	err := storage.Cassandra.Session.Query(q, s.PostId, s.Content, s.Author, s.Timestamp).Exec()
	if err != nil {
		log.Printf("ERROR: fail create post, %s", err.Error())
	}

	return err
}

// listing top posts per category
func (s *Post) ListPosts() ([]*Post, error) {
	q := "SELECT * FROM POSTS LIMIT 20"

	rawposts := storage.Cassandra.Session.Query(q, s.PostId, s.Content, s.Author, s.Timestamp).Iter()
	defer rawposts.Close()

	posts := make([]*Post, rawposts.NumRows())
	scanner := rawposts.Scanner()

	for scanner.Next() {
		scanner.Scan()
		post := new(Post)
		err := scanner.Scan(&post.PostId, &post.Content, &post.Timestamp, &post.Tags)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Post) GetPost() (*Post, error) {
	q := "SELECT * FROM POSTS WHERE ID = ? LIMIT 1"

	iter := storage.Cassandra.Session.Query(q, s.PostId).Iter()

	post := Post{}
	// post.Title = m["title"].(string)
	for iter.Scan(&post.PostId, &post.Content, &post.Timestamp, &post.Tags) {
		break
	}

	return &post, nil
}
