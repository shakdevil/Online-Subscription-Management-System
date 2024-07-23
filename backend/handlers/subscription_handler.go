// handlers/subscription_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/db"
	"backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

// AddSubscriptionHandler handles the endpoint for adding a new subscription.
func AddSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var newSubscription models.Subscription

	// Decode JSON request body into Subscription struct
	err := json.NewDecoder(r.Body).Decode(&newSubscription)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate user permissions and ensure the user exists
	// For simplicity, assume user ID is provided in the request body.
	newSubscription.UserID = "user_id_here"

	// Set default values or perform additional validation if needed
	newSubscription.CreatedAt = time.Now()

	// Insert the new subscription into the database
	_, err = db.SubscribeCollection.InsertOne(db.Ctx, newSubscription)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add subscription: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSubscriptionsHandler handles the endpoint for fetching user subscriptions.
func GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from authentication token or session
	// For simplicity, assume user ID is provided in the request body.
	userID := "user_id_here"

	cursor, err := db.SubscribeCollection.Find(db.Ctx, bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch subscriptions: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(db.Ctx)

	var subscriptions []models.Subscription
	if err := cursor.All(db.Ctx, &subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(subscriptions)
	if err != nil {
		http.Error(w, "Failed to serialize subscriptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
