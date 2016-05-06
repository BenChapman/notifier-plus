package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	slackGroupURL = "https://hooks.slack.com/services/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"
)

type FailureInfo struct {
	Pipeline string
	Job      string
	Build    string
}

func pipelineFailure(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var failureData FailureInfo

	err := decoder.Decode(&failureData)
	if err != nil {
		panic(fmt.Sprintf("Decoding failed %s", err))
	}

	notifyHumansOfFailure(failureData)
}

func launchTmate(w http.ResponseWriter, request *http.Request) {

}

func notifyHumansOfFailure(data FailureInfo) {
	message := "Hello lovely people"
	messageReader := strings.NewReader(message)
	response, err := http.Post(slackGroupURL, "text", messageReader)

	fmt.Print(response)
	if err != nil {
		panic(fmt.Sprintf("Failed to post to slack channel %s", err))
	}

	fmt.Printf("FAILUREDATA: %#v\n", data)
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/failure", pipelineFailure).Methods("POST")

	r.HandleFunc("/tmate", launchTmate).Methods("POST")
	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}
