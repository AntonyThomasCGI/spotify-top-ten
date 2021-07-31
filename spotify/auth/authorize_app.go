// Does the initial authorization of the app for a user to allow access.
package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

const (
	authorizeURL = "https://accounts.spotify.com/authorize"
)

var scopes = []string{"user-top-read", "playlist-modify-public", "playlist-modify-private", "playlist-read-private"}

// getAuthCode follows the spotify auth flow to get an auth code:
// https://developer.spotify.com/documentation/general/guides/authorization-guide/
func getAuthCode(spotifyId string, port string) (*string, error) {
	codeCh, serverErr := startRedirectServer(port)
	if serverErr != nil {
		return nil, serverErr
	}
	url, parseErr := url.Parse(authorizeURL)
	if parseErr != nil {
		return nil, parseErr
	}
	q := url.Query()
	q.Add("response_type", "code")
	q.Add("client_id", spotifyId)
	q.Add("redirect_uri", "http://localhost:"+port)
	q.Add("scope", strings.Join(scopes, " "))
	url.RawQuery = q.Encode()

	fmt.Println("Please authorize this app via browser.")
	openErr := openURL(url.String())
	if openErr != nil {
		return nil, openErr
	}

	code := <-codeCh

	return &code, nil
}

// startRedirectServer starts a local server to obtain the redirected code from spotify.
func startRedirectServer(port string) (code chan string, err error) {
	code = make(chan string)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code <- r.FormValue("code")
		})
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	return code, nil
}

// openURL opens the spotify authorization url in a browser window.
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}
