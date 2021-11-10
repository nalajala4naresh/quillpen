package login

import "net/http"
import "github.com/gorilla/schema"


type account struct{

	UserName string `schema: "username,required"`
	Password string `schema: "password,required"`
}


func LoginHandler(resp http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic("Unable to parse the form")
	}

	var account account

	decoder := schema.NewDecoder()
	decoder.Decode(&account, req.Form)

	// go to userdatabase and authenticate the user
    





}