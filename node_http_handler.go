package json_log_monitoring

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
)

func GetCounter(c *gin.Context) {
	valid, counter := Analyzer.getResult()

	httpStatusCode := http.StatusOK
	if !valid {
		httpStatusCode = 480
	}
	jsonRes, err := json.Marshal(counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(httpStatusCode, jsonRes)
}

func GetCounterAndReset(c *gin.Context) {
	valid, counter := Analyzer.getResult()
	Analyzer.ResetCounting()

	httpStatusCode := http.StatusOK
	if !valid {
		httpStatusCode = 480
	}
	jsonRes, err := json.Marshal(counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(httpStatusCode, jsonRes)
}
