package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	"mediPath-backend/internal/config"
)

type AIService struct {
	cfg *config.Config
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{cfg: cfg}
}

func (a *AIService) Explain(disease string) (string, string, error) {

	prompt := "Explain the disease " + disease + " in simple language for a patient and suggest daily routine."

	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Authorization", "Bearer "+a.cfg.OpenAIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "AI explanation unavailable", "Consult doctor for advice.", nil
	}

	choice := choices[0].(map[string]interface{})

	message := choice["message"].(map[string]interface{})

	content := message["content"].(string)

	return content, "Maintain healthy lifestyle and consult doctor if symptoms persist.", nil
}
