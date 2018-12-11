package models

import "time"

type Stats struct {
	UpdatedAt                    time.Time `json:"updated_at" bson:"updated_at"`
	NumberOfTotalTransactions    int64     `json:"total_transactions_count" bson:"total_transactions_count"`
	NumberOfLastWeekTransactions int64     `json:"last_week_transactions_count" bson:"last_week_transactions_count"`
	NumberOfLastDayTransactions  int64     `json:"last_day_transactions_count" bson:"last_day_transactions_count"`
}
