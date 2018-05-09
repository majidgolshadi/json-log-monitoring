package json_log_monitoring

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCounter(c *gin.Context) {
	valid, counter := Analyzer.getResult()
	if !valid {
		c.JSON(480, gin.H{"error": "invalid data"})
	}

	c.JSON(http.StatusOK, counter)
}

func GetCounterAndReset(c *gin.Context) {
	valid, counter := Analyzer.getResult()
	Analyzer.ResetCounting()

	httpStatusCode := http.StatusOK
	if !valid {
		c.JSON(480, gin.H{"error": "invalid data"})
	}

	c.JSON(httpStatusCode, counter)
}
