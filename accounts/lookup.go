package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func AccountLookUpHandler(w http.ResponseWriter, r *http.Request) {

	emailMap := mux.Vars(r)
	email := emailMap["email"]

	account := Account{Email: email}
	caccount, err := account.GetAccount()
	if errors.Is(err, ACCOUNT_NOT_FOUND) {
		w.WriteHeader(http.StatusNotFound)
		return

	}
	caccount.Password = ""

	data, err := json.Marshal(caccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
