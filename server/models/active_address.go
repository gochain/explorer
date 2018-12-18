package models

import "time"

type ActiveAddress struct {
	Address   string    `json:"address" bson:"address"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
