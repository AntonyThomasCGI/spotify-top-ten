package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const (
	authenticationURL = "https://accounts.spotify.com/api/token"
)

type SpotifyAuth struct {
	AccessToken  string `json:"access_token" yaml:"access_token"`
	RefreshToken string `json:"refresh_token" yaml:"refresh_token"`
}

func (s SpotifyAuth) Bearer() *[]string {
	return &[]string{"Bearer " + s.AccessToken}
}

// GetAuthToken attempts to read access tokens from an auth file or authenticates a new user.
func GetAuthToken() (*SpotifyAuth, error) {
	spotifyId := os.Getenv("CLIENT_ID")
	spotifySecret := os.Getenv("CLIENT_SECRET")

	auth, fileErr := authFromCredentialFile()
	if fileErr != nil {
		logger.Debug("Failed to read from auth file:", fileErr)
	}
	if auth != nil { // Got auth from file.
		// TODO: First check if refresh required, for now just always refresh.
		refreshErr := refreshAuthToken(spotifyId, spotifySecret, auth)
		if refreshErr != nil {
			return nil, refreshErr
		}
		writeErr := writeCredentialFile(auth)
		if writeErr != nil {
			logger.Error("Could not write cred file:", writeErr)
		}

		return auth, nil
	}

	// Authenticate new user.
	port := os.Getenv("PORT")

	code, err := getAuthCode(spotifyId, port)
	if err != nil {
		return nil, err
	}

	newAuth, authErr := newAuthTokenFromCode(spotifyId, spotifySecret, code, port)
	if authErr != nil {
		return nil, authErr
	}

	writeErr := writeCredentialFile(newAuth)
	if writeErr != nil {
		logger.Error("Could not write cred file:", writeErr)
	}

	return newAuth, nil
}

// refreshAuthToken requests a refreshed access token using the refresh token.
func refreshAuthToken(spotifyId string, spotifySecret string, auth *SpotifyAuth) error {
	postBody := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{auth.RefreshToken},
		"client_id":     []string{spotifyId},
		"client_secret": []string{spotifySecret},
	}

	request, reqErr := http.NewRequest("POST", authenticationURL, strings.NewReader(postBody.Encode()))
	if reqErr != nil {
		return reqErr
	}
	request.Header = http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Accept":       []string{"application/json"},
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		return respErr
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			errNoBody := fmt.Errorf("Non-OK HTTP status: %d", response.StatusCode)
			return errNoBody
		}
		errF := fmt.Errorf("Non-OK HTTP status: %d Body: %s", response.StatusCode, string(body))
		return errF
	}

	decodeErr := json.NewDecoder(response.Body).Decode(auth)
	if decodeErr != nil {
		return decodeErr
	}

	return nil
}

// newAuthTokenFromCode requests a new auth token from the spotify accounts service with an authorization code.
func newAuthTokenFromCode(spotifyId string, spotifySecret string, code *string, port string) (*SpotifyAuth, error) {
	postBody := url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{*code},
		"redirect_uri":  []string{"http://localhost:" + port},
		"client_id":     []string{spotifyId},
		"client_secret": []string{spotifySecret},
	}

	request, reqErr := http.NewRequest("POST", authenticationURL, strings.NewReader(postBody.Encode()))
	if reqErr != nil {
		return nil, reqErr
	}
	request.Header = http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
		"Accept":       []string{"application/json"},
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		return nil, respErr
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			errNoBody := fmt.Errorf("Non-OK HTTP status: %d", response.StatusCode)
			return nil, errNoBody
		}
		errF := fmt.Errorf("Non-OK HTTP status: %d Body: %s", response.StatusCode, string(body))
		return nil, errF
	}

	auth := &SpotifyAuth{}
	decodeErr := json.NewDecoder(response.Body).Decode(auth)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return auth, nil
}
