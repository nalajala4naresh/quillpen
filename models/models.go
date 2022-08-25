package models

import (
	_ "github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Profile struct {
	Account
	Fullname string `json:"fullname" cql:"fullname,required"`
	Userhandle string `json:"userhandle" cql:"userhandle,required"`
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
	Email    string `json:"email" schema:"email,required" cql:"email,required"`
	Password string `json:"password" schema:"password,required" cql:"password,required"`
}

func (s *Account) ModelType() string {
	return "Account"

}

type Post struct {
	PostId    string `json:"id" cql:"id"`
	// Title     string     `json:"" cql:"title"`
	Content   string     `json:"content" cql:"content"`
	Author    string     `json:"author" cql:"-"`
	Timestamp time.Time  `json:"-" cql:"timestamp"`
	Tags      []string   `json:"-" cql:"tags"`
}

func (s *Post) ModelType() string {

	return "Post"

}
