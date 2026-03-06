package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectRoot string
	MLDir       string
	MLScript    string
	OpenAIKey   string
	Port        string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		ProjectRoot: os.Getenv("PROJECT_ROOT"),
		MLDir:       os.Getenv("ML_DIR"),
		MLScript:    os.Getenv("ML_SCRIPT"),
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
		Port:        os.Getenv("PORT"),
	}
}

func (c *Config) GetMLScriptPath() string {

	return filepath.Join(
		c.ProjectRoot,
		c.MLDir,
		c.MLScript,
	)
}
