package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Sum struct{
	A int `json:"a"`
	B int `json:"b"`
}

func main(){
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		start := time.Now()
		c.JSON(200, gin.H{
			"hello": "Hello world",
		})
		fmt.Println("Время", time.Since(start))
	})

	r.POST("/sum", func(c *gin.Context){
		start := time.Now()
		var req Sum

		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(400, gin.H{
				"error": "Неверный json",
			})
			return
		}
		
		result := req.A + req.B

		c.JSON(200, gin.H{
			"result": result, 
		})

		c.Next()

		fmt.Println("Время", time.Since(start))
	})


		
	

	fmt.Println("Успешно")
	r.Run()
	
}