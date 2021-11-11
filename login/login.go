package login

import "net/http"
import "quillpen/models"
import "github.com/gorilla/schema"





func LoginHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")
	}

	var account models.Account

	decoder := schema.NewDecoder()
	decoder.Decode(&account, req.Form)

	// go to userdatabase and authenticate the user
    





}