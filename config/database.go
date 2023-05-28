package config

import "os"

type Database interface{}

type PsqlDbConnection struct {
	DbHost     string
	DbPort     string
	DbDatabase string
	DbUsername string
	DbPassword string
}

type DatabaseConfig struct {
	Psql PsqlDbConnection
}

func DatabaseNew() Database {
	return &DatabaseConfig{
		Psql: PsqlDbConnection{
			DbHost:     os.Getenv("DB_HOST"),
			DbPort:     os.Getenv("DB_PORT"),
			DbDatabase: os.Getenv("DB_DATABASE"),
			DbUsername: os.Getenv("DB_USERNAME"),
			DbPassword: os.Getenv("DB_PASSWORD"),
		},
	}
}
