package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"webber"
	"webber/db"
	"webber/logging"

	"go.uber.org/zap"
)

func main() {
	var neo4jDao *db.Neo4jDao
	var err error

	logger := logging.NewWebberLogger()

	if len(os.Args) == 3 {
		target := os.Args[2]
		neo4jDao, err = db.NewNeo4jDao(target, "", "", logger)
		if err != nil {
			logger.Fatalw("error while connecting to Neo4j", zap.Error(err))
		}
	} else if len(os.Args) == 5 {
		target := os.Args[2]
		username := os.Args[3]
		password := os.Args[4]
		neo4jDao, err = db.NewNeo4jDao(target, username, password, logger)
		if err != nil {
			logger.Fatalw("error while connecting to Neo4j", zap.Error(err))
		}
	} else {
		logger.Fatalw("no Neo4j configs specified.\n" +
			"Usage: webber <http port> <neo4j target uri> [neo4j username] [neo4j password]")
	}

	defer neo4jDao.Close()

	router := webber.GetGinEngine(neo4jDao, true)
	bindAddress := os.Args[1]

	httpSrv := http.Server{
		Addr:    bindAddress,
		Handler: router,
	}

	ctx, cancelFx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFx()
	defer func() {
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal("error while shutting down http server", zap.Error(err))
		}
	}()
	log.Fatal(httpSrv.ListenAndServe())
}
