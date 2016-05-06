package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/BenChapman/notifier-plus/core/launchTmate"

	"github.com/gorilla/mux"
)

var (
	slackGroupURL = SECRET
)

func pipelineFailure(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&launchTmate.FailureData)
	if err != nil {
		panic(fmt.Sprintf("Decoding failed %s", err))
	}

	notifyHumansOfFailure()
}

func notifyHumansOfFailure() {
	message := fmt.Sprintf("{\"channel\": \"@somebody\", \"username\": \"webhookbot\", \"text\": \"%s\", \"icon_emoji\": \":ghost:\"}", launchTmate.FailureData.Job)
	messageReader := strings.NewReader(message)
	response, err := http.Post(slackGroupURL, "text/json", messageReader)

	fmt.Print(response)
	if err != nil {
		panic(fmt.Sprintf("Failed to post to slack channel %s", err))
	}

	fmt.Printf("FAILUREDATA: %#v\n", launchTmate.FailureData)
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/failure", pipelineFailure).Methods("POST")

	r.HandleFunc("/tmate", launchTmate.Launch()).Methods("POST")
	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}
