package models

type Config struct {
	DBDriver string `envconfig:"DB_DRIVER" default:"postgres"`
	DBUser   string `envconfig:"DB_USER" default:"postgres"`
	DBPass   string `envconfig:"DB_PASS" default:"postgres"`
	DBName   string `envconfig:"DB_NAME" default:"cli_wikis"`
	DBHost   string `envconfig:"DB_HOST" default:"localhost"`
	DBPort   int    `envconfig:"DB_PORT" default:"5432"`
}
