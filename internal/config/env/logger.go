package env

const logFilepath = "LOG_FILEPATH"

const defaultLogFilePath = "bot.log"

type LoggerConfig struct {
	logFilepath string
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		logFilepath: readEnvAsString(logFilepath, defaultLogFilePath),
	}
}

func (l *LoggerConfig) LogFilepath() string {
	return l.logFilepath
}
