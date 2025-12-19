package config

import "os"

type Config struct {
	DB_DSN     string
	BcryptCost int // например, 12
	Addr       string
}

func Load() Config {
	cost := 12
	if v := os.Getenv("BCRYPT_COST"); v != "" {
		// необязательно: распарсить int, при ошибке оставить 12
	}
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	dns_bd := os.Getenv("DB_DSN")
	if dns_bd == "" {
		dns_bd = "postgres://postgres:postgres@localhost:5432/prak_9?sslmode=disable"
	}
	return Config{
		DB_DSN:     dns_bd,
		BcryptCost: cost,
		Addr:       addr,
	}
}
