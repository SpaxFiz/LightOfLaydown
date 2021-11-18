package main

import (
	"SpaxFiz/LaydownLight/core/model"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	route := r.Group("/api")
	{
		route.GET("/em_account_data", func(c *gin.Context) {
			em := &model.EmAccount{}
			if err := em.Fetch(); err != nil {
				c.JSON(500, gin.H{"message": err.Error()})
				return
			}
			if data, e := em.Render(); e != nil {
				c.JSON(500, gin.H{"message": e.Error()})
				return
			} else {
				c.JSON(200, data)
			}
		})
	}

	//r.GET("/", func(context *gin.Context) {
	//	context.JSON(200, gin.H{"message": "hello world"})
	//})

	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
