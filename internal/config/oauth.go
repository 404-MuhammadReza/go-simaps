package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func AzureConfig(config *Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.AzureClientID,
		ClientSecret: config.AzureSecret,
		RedirectURL:  config.AzureRedirect,
		Scopes:       []string{"User.Read"},
		Endpoint:     microsoft.AzureADEndpoint(config.AzureTenantID),
	}
}