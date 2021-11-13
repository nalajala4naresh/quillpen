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




type Post struct{
	PostId string `bson:"post_id"`
	Summary string `bson:"summary"`
	Title string `bson:"title"`
	Content string `bson:"content"`
	Created_by string `bson: "created_by"`
	Created_at time.Time `bson: "created_time"`
}

