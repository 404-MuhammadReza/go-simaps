package config

import (
	"log"

	"googlemaps.github.io/maps"
)

func GoogleConfig(config *Config) *maps.Client {
	client, err := maps.NewClient(maps.WithAPIKey(config.GoogleAPIKey))
	if err != nil { log.Fatal("Failed to create Google Maps client: ", err) }

	return client
}