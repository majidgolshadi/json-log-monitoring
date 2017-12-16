package rest

import "github.com/gin-gonic/gin"

func RunHttpServer(port string) error {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/is-json", IsJson)
		v1.GET("/is-json/reset", IsJsonAndReset)
	}

	return router.Run(port)
}
