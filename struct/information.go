package _struct

import "time"

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

func Uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)
	return now.Sub(startTime).String()
}