package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{
				"200": "OK",
			},
		)
	},
	)
	router.Run()
}
