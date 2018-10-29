//track struct and correlating functions
package _struct

import (
	"github.com/marni/goigc"
	"time"
)

//struct
type Track struct {
	UniqueID     string
	Pilot       string    `json:"Pilot"`
	Glider      string    `json:"Glider"`
	GliderId    string    `json:"Glider id"`
	TrackLength float64   `json:"Track length"`
	HeaderDate  time.Time `json:"Header date"`
	URL          string
	TimeRecorded time.Time
}

//struct
type TrackDB struct {
	tracks map[string]Track
}

//struct
type ID struct {
	ID string `json:"id"`
}

//struct
type URL struct {
	URL string `json:"url"`
}

//var
var IDs []string
var Db TrackDB
var LastUsed int

//function to initialize database
func (db *TrackDB) Init() {
	db.tracks = make(map[string]Track)
}

//function for adding to database
func (db *TrackDB) Add(t Track, i ID) {
	db.tracks[i.ID] = t
	IDs = append(IDs, i.ID)
}

//function for getting id from database
func (db *TrackDB) Get(keyID string) (Track, bool) {
	t, err := db.tracks[keyID]
	return t, err
}

//function for calculating track distance
func CalculatedDistance(track igc.Track) float64 {
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}
	return distance
}
