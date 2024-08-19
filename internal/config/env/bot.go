package env

const botToken = "BOT_TOKEN"
const defaultBotToken = ""

type BotConfig struct {
	token string
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		token: readEnvAsString(botToken, defaultBotToken),
	}
}

func (c *BotConfig) Token() string {
	return c.token
}
