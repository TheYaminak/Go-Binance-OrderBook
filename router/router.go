package router

import (
	//"github.com/TheYaminaK/Go-Binance-OrderBook/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/api/fill-book", middleware.BinanceData).Methods("POST")
	return r
}
