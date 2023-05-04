package google_oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func NewHttpClient(credPath, tokenPath string) (*http.Client, error) {

	credData, err := os.ReadFile(credPath)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credData, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer tokenFile.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), token), nil
}
