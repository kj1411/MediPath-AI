package services

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"mediPath-backend/internal/config"
)

type MLService struct {
	cfg *config.Config
}

func NewMLService(cfg *config.Config) *MLService {
	return &MLService{cfg: cfg}
}

func (m *MLService) Predict(drugs []string) (map[string]interface{}, error) {

	script := m.cfg.GetMLScriptPath()

	args := append([]string{script}, drugs...)

	cmd := exec.Command("python3", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf(
			"python execution failed: %v, output: %s",
			err,
			string(output),
		)
	}

	var result map[string]interface{}

	err = json.Unmarshal(output, &result)

	return result, err
}
