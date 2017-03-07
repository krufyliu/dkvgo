package scheduler

import (
	"net/http"
)

func getJobDetail(w http.ResponseWriter, r *http.Request) {

}

func stopJob(w http.ResponseWriter, r *http.Request) {

}

http.Handle("/api/jobs", getJobDetail)
http.Handle("/api/job/action/stop", stopJob)

