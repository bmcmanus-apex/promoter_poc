package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func main() {
	// Read values from environment variables
	projectID := os.Getenv("GCP_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC_NAME")
	keyFilePath := os.Getenv("GCP_KEY_FILE_PATH")

	// Check if environment variables are set
	if projectID == "" || topicName == "" || keyFilePath == "" {
		log.Fatalf("Ensure all environment variables (GCP_PROJECT_ID, PUBSUB_TOPIC_NAME, GCP_KEY_FILE_PATH) are set.")
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

	message := "Hello, AFS!"

	// Publish a message
	log.Print("Setting topic")
	topic := client.Topic(topicName)
	log.Print("Topic set")
	defer topic.Stop()

	publishResult := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
		Attributes: map[string]string{
			"origin": "local",
		},
	})
	log.Print("Message published")

	// Do we need this?
	msgID, err := publishResult.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
	fmt.Printf("Published message with ID: %v\n", msgID)
}
