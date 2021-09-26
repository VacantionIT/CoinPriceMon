package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store - Структура для работы с БД
type Store struct {
	config   *Config
	DBClient *mongo.Client
}

// New - Функция создания структуры для работы с БД
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open - Открытие и проверка соединения с БД
func (s *Store) Open() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(s.config.DatabaseURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	s.DBClient = client

	return nil
}

// Close - Закрытие соединения с БД
func (s *Store) Close() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	s.DBClient.Disconnect(ctx)
}
