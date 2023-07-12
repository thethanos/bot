// for linter
package main

func main() {

}

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type HTTPHandler struct {
	signalHandler chan os.Signal
	config        *oauth2.Config
}

func (h *HTTPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	authCode := req.URL.Query().Get("code")

	token, err := h.config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	fmt.Printf("Saving credential file to: %s\n", "token.json")
	file, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Token is saved to the file, you may close the window."))
	h.signalHandler <- syscall.SIGTERM
}

func createConfig(path string) (*oauth2.Config, error) {
	credData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credData, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
}

func setupSignalHandler() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}

func main() {

	signalHandler := setupSignalHandler()
	config, err := createConfig("credentials.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: &HTTPHandler{signalHandler: signalHandler, config: config},
	}

	go func() {
		fmt.Println(server.ListenAndServe())
	}()

	getTokenFromWeb(config)

	<-signalHandler

	server.Shutdown(context.Background())
}
*/
