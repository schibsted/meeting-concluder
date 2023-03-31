package concluder

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/xyproto/env/v2"
)

type APIConfig struct {
	OpenAIKey    string `toml:"openai_api_key"`
	SlackWebhook string `toml:"slack_webhook"`
}

// global configuration
var Config APIConfig

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting user's home directory: %v", err)
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "concluder.toml")
	data, err := os.ReadFile(configPath)
	if err == nil {
		err = toml.Unmarshal(data, &Config)
		if err != nil {
			log.Fatalf("Error unmarshaling TOML: %v\n", err)
		}
	}

	if val := env.StrAlt("OPENAI_API_KEY", "OPENAI_KEY", ""); val != "" {
		Config.OpenAIKey = val
	}
	if val := env.Str("SLACK_WEBHOOK_URL"); val != "" {
		Config.SlackWebhook = val
	}

	if Config.SlackWebhook == "" {
		log.Println("WARNING: openai_api_key must be set in ~/.config/concluder.toml or as $OPENAI_API_KEY or $OPENAI_KEY")
	}
	if Config.SlackWebhook == "" {
		log.Println("WARNING: slack_webhook must be set in ~/.config/concluder.toml or as $SLACK_WEBHOOK_URL")
	}
}
