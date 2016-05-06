package launchTmate

type FailureInfo struct {
	Pipeline string
	Job      string
	Build    string
}

var FailureData FailureInfo

func Launch(tmateUrl string) string {
	return tmateUrl
}
