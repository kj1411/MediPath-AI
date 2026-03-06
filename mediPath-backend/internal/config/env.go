package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port         string
	OpenAIKey    string
	PythonPath   string
	MLScriptPath string
}

func LoadEnv() *EnvConfig {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	mlScript := filepath.Join(
		rootDir,
		"..",
		"mediPath-core",
		"mediPath-ml",
		"predict.py",
	)

	return &EnvConfig{
		Port:         getEnv("PORT", "8080"),
		OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		PythonPath:   getEnv("PYTHON_PATH", "python3"),
		MLScriptPath: mlScript,
	}
}

func getEnv(key string, fallback string) string {

	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}
