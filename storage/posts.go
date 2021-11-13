package storage

import (
	"context"
	"fmt"
	"quillpen/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var POSTS_COLLECTION = "posts"

func ListPosts() []*models.Post {

	connect, cancel := context.WithTimeout(context.Background(),2 * time.Second)
	defer cancel()	
	client, err := mongo.Connect(connect,MongoOptions(DOCDB_ENDPOINT))

	if err != nil {
		panic("Unable to connect to mongoDb")


	}
	defer client.Disconnect(context.Background())

	data_collection := client.Database(DOCDB_DB).Collection(POSTS_COLLECTION)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2 * time.Second)
    defer cancel2()

	cursor , lerr := data_collection.Find(ctx2, bson.D{})
	if err != nil {
		panic(lerr.Error())
	}

	

	var posts []*models.Post
	for cursor.Next(ctx2) {
		record := new(models.Post)
		derr := cursor.Decode(record)
		if derr != nil {
			fmt.Println(derr.Error())
			continue
		}
		posts = append(posts,record)

	}

	return posts


}



func FindAPost(post_id string) *models.Post {

	connect, cancel := context.WithTimeout(context.Background(),2 * time.Second)
	defer cancel()	
	client, err := mongo.Connect(connect,MongoOptions(DOCDB_ENDPOINT))

	if err != nil {
		panic("Unable to connect to mongoDb")


	}
	defer client.Disconnect(context.Background())

	data_collection := client.Database(DOCDB_DB).Collection(POSTS_COLLECTION)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2 * time.Second)
    defer cancel2()

	var post models.Post

	error := data_collection.FindOne(ctx2, bson.M{"post_id":post_id}).Decode(&post)
    
	if error != nil {
		println("unable to deocde a post")

	}

	return &post


}
