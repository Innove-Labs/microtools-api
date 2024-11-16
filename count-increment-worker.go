package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func updateApiHitCounts() {
	log.Println("Starting API hit count worker...")
	ticker := time.NewTicker(1 * time.Minute) // Update every minute
	defer ticker.Stop()

	for {
		<-ticker.C // Wait for the next tick

		// Update the count in MongoDB
		err := incrementApiHitCount()
		if err != nil {
			log.Printf("Error updating count: %v", err)
		} else {
			fmt.Println("API hit count updated successfully.")
		}
	}
}

func incrementApiHitCount() error {
	// Create a filter to find the document (you can have a single document for counts)
	filter := bson.M{"name": "api_hit_count"}

	// Update the count by incrementing it
	update := bson.M{
		"$inc": bson.M{"count": 1}, // Increment the 'count' field by 1
	}

	// Perform the update
	collection := MongoClient.Database("microapis").Collection("api_analytics")
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment count: %v", err)
	}

	return nil
}
