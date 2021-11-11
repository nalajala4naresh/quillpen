package login

import (
	"context"
	"html/template"
	"net/http"
	"quillpen/models"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("/Users/nareshnalajala/go/src/quillpen/templates/login.html"))

}

func LoginHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")
	}

	var given_account models.Account

	var existing_account models.Account

	decoder := schema.NewDecoder()
	decoder.Decode(&given_account, req.Form)

	// go to userdatabase and authenticate the user
	// check for existing account using email address
	options := options.Client().ApplyURI("mongodb://localhost:27017")
	time_out_context, t_cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer t_cancel()
	client, c_error := mongo.Connect(time_out_context, options)

	if c_error != nil {
		// Render server error in the UI
		panic("unable to establish the connection")

	}

	defer client.Disconnect(context.TODO())

	accounts := client.Database("quillpen").Collection("accounts")
	finder_context, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	// if account exists return an error mesage
	resultOne := accounts.FindOne(finder_context, bson.D{{"email", given_account.Email}})
	decode_err := resultOne.Decode(&existing_account)

	if decode_err != nil {
		templates.ExecuteTemplate(resp,"ISE", nil)
	return	
	}

	password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))

	if password_check_err != nil {
		// write Inavlid password to HTML

		templates.ExecuteTemplate(resp, "InvalidPassword", nil)
		return
	}
	// return session token
	templates.ExecuteTemplate(resp, "loginsuccess", nil)

}
