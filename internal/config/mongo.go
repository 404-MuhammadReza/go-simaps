package config

import (
	"log"
	"time"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConfig(config *Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil { log.Fatal("Can't Connect to MongoDB: ", err) }

	err = client.Ping(ctx, nil)
	if err != nil { log.Fatal("Can't Ping MongoDB: ", err) }
	
	db := client.Database(config.MongoName)
	return db
}