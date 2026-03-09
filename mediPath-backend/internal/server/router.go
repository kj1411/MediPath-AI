package server

import (
	"mediPath-backend/internal/server/endpoints"
	"mediPath-backend/internal/server/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter(handler *endpoints.PredictHandler) *mux.Router {

	r := mux.NewRouter()

	r.Use(middleware.EnableCORS)

	r.HandleFunc("/predict", handler.Handle).Methods("POST", "OPTIONS")

	return r
}
