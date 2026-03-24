package main 

import (
	"context"
	"fmt"
	"os"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func connectDB(){
	url := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil{
		panic(err)
	}
	db = pool
	fmt.Println("База данных подключена")
}

