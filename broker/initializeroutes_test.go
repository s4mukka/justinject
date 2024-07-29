package broker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntializeRoutes(t *testing.T) {
	router := gin.Default()
	intializeRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `"pong"`, w.Body.String())
}

func exampleHandler(ctx IRestContext) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func TestRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", route(exampleHandler))

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"success"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
