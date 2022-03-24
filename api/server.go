package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kalyansarwa/stocksapi/api/controllers"
	"github.com/kalyansarwa/stocksapi/api/seed"
)

var server = controllers.Server{}

func Run() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not coming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"),os.Getenv("DB_USER"),os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_URL"))

	seed.LoadDB(server.DB)
	
	server.Run(":7070")
}