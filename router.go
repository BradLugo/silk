package webber

import (
	"net/http"
	"webber/db"
	"webber/graph"
	"webber/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

func GetGinEngine(dao db.Dao, debug bool) (r *gin.Engine) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r = gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", routeHealth)

	r.POST("/query", graphqlHandler(dao))
	r.GET("/", playgroundHandler())

	if debug {
		ginpprof.Wrapper(r)
	}

	return r
}

func routeHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "imok",
	})
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Graphql handler
func graphqlHandler(dao db.Dao) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	r := generated.Config{Resolvers: &graph.Resolver{Dao: dao}}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(r))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
