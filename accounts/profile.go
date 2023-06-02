package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	_ "github.com/quillpen/sessionManager"
)

func ProfileHandler(resp http.ResponseWriter, req *http.Request) {

	// session, _ := sessionManager.Store.Get(req, sessionManager.SessionName)
	// if session.IsNew {
	// 	resp.WriteHeader(http.StatusUnauthorized)
	// 	return

	// }
	// isauthenticated := session.Values[sessionManager.SessionIsAuthenticated].(bool)
	// if !isauthenticated {

	// 	resp.WriteHeader(http.StatusUnauthorized)
	// 	return

	// }
	uri_params := mux.Vars(req)
	userId := uri_params["id"]

	userUuid, err := gocql.ParseUUID(userId)
	fmt.Printf("user id is %s", userId)
	if err != nil {

		resp.WriteHeader(http.StatusBadRequest)
		return

	}

	user := User{UserId: userUuid}
	fuser, err := user.GetUser()
	fmt.Printf(" data from get user is %s, %s, %s", fuser.Email, fuser.UserId, fuser.Username)
	if err != nil {
		fmt.Printf("unable to Get User %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return

	}

	data, err := json.Marshal(fuser)
	if err != nil {
		fmt.Printf("unable to marshal User %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return

	}
	resp.WriteHeader(http.StatusOK)
	fmt.Println(string(data))
	resp.Write(data)

}
