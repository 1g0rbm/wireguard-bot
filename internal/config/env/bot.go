package env

const BotToken = "BOT_TOKEN"
const defaultBotToken = ""

type BotConfig struct {
	token string
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		token: readEnvAsString(BotToken, defaultBotToken),
	}
}

func (c *BotConfig) Token() string {
	return c.token
}
