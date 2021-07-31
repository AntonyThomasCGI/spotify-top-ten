package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// getTopSongs returns a numSongs length TopSongs object containing the users top songs from the last 4 weeks.
func getTopSongs(numSongs int, auth *[]string) ([]*Song, error) {
	topTrackEndpoint := fmt.Sprintf(
		"https://api.spotify.com/v1/me/top/tracks?time_range=short_term&limit=%d",
		numSongs,
	)

	request, reqErr := http.NewRequest("GET", topTrackEndpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	request.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": *auth,
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

	respJson := struct {
		Items []*Song `json:"items"`
	}{}
	decodeErr := json.NewDecoder(response.Body).Decode(&respJson)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return respJson.Items, nil
}
