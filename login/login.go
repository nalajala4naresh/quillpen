package login

import (
	"html/template"
	"net/http"
	"quillpen/models"
	"quillpen/storage"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/login.html"))

}

func LoginHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")
	}

	var given_account models.Account

	decoder := schema.NewDecoder()
	decoder.Decode(&given_account, req.Form)

	


	// if account exists return an error mesage
	existing_account := storage.FindAccount(bson.D{{"email", given_account.Email}})

	if existing_account != nil {

		// Password comparison

		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Inavlid password to HTML
	
			templates.ExecuteTemplate(resp, "InvalidPassword", nil)
			return
		}
		templates.ExecuteTemplate(resp, "loginsuccess", nil)


	} else{

		templates.ExecuteTemplate(resp, "MissingAccount", nil)

	}
	

	// return session token

	

}
