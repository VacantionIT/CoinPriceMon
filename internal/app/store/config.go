package store

// Config - настройки базы
type Config struct {
	DatabaseURL string `toml:"database_url"`
}

// NewConfig - Создание новой кофигурации БД
func NewConfig() *Config {
	return &Config{}
}
