package signup

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"quillpen/models"
	"quillpen/storage"

	"github.com/gorilla/schema"
)




var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("templates/postsignup.html"))

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


	// Try to create Account
	result, create_err := storage.CreateAccount(new_account)
    
	_, ok := result.(int)
	if ok {
		println(create_err.Error())
		log.Default().Printf("Lookup failed for email")
		templates.ExecuteTemplate(resp,"conflict",nil)
		return

	}
	fmt.Println(result)
	templates.ExecuteTemplate(resp,"thankyou",nil)
}