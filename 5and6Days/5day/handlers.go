package main

import (
	"context"
	"github.com/gin-gonic/gin"
)



func CreateUser(c *gin.Context){
	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Ошибка неверный ввод"})
		return
	}
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO users (name, email) VALUES ($1, $2)",
		user.Name,
		user.Email,
	)

	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(200, gin.H{"message": "Данные успешно были добавлены в баззу данных"})
}