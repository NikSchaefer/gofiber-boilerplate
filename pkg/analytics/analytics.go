package analytics

import (
	"fmt"
	"log"

	"github.com/NikSchaefer/go-fiber/config"
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/posthog/posthog-go"
)

// analyticsClient is the global PostHog client instance
var analyticsClient posthog.Client

var isProduction bool = false

// InitAnalytics initializes the PostHog client with the configured token
func InitAnalytics() {
	key := config.GetPosthogKey()

	client, err := posthog.NewWithConfig(key,
		posthog.Config{Endpoint: "https://us.i.posthog.com"},
	)
	if err != nil {
		log.Fatal("Failed to initialize PostHog client:", err)
	}

	fmt.Println("Initializing PostHog client")
	isProduction = config.GetIsProduction()

	analyticsClient = client
}

// TrackEvent sends an event to PostHog with the given properties
// Returns an error if the tracking failed
func TrackEvent(name string, properties map[string]interface{}, distinctID string) error {
	if !isProduction {
		return nil
	}

	if analyticsClient == nil {
		return fmt.Errorf("analytics client not initialized")
	}

	return analyticsClient.Enqueue(posthog.Capture{
		DistinctId: distinctID,
		Event:      name,
		Properties: properties,
	})
}

// TrackEventWithUser sends an event to PostHog associated with a specific user
// Returns an error if the tracking failed or nil if the user is blacklisted
func TrackEventWithUser(name string, properties map[string]interface{}, user *ent.User) error {
	if !isProduction {
		return nil
	}

	if analyticsClient == nil {
		return fmt.Errorf("analytics client not initialized")
	}

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	distinctID := user.ID.String()
	return TrackEvent(name, properties, distinctID)
}

// IdentifyUser updates user properties in PostHog
func IdentifyUser(distinctID string, properties map[string]interface{}) error {
	if !isProduction {
		return nil
	}

	if analyticsClient == nil {
		return fmt.Errorf("analytics client not initialized")
	}

	return analyticsClient.Enqueue(posthog.Identify{
		DistinctId: distinctID,
		Properties: properties,
	})
}
