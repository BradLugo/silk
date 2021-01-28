package db

import (
	"encoding/json"
	"fmt"
	"webber/graph/model"

	"github.com/google/uuid"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Dao interface {
	CreateNote(n *model.NewNote) (*model.Note, error)
	Close()
}

type Neo4jDao struct {
	driver neo4j.Driver
}

func NewNeo4jDao(target, username, password string) (*Neo4jDao, error) {
	driver, err := neo4j.NewDriver(target, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		return nil, err
	}

	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session := driver.NewSession(sessionConfig)
	defer session.Close()

	_, err = session.Run("CALL db.ping()", nil)
	if err != nil {
		return nil, err
	}

	return &Neo4jDao{
		driver: driver,
	}, nil
}

func (d *Neo4jDao) CreateNote(nn *model.NewNote) (*model.Note, error) {
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session := d.driver.NewSession(sessionConfig)
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var result neo4j.Result
		var err error

		if nn.RelatedTo != nil && len(nn.RelatedTo) > 0 {
			result, err = tx.Run(
				"CREATE (n:Note { uuid: apoc.create.uuid(), text: $text, citation: $citation }) WITH n "+
					"UNWIND $relatedTo AS uuid "+
					"OPTIONAL MATCH (m:Note { uuid: uuid })"+
					"CREATE r=((n)-[:RELATED_TO]->(m)) "+
					"RETURN n as new_node, collect(r) as associations",
				map[string]interface{}{
					"text":      nn.Text,
					"citation":  nn.Citation,
					"relatedTo": nn.RelatedTo,
				})
		} else {
			result, err = tx.Run(
				"CREATE (nn:Note { uuid: apoc.create.uuid(), text: $text, citation: $citation }) "+
					"RETURN nn",
				map[string]interface{}{
					"text":     nn.Text,
					"citation": nn.Citation,
				})
		}

		if err != nil {
			return nil, err
		}

		if result.Next() {
			node := result.Record().Values[0].(neo4j.Node)
			jsonbody, err := json.Marshal(node.Props)
			if err != nil {
				return nil, err
			}

			n := &model.Note{}
			if err := json.Unmarshal(jsonbody, n); err != nil {
				return nil, err
			}

			if suuid, ok := node.Props["uuid"]; ok {
				if _, err := uuid.Parse(suuid.(string)); err != nil {
					return nil, err
				}
				n.ID = suuid.(string)
			} else {
				return nil, fmt.Errorf("could not find uuid")
			}

			if nn.RelatedTo != nil && len(nn.RelatedTo) > 0 {
				pl := result.Record().Values[1].([]interface{})
				for _, up := range pl {
					p := up.(neo4j.Path)

					rids := make(map[int64]struct{})
					for _, pr := range p.Relationships {
						if pr.Type == "RELATED_TO" {
							if pr.StartId != node.Id && pr.EndId == node.Id {
								rids[pr.StartId] = struct{}{}
							} else if pr.StartId == node.Id && pr.EndId != node.Id {
								rids[pr.EndId] = struct{}{}
							} else if pr.StartId == node.Id && pr.EndId == node.Id {
								return nil, fmt.Errorf("circular relationship on new node")
							} else {
								return nil, fmt.Errorf("relationship does contain new node")
							}
						}
					}

					for _, pn := range p.Nodes {
						if pn.Id == node.Id {
							continue
						}

						if _, found := rids[pn.Id]; found {
							jsonbody, err := json.Marshal(pn.Props)
							if err != nil {
								return nil, err
							}

							rn := &model.Note{}
							if err := json.Unmarshal(jsonbody, rn); err != nil {
								return nil, err
							}

							if suuid, ok := pn.Props["uuid"]; ok {
								if _, err := uuid.Parse(suuid.(string)); err != nil {
									return nil, err
								}
								rn.ID = suuid.(string)
							} else {
								return nil, fmt.Errorf("could not find uuid")
							}

							n.RelatedTo = append(n.RelatedTo, rn)
						} else {
							return nil, fmt.Errorf("related node in path not found in relationships")
						}
					}
				}
			}

			return n, nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Note), nil
}

func (d *Neo4jDao) Close() {
	d.driver.Close()
}
