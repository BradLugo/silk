package main

import (
	"log"
	"os"
	"webber"
)

func main() {
	var neo4jDao *webber.Neo4jDao
	var err error

	if len(os.Args) == 2 {
		target := os.Args[1]
		neo4jDao, err = webber.NewNeo4jDao(target, "", "")
		if err != nil {
			log.Panic("error while connecting to Neo4j", err)
		}
	} else if len(os.Args) == 4 {
		target := os.Args[1]
		username := os.Args[2]
		password := os.Args[3]
		neo4jDao, err = webber.NewNeo4jDao(target, username, password)
		if err != nil {
			log.Panicln("error while connecting to Neo4j", err)
		}
	} else {
		log.Fatalln("No Neo4j configs specified.\n" +
			"Usage: webber <neo4j target uri> [neo4j username] [neo4j password]")
	}

	defer neo4jDao.Close()
	log.Println("Exiting successfully")
}
