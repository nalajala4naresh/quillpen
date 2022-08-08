package accounts

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"quillpen/models"
	"quillpen/storage"

	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
)

var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/postsignup.html", "templates/signup.html", "templates/login.html"))

}

func SignUpForm(resp http.ResponseWriter, req *http.Request) {

	templates.ExecuteTemplate(resp, "signupview", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
	})

}

func SignUpHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	var new_account models.Profile
	if err != nil {
		panic("Unable to parse the form")

	}
	decoder := schema.NewDecoder()
	decoder.Decode(&new_account, req.Form)
	// bcrypt the password
	new_account.Hash()

	// check if account exists based on email address
	_, lerr := storage.GetAccount(new_account.Email)
	if errors.Is(lerr, storage.ACCOUNT_NOT_FOUND) {
		cerr := storage.CreateAccount(new_account)
		if cerr != nil {
			resp.Write([]byte("Please retry again Later !"))

			return
		}
		http.Redirect(resp, req, "/posts", http.StatusSeeOther)
		templates.ExecuteTemplate(resp, "thankyou", nil)

	} else {
		log.Default().Printf("Lookup failed for email")
		templates.ExecuteTemplate(resp, "conflict", nil)
		return

	}

}
