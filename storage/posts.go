package storage

import "quillpen/models"

// publishing new post
func createPost(post models.Post) {


}


// listing top posts per category
func ListPosts() []models.Post {
var posts []models.Post
posts = append(posts, models.Post{})
return posts
}

func GetPost(postid string) models.Post {
	var post models.Post
	
	return post
	}