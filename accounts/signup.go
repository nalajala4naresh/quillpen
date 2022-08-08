package accounts

import (
	"fmt"
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

	templates = template.Must(template.ParseFiles("templates/postsignup.html","templates/signup.html","templates/login.html"))

}


func SignUpForm(resp http.ResponseWriter , req *http.Request) {


	templates.ExecuteTemplate(resp,"signupview",map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
	})


}

func SignUpHandler(resp http.ResponseWriter , req *http.Request) {

	err := req.ParseForm()

	var new_account models.Signupform
	if err != nil {
		panic("Unable to parse the form")

	}
	decoder := schema.NewDecoder()
	decoder.Decode(&new_account, req.Form)
	// bcrypt the password
	new_account.Hash()

    // check if account exists based on email address
	
	// Try to create Account
	cerr := storage.CreateAccount(models.Account{Email: new_account.Email,Password: new_account.Password})
    
	//  use errors as is methods
	if cerr !=nil {
		
		log.Default().Printf("Lookup failed for email")
		templates.ExecuteTemplate(resp,"conflict",nil)
		return

	}
	http.Redirect(resp,req,"/posts",http.StatusSeeOther)
	templates.ExecuteTemplate(resp,"thankyou",nil)
}