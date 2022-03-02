package config

import (
	"flag"
	"fmt"
	"os")

const (
	EnvDbHost = "DB_HOST"
	EnvDbName = "DB_NAME"
	EnvDbPort = "DB_PORT"
	EnvDbUser = "DB_USER"
	EnvDbPswd = "DB_PASSWORD"
)

type Config struct {
	DbHost string
	DbName string
	DbPort string
	DbUser string
	DbPswd string
}

func GetFromEnv() *Config {
	conf := &Config{}

	flag.StringVar(&conf.DbHost, EnvDbHost, os.Getenv(EnvDbHost), "database host")
	flag.StringVar(&conf.DbName, EnvDbName, os.Getenv(EnvDbName), "database name")
	flag.StringVar(&conf.DbPort, EnvDbPort, os.Getenv(EnvDbPort), "database port")
	flag.StringVar(&conf.DbUser, EnvDbUser, os.Getenv(EnvDbUser), "database user name")
	flag.StringVar(&conf.DbPswd, EnvDbPswd, os.Getenv(EnvDbPswd), "database user password")

	flag.Parse()
	return conf
}

func (cfg *Config) ConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPswd,
		cfg.DbName,
		cfg.DbPort,
	)
}
