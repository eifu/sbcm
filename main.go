package main

import (
	"./lib"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"os"
	"html/template"
	"errors"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {

	filename := "data/"+p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}

func loadPage(title string) (*Page, error) {
	filename := "data/"+title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title,Body:  body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil{
    	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    	return 
    }
    renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	// The value returned by FormValue is of type string.
	p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil{
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return 
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
    // redirect to the view/*title*
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func init(){
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}
}

func main(){
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
    http.ListenAndServe(":8080", nil)
}
func _main() {

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
	client := lib.GetGoogleClient(ctx, config)

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

	now := time.Now()

	start = &calendar.EventDateTime{
		// Formatted according to RFC3339
		// the last time thingy is TIME OFFSET,
		// in NY, its -04:00:00. EDT
		DateTime: now.Format(time.RFC3339),
		// Formatted as an IANA Time Zone Database name, e.g. "Europe/Zurich"
		TimeZone: "America/New_York",
	}
	end = &calendar.EventDateTime{
		DateTime: now.Add(time.Duration(30) * time.Minute).Format(time.RFC3339),
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
