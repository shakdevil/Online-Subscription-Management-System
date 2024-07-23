// handlers/spend_handler.go
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

// AddSpendHandler handles the endpoint for adding spend data.
func AddSpendHandler(w http.ResponseWriter, r *http.Request) {
	var newSpend models.Spend

	// Decode JSON request body into Spend struct
	err := json.NewDecoder(r.Body).Decode(&newSpend)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate user permissions and ensure the user exists
	// For simplicity, assume user ID is provided in the request body.
	newSpend.UserID = "user_id_here"

	// Set default values or perform additional validation if needed
	newSpend.Date = time.Now()

	// Insert the new spend data into the database
	_, err = db.SpendCollection.InsertOne(db.Ctx, newSpend)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add spend data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSpendHandler handles the endpoint for fetching spend data.
func GetSpendHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from authentication token or session
	// For simplicity, assume user ID is provided in the request body.
	userID := "user_id_here"

	cursor, err := db.SpendCollection.Find(db.Ctx, bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch spend data: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(db.Ctx)

	var spendData []models.Spend
	if err := cursor.All(db.Ctx, &spendData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode spend data: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(spendData)
	if err != nil {
		http.Error(w, "Failed to serialize spend data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
