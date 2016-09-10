package lib

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetGoogleClient(ctx context.Context, config *oauth2.Config) *http.Client {

	// cacheFile is the file path for credential file.
	cacheFile, err := getFilepathTokenCache()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	// tok is a token Object that reads cacheFile (credential file)
	tok, err := getGoogleTokenFromFile(cacheFile)
	if err != nil {
		// prompts users to retrieve a token from browser.
		// retruns a oauth2.Token object.
		tok = getGoogleTokenFromWeb(config)

		// save a oauth2.Token object to the file path.
		saveGoogleToken(cacheFile, tok)
	}

	// Client returns an HTTP client using the provided token. The token
	// will auto-refresh as necessary. The underlying HTTP transport will
	// be obtained using the provided context. The returned client and
	// its Transport should not be modified.
	return config.Client(ctx, tok)
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func getFilepathTokenCache() (string, error) {

	// user.Current returns the current user.
	// usr is type User
	/*
		type User struct {
			Uid string // user ID
			Gid string // primary group ID
			Username string
			Name string
			HomeDir string
		}
	*/
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// make "$HOME/.credentials hidden directory."
	// filepath.Join creates string name of hidden directory.
	// tokenCacheDir string
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")

	// os.MkdirAll creates directory.
	os.MkdirAll(tokenCacheDir, 0700)

	// url.QueryEscape escapes the string so it can be safely placed
	// inside a URL query.
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func getGoogleTokenFromFile(file string) (*oauth2.Token, error) {

	// os.Open opens the credential file.
	// creates File object
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// initialize Token object
	t := &oauth2.Token{}

	// json.NewDecoder returns a new decoder that reads from r.
	// the decoder introduces its own buddering and may read data
	// from r beyond the JSON values requested.
	d := json.NewDecoder(f)

	// Decode reads the next JSON-encoded value from its input
	// and stores it in the value pointed to by v.
	err = d.Decode(t)
	defer f.Close()
	return t, err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token object.
func getGoogleTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
	// that asks for permissions for the required scopes explicitly.
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	// config.Exchange converts an authorization code into a Token object.
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveGoogleToken(filepath string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", filepath)
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	// json.NewEncoder returns a new encoder that writes to f, io.Writer.
	e := json.NewEncoder(f)

	// Encode writes the JSON encoding of v to the stream, following by a new
	// character.
	e.Encode(token)
}
