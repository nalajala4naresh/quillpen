package storage

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
)

const (
	caFilePath     = "rds-combined-ca-bundle.pem"
	connectTimeout = 5
	queryTimeout   = 30
)

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
