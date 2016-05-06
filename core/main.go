package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func notifyHumansOfFailure(data FailureInfo) {
	failureUrl := fmt.Sprintf("http://192.168.100.4:8080/pipelines/%s/jobs/%s/builds/%s", data.Pipeline, data.Job, data.Build)
	messageText := fmt.Sprintf("There was an error on %s/%s, build #%s. %s", data.Pipeline, data.Job, data.Build, failureUrl)
	message := fmt.Sprintf("{\"channel\": \"#notifier-plus\", \"username\": \"notifier-plus-bot\", \"text\": \"%s\", \"icon_emoji\": \":ghost:\"}", messageText)
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

	r.HandleFunc("/tmate", launchTmate.Launch).Methods("POST")
	// Bind to a port and pass our router in
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Listening at :" + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
