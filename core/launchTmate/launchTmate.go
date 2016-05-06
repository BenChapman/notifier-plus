package launchTmate

import "net/http"

type FailureInfo struct {
	Pipeline string
	Job      string
	Build    string
}

var FailureData FailureInfo

func Launch(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(FailureData.Job))
}
