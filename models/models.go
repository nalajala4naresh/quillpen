package models

import (

	"time"
	"golang.org/x/crypto/bcrypt"
)

type Signupform struct {

	Fullname string `bson:"fullname,omitempty", schema:"fullname,omitempty"`
	Email string      `bson:"email", schema:"email,required"`
	Password string   `bson:"password", schema:"password,required"`


}

func (s *Signupform) ModelType() string {
	return "SignupForm"

}

type Model interface {

	ModelType() string

}

func (s *Signupform) Hash(){

	hashed_pass, error := bcrypt.GenerateFromPassword([]byte(s.Password),bcrypt.DefaultCost)
	if error !=nil {
		panic("unable to hash password")
	}
	
	s.Password = string(hashed_pass)

}

type Account struct {


	Email string `schema: "email,required", cql:"email,required"`
	Password string `schema: "password,required" , cql:"password,required"`
	
}

func (s *Account) ModelType() string {
	return "Account"

}



type Post struct{




	PostId string `cql:"id,omitempty"`
	Title string `cql:"title"`
	Content []byte `cql:"content"`
	Author string `cql:"author"`
	Timestamp  time.Time   `cql:"timestamp"`
	Tags       []string    `cql:"tags"`
}

func (s *Post) ModelType() string {


	return "Post"

}
