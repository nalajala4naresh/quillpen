package storage

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
)

var Cassandra *CassandraStore

type CassandraConfig struct {
	host       string
	port       string
	keyspace   string
	conistency string
	tls        *tls.Config
}

var cassandraConfig = CassandraConfig{
	host:       getEnv("CASSANDRA_HOST", "cassandra.us-east-1.amazonaws.com"),
	port:       getEnv("CASSANDRA_PORT", "9142"),
	keyspace:   getEnv("CASSANDRA_KEYSPACE", "quillpen"),
	conistency: getEnv("CASSANDRA_CONSISTANCY", "LOCAL_QUORUM"),
}

func init() {
	// initialize chat DB

	var err error
	Cassandra, err = NewCassandraStore(&cassandraConfig)
	if err != nil {
		panic(err.Error())
	}
	initCaassandra()
}

const (
	caFilePath     = "sf-class2-root.crt"
	connectTimeout = 5
	queryTimeout   = 30
)

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
