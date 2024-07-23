// handlers/graph_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/db"
	"backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

// GetGraphDataHandler handles the endpoint for fetching graph data.
func GetGraphDataHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from authentication token or session
	// For simplicity, assume user ID is provided in the request body.
	userID := "user_id_here"

	// TODO: Add logic to fetch and aggregate graph data based on user ID.

	// Mocked data for illustration purposes
	mockGraphData := []models.Graph{
		{
			Month:     "January",
			Amounts:   1000.0,
			TotalCost: 500.0,
		},
		{
			Month:     "February",
			Amounts:   1500.0,
			TotalCost: 700.0,
		},
	}


	// Define the pipeline to group and sum the amounts by month for subscriptions
	subscriptionPipeline := []bson.D{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "month", Value: bson.D{{Key: "$month", Value: "$date"}}},
			}},
			{Key: "totalAmount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "month", Value: "$_id.month"},
			{Key: "totalAmount", Value: 1},
		}}},
	}

	// Aggregate data using pipeline for subscriptions
	subscriptionCursor, err := db.SubscribeCollection.Aggregate(db.Ctx, subscriptionPipeline)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error aggregating subscription data: %v", err), http.StatusInternalServerError)
		return
	}
	defer subscriptionCursor.Close(db.Ctx)

	// Store aggregated data in graph collection for subscriptions
	var subscriptionGraphData []models.Graph
	for subscriptionCursor.Next(db.Ctx) {
		var result models.Graph
		err := subscriptionCursor.Decode(&result)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding subscription cursor data: %v", err), http.StatusInternalServerError)
			return
		}
		subscriptionGraphData = append(subscriptionGraphData, result)
	}

	// TODO: Add similar pipeline for spend data and merge the results if needed

	response, err := json.Marshal(subscriptionGraphData)
	if err != nil {
		http.Error(w, "Failed to serialize graph data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
