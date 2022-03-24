package controllers

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/kalyansarwa/stocksapi/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
	healthy int64
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, DbUrl string) {

	var err error

	DBURL := DbUrl
	server.DB, err = gorm.Open(Dbdriver, DBURL)

	if err != nil {
		log.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		log.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.Debug().AutoMigrate(&models.StockEntry{}) //database migration

	server.Router = mux.NewRouter()

	//router := http.NewServeMux()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Println("Listening to port 7070")
	atomic.StoreInt64(&server.healthy, time.Now().UnixNano())
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
