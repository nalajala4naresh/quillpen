package storage

import (
	"log"
	"strconv"

	"github.com/gocql/gocql"
)

func NewCassandraStore(config *CassandraConfig) (*CassandraStore, error) {
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
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "quillpen-cassandra-at-783686645338", Password: "wGycWm0Erfc3u7nFBdk1oB8T0hiiVWeStzO7lI/2l2E="}
	cluster.SslOpts = &gocql.SslOptions{
		CaPath:                 "sf-class2-root.crt",
		EnableHostVerification: false,
	}
	s, err := cluster.CreateSession()
	if err != nil {
		log.Printf("ERROR: fail create cassandra session, %s", err.Error())
		return nil, err
	}

	return &CassandraStore{Session: s}, nil
}

type CassandraStore struct {
	Session *gocql.Session
}
