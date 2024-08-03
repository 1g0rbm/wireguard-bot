package env

import "os"

func readEnvAsString(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return value
}
