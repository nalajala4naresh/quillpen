package storage

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"quillpen/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const caFilePath = "rds-combined-ca-bundle.pem"
const connectTimeout=5
const queryTimeout = 30

const connectionStringTemplate = "mongodb://%s:%s@%s/quillpen"
// var username = "myUserAdmin"
// var password = "abc123"
// var DOCDB_ENDPOINT = "localhost"
// var DOCDB_DB = "quillpen"
// var ACCOUNTS_COLLECTION = "accounts"

var username = os.Getenv("DOCDB_USER")
var password = os.Getenv("DOCDB_PASS")
var DOCDB_ENDPOINT = os.Getenv("DOCDB_ENDPOINT")
var DOCDB_DB = os.Getenv("DOCDB_DB")
var ACCOUNTS_COLLECTION = os.Getenv("DOCDB_ACCOUNTS")


func init() {

	optios := MongoOptions(DOCDB_ENDPOINT)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
    defer cancel()

	client, err := mongo.Connect(ctx,optios)
	
	
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to DocumentDB!")



}
func MongoOptions(endpoint string) *options.ClientOptions {


	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, endpoint)

	// tlsConfig, err := getCustomTLSConfig(caFilePath)
	// if err != nil {
		// log.Fatalf("Failed getting TLS configuration: %v", err)
	// }

	c_options := options.Client().ApplyURI(connectionURI)

	// if os.Getenv("STAGE") != "local" {
	// 	c_options = c_options.SetTLSConfig(tlsConfig)
		
		
	// }

	return c_options


	

}


func FindAccount(query interface{}) *models.Account{
    query , ok := query.(bson.D)
	if !ok {
		panic("The query should be of bson.D type")

	}

    connect, cancel := context.WithTimeout(context.Background(),2 * time.Second)
	defer cancel()
    client, err := mongo.Connect(connect,MongoOptions(DOCDB_ENDPOINT))
	if err != nil {
		panic("Mongo comnnection failed")
	}
	defer client.Disconnect(context.Background())

	data_collection := client.Database(DOCDB_DB).Collection(ACCOUNTS_COLLECTION)

	ctx2, cancel2 := context.WithTimeout(context.Background(), queryTimeout*time.Second)
    defer cancel2()

	var account models.Account

	find_err := data_collection.FindOne(ctx2, query).Decode(account)

	if find_err !=nil {
		return nil

	}
	return &account

	

}

func CreateAccount(account models.Signupform) (interface{}, error) {

	connect, cancel := context.WithTimeout(context.Background(),2 * time.Second)
	defer cancel()
    client, err := mongo.Connect(connect,MongoOptions(DOCDB_ENDPOINT))
	if err != nil {
		panic("Mongo comnnection failed")
	}
	defer client.Disconnect(context.Background())
	
	data_collection := client.Database(DOCDB_DB).Collection(ACCOUNTS_COLLECTION)
	ctx2, cancel2 := context.WithTimeout(context.Background(), queryTimeout*time.Second)
    defer cancel2()

	insert_result, err := data_collection.InsertOne(ctx2, account)
	if err != nil {
		return -1, err
	}
	return insert_result.InsertedID , nil


}
func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}