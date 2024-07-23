// models/graph.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Graph represents the data structure for graph data
type Graph struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Month     string             `json:"month" bson:"month"`
	Amounts   float64            `json:"totalAmount" bson:"totalAmount"`
	TotalCost float64            `json:"totalCost" bson:"totalCost"`
}
