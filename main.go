package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
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
	HeaderDate  time.Time 	`json:"Header date"`
	Pilot       string 		`json:"Pilot"`
	Glider      string 		`json:"Glider"`
	GliderId    string 		`json:"Glider id"`
	TrackLength float64		`json:"Track length"`
}

type TrackDB struct {
	tracks map[string]Track
}

type ID struct {
	ID	string	`json:"id"`
}

type URL struct {
	URL string `json:"url"`
}

var startTime time.Time
var tracks map[int]Track
var IDs []string
var db TrackDB
var lastUsed int

func init(){
	startTime = time.Now()
}

func (db *TrackDB) Init() {
	db.tracks = make(map[string]Track)
}

func (db *TrackDB) Add(t Track, i ID) {
	db.tracks[i.ID] = t
	IDs = append(IDs, i.ID)
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
	case "POST":
		if r.Body == nil {
			http.Error(w, "Missing body", http.StatusBadRequest)
			return
		}
		var u URL
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if checkURL(u.URL) == false {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		track, err := igc.ParseLocation(u.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		totalDistance := 0.0
		for i := 0; i < len(track.Points)-1; i++ {
			totalDistance += track.Points[i].Distance(track.Points[i+1])
		}
		var i ID
		i.ID = ("ID" + strconv.Itoa(lastUsed))
		t := Track{track.Header.Date,
			track.Pilot,
			track.GliderType,
			track.GliderID,
			totalDistance}
		lastUsed++
		if db.tracks == nil {
			db.Init()
		}
		db.Add(t, i)
		return
	case "GET":
		if len(IDs) == 0 {
			IDs = make([]string, 0)
		}
		json.NewEncoder(w).Encode(IDs)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func checkURL(u string) bool {
	check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
	if check == true {
		return true
	}
	return false
}
/*  ta i bruk denne senere
func CalculatedDistance(track igc.Track) float64 {
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}
	return distance
}
*/
func main(){
	http.HandleFunc("/igcinfo/api/", handlerApi)
	http.HandleFunc("/igcinfo/api/igc/", handlerIgc)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
