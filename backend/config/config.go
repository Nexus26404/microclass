package config

import (
	"os"
	"sync"
)

type Config struct {
	Port        string
	MockMode    bool
	OpenAIKey   string
	OpenAIURL   string
	OpenAIModel string
	DBPath      string
	mu          sync.RWMutex
}

var cfg *Config

func Load() *Config {
	cfg = &Config{
		Port:        getEnv("PORT", "8080"),
		MockMode:    getEnv("MOCK_MODE", "true") == "true",
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		OpenAIURL:   getEnv("OPENAI_BASE_URL", "https://api.minimax.chat"),
		OpenAIModel: getEnv("OPENAI_MODEL", "MiniMax-Text-01"),
		DBPath:      getEnv("DB_PATH", "./data/microclass.db"),
	}
	return cfg
}

func Get() *Config {
	return cfg
}

func (c *Config) Update(mockMode *bool, openAIKey string, openAIURL, model *string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if mockMode != nil {
		c.MockMode = *mockMode
	}
	if openAIKey != "" {
		c.OpenAIKey = openAIKey
	}
	if openAIURL != nil {
		c.OpenAIURL = *openAIURL
	}
	if model != nil {
		c.OpenAIModel = *model
	}
}

func InitFromDB(mockMode bool, openAIKey, openAIURL, openAIModel string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if !mockMode || openAIKey != "" || openAIURL != "" || openAIModel != "" {
		if !mockMode {
			cfg.MockMode = false
		}
		if openAIKey != "" {
			cfg.OpenAIKey = openAIKey
		}
		if openAIURL != "" {
			cfg.OpenAIURL = openAIURL
		}
		if openAIModel != "" {
			cfg.OpenAIModel = openAIModel
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
