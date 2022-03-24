package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kalyansarwa/stocksapi/api/models"
	"github.com/kalyansarwa/stocksapi/api/responses"
)

func (server *Server) GetStockEntries(w http.ResponseWriter, r *http.Request) {
	stockEntry := models.StockEntry{}

	stockEntries, err := stockEntry.FindAllStockEntries(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, stockEntries)
}

func (server *Server) GetStockEntriesByPortfolioId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	portfolioId := vars["portfolioId"]

	stockEntry := models.StockEntry{}

	stockEntries, err := stockEntry.FindAllStockEntriesByPortfolioId(server.DB, portfolioId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, stockEntries)

}

func (server *Server) GetStockEntryBySymbol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioId := vars["portfolioId"]
	symbol := vars["symbol"]

	stockEntry := models.StockEntry{}

	stockEntryFound, err := stockEntry.FindStockEntryBySymbol(server.DB, portfolioId, symbol)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	responses.JSON(w, http.StatusOK, stockEntryFound)
}

func (server *Server) CreateStockEntry(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	se := models.StockEntry{}

	err = json.Unmarshal(body, &se)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	se.Prepare()
	err = se.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	seNew, err := se.SaveStockEntry(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%s-%s", r.Host, r.RequestURI, seNew.PortfolioId, seNew.Symbol))
	responses.JSON(w, http.StatusCreated, seNew)
}

func (server *Server) UpdateStockEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioId := vars["portfolioId"]
	symbol := vars["symbol"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	se := models.StockEntry{}

	err = json.Unmarshal(body, &se)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	se.Prepare()
	//err = se.Valuedate("update")
	updatedEntry, err := se.UpdateStockEntry(server.DB, portfolioId, symbol)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, updatedEntry)
}

func (server *Server) DeleteStockEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioId := vars["portfolioId"]
	symbol := vars["symbol"]

	se := models.StockEntry{}
	var err error
	_, err = se.DeleteStockEntry(server.DB, portfolioId, symbol)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%s-%s", portfolioId, symbol))
	responses.JSON(w, http.StatusNoContent, "")

}
