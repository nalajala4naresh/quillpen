package signup

import "github.com/gorilla/schema"
import "net/http"


type Signupform struct {

	FirstName string
	LastName string
	Email string
	Phone string
	Password string


}

func SignUpHandler(resp http.ResponseWriter , req *http.Request) {

	err := req.ParseForm()

	var account Signupform
	if err != nil {
		panic("Unable to parse the form")

	}
	decoder := schema.NewDecoder()
	decoder.Decode(&account, req.Form)



}