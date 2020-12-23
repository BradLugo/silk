package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"webber"
)

func main() {
	var neo4jDao *webber.Neo4jDao
	var err error

	if len(os.Args) == 3 {
		target := os.Args[2]
		neo4jDao, err = webber.NewNeo4jDao(target, "", "")
		if err != nil {
			log.Panic("error while connecting to Neo4j", err)
		}
	} else if len(os.Args) == 5 {
		target := os.Args[2]
		username := os.Args[3]
		password := os.Args[4]
		neo4jDao, err = webber.NewNeo4jDao(target, username, password)
		if err != nil {
			log.Panicln("error while connecting to Neo4j", err)
		}
	} else {
		log.Fatalln("no Neo4j configs specified.\n" +
			"Usage: webber <http port> <neo4j target uri> [neo4j username] [neo4j password]")
	}

	defer neo4jDao.Close()

	router := webber.GetGinEngine(true)
	bindAddress := os.Args[1]

	httpSrv := http.Server{
		Addr:    bindAddress,
		Handler: router,
	}

	ctx, cancelFx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFx()
	defer log.Fatal(httpSrv.Shutdown(ctx))
	log.Fatal(httpSrv.ListenAndServe())
}
