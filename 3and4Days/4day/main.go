package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
)


type Tasks struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
}

type TaskU struct{
	Title string `json:"title" binding:"required"`
}

var tasks []Tasks
var nextID = 1

func main(){
	r := gin.Default()
	r.POST("/tasks", func(c *gin.Context){
		var input TaskU
		if err := c.ShouldBindJSON(&input); err != nil{
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		task := Tasks{
			ID: nextID,
			Title: input.Title,
			Done: false,
		}
		tasks = append(tasks, task)
		nextID++
		c.JSON(201, task)
	})

	r.GET("/tasks", func (c *gin.Context){
		c.JSON(200, tasks)
	})

	r.GET("/tasks/:id", func (c *gin.Context){
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil{
			c.JSON(400, gin.H{
        		"error": "invalid id",
    		})
    		return
		}
		for _, v := range tasks{
			if v.ID == id{
				c.JSON(200, v)
				return
			}
		}
		c.JSON(404, gin.H{
			"error": "id не был найден",
		})
	})


	r.PUT("/tasks/:id", func (c *gin.Context){
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil{
			c.JSON(400, gin.H{
				"error": "invalid id",
			})
		}

		for i := range tasks{
			if tasks[i].ID == id{
				tasks[i].Done = true
				c.JSON(200, tasks[i])
				return
			}
		}

		c.JSON(404, gin.H{
			"error": "Такого id нету",
		})
		return
	})

	r.Run()
}