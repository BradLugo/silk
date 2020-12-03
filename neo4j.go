package webber

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Dao interface {
	Close()
}

type Neo4jDao struct {
	driver neo4j.Driver
}

func NewNeo4jDao(target, username, password string) (*Neo4jDao, error) {
	driver, err := neo4j.NewDriver(target, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
		config.Encrypted = false
	})

	if err != nil {
		return nil, err
	}

	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	_, err = session.Run("CALL db.ping()", nil)
	if err != nil {
		return nil, err
	}

	return &Neo4jDao{
		driver: driver,
	}, nil
}

func (d *Neo4jDao) Close() {
	d.driver.Close()
}
