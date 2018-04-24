package json_log_monitoring

import "github.com/gin-gonic/gin"

var Analyzer *analyzer
func RunHttpServer(analyzer *analyzer, port string) error {
	Analyzer = analyzer

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/is-json/reset", IsValidAndReset)
		v1.GET("/counter", GetCounter)
		v1.GET("/counter/reset", GetCounterAndReset)
	}

	return router.Run(port)
}
