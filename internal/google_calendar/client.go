package google_calendar

import (
	"context"

	"google.golang.org/api/calendar/v3"
)

type GoogleCalendarClient struct {
	client *calendar.Service
}

func NewGoogleCalendarClient() (*GoogleCalendarClient, error) {
	client, err := calendar.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	return &GoogleCalendarClient{client: client}, nil
}
