package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"github.com/quillpen/models"
	"github.com/quillpen/storage"

	
)

var templates *template.Template


func SignUpHandler(resp http.ResponseWriter, req *http.Request) {

	
    
	var new_account models.Profile
	 
	   defer req.Body.Close()
	   jb, _ := ioutil.ReadAll(req.Body)
       err := json.Unmarshal(jb,&new_account)
	   
	   if err != nil {
		panic("Unable to decode  the json")}
	
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
