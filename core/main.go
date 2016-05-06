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
	slackGroupURL string
	failure       launchTmate.FailureInfo
	tmateUser     string
	tmateHost     string
)

func pipelineFailure(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&failure)
	if err != nil {
		panic(fmt.Sprintf("Decoding failed %s", err))
	}

	notifyHumansOfFailure(failure)
}

func notifyHumansOfFailure(data launchTmate.FailureInfo) {
	failureUrl := fmt.Sprintf("http://192.168.100.4:8080/pipelines/%s/jobs/%s/builds/%s", data.Pipeline, data.Job, data.Build)
	messageText := fmt.Sprintf("Human, there was an error on %s/%s, build #%s. %s\nType /helpplz if you'd like some help. Meow.", data.Pipeline, data.Job, data.Build, failureUrl)
	message := fmt.Sprintf("{\"channel\": \"#concourse-cat\", \"username\": \"concourse-cat\", \"text\": \"%s\", \"icon_emoji\": \":smirk_cat:\"}", messageText)
	messageReader := strings.NewReader(message)
	response, err := http.Post(slackGroupURL, "text/json", messageReader)

	fmt.Print(response)
	if err != nil {
		panic(fmt.Sprintf("Failed to post to slack channel %s", err))
	}

	fmt.Printf("FAILUREDATA: %#v\n", launchTmate.FailureData)
}

func incomingCommand(w http.ResponseWriter, request *http.Request) {
	sessionUrl := launchTmate.Launch("https://tmate.io/t/" + tmateUser)
	w.Write([]byte(fmt.Sprintf("Your hijacked session for pipeline '%s' is ready to use - %s\nUse the command `fly -t lite hijack -j %s/%s`", failure.Pipeline, sessionUrl, failure.Pipeline, failure.Job)))
}

func main() {
	slackGroupURL = os.Getenv("SLACK_HOOK_URL")
	if slackGroupURL == "" {
		log.Fatal("You must specify a SLACK_HOOK_URL environment variable for your hook")
	}

	tmateUser = os.Getenv("TMATE_USER")
	if tmateUser == "" {
		log.Fatal("You must specify a TMATE_USER environment variable to ssh into a tmate session.")
	}

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/failure", pipelineFailure).Methods("POST")
	r.HandleFunc("/helpmeout", incomingCommand).Methods("POST")

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
