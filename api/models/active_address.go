package models

import "time"

type ActiveAddress struct {
	UpdatedAt time.Time `json:"last_updated_at" firestore:"updated_at"`
}
