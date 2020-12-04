package webber

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetGinEngine(debug bool) (r *gin.Engine) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r = gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", routeHealth)

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
