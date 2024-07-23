// db/mongodb.go
package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    Client *mongo.Client
    ctx    context.Context
)

func ConnectDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    var err error
    Client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = Client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    ctx = context.TODO()
}
