package models

type Config struct {
	DBDriver string `envconfig:"DB_DRIVER" default:"postgres"`
	DBUser   string `envconfig:"DB_USER" default:"postgres"`
	DBPass   string `envconfig:"DB_PASS" default:"1903"`
	DBName   string `envconfig:"DB_NAME" default:"CLI_Wikis"`
	DBHost   string `envconfig:"DB_HOST" default:"localhost"`
	DBPort   int    `envconfig:"DB_PORT" default:"5432"`
}
