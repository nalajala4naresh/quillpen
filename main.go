package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//go:embed index.html
var Templates embed.FS


var cookie_store = sessions.NewCookieStore([]byte("NareshPass"))

type Qsession struct {

	Id primitive.ObjectID `bson: "_id,omitempty"`

}

type IndexCard struct {
	Title string
	Content string
}
func ArticleHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    
	fmt.Fprintf(resp, "Category: %v\n , Id: %v\n", vars["category"],vars["id"])
	
}




func main() {
    

	connect_context, cancel := context.WithTimeout(context.Background(), 3* time.Second)
	defer cancel()

	mongo_client, err := mongo.Connect(connect_context,options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic("Unable to Connect to mongo DB")
	}
	defer func(){
		if err = mongo_client.Disconnect(connect_context); err != nil {
			panic(err)
		}
	}()

	if err := mongo_client.Ping(context.TODO(),readpref.Primary()); err != nil {
		panic("Pinged the primary and Connected")
	}

	fmt.Println("Successfully connected and pinged.")

	databases, derr := mongo_client.ListDatabaseNames(connect_context,bson.M{})
    
	if derr != nil {
		println(derr.Error())
	}

	fmt.Println(databases)

	quickstartDb := mongo_client.Database("quickstart")
	podcasts := quickstartDb.Collection("podcasts")
	result, err := podcasts.InsertOne(connect_context,bson.D{{"Name","Naresh"},{"age",25}})
	if err != nil {
		log.Fatal("Unable to insert the data")

	}
	println(result.InsertedID)




	
	





	router := mux.NewRouter()
	router.Schemes("https")


	router.HandleFunc("/",IndexHandler).Methods("GET")
	router.HandleFunc("/signup",SignUpHandler).Methods("POST")

	router.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
    
	logged_handlers := handlers.LoggingHandler(os.Stdout,router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers,"application/json","application/x-www-form-urlencoded")
    compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServe(":8080",compressedHandlers)

	


}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := cookie_store.Get(req, "_first")

	session.Values["name"] = "naresh"
	session.Values["age"]= 28
	err := session.Save(req, resp)
	if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
 	
	index := template.Must(template.ParseFS(Templates, "index.html"))
	index.Execute(resp,nil)


}

type SignUpForm struct {

	FirstName string
	LastName string
	Phone string
	Email string
	Password string

}

func SignUpHandler(resp http.ResponseWriter, req *http.Request) {


	err := req.ParseForm()
	if err != nil {
    //    Invalid form data

	}
    var account_info SignUpForm
	decoder:= schema.NewDecoder()
	decoder.Decode(&account_info,req.Form)

	fmt.Println("Got the login deets as follows %s, %s, %s",account_info.Phone, account_info.Email,
account_info.Password)



}

func readDocuments() {



}
