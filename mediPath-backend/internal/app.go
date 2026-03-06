package internal

import (
	"mediPath-backend/internal/config"
	"mediPath-backend/internal/server"
)

func StartApp() {

	cfg := config.Load()

	s := server.NewServer(cfg)

	s.Start()
}
