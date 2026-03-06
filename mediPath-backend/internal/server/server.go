package server

import (
	"log"
	"net/http"

	"mediPath-backend/internal/config"
	"mediPath-backend/internal/server/endpoints"
	"mediPath-backend/internal/services"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start() {

	ml := services.NewMLService(s.cfg)
	ai := services.NewAIService(s.cfg)
	drug := services.NewDrugService()
	agent := services.NewAgentService(ai)

	handler := &endpoints.PredictHandler{
		ML:    ml,
		Agent: agent,
		Drug:  drug,
	}

	router := SetupRouter(handler)

	log.Println("Server running on port", s.cfg.Port)

	http.ListenAndServe(":"+s.cfg.Port, router)
}
