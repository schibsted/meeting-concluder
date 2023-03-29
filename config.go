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
	OpenAI_APIKey string `toml:"openai_api_key"`
	Slack_Webhook string `toml:"slack_webhook"`
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
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error unmarshaling TOML: %v", err)
	}

	if val := env.StrAlt("OPENAI_API_KEY", "OPENAI_KEY", ""); val != "" {
		Config.OpenAI_APIKey = val
	}
	if val := env.Str("SLACK_WEBHOOK_URL"); val != "" {
		Config.Slack_Webhook = val
	}

	if Config.Slack_Webhook == "" {
		log.Println("openai_api_key can be set in ~/.config/concluder.toml or as $OPENAI_API_KEY")
	}
	if Config.Slack_Webhook == "" {
		log.Println("slack_webhook can be set in ~/.config/concluder.toml or as $SLACK_WEBHOOK_URL")
	}
}
