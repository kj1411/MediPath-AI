package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"

	appconfig "mediPath-backend/internal/config"
)

type AIService struct {
	cfg    *appconfig.Config
	client *bedrockruntime.Client
}

func NewAIService(cfg *appconfig.Config) *AIService {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Println("WARNING: Failed to load AWS config:", err)
		return &AIService{cfg: cfg}
	}

	client := bedrockruntime.NewFromConfig(awsCfg)
	return &AIService{cfg: cfg, client: client}
}

type bedrockMessage struct {
	Role    string           `json:"role"`
	Content []bedrockContent `json:"content"`
}

type bedrockContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type bedrockRequest struct {
	AnthropicVersion string           `json:"anthropic_version"`
	MaxTokens        int              `json:"max_tokens"`
	Messages         []bedrockMessage `json:"messages"`
}

type bedrockResponse struct {
	Content []bedrockContent `json:"content"`
}

func (a *AIService) Explain(disease string) (string, string, error) {

	if a.client == nil {
		log.Println("Bedrock client not initialized — check AWS credentials")
		return "AI explanation unavailable", "Consult doctor for advice.",
			fmt.Errorf("bedrock client not initialized")
	}

	prompt := fmt.Sprintf(
		`You are a caring medical assistant helping patients understand their health.

The patient has been identified with: %s

Please provide:
1. A clear, simple explanation of this condition in 3-4 sentences that a non-medical person can understand. Avoid jargon.
2. A practical daily care routine with 4-5 actionable tips.

Format your response EXACTLY as:
EXPLANATION:
<your explanation here>

ROUTINE:
<your routine here>`, disease)

	reqBody := bedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        1024,
		Messages: []bedrockMessage{
			{
				Role: "user",
				Content: []bedrockContent{
					{Type: "text", Text: prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	modelID := a.cfg.BedrockModel

	output, err := a.client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		ContentType: strPtr("application/json"),
		Accept:      strPtr("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		log.Println("Bedrock InvokeModel error:", err)
		return "AI explanation unavailable", "Consult doctor for advice.",
			fmt.Errorf("bedrock error: %w", err)
	}

	var result bedrockResponse
	if err := json.Unmarshal(output.Body, &result); err != nil {
		return "AI explanation unavailable", "Consult doctor for advice.",
			fmt.Errorf("failed to parse bedrock response: %w", err)
	}

	if len(result.Content) == 0 {
		return "AI explanation unavailable", "Consult doctor for advice.", nil
	}

	content := result.Content[0].Text

	explanation, routine := parseResponse(content)

	return explanation, routine, nil
}

func parseResponse(content string) (string, string) {
	parts := strings.SplitN(content, "ROUTINE:", 2)

	explanation := content
	routine := "Maintain a healthy lifestyle and consult your doctor if symptoms persist."

	if len(parts) == 2 {
		explanation = strings.TrimPrefix(parts[0], "EXPLANATION:")
		explanation = strings.TrimSpace(explanation)
		routine = strings.TrimSpace(parts[1])
	} else if idx := strings.Index(content, "EXPLANATION:"); idx >= 0 {
		explanation = strings.TrimSpace(content[idx+len("EXPLANATION:"):])
	}

	return explanation, routine
}

func strPtr(s string) *string {
	return &s
}
