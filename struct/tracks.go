package _struct

import ("time"
		"github.com/marni/goigc"
)

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


var IDs []string
var Db TrackDB
var LastUsed int



func (db *TrackDB) Init() {
	db.tracks = make(map[string]Track)
}

func (db *TrackDB) Add(t Track, i ID) {
	db.tracks[i.ID] = t
	IDs = append(IDs, i.ID)
}

func (db *TrackDB) Get(keyID string) (Track, bool) {
	t, err := db.tracks[keyID]
	return t, err
}


func CalculatedDistance(track igc.Track) float64 {
	distance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		distance += track.Points[i].Distance(track.Points[i+1])
	}
	return distance
}