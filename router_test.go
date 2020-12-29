package webber_test

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"webber"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type RestSuite struct {
	suite.Suite
	engine *gin.Engine
}

func TestRestSuite(t *testing.T) {
	suite.Run(t, new(RestSuite))
}

func (s *RestSuite) SetupSuite() {
	s.engine = webber.GetGinEngine(nil, false)
}

func (s *RestSuite) TestBasicEngine() {
	health := serveRequest(s.engine, "GET", ``, "/health")
	s.Equal(`{"status":"imok"}`, string(health))
}

func serveRequest(engine *gin.Engine, method, body, url string) []byte {
	// generate and serve request
	requestBody := bytes.NewBufferString(body)
	req := httptest.NewRequest(method, url, requestBody)
	defer req.Body.Close()
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// review results
	resp := w.Result()
	result, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// http response may end in newline, does not need to propagate
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}
