package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"webber/models"
)

type Dao interface {
	CreateNote(n *models.Note) (uuid.UUID, error)
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

func (d *Neo4jDao) CreateNote(n *models.Note) (uuid.UUID, error) {
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := d.driver.NewSession(sessionConfig)
	if err != nil {
		return uuid.Nil, err
	}
	defer session.Close()

	id, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (n:Note { uuid: apoc.create.uuid(), text: $text, citation: $citation }) "+
				"RETURN n.uuid",
			map[string]interface{}{
				"text":     n.Text,
				"citation": n.Citation,
			})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return uuid.Nil, err
	}

	sid, ok := id.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("could not convert following interface{} to a string: %+v", id)
	}

	pid, err := uuid.Parse(sid)
	if err != nil {
		return uuid.Nil, err
	}

	return pid, nil
}

func (d *Neo4jDao) Close() {
	d.driver.Close()
}
