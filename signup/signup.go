package signup

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"quillpen/models"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)




var templates *template.Template

func init() {

	templates = template.Must(template.ParseFiles("/Users/nareshnalajala/go/src/quillpen/templates/postsignup.html"))

}

func SignUpHandler(resp http.ResponseWriter , req *http.Request) {

	err := req.ParseForm()

	var new_account models.Signupform
	if err != nil {
		panic("Unable to parse the form")

	}
	decoder := schema.NewDecoder()
	decoder.Decode(&new_account, req.Form)
	new_account.Hash()

	// check for existing account using email address
	options:= options.Client().ApplyURI("mongodb://localhost:27017")
    time_out_context , t_cancel := context.WithTimeout(context.Background(),1 * time.Second)
	defer t_cancel()
	client, c_error := mongo.Connect(time_out_context,options)
	

	if c_error != nil {
		// Render server error in the UI
		panic("unable to establish the connection")

	}

	defer client.Disconnect(context.TODO())

	accounts := client.Database("quillpen").Collection("accounts")
    upsert_context, cancel:= context.WithTimeout(context.TODO(), 1 * time.Second)
	defer cancel()
    // if account exists return an error mesage
	fmt.Printf("data after decoding with schema %s, %s",new_account.Email, new_account.Fullname)
	result, create_err := accounts.InsertOne(upsert_context, &new_account)
    
	
	if create_err != nil {
		println(create_err.Error())
		log.Default().Printf("Lookup failed for email")
		templates.ExecuteTemplate(resp,"conflict",nil)
		return
	}

	fmt.Println(result.InsertedID)
	templates.ExecuteTemplate(resp,"thankyou",nil)
	return
}