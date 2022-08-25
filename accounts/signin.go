package accounts

import "net/http"
import "encoding/json"


import "golang.org/x/crypto/bcrypt"
import "quillpen/storage"
import "quillpen/models"



func SignInHandler(resp http.ResponseWriter, req *http.Request) {
    
	var given_account models.Account
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&given_account)
		
	

	// if account exists return an error mesage
	account, err := storage.GetAccount(given_account.Email)
	if err == nil {
		existing_account := account
		// Password comparison

		password_check_err := bcrypt.CompareHashAndPassword([]byte(existing_account.Password), []byte(given_account.Password))
		if password_check_err != nil {
			// write Inavlid password to HTML

			resp.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {

		resp.WriteHeader(http.StatusNotFound)
	}

}
