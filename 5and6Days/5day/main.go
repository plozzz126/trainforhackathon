package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func main(){
	godotenv.Load()
	connectDB()
	r := gin.Default()

	r.POST("/users", CreateUser)

	r.Run(":8080")
}