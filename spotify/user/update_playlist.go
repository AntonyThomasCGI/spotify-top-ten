package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
