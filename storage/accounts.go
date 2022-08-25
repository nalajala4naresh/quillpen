package storage

import (
	"errors"
	"fmt"
	"quillpen/models"

	"github.com/gocql/gocql"
)

var ACCOUNT_NOT_FOUND = errors.New("Account not found")
var CAN_NOT_CREATE_ACCOUNT = errors.New("Account creation failed")

func CreateAccount(account models.Profile) error {

	q := "INSERT INTO ACCOUNTS (fullname,userhandle, email, password) VALUES(?,?,?,?)"
	query := Session.Query(q, account.Fullname, account.Userhandle, account.Email, account.Password)
	fmt.Println(query.String())
	err := query.Consistency(gocql.Quorum).Exec()
	if err != nil {
		return CAN_NOT_CREATE_ACCOUNT

	}
	return nil

}

func GetAccount(id string) (*models.Profile, error) {
	q := "SELECT * FROM ACCOUNTS WHERE email = ? LIMIT 1"
	iter := Session.Query(q, id).Consistency(gocql.Quorum).Iter()
	m := map[string]interface{}{}
	for iter.MapScan(m) {

		account := &models.Profile{}
		account.Fullname = m["fullname"].(string)
		account.Email = m["email"].(string)
		account.Password = m["password"].(string)

		return account, nil

	}
	return nil, ACCOUNT_NOT_FOUND

}
