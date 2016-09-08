package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {

	// cacheFile is the file path for credential file.
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	// tok is a token Object that reads cacheFile (credential file)
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		// prompts users to retrieve a token from browser.
		// retruns a oauth2.Token object.
		tok = getTokenFromWeb(config)

		// save a oauth2.Token object to the file path.
		saveToken(cacheFile, tok)
	}

	// Client returns an HTTP client using the provided token. The token
	// will auto-refresh as necessary. The underlying HTTP transport will
	// be obtained using the provided context. The returned client and
	// its Transport should not be modified.
	return config.Client(ctx, tok)
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {

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
func tokenFromFile(file string) (*oauth2.Token, error) {

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
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

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
func saveToken(filepath string, token *oauth2.Token) {
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

func main() {

	// context.Background returns a non-nil, empty Context. It is never canceled,
	// has no values, and has no deadline.  It is typically used by the main function,
	// initialization, and tests, and as the top-level Context for incoming
	// requests.
	ctx := context.Background()

	// read JSON files, and create []uint8 data which  is the type expected by the io libraries
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		fmt.Println(err)
	}

	// ConfigFromJSON uses a Google Developers Console client_credentials.json
	// file to construct a config.
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		fmt.Println(err)
	}

	// getClient returns *http.Client.
	client := getClient(ctx, config)

	// calendar.New return a Service object.
	/*	type Service struct {
			BasePath  string // API endpoint base URL
			UserAgent string // optional additional User-Agent fragment
			Acl *AclService
			CalendarList *CalendarListService
			Calendars *CalendarsService
			Channels *ChannelsService
			Colors *ColorsService
			Events *EventsService
			Freebusy *FreebusyService
			Settings *SettingsService
			// contains filtered or unexported fields
		}
	*/
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: "2015-05-28T09:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2015-05-28T17:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: "lpage@example.com"},
			&calendar.EventAttendee{Email: "sbrin@example.com"},
		},
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		fmt.Printf("%#v", srv)
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
