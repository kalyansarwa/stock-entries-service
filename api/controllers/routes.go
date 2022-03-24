package controllers

import "github.com/kalyansarwa/stocksapi/api/middlewares"

func (s *Server) initializeRoutes() {

	s.Router.Use(middlewares.TracingMiddleware, middlewares.LoggingMiddleware)
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Health Route
	s.Router.HandleFunc("/health", s.HealthCheckHandler)

	// StockEntry routes
	s.Router.HandleFunc("/stockEntries", middlewares.SetMiddlewareJSON(s.CreateStockEntry)).Methods("POST")
	s.Router.HandleFunc("/stockEntries", middlewares.SetMiddlewareJSON(s.GetStockEntries)).Methods("GET")
	s.Router.HandleFunc("/stockEntries/{portfolioId}", middlewares.SetMiddlewareJSON(s.GetStockEntriesByPortfolioId)).Methods("GET")
	s.Router.HandleFunc("/stockEntries/{portfolioId}/{symbol}", middlewares.SetMiddlewareJSON(s.GetStockEntryBySymbol)).Methods("GET")
	s.Router.HandleFunc("/stockEntries/{portfolioId}/{symbol}", middlewares.SetMiddlewareJSON(s.UpdateStockEntry)).Methods("PUT")
	s.Router.HandleFunc("/stockEntries/{portfolioId}/{symbol}", middlewares.SetMiddlewareJSON(s.DeleteStockEntry)).Methods("DELETE")

}
