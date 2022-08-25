package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
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

	
    
	var new_account models.Profile
	if req.Header["Content-Type"][0] == "application/json" {
	   defer req.Body.Close()
	   jb, _ := ioutil.ReadAll(req.Body)
       err := json.Unmarshal(jb,&new_account)
	   
	   if err != nil {
		panic("Unable to decode  the json")

	}
       

	} else {
		err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")

	}
	decoder := schema.NewDecoder()
	derr := decoder.Decode(&new_account, req.Form)
	if derr != nil {
		panic("Unable to parse the form")

	}
	}
	
    fmt.Printf("%s, %s, %s, %s",new_account.Email,new_account.Fullname,new_account.Password,new_account.Userhandle)
	

	// check if account exists based on email address
	existing_account, lerr := storage.GetAccount(new_account.Email)
	if errors.Is(lerr, storage.ACCOUNT_NOT_FOUND) {
		// bcrypt the password
	new_account.Hash()
		cerr := storage.CreateAccount(new_account)
		if cerr != nil {
            resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte("Please retry again Later !"))
            
			return
		}

		
		// http.Redirect(resp, req, "/posts", http.StatusSeeOther)
		// templates.ExecuteTemplate(resp, "thankyou", nil)

	} else if existing_account != nil {
		resp.WriteHeader(http.StatusConflict)
		return

	}

}
