package accounts

import "net/http"

import "github.com/gorilla/csrf"
import "github.com/gorilla/schema"
import "golang.org/x/crypto/bcrypt"
import "quillpen/storage"
import "quillpen/models"

func SignInForm(resp http.ResponseWriter, req *http.Request) {

	templates.ExecuteTemplate(resp, "signinview", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
	})

}

func SignInHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")
	}

	var given_account models.Account

	decoder := schema.NewDecoder()
	decoder.Decode(&given_account, req.Form)

	// if account exists return an error mesage
	account, err := storage.GetAccount(given_account.Email)
	if err == nil {
		existing_account := account

		// Password comparison

		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Inavlid password to HTML

			templates.ExecuteTemplate(resp, "InvalidPassword", nil)
			return
		}
		http.Redirect(resp, req, "/posts", http.StatusSeeOther)

	} else {

		templates.ExecuteTemplate(resp, "MissingAccount", nil)

	}

}
