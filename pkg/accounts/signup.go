package accounts

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/quillpen/pkg/sessionManager"
)

var templates *template.Template

func SignUpHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := sessionManager.Store.Get(req, sessionManager.SessionName)

	var new_account Account

	defer req.Body.Close()
	jb, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(jb, &new_account)
	if err != nil {

		resp.WriteHeader(http.StatusBadRequest)
		log.Printf("Json marshalling error %s", err)
		return
	}

	// check if account exists based on email address
	existing_account, lerr := new_account.GetAccount()
	if errors.Is(lerr, ACCOUNT_NOT_FOUND) {
		// bcrypt the password
		cerr := new_account.CreateAccount()
		if cerr != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte("Please retry again Later !"))

			return
		} else {
			session.Values[sessionManager.SessionUserId] = new_account.UserId.String()
			session.Values[sessionManager.SessionIsAuthenticated] = true
			err = session.Save(req, resp)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			resp.WriteHeader(http.StatusOK)
			data, _ := json.Marshal(new_account)
			resp.Write(data)

		}

		// http.Redirect(resp, req, "/posts", http.StatusSeeOther)
		// templates.ExecuteTemplate(resp, "thankyou", nil)

	} else if existing_account != nil {
		resp.WriteHeader(http.StatusConflict)
	}
}
