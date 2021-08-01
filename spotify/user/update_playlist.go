package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

type Song struct {
	URI     string    `json:"uri"`
	Name    string    `json:"name"`
	Artists []*Artist `json:"artists"`
}

type Artist struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// UpdatePlaylist gets a users top songs and replaces the songs in the target playlist.
func UpdatePlaylist(auth *[]string) error {
	// Get env variables.
	playlistId := os.Getenv("PLAYLIST_ID")
	numTopSongs, convErr := strconv.Atoi(os.Getenv("N_TOP_SONGS"))
	if convErr != nil {
		errF := fmt.Errorf("Could not convert N_TOP_SONGS env variable: %v", convErr)
		return errF
	}

	topSongs, topErr := getTopSongs(numTopSongs, auth)
	if topErr != nil {
		return topErr
	}

	replaceErr := replacePlaylistSongs(playlistId, topSongs, auth)
	if replaceErr != nil {
		return replaceErr
	}

	now := time.Now()
	description := fmt.Sprintf("--- last updated: %s --- beep boop im a bot.", now.Format("02/01 03:04pm"))

	descriptionErr := updateDescription(playlistId, description, auth)
	if descriptionErr != nil {
		logger.Error(descriptionErr)
	}

	return nil
}

// replacePlaylistSongs takes a list of songs and replaces the songs in the playlist.
func replacePlaylistSongs(playlistId string, songs []*Song, auth *[]string) error {
	uris := []string{}
	for _, song := range songs {
		uris = append(uris, song.URI)
	}

	replaceEndpoint := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s/tracks?uris=%s",
		playlistId,
		strings.Join(uris, ","),
	)
	request, reqErr := http.NewRequest("PUT", replaceEndpoint, nil)
	if reqErr != nil {
		return reqErr
	}
	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": *auth,
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		return respErr
	}
	defer response.Body.Close()

	if response.StatusCode > 299 { // 201 considered success.
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			errNoBody := fmt.Errorf("Non-OK HTTP status: %d", response.StatusCode)
			return errNoBody
		}
		errF := fmt.Errorf("Non-OK HTTP status: %d Body: %s", response.StatusCode, string(body))
		return errF
	}

	return nil
}

// updateDescription updates the playlists description.
func updateDescription(playlistId string, description string, auth *[]string) error {
	playlistEndpoint := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s",
		playlistId,
	)

	body, marshalErr := json.Marshal(map[string]string{
		"description": description,
	})
	if marshalErr != nil {
		return marshalErr
	}

	request, reqErr := http.NewRequest("PUT", playlistEndpoint, bytes.NewBuffer(body))
	if reqErr != nil {
		return reqErr
	}
	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": *auth,
	}

	response, respErr := client.Do(request)
	if respErr != nil {
		return respErr
	}
	defer response.Body.Close()

	if response.StatusCode > 299 { // 201 considered success.
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			errNoBody := fmt.Errorf("Non-OK HTTP status: %d", response.StatusCode)
			return errNoBody
		}
		errF := fmt.Errorf("Non-OK HTTP status: %d Body: %s", response.StatusCode, string(body))
		return errF
	}

	return nil
}
