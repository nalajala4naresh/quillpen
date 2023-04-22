package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/quillpen/models"

	"golang.org/x/crypto/bcrypt"
)



func SignInHandler(resp http.ResponseWriter, req *http.Request) {
    
	var given_account models.Account
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&given_account)
		
	
	account, err := getAccount(given_account.Email)
	if err == nil {
		existing_account := account
		// Password comparison

		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Unauthorized header
			resp.WriteHeader(http.StatusUnauthorized)
			return
			
		} else {
			// nullfying password
			existing_account.Password = "Unknown"
			data , merr := json.Marshal(existing_account)
			if merr != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return

			}

			resp.Write(data)

		}
		
	} else if errors.Is(err,ACCOUNT_NOT_FOUND) {

		resp.WriteHeader(http.StatusNotFound)
		
	}

}
