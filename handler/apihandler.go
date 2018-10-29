package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hagelund96/Kejd/struct"
	"github.com/marni/goigc"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//checks if the given id exists
func checkId(id string) bool {
	idExists := false
	for i := 0; i < len(_struct.IDs); i++ {
		if _struct.IDs[i] == strings.ToUpper(id) {
			idExists = true
			break
		}
	}
	return idExists
}

//checks if url given is valid
func checkURL(u string) bool {
	check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
	if check == true {
		return true
	}
	return false
}

//parses ids into json, and encodes and outputs whole array of ids
func replyWithAllTracksId(w http.ResponseWriter, db _struct.TrackDB) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	if len(_struct.IDs) == 0 {
		_struct.IDs = make([]string, 0)
	}
	json.NewEncoder(w).Encode(_struct.IDs)
	return
}

//parses id into json, and encodes and outputs whole track mapped to id
func replyWithTracksId(w http.ResponseWriter, db _struct.TrackDB, id string) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	t, _ := db.Get(strings.ToUpper(id))
	api := _struct.Track{t.UniqueID, t.Pilot, t.Glider, t.GliderId, t.TrackLength, t.HeaderDate, t.URL, t.TimeRecorded}
	json.NewEncoder(w).Encode(api)
}

//parses field into json, and encodes and outputs it
func replyWithField(w http.ResponseWriter, db _struct.TrackDB, id string, field string) {
	http.Header.Set(w.Header(), "content-type", "application/json")
	t, _ := db.Get(strings.ToUpper(id))

	api := _struct.Track{t.UniqueID, t.Pilot, t.Glider, t.GliderId, t.TrackLength, t.HeaderDate, t.URL, t.TimeRecorded}

	switch strings.ToUpper(field) {
	case "PILOT":
		json.NewEncoder(w).Encode(api.Pilot)
	case "GLIDER":
		json.NewEncoder(w).Encode(api.Glider)
	case "GLIDER_ID":
		json.NewEncoder(w).Encode(api.GliderId)
	case "TRACK_LENGTH":
		json.NewEncoder(w).Encode(api.TrackLength)
	case "H_DATE":
		json.NewEncoder(w).Encode(api.HeaderDate)
	default:
		http.Error(w, "Not a valid option", http.StatusNotFound)
		return
	}
}

func HandlerIgc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//handling POST /api/igc/   **NEED TO HAVE SLASH, DOES NOT WORK WIHTOUT**
	case "POST":
		//checks that input is not empty
		if r.Body == nil {
			http.Error(w, "Missing body", http.StatusBadRequest)
			return
		}
		var u _struct.URL
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//checks if url is valid
		if checkURL(u.URL) == false {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		track, err := igc.ParseLocation(u.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//calculates total distance
		totalDistance := _struct.CalculatedDistance(track)

		URL := _struct.URL{}

		ID := rand.Intn(1000)

		track.UniqueID = strconv.Itoa(ID)

		trackFile := _struct.Track{}

		timestamp := time.Now().Second()
		timestamp = timestamp * 1000

		client := mongoConnect()

		collection := client.Database("paragliding").Collection("tracks")

		// Checking for duplicates so that the user doesn't add into the database igc files with the same URL
		duplicate := urlInMongo(URL.URL, collection)

		if !duplicate {

			trackFile = _struct.Track{
				track.UniqueID,
				track.Pilot,
				track.GliderType,
				track.GliderID,
				totalDistance,
				track.Date,
				URL.URL, time.Now()}

			res, err := collection.InsertOne(context.Background(), trackFile)
			if err != nil {
				log.Fatal(err)
			}

			id := res.InsertedID

			if id == nil {
				http.Error(w, "", 300)
			}

			// Encoding the ID of the track that was just added to DB
			fmt.Fprint(w, "{\n\"id\":\""+track.UniqueID+"\"\n}")

			//triggerWhenTrackIsAdded()

		} else {

			trackInDB := getTrack(client, URL.URL)
			// If there is another file in igcFilesDB with that URL return and tell the user that that IGC FILE is already in the database
			http.Error(w, "409 Conflict - The Igc File you entered is already in our database!", http.StatusConflict)
			fmt.Fprintln(w, "\nThe file you entered has the following ID: ", trackInDB.UniqueID)
			return

		}



		return
		//Handling all GETs after /api/
	case "GET":
		client := mongoConnect()

		ids := getTrackID(client)

		fmt.Fprint(w, ids)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Redirect to /paragliding/api
	http.Redirect(w, r, "/paragliding/api", 302)
}

//handler for /api shows uptime description and version in json
func HandlerApi(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 4 && parts[3] == "" {
		api := _struct.Information{_struct.Uptime(), _struct.Description, _struct.Version}
		json.NewEncoder(w).Encode(api)
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

//Old functions that we went away from, left behind to show working progress
/*func handlerIdAndField(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idExists := false
	for i := 0; i < len(IDs); i++ {
		if IDs[i] == strings.ToUpper(parts[4]) {
			idExists = true
			break
		}
	}
	if !idExists {
		http.Error(w, "ID out of range.", http.StatusNotFound)
		return
	}
	t, _ := db.Get(strings.ToUpper(parts[4]))
	if len(parts) == 5 {
		http.Header.Set(w.Header(), "content.type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
	if len(parts) == 6 {
		switch strings.ToUpper(parts[5]) {
		case "PILOT":
			fmt.Fprint(w, t.Pilot)
		case "GLIDER":
			fmt.Fprint(w, t.Glider)
		case "GLIDER_ID":
			fmt.Fprint(w, t.GliderId)
		case "TRACK_LENGTH":
			fmt.Fprint(w, t.TrackLength)
		case "H_DATE":
			fmt.Fprint(w, t.HeaderDate)
		default:
			http.Error(w, "Not a valid option", http.StatusNotFound)
			return
		}

	}
	if len(parts) > 6 {
		http.Error(w, "Too many /.", http.StatusNotFound)
	}
}*/
