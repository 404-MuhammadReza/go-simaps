package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort	       	 string
	AppDomain		 string
	FrontendURL    	 string

	GoogleAPIKey   	 string

	MongoURI       	 string
	MongoName      	 string

	MinioEndpoint  	 string
	MinioAccessKey 	 string
	MinioSecretKey 	 string
	MinioBucket	   	 string
	MinioUseSSL	   	 bool

	AzureClientID  	 string
	AzureSecret    	 string
	AzureTenantID  	 string
	AzureRedirect  	 string

	JWTSecret      	 []byte
	JWTIssuer      	 string
	JWTAccessExpiry  int
	JWTRefreshExpiry int

	JWTAccessMaxAge	 int
	JWTRefreshMaxAge int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	return &Config{
		AppPort: getEnv("APP_PORT"),
		AppDomain: getEnv("APP_DOMAIN"),
		FrontendURL: getEnv("FRONTEND_URL"),

		GoogleAPIKey: getEnv("GOOGLE_API_KEY"),

		MongoURI: getEnv("MONGODB_URI"),
		MongoName: getEnv("MONGODB_NAME"),

		MinioEndpoint: getEnv("MINIO_ENDPOINT"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY"),
		MinioBucket: getEnv("MINIO_BUCKET"),
		MinioUseSSL: getEnvBool("MINIO_USE_SSL"),

		AzureClientID: getEnv("AZURE_CLIENT_ID"),
		AzureSecret: getEnv("AZURE_CLIENT_SECRET"),
		AzureTenantID: getEnv("AZURE_TENANT_ID"),
		AzureRedirect: getEnv("AZURE_REDIRECT_URL"),

		JWTSecret: []byte(getEnv("JWT_SECRET")),
		JWTIssuer: getEnv("JWT_ISSUER"),
		JWTAccessExpiry: getEnvInt("JWT_ACCESS_EXPIRY"), // In minutes
		JWTRefreshExpiry: getEnvInt("JWT_REFRESH_EXPIRY"), // In hours

		JWTAccessMaxAge: getEnvInt("JWT_ACCESS_EXPIRY") * 60, // Convert minutes to seconds
		JWTRefreshMaxAge: getEnvInt("JWT_REFRESH_EXPIRY") * 3600, // Convert hours to seconds
	}
}

// Helper: Get ENV (string)
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable '%s' is missing or empty.", key)
	}

	return value
}

// Helper: Get ENV (int)
func getEnvInt(key string) int {
	strValue := getEnv(key)
	value, err := strconv.Atoi(strValue)
	if err != nil {
		log.Fatalf("Environment variable '%s' must be a valid integer. Got: '%s'", key, strValue)
	}

	return value
}

// Helper: Get ENV (bool)
func getEnvBool(key string) bool {
	strValue := getEnv(key)
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		log.Fatalf("Environment variable '%s' must be 'true' or 'false'. Got: '%s'", key, strValue)
	}
	return value
}