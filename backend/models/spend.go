// models/spend.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Spend represents the data structure for spend analysis
type Spend struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID string             `json:"user_id" bson:"user_id"`
	Date   time.Time          `json:"date" bson:"date"`
	Amount float64            `json:"amount" bson:"amount"`
}
