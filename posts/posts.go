package posts

import (
	"context"
	"fmt"
	
	"time"
	"quillpen/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)



func init() {
    uri := "mongodb://0.0.0.0:27017/?maxPoolSize=20&w=majority"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetConnectTimeout(2 * time.Second)

	client, err := mongo.Connect(context.TODO(),opts)

	if err!= nil {
		println(err.Error())

		panic("unable to connect to Mongodb")
	}


	defer func(){
		if err:= client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(context.TODO(),readpref.Primary()); err != nil {
		panic(err)

	}
	fmt.Println("Successfully connected and pinged.")


}

func Write_posts(posts_to_write []*models.Post) {
	uri := "mongodb://0.0.0.0:27017/?maxPoolSize=20&w=majority"
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetConnectTimeout(2 * time.Second)

	
	if posts_to_write == nil {
		fmt.Println("the post can not be empty and a bad call was done")
	}
	
	post := models.Post{PostId: "1234",Content: "this is my content ", Created_by: "naresh"}
	
	client, err := mongo.Connect(context.TODO(),opts)
	if err != nil{
		panic("connection failed")
	}
	posts := client.Database("test").Collection("posts")
	insert_context , cancel := context.WithTimeout(context.Background(),1* time.Second)
	defer cancel()

	posts.InsertOne(insert_context, post)




}
