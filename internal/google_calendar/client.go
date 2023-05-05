package google_calendar

import (
	"context"
	"os"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarClient struct {
	client *calendar.Service
}

func NewGoogleCalendarClient(credPath string) (*GoogleCalendarClient, error) {
	file, err := os.ReadFile(credPath)
	if err != nil {
		return nil, err
	}
	client, err := calendar.NewService(context.Background(), option.WithCredentialsJSON(file))
	if err != nil {
		return nil, err
	}
	return &GoogleCalendarClient{client: client}, nil
}

func (gc *GoogleCalendarClient) CreateAclRule(user, role, calendarId string) (*calendar.AclRule, error) {
	rule := &calendar.AclRule{
		Role: role,
		Scope: &calendar.AclRuleScope{
			Type:  "user",
			Value: user,
		},
	}
	return gc.client.Acl.Insert(calendarId, rule).Do()
}

func (gc *GoogleCalendarClient) CreateCalendar(summary string) (*calendar.Calendar, error) {
	newCalendar := &calendar.Calendar{
		Summary: summary,
	}
	return gc.client.Calendars.Insert(newCalendar).Do()
}

func (gc *GoogleCalendarClient) CreateEvent(summary, date, calendarId string) (*calendar.Event, error) {
	event := calendar.Event{
		Summary:            summary,
		EndTimeUnspecified: true,
		Start:              &calendar.EventDateTime{Date: date},
	}
	return gc.client.Events.Insert(calendarId, &event).Do()
}

func (gc *GoogleCalendarClient) GetCalendars() (*calendar.CalendarList, error) {
	call := gc.client.CalendarList.List()
	return call.Do()
}

func (gc *GoogleCalendarClient) GetEvents(calendarId string) (*calendar.Events, error) {
	return gc.client.Events.List(calendarId).Do()
}

func (gc *GoogleCalendarClient) DeleteCalendar(calendarId string) error {
	return gc.client.Calendars.Delete(calendarId).Do()
}
