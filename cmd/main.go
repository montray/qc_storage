package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/montray/storage"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL") )

	if err != nil {
		log.Fatalf("postgres connection error: %s", err.Error())
	}

	r := gin.Default()

	storeHandler := storage.NewStoreHandler(storage.NewMStorageService(storage.NewPgStorageRepository(db)))

	r.POST("/store/add", storeHandler.Add)
	r.POST("/store/order", storeHandler.Store)
	r.GET("/store/:product_id", storeHandler.Get)

	r.Run(":8888")
}


