package main

import (
	"topten/spotify/auth"
	"topten/spotify/user"

	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

func init() {
	// loads values from .env into the system.
	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found!")
	}
}

func main() {
	authToken, authErr := auth.GetAuthToken()
	if authErr != nil {
		logger.Fatal(authErr)
	}

	updateErr := user.UpdatePlaylist(authToken.Bearer())
	if updateErr != nil {
		logger.Fatal(updateErr)
	}

	logger.Info("Successfully updated playlist.")
}
