package models


import "time"

type Signupform struct {

	Fullname string `bson:"fullname,omitempty", schema:"fullname,omitempty"`
	Email string      `bson:"email", schema:"email,required"`
	Password string   `bson:"password", schema:"password,required"`


}

type Account struct {

	UserName string `schema: "username,required"`
	Password string `schema: "password,required"`
	
}

type Post struct{
	PostId string `bson:"post_id"`
	Content string `bson:"content"`
	Created_by string `bson: "created_by"`
	Created_at time.Time `bson: "created_time"`
}
