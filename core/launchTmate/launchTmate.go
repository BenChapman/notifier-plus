package launchTmate

import "net/http"

type FailureInfo struct {
	Pipeline string
	Job      string
	Build    string
}

var FailureData FailureInfo

func Launch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(FailureData.Job))
	}
}
