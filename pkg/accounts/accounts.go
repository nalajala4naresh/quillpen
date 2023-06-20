package accounts

import (
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"

	"github.com/quillpen/pkg/storage"
)

var (
	ACCOUNT_NOT_FOUND      = errors.New("Account not found")
	CAN_NOT_CREATE_ACCOUNT = errors.New("Account creation failed")
)

type User struct {
	Email         string                    `json:"email" cql:"email,required"`
	Username      string                    `json:"fullname" cql:"username,required"`
	UserId        gocql.UUID                `json:"userhandle,inline" cql:"user_id,required"`
	Conversations map[gocql.UUID]gocql.UUID `json:"conversations" cql:"conversations"`
}

func (u *User) ModelType() string {
	return "User"
}

func (u *User) GetUser() (*User, error) {
	query := `SELECT user_id, username, email FROM  users WHERE user_id= ? ;`
	iter := storage.Cassandra.Session.Query(query, u.UserId).Iter()

	var user User
	scanner := iter.Scanner()
	for scanner.Next() {

		err := scanner.Scan(&user.UserId, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}

	}

	return &user, nil
}

func (u *User) CreateUser() error {
	// insert into table
	query := `INSERT INTO quillpen.users(email, username, user_id) VALUES(?,?,?);`
	err := storage.Cassandra.Session.Query(query, u.Email, u.Username, u.UserId).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateLastRead(conversation, messageId gocql.UUID) error {
	u.Conversations[conversation] = messageId
	query := `UPDATE quillpen.users
	SET conversations  = conversations + {?:?}
	WHERE user_id = ?;`

	err := storage.Cassandra.Session.Query(query, conversation, messageId, u.UserId).Exec()

	return err
}

func (u *User) DeleteUser() error {
	conv_delete := `DELETE FROM quillpen.messages
	WHERE conversation_id = ? ;`
	for conversation := range u.Conversations {
		err := storage.Cassandra.Session.Query(conv_delete, conversation).Exec()
		if err != nil {
			return err
		}
	}
	query := `DELETE FROM quillpen.users
			WHERE user_id = ? ;`
	err := storage.Cassandra.Session.Query(query, u.UserId).Exec()
	return err
}

type Account struct {
	Email    string     `json:"email" schema:"email,required" cql:"email,required"`
	Password string     `json:"password,omitempty" schema:"password,required" cql:"password,required"`
	Username string     `json:"username" schema:"username,required" cql:"username,required"`
	UserId   gocql.UUID `json:"userid" cql:"user_id,required"`
	PublicKey string     `json:"publickey" cql:"public_key,required"` 
}

func (a *Account) GetAccount() (*Account, error) {
	q := "SELECT * FROM  accounts WHERE email = ?"
	iter := storage.Cassandra.Session.Query(q, a.Email).Iter()
	var account Account

	for iter.Scan(&account.Email, &account.Password, &account.UserId, &account.Username) {
		// Process each row of the result
		return &account, nil
	}

	return nil, ACCOUNT_NOT_FOUND
}

func (a *Account) CreateAccount() error {
	// generate Hash of the password
	a.Hash()
	user_id := gocql.MustRandomUUID()
	// add a row to accounts table
	q := "INSERT INTO accounts (email, password, user_id, username,publickey) VALUES(?,?,?,?,?);"
	err := storage.Cassandra.Session.Query(q, a.Email, a.Password, user_id, a.Username,a.PublicKey).Exec()
	if err != nil {
		fmt.Printf("user create error%s", err)
		return CAN_NOT_CREATE_ACCOUNT
	}
	q = "SELECT  email, user_id, username FROM accounts WHERE user_id = ?;"
	iter := storage.Cassandra.Session.Query(q, user_id).Iter()

	if !iter.Scan(&a.Email, &a.UserId, &a.Username) {
		return fmt.Errorf("unable to read created account  ")

	}
	// add a row to users table
	user := User{Email: a.Email, UserId: user_id,Username: a.Username}

	err = user.CreateUser()
	if err != nil {
		q := "DELETE FROM accounts WHERE email = ?;"
		derr := storage.Cassandra.Session.Query(q, a.Email).Exec()
		return fmt.Errorf("user create error%s and delete error %s", err, derr)
	}

	return nil
}

func (a *Account) UpdateAccount() error {
	// generate Hash of the password
	a.Hash()

	q := "UPDATE  accounts SET (email, password) VALUES(?,?)"
	err := storage.Cassandra.Session.Query(q, a.Email, a.Password).Exec()
	if err != nil {
		return CAN_NOT_CREATE_ACCOUNT
	}
	return nil
}

func (a *Account) Hash() {
	hashed_pass, error := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if error != nil {
		panic("unable to hash password")
	}

	a.Password = string(hashed_pass)
}

func (s *Account) ModelType() string {
	return "Account"
}
