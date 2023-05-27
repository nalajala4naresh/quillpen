package storage

import (
	"log"
	"strconv"

	"github.com/gocql/gocql"
)

type Store interface {
	Get(query string,values ...interface{}) (map[string]interface{}, error)
	List(query string,values ...interface{}) ([]map[string]interface{}, error)
	Delete(query string,values ...interface{}) error
	Create(query string,values ...interface{}) (map[string]interface{}, error)
	Update(query string, values ...interface{}) (map[string]interface{}, error)
}

func NewCassandraStore(config *CassandraConfig) (Store, error) {
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
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra",Password: "cassandra"}

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

func (c *CassandraStore) Get(query string, values ...interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	itr := c.session.Query(query,values...).Consistency(gocql.EachQuorum).Iter()

	for itr.MapScan(m) {
		break
	}
	err := itr.Close()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *CassandraStore) List(query string, values ...interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 10)

	itr := c.session.Query(query,values...).Consistency(gocql.EachQuorum).Iter()

	for {
		entity := make(map[string]interface{})
		isend := itr.Scan(entity)
		if isend {
			break
		}
		result = append(result, entity)

	}
	err := itr.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CassandraStore) Create(query string, values ...interface{}) (map[string]interface{}, error) {
	itr := c.session.Query(query,values...).Consistency(gocql.EachQuorum).Iter()
	entity := make(map[string]interface{})
	for itr.MapScan(entity) {
		break
	}
	err := itr.Close()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (c *CassandraStore) Update(query string, values ...interface{}) (map[string]interface{}, error) {
	itr := c.session.Query(query,values...).Consistency(gocql.EachQuorum).Iter()
	entity := make(map[string]interface{})
	for itr.MapScan(entity) {
		break
	}
	err := itr.Close()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (c *CassandraStore) Delete(query string, values ...interface{}) error {
	return nil
}
