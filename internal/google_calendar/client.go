package google_calendar

import (
	"context"
	"net/http"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarClient struct {
	client *calendar.Service
}

func NewGoogleCalendarClient(oauthClient *http.Client) (*GoogleCalendarClient, error) {

	client, err := calendar.NewService(context.Background(), option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, err
	}
	return &GoogleCalendarClient{client: client}, nil
}

func (gc *GoogleCalendarClient) CreateEvent() {
	//event := calendar.Event{
	//	AttendeesOmitted:   true,
	//	EndTimeUnspecified: true,
	//	Summary:            "test event",
	//	Start:              &calendar.EventDateTime{Date: "2023-05-04"},
	//}

	//call := gc.client.Events.Insert("03ce5d8053cc70eb468fa7b3365fc0f17c6c04d41155f50cf662846356d24349@group.calendar.google.com", &event)
	//res, err := call.Do()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(res.HTTPStatusCode)
	//}
}
