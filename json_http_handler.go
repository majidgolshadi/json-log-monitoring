package json_log_monitoring

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func IsValidAndReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": strconv.FormatBool(Analyzer.alwaysValid)})
	Analyzer.ResetCounting()
	return
}
