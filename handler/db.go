package handler

import (
	"context"
	"github.com/Hagelund96/Kejd/struct"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"net/http"
)

func mongoConnect() *mongo.Client {
	// Connect to MongoDB
	conn, err := mongo.Connect(context.Background(), "mongodb://hagelund96:hagelund123@ds145053.mlab.com:45053/paragliding", nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
}

// Check if the track already exists in the database
func urlInMongo(url string, trackColl *mongo.Collection) bool {

	// Read the documents where the trackurl field is equal to url parameter
	cursor, err := trackColl.Find(context.Background(),
		bson.NewDocument(bson.EC.String("url", url)))
	if err != nil {
		log.Fatal(err)
	}

	// 'Close' the (cursor A pointer to the result set of a query. Clients can iterate through a cursor to retrieve results).
	defer cursor.Close(context.Background())

	track := _struct.Track{}

	// Point the cursor at whatever is found
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}
	}

	if track.URL == "" { // If there is an empty field, in this case, `url`, it means the track is not on the database
		return false
	}
	return true
}

// Get track
func getTrack(client *mongo.Client, url string) _struct.Track {
	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection

	cursor, err := collection.Find(context.Background(), bson.NewDocument(bson.EC.String("url", url)))

	if err != nil {
		log.Fatal(err)
	}

	resTrack := _struct.Track{}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
	}

	return resTrack

}

func getTrackID(client *mongo.Client) string {

	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection

	cursor, err := collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	resTrack := _struct.Track{}
	length, error := collection.Count(context.Background(), nil)
	if error != nil {
		log.Fatal(error)
	}
	ids := "["
	i := int64(0)
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&resTrack)
		if err != nil {
			log.Fatal(err)
		}
		ids += resTrack.UniqueID
		if i == length-1 {
			break
		}
		ids += ","
		i++
	}
	ids += "]"
	return ids
}

func getTrack1(client *mongo.Client, id string, w http.ResponseWriter) _struct.Track {
	db := client.Database("paragliding")     // `paragliding` Database
	collection := db.Collection("tracks") // `track` Collection
	filter := bson.NewDocument(bson.EC.String("uniqueid", ""+id+""))
	resTrack := _struct.Track{}
	err := collection.FindOne(context.Background(), filter).Decode(&resTrack)
	if err != nil {
		http.Error(w, "File not found!", 404)
	}
	return resTrack

}