package main

import (
	"github.com/Hagelund96/Kejd/handler"
	"github.com/Hagelund96/Kejd/struct"
	"net/http"
	"os"
)

//main function for application. Initialises storage database
func main() {
	_struct.Db.Init()

	var p string
	if port := os.Getenv("PORT"); port != "" {
		p = ":" + port
	} else {
		p = ":8080"
	}

	//different handlers for urls
	http.HandleFunc("/paragliding/", handler.Handler)
	http.HandleFunc("/paragliding/api/", handler.HandlerApi)
	http.HandleFunc("/paragliding/api/track/", handler.HandlerTrack)
	http.HandleFunc("/paragliding/api/track/{id}", handler.HandlerTrackId)
	http.HandleFunc("/paragliding/api/track/{id}/{field}", handler.HandlerTrackIdFIeld)


	http.ListenAndServe(p, nil)
}
