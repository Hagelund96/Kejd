package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const (
	Version = "1.0"
	Description = "Service for IGC tracks."
)

type Information struct {
	Uptime 		string 	`json:"uptime"`
	Info   		string 	`json:"info"`
	Version		string	`json:"version"`
}

var startTime time.Time

func init(){
	startTime = time.Now()
}

func uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)
	return now.Sub(startTime).String()
}

func handlerApi(w http.ResponseWriter, r *http.Request){
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if parts[2] == "api" && parts[3] == "" {
		api := Information{uptime(), Description, Version}
		json.NewEncoder(w).Encode(api)
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

func main(){
	http.HandleFunc("/igcinfo/api/", handlerApi)


	http.ListenAndServe("127.0.0.1:8080", nil)
}
