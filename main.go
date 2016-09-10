package main

import (
	"./lib/googleCalendar"
	"fmt"
	"golang.org/x/net/context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
)

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
	client := googleCalendar.GetGoogleClient(ctx, config)

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

	/*	Event — An event on a calendar containing information such as the title,
		start and end times, and attendees. Events can be either single events or recurring
		events. An event is represented by an Event resource. The Events collection for a given
		calendar contains all event resources for that calendar.
	*/

	var summary, location, description string
	//var recurrence []string
	var start, end *calendar.EventDateTime

	summary = "test"

	location = "Stony Brook"

	description = "TEST"

	start = &calendar.EventDateTime{
		// Formatted according to RFC3339
		// the last time thingy is TIME OFFSET,
		// in NY, its -04:00:00. EDT
		DateTime: "2016-09-08T09:00:00-04:00:00",
		// Formatted as an IANA Time Zone Database name, e.g. "Europe/Zurich"
		TimeZone: "America/New_York",
	}
	end = &calendar.EventDateTime{
		DateTime: "2016-09-08T11:00:00-04:00:00",
		TimeZone: "America/New_York",
	}

	event := &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start:       start,
		End:         end,
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		fmt.Printf("%#v", srv)
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
