package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// authFromCredentialFile attempts to get auth tokens from a file.
func authFromCredentialFile() (*SpotifyAuth, error) {
	homeDir, dirErr := os.UserHomeDir()
	if dirErr != nil {
		return nil, dirErr
	}

	credsPath := filepath.Join(homeDir, ".spotify", "auth.yaml")

	credsFile, openErr := os.Open(credsPath)
	if openErr != nil {
		return nil, openErr
	}
	defer credsFile.Close()

	credsBytes, readErr := ioutil.ReadAll(credsFile)
	if readErr != nil {
		return nil, readErr
	}

	auth := &SpotifyAuth{}
	unmarshalErr := yaml.Unmarshal(credsBytes, auth)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return auth, nil
}

// writeCrednetialFile attempts to write a yaml file of auth tokens.
func writeCredentialFile(auth *SpotifyAuth) error {
	homeDir, dirErr := os.UserHomeDir()
	if dirErr != nil {
		return dirErr
	}

	credsPath := filepath.Join(homeDir, ".spotify", "auth.yaml")

	credsJson, marshalErr := yaml.Marshal(auth)
	if marshalErr != nil {
		return marshalErr
	}

	// Check if the .spotify dir exists, if not create.
	baseDir := filepath.Dir(credsPath)
	if _, statErr := os.Stat(baseDir); os.IsNotExist(statErr) {
		dirErr := os.Mkdir(baseDir, 0700)
		if dirErr != nil {
			return (dirErr)
		}
	}

	// Write the file.
	credsFile, createErr := os.Create(credsPath)
	if createErr != nil {
		return createErr
	}
	defer credsFile.Close()

	_, writeErr := credsFile.Write(credsJson)
	if writeErr != nil {
		return writeErr
	}
	// Read write to user only.
	chmodErr := credsFile.Chmod(0600)
	if chmodErr != nil {
		return chmodErr
	}

	return nil
}
