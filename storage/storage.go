package storage

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gocql/gocql"
)

type Store interface {
	Get(query string) (map[string]interface{}, error)
	List(query string) ([]map[string]interface{}, error)
	Delete(query string) (map[string]interface{}, error)
	Create(query string) (map[string]interface{}, error)
}

func New(config *CassandraConfig) (Store, error) {
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
		return nil, err
	}

	return &CassandraStore{session: s}, nil
}

type CassandraStore struct {
	session *gocql.Session
}

func (c *CassandraStore) Get(query string) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	itr := c.session.Query(query).Consistency(gocql.EachQuorum).Iter()

	for itr.MapScan(m) {
		return m, nil
	}
	return nil, errors.New(fmt.Sprintf("Get Query Failed %s ", query))
}

func (c *CassandraStore) List(query string) ([]map[string]interface{}, error) {
	return nil, nil
}

func (c *CassandraStore) Create(query string) (map[string]interface{}, error) {
	return nil, nil
}

func (c *CassandraStore) Delete(query string) (map[string]interface{}, error) {
	return nil, nil
}
