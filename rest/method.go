package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var JsonStructErrorflag bool

func init()  {
	JsonStructErrorflag = false
}

func IsJson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": strconv.FormatBool(JsonStructErrorflag)})
	return
}

func IsJsonAndReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": strconv.FormatBool(JsonStructErrorflag)})
	JsonStructErrorflag = false
	return
}
