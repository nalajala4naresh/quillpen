package main

import (
	"net/http"
	"os"

	"github.com/quillpen/pkg/accounts"
	"github.com/quillpen/pkg/chat"
	"github.com/quillpen/pkg/posts"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/gorilla/sessions"
)

func main() {
	router := mux.NewRouter()
    
	serverAddr := ":443"
	certFile := "/etc/tls/cert.pem" // Path to your TLS certificate file
	keyFile := "/etc/tls/key.pem"

	router.HandleFunc("/signin", accounts.SignInHandler).Methods("POST")
	router.HandleFunc("/users/{id}", accounts.UserHandler).Methods("GET", "PATCH")
	router.HandleFunc("/signup", accounts.SignUpHandler).Methods("POST")
	router.HandleFunc("/posts", posts.ListPosts).Methods("GET")
	router.HandleFunc("/post", posts.CreatePost).Methods("POST")
	router.HandleFunc("/post/{postid}", posts.GetPost).Methods("GET")
	router.HandleFunc("/health", IndexHandler).Methods("GET")
	router.HandleFunc("/accounts/{email}", accounts.AccountLookUpHandler).Methods("GET")

	router.HandleFunc("/conversations", chat.ConversationsHandler).Methods("POST")
	router.HandleFunc("/conversations/{conversationId}", chat.DeleteConversationHandler).Methods("DELETE")
	router.HandleFunc("/conversations/{userId}", chat.ListConversationsHandler).Methods("GET")
	router.HandleFunc("/conversation/{id}/{userid}", chat.ChatHandler)
	logged_handlers := handlers.LoggingHandler(os.Stdout, router)
	contetTypeHandler := handlers.ContentTypeHandler(logged_handlers, "application/json")
	compressedHandlers := handlers.CompressHandler(contetTypeHandler)
	http.ListenAndServeTLS(serverAddr, certFile, keyFile, compressedHandlers)
}

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}
