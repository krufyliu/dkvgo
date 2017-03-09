package scheduler

import (
	"net/http"
)

func getJobDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)	
		return
	}

}

func stopJob(w http.ResponseWriter, r *http.Request) {

}

func APIGate() {
 	http.HandleFunc("/api/jobs", getJobDetail)
 	http.HandleFunc("/api/job/action/stop", stopJob)
}

type HTTPServer struct {
	ctx *DkvScheduler
}