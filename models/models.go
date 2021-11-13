package models

import (
	"html/template"
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

	Email string `schema: "email,required", bson:"email,required"`
	Password string `schema: "password,required" , bson:"password,required"`
	
}

func (s *Account) ModelType() string {
	return "Account"

}



type Post struct{
	PostId string `bson:"post_id,omitempty"`
	Summary string `bson:"summary"`
	Title string `bson:"title"`
	MD_Content []byte `bson:"md_content", schema: "md_content,required"`
	HTML_Content template.HTML `bson:"-", schema: "-"`
	Created_by string `bson: "created_by",schema: "-"`
	Created_at time.Time `bson: "created_time",schema: "-"`
}


func (s *Post) ModelType() string {


	return "Post"

}

type ResultSet struct{
    Posts []*Post
	Accounts []*Account

}
type Result struct {

	Post *Post
	Account *Account
}
