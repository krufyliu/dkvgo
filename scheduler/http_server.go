package scheduler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/krufyliu/dkvgo/scheduler/tracker"
)

func stopJob(w http.ResponseWriter, r *http.Request) {
	jobId, _ := strconv.Atoi(mux.Vars(r)["id"])
	retMap := make(map[string]interface{})
	ok := tracker.StopJobById(jobId)
	retMap["success"] = ok
	if ok {
		retMap["message"] = "SUCCESS"
	} else {
		retMap["message"] = "job is not running"
	}
	retJson, _ := json.Marshal(retMap)
	log.Printf("API: stop task %d, result: %s", jobId, string(retJson))
	w.Write(retJson)
}

func APIServer(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/jobs/{id:[0-9]+}/action/stop", stopJob).Methods("POST")
	return &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}
