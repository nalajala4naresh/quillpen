package storage

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"strconv"

	"github.com/gocql/gocql"
)

const caFilePath = "rds-combined-ca-bundle.pem"
const connectTimeout = 5
const queryTimeout = 30

var Session *gocql.Session

func init() {

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

	s, err := cluster.CreateSession()
	if err != nil {
		log.Printf("ERROR: fail create cassandra session, %s", err.Error())
		os.Exit(1)
	}
	Session = s

}

func clearSession() {
	Session.Close()
}

type CassandraConfig struct {
	host       string
	port       string
	keyspace   string
	conistency string
}

var cassandraConfig = CassandraConfig{
	host:       getEnv("CASSANDRA_HOST", "localhost"),
	port:       getEnv("CASSANDRA_PORT", "9042"),
	keyspace:   getEnv("CASSANDRA_KEYSPACE", "quillpen"),
	conistency: getEnv("CASSANDRA_CONSISTANCY", "LOCAL_QUORUM"),
}

func getEnv(key, fallback string) string {

	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
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
