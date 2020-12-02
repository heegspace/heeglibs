package youyoulibs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func HandleErr(c *gin.Context, code float64, msg string, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"error":   err,
	})
}
