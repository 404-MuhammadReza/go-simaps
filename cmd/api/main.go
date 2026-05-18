package main

import (
	"log"

	"go-simaps/internal/config"
	"go-simaps/internal/handler"
	"go-simaps/internal/repository"
	"go-simaps/internal/router"
	"go-simaps/internal/service"
)

func main() {
	baseConfig := config.LoadConfig()
	mongoConfig := config.MongoConfig(baseConfig)
	minioConfig := config.MinioConfig(baseConfig)
	azureConfig := config.AzureConfig(baseConfig)
	googleConfig := config.GoogleConfig(baseConfig)

	
	storageRepository := repository.NewStorageRepository(minioConfig, baseConfig.MinioBucket)

	geocodingService := service.NewGeocodingService(googleConfig)
	geocodingHandler := handler.NewGeocodingHandler(geocodingService)

	logRepository := repository.NewLogRepository(mongoConfig)
	loggingService := service.NewLoggingService(logRepository)
	loggingHandler := handler.NewLoggingHandler(loggingService)

	userRepository := repository.NewUserRepository(mongoConfig)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userRepository, baseConfig, azureConfig)
	authHandler := handler.NewAuthHandler(authService, baseConfig)

	employeeRepository := repository.NewEmployeeRepository(mongoConfig)
	employeeService := service.NewEmployeeService(employeeRepository, storageRepository, geocodingService, loggingService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	hospitalRepository := repository.NewHospitalRepository(mongoConfig)
	hospitalService := service.NewHospitalService(hospitalRepository, geocodingService, loggingService)
	hospitalHandler := handler.NewHospitalHandler(hospitalService)

	handler := router.Handler{
		User:     userHandler,
		Auth:     authHandler,
		Hospital: hospitalHandler,
		Employee: employeeHandler,
		Logging:  loggingHandler,
		Geocoding: geocodingHandler,
	}

	r := router.Setup(handler, baseConfig)
	if err := r.Run(baseConfig.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
