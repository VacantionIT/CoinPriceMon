package store

import "go.mongodb.org/mongo-driver/bson/primitive"

// Monitoring - структура с описанием параметров мониторинга
type Monitoring struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CoinID   string             `json:"coinid,omitempty" bson:"coinid,omitempty"`
	Interval int                `json:"interval,omitempty" bson:"interval,omitempty"`
}
