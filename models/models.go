package models

import (
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Profile struct {
	Account
	Fullname string `json: "fullname,required" cql:"fullname,required" schema:"fullname,omitempty"`
	Userhandle string `json: "userhandle,required" cql:"userhandle,required"`
}

func (p *Profile) ModelType() string {
	return "SignupForm"

}

type Model interface {
	ModelType() string
}

func (p *Profile) Hash() {

	hashed_pass, error := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if error != nil {
		panic("unable to hash password")
	}

	p.Password = string(hashed_pass)

}

type Account struct {
	Email    string `json: "email,required" schema: "email,required", cql:"email,required"`
	Password string `json: "password,required" schema: "password,required" , cql:"password,required"`
}

func (s *Account) ModelType() string {
	return "Account"

}

type Post struct {
	PostId    gocql.UUID `cql:"id,omitempty"`
	Title     string     `cql:"title"`
	Content   string     `cql:"content"`
	Author    string     `cql:"-"`
	Timestamp time.Time  `cql:"timestamp"`
	Tags      []string   `cql:"tags"`
}

func (s *Post) ModelType() string {

	return "Post"

}
