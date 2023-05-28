package storage

import (
	"log"
	"strconv"

	"github.com/gocql/gocql"
)

func initCaassandra() {
	port := func(p string) int {
		iport, err := strconv.Atoi(p)
		if err != nil {
			return 9042
		}
		return iport
	}

	consistency := func(con string) gocql.Consistency {
		cc, err := gocql.MustParseConsistency(con)
		if err != nil {
			return gocql.Quorum
		}
		return cc
	}

	cluster := gocql.NewCluster(cassandraConfig.host)
	cluster.Port = port(cassandraConfig.port)
	cluster.Keyspace = cassandraConfig.keyspace
	cluster.Consistency = consistency(cassandraConfig.conistency)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}

	s, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("ERROR: fail create cassandra session, %s", err.Error())
	}

	// create keyspace if it does not exists
	if err := s.Query(`CREATE KEYSPACE IF NOT EXISTS quillpen 
	WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace quillpen: %s", err)
	}
	// create users if not exists

	// user_id is still the UUID that uniquely identifies each user, serving as the primary key.
	// username and email columns store user information, as before.
	// conversations is  a map data type where the key represents the conversation ID and the value represents the ID of the last read message.
	if err := s.Query(`CREATE TABLE IF NOT EXISTS quillpen.users ( 
		user_id UUID PRIMARY KEY,
		username text,
		email text,
		conversations map<UUID, UUID>
	);`).Exec(); err != nil {
		log.Fatalf("Failed to create users  table %s:", err)
	}

	// create Accounts

	if err := s.Query(`CREATE TABLE IF NOT EXISTS quillpen.accounts ( 
		email text PRIMARY KEY,
		password text 
		
	);`).Exec(); err != nil {
		log.Fatalf("Failed to create accounts  table %s:", err)
	}

	// create message table if does not exist
	if err := s.Query(`CREATE TABLE IF NOT EXISTS quillpen.messages (
		conversation_id UUID,
		message_id UUID,
		sender_id UUID,
		reciepent_id UUID,
		message text, 
		timestamp timestamp,
		PRIMARY KEY (conversation_id, timestamp,message_id ) 

	) WITH CLUSTERING ORDER BY (timestamp DESC);`); err != nil {
		log.Fatalf("Failed to create messages table %s:", err)
	}
}
