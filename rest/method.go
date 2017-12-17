package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func IsJson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": strconv.FormatBool(getStatus())})
	return
}

func IsJsonAndReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": strconv.FormatBool(getStatus())})
	SetIsJSON()
	return
}
