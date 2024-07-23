// models/subscription.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscription represents a subscription entity
type Subscription struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      string             `json:"user_id" bson:"user_id"`
	Service     string             `json:"service" bson:"service"`
	Amount      float64            `json:"amount" bson:"amount"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	IsAutoRenew bool               `json:"is_auto_renew" bson:"is_auto_renew"`
}
