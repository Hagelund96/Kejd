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
	http.HandleFunc("/igcinfo/api/", handler.HandlerApi)
	http.HandleFunc("/igcinfo/api/igc/", handler.HandlerIgc)

	http.ListenAndServe(p, nil)
}