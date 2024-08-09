package config

type PGConfig interface {
	DSN() string
}

type BotConfig interface {
	Token() string
}

type LoggerConfig interface {
	LogFilepath() string
}
