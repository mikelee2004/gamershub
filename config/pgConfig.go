package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHOST     string
	DBPORT     string
	DBNAME     string
	DBUSER     string
	DBPASSWORD string
	JWTSecret  string
	JWTExpire  int
}

func LoadConfig() *Config {
	ttl, _ := strconv.Atoi(os.Getenv("JWWT_TTL"))
	if ttl == 0 {
		ttl = 24
	}
	return &Config{
		DBHOST:     os.Getenv("DBHOST"),
		DBPORT:     os.Getenv("DBPORT"),
		DBNAME:     os.Getenv("DBNAME"),
		DBUSER:     os.Getenv("DBUSER"),
		DBPASSWORD: os.Getenv("DBPASSWORD"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JWTExpire:  ttl,
	}
}
