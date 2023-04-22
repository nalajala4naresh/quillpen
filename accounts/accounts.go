package accounts

import (
	"errors"
	"fmt"

	"github.com/quillpen/models"
	"github.com/quillpen/storage"

	"github.com/gocql/gocql"
)

var (
	ACCOUNT_NOT_FOUND      = errors.New("Account not found")
	CAN_NOT_CREATE_ACCOUNT = errors.New("Account creation failed")
)

func createAccount(account models.Profile) error {
	q := "INSERT INTO ACCOUNTS (fullname,userhandle, email, password) VALUES(?,?,?,?)"
	query := storage.Session.Query(q, account.Fullname, account.Userhandle, account.Email, account.Password)
	fmt.Println(query.String())
	err := query.Consistency(gocql.Quorum).Exec()
	if err != nil {
		return CAN_NOT_CREATE_ACCOUNT
	}
	return nil
}

func getAccount(id string) (*models.Profile, error) {
	q := "SELECT * FROM ACCOUNTS WHERE email = ? LIMIT 1"
	iter := storage.Session.Query(q, id).Consistency(gocql.Quorum).Iter()
	m := map[string]interface{}{}
	for iter.MapScan(m) {

		account := &models.Profile{}
		account.Fullname = m["fullname"].(string)
		account.Email = m["email"].(string)
		account.Password = m["password"].(string)
		account.Userhandle = m["userhandle"].(string)

		return account, nil

	}
	return nil, ACCOUNT_NOT_FOUND
}
