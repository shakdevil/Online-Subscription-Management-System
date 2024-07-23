// main.go
package main

import (
    "fmt"
    "net/http"

    "backend/db"
    "backend/handlers"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    db.ConnectDB()
    defer func() {
        if db.Client != nil {
            db.Client.Disconnect(db.ctx)
            fmt.Println("Disconnected from MongoDB.")
        }
    }()

    // Set up CORS middleware
    corsHandler := cors.Default().Handler

    // Create a new router and register handlers
    router := mux.NewRouter()
    router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
    router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    router.Handle("/add-subscription", corsHandler(http.HandlerFunc(handlers.AddSubscriptionHandler))).Methods("POST")
    router.Handle("/get-subscriptions", corsHandler(http.HandlerFunc(handlers.GetSubscriptionsHandler))).Methods("GET")
    router.Handle("/add-spend", corsHandler(http.HandlerFunc(handlers.AddSpendHandler))).Methods("POST")
    router.Handle("/get-graph-data", corsHandler(http.HandlerFunc(handlers.GetGraphDataHandler))).Methods("GET")

    // Start HTTP server with CORS middleware and router
    fmt.Println("Server listening on :8080...")
    http.ListenAndServe(":8080", router)
}
