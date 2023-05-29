package accounts

import (
	"errors"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"

	"github.com/quillpen/storage"
)

var (
	ACCOUNT_NOT_FOUND      = errors.New("Account not found")
	CAN_NOT_CREATE_ACCOUNT = errors.New("Account creation failed")
)

type User struct {
	Email         string                    `json:"email" cql:"email,required"`
	Username      string                    `json:"fullname" cql:"username,required"`
	UserId        gocql.UUID                `json:"userhandle" cql:"user_id,required"`
	Conversations map[gocql.UUID]gocql.UUID `json:"conversations" cql:"conversations"`
}

func (u *User) ModelType() string {
	return "User"
}

func (u *User) GetUser() (*User, error) {
	query := `SELECT * FROM  users WHERE user_id= ? ;`
	iter := storage.Cassandra.Session.Query(query, u.UserId).Iter()

	var user User
	for iter.Scanner().Next() {

		err := iter.Scanner().Scan(&u.UserId, &u.Username, &u.Email, &u.Conversations)
		if err != nil {
			return nil, err
		}

	}
	return &user, nil
}

func (u *User) CreateUser() error {
	// create uuid
	user_id := gocql.MustRandomUUID()
	// insert into table
	query := `INSERT INTO quillpen.users(email, username, user_id) VALUES(?,?,?);`
	err := storage.Cassandra.Session.Query(query, u.Email, u.Username, user_id).Exec()
	if err != nil {
		return err
	}
	// set the value upon success
	u.UserId = user_id
	return nil
}

func (u *User) AddConversation() (gocql.UUID, error) {
	if len(u.UserId) == 0 {
		return [16]byte{}, errors.New("UUID can not be empty")
	}
	conv_id := gocql.MustRandomUUID()
	query := `UPDATE quillpen.users
	SET conversations = conversations + {?}
	WHERE user_id = ?;`
	storage.Cassandra.Session.Query(query, conv_id, u.UserId)

	return conv_id, nil
}

func (u *User) UpdateLastRead(messageId gocql.UUID) error {
	return nil
}

func (u *User) DeleteUser() error {
	return nil
}

type Account struct {
	Email    string `json:"email" schema:"email,required" cql:"email,required"`
	Password string `json:"password" schema:"password,required" cql:"password,required"`
	FullName string `json:"username" schema:"username,required" cql:"username,required"`
}

func (a *Account) GetAccount() (*Account, error) {
	// generate Hash of the password
	a.Hash()

	q := "SELECT * FROM  accounts WHERE email = ?"
	iter := storage.Cassandra.Session.Query(q, a.Email).Iter()
	var account Account

	for iter.Scan(&account.Email, &account.Password, &account.FullName) {
		// Process each row of the result
		return &account, nil
	}

	return nil, ACCOUNT_NOT_FOUND
}

func (a *Account) CreateAccount() error {
	// generate Hash of the password
	a.Hash()

	q := "INSERT INTO accounts (email, password) VALUES(?,?)"
	err := storage.Cassandra.Session.Query(q, a.Email, a.Password).Exec()
	if err != nil {
		return CAN_NOT_CREATE_ACCOUNT
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
