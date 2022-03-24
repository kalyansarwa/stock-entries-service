package controllers

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/kalyansarwa/stocksapi/api/responses"
)

type Health struct {
	Alive bool   `json:"alive"`
	Since string `json:"since"`
}

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}

func (server *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	if h := atomic.LoadInt64(&server.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		health := Health{}
		health.Alive = true
		since := time.Since(time.Unix(0, h))
		health.Since = since.String()
		responses.JSON(w, http.StatusOK, health)
	}

}
