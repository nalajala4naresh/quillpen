package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/quillpen/sessionManager"
)

func SignInHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := sessionManager.Store.Get(req, sessionManager.SessionName)
	if !session.IsNew {
		isauthenticated := session.Values[sessionManager.SessionIsAuthenticated].(bool)
		if isauthenticated {
			http.Redirect(resp, req, "/chat", http.StatusSeeOther)
		}

	}

	var given_account Account
	defer req.Body.Close()

	jb, _ := ioutil.ReadAll(req.Body)

	jerr := json.Unmarshal(jb, &given_account)
	if jerr != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}

	account, err := given_account.GetAccount()
	if err == nil {
		existing_account := account
		// Password comparison
		fmt.Println(existing_account.Password, given_account.Password)
		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Unauthorized header

			resp.WriteHeader(http.StatusUnauthorized)
			return

		} else {

			existing_account.Password = "Unknown"
			session.Values[sessionManager.SessionUserId] = existing_account.UserId.String()
			session.Values[sessionManager.SessionIsAuthenticated] = true
			err = session.Save(req, resp)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			data, merr := json.Marshal(existing_account)
			if merr != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return

			}
			resp.WriteHeader(http.StatusAccepted)
			resp.Write(data)

		}

	} else if errors.Is(err, ACCOUNT_NOT_FOUND) {
		resp.WriteHeader(http.StatusNotFound)
	}
}
