package login

import (
	"html/template"
	"net/http"
	"quillpen/models"
	"quillpen/storage"

	"github.com/gorilla/schema"
	"github.com/gorilla/csrf"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/login.html"))

}

func LoginForm(resp http.ResponseWriter, req *http.Request) {

	templates.ExecuteTemplate(resp,"loginview",map[string]interface{}{
        csrf.TemplateTag: csrf.TemplateField(req),
    })

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
	result := storage.FindOne(bson.M{"email": given_account.Email},storage.ACCOUNTS_COLLECTION)
	if result.Account != nil {
		existing_account := result.Account

		// Password comparison

		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Inavlid password to HTML
	
			templates.ExecuteTemplate(resp, "InvalidPassword", nil)
			return
		}
		http.Redirect(resp,req,"/posts",http.StatusSeeOther)


	} else{

		templates.ExecuteTemplate(resp, "MissingAccount", nil)

	}

	

}
