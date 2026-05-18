package service

import (
	"log"
	"fmt"
	"context"
	
	"go-simaps/internal/model"
	"go-simaps/internal/apperror"

	"googlemaps.github.io/maps"
)

type GeocodingService interface {
	ToCoordinate(ctx context.Context, address string) (*model.Coordinate, error)
	ToAddress(ctx context.Context, coordinate *model.Coordinate) (string, error)
}

type geocodingService struct {
	mapsClient *maps.Client
}

func NewGeocodingService(mapsClient *maps.Client) GeocodingService {
	return &geocodingService{mapsClient}
}

func (s *geocodingService) ToCoordinate(ctx context.Context, address string) (*model.Coordinate, error) {
	log.Printf("Geocoding Running for Address: %s", address)

	request := &maps.GeocodingRequest{Address: address}
	response, err := s.mapsClient.Geocode(ctx, request)
	if err != nil {
		return nil, apperror.ErrExternal("Failed to geocode address: " + err.Error())
	}

	if len(response) == 0 {
		return nil, apperror.ErrNotFound("Coordinates not found for address: " + address)
	}
	lat := response[0].Geometry.Location.Lat
	lng := response[0].Geometry.Location.Lng

	return &model.Coordinate{Lat: lat, Lng: lng}, nil
}

func (s *geocodingService) ToAddress(ctx context.Context, coordinate *model.Coordinate) (string, error) {
	log.Printf("Reverse Geocoding Running for Coordinates: (%f, %f)", coordinate.Lat, coordinate.Lng)

	request := &maps.GeocodingRequest{LatLng: &maps.LatLng{Lat: coordinate.Lat, Lng: coordinate.Lng}}
	response, err := s.mapsClient.ReverseGeocode(ctx, request)
	if err != nil {
		return "", apperror.ErrExternal("Failed to reverse geocode coordinates: " + err.Error())
	}

	if len(response) == 0 {
		message := fmt.Sprintf("Address not found for coordinates: (%f, %f)", coordinate.Lat, coordinate.Lng)
		return "", apperror.ErrNotFound(message)
	}
	
	return response[0].FormattedAddress, nil
}
