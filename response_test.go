package youyoulibs

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleOK(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	HandleOK(c, "TestHandleOK")
}

func TestHandleErr(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	HandleErr(c, 200, "Test Success", nil)
}
