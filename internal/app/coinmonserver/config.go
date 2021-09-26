package coinmonserver

import "github.com/VacantionIT/coin-price-mon/internal/app/store"

// Config - струкура с настройками сервера, с указанием ключей в настроечном файле
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	HashSalt    string `toml:"hash_salt"`
	ExpDuration int    `toml:"token_exp_duration"`
	TokenKey    string `toml:"token_key"`
	Store       *store.Config
}

// NewConfig - создание конфига с параметрами по умолчанию
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		HashSalt:    "some_hash_salt",
		ExpDuration: 86400,
		TokenKey:    "skey_for_token",
		Store:       store.NewConfig(),
	}
}
