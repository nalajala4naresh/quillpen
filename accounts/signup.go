package accounts

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
)

var templates *template.Template

func SignUpHandler(resp http.ResponseWriter, req *http.Request) {
	var new_account Account

	defer req.Body.Close()
	jb, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(jb, &new_account)
	if err != nil {
		panic("Unable to decode  the json")
	}

	// check if account exists based on email address
	existing_account, lerr := new_account.GetAccount()
	if errors.Is(lerr, ACCOUNT_NOT_FOUND) {
		// bcrypt the password
		new_account.Hash()
		cerr := new_account.CreateAccount()
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
