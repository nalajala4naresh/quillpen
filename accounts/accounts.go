package accounts

import (
	"errors"

	"github.com/quillpen/models"
	"github.com/quillpen/storage"

)

var (
	ACCOUNT_NOT_FOUND      = errors.New("Account not found")
	CAN_NOT_CREATE_ACCOUNT = errors.New("Account creation failed")
)

func createAccount(account models.Profile) error {
	q := "INSERT INTO ACCOUNTS (fullname,userhandle, email, password) VALUES(?,?,?,?)"
	_, err := storage.Cassandra.Create(q, account.Fullname, account.Userhandle, account.Email, account.Password)
	if err != nil {
		return CAN_NOT_CREATE_ACCOUNT
	}
	return nil
}

func getAccount(id string) (*models.Profile, error) {
	q := "SELECT * FROM ACCOUNTS WHERE email = ? LIMIT 1"

	raccount, err := storage.Cassandra.Get(q, id)
	if err != nil {
		return nil, err
	}

	account := &models.Profile{}
	account.Fullname = raccount["fullname"].(string)
	account.Email = raccount["email"].(string)
	account.Password = raccount["password"].(string)
	account.Userhandle = raccount["userhandle"].(string)

	return account, nil
}
