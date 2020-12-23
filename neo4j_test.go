package webber_test

import (
	"testing"
	"webber"

	"github.com/stretchr/testify/suite"
)

var (
	neo4jLocalBoltTarget = "bolt://localhost:7687"
)

type Neo4jDaoTestSuite struct {
	suite.Suite
	dao *webber.Neo4jDao
}

func TestCacheDaoTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode.")
	}
	suite.Run(t, new(Neo4jDaoTestSuite))
}

func (s *Neo4jDaoTestSuite) SetupSuite() {
	neo4jDao, err := webber.NewNeo4jDao(neo4jLocalBoltTarget, "", "")
	if err != nil {
		s.FailNow("Could not connect to local Neo4j database: ", err.Error())
	}
	s.dao = neo4jDao
}

func (s *Neo4jDaoTestSuite) TeardownSuite() {
	s.dao.Close()
}

func (s *Neo4jDaoTestSuite) TestInvalidConnection() {
	neo4jDao, err := webber.NewNeo4jDao("127.0.0.1:777", "doesNotExist", "fancyPassword")
	s.Nil(neo4jDao)
	s.Error(err)

	_, err = webber.NewNeo4jDao("bolt://127.0.0.1:777", "doesNotExist", "fancyPassword")
	s.Nil(neo4jDao)
	s.Error(err)
}
