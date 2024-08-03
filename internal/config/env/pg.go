package env

const PgDsn = "PG_DSN"

const defaultPgDSN = "host=localhost port=54321 dbname=wireguard user=user password=qwerty sslmode=disable"

type PgConfig struct {
	dsn string
}

func NewPgConfig() *PgConfig {
	return &PgConfig{
		dsn: readEnvAsString(PgDsn, defaultPgDSN),
	}
}

func (c *PgConfig) DSN() string {
	return c.dsn
}
