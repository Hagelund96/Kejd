package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/marni/goigc"
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

type Track struct {
	ID 			int			`json:"ID"`
	HeaderDate  time.Time 	`json:"Header date"`
	Pilot       string 		`json:"Pilot"`
	Glider      string 		`json:"Glider"`
	GliderId    string 		`json:"Glider id"`
	TrackLength float64		`json:"Track length"`
}

var startTime time.Time
var tracks map[int]Track
var ID int

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

func handlerIgc(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case("POST"):
		var url string
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		track, err := igc.ParseLocation(url)
		tracks[ID] = Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, CalculatedDistance(track)}
		post, err := json.Marshal(tracks[ID].ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		http.Header.Add(w.Header(), "content-type", "application/json")
		json.NewEncoder(w).Encode(post)
		ID++
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CalculatedDistance(track igc.Track) float64 {
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}
	return distance
}

func main(){
	http.HandleFunc("/igcinfo/api/", handlerApi)
	http.HandleFunc("/igcinfo/api/igc", handlerIgc)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
