package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectRoot  string
	MLDir        string
	MLScript     string
	OpenAIKey    string
	Port         string
	AWSRegion    string
	BedrockModel string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	bedrockModel := os.Getenv("BEDROCK_MODEL")
	if bedrockModel == "" {
		bedrockModel = "anthropic.claude-3-haiku-20240307-v1:0"
	}

	return &Config{
		ProjectRoot:  os.Getenv("PROJECT_ROOT"),
		MLDir:        os.Getenv("ML_DIR"),
		MLScript:     os.Getenv("ML_SCRIPT"),
		OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		Port:         os.Getenv("PORT"),
		AWSRegion:    awsRegion,
		BedrockModel: bedrockModel,
	}
}

func (c *Config) GetMLScriptPath() string {

	return filepath.Join(
		c.ProjectRoot,
		c.MLDir,
		c.MLScript,
	)
}
