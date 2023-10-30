package main

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func main() {
	// Read values from environment variables
	projectID := os.Getenv("GCP_PROJECT_ID")
	subscriptionName := "sub"
	keyFilePath := os.Getenv("GCP_KEY_FILE_PATH")

	// Check if environment variables are set
	if projectID == "" || subscriptionName == "" || keyFilePath == "" {
		log.Fatalf("Ensure all environment variables (GCP_PROJECT_ID, PUBSUB_SUBSCRIPTION_NAME, GCP_KEY_FILE_PATH) are set.")
	} else {
		log.Print("Found all environment variables")
	}

	// Set up the context and Pub/Sub client
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(keyFilePath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	} else {
		log.Print("Client created")
	}
	defer client.Close()

	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error getting subscriptions: %v", err)
		}
		fmt.Println(s)
	}

	// Create a subscription object
	sub := client.SubscriptionInProject(projectID, subscriptionName)
	log.Print("Bound subscription")

	// Pull messages in a loop
	//for {
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// Handle the received message
		log.Printf("Received message: %s\n", string(msg.Data))

		// Acknowledge the message to mark it as processed
		msg.Ack()
	})
	if err != nil {
		log.Printf("Error receiving message: %v", err)
		// You can add error handling and retry logic here if needed
	}

	// Add a short delay to control the rate of message pulling
	time.Sleep(1 * time.Second)
	//}
}
