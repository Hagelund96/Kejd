//Information struct and correlating functions
package _struct

import "time"

//const
const (
	Version     = "1.0"
	Description = "Service for IGC tracks."
)

//struct
type Information struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//var
var startTime time.Time

//function to get server start time
func init() {
	startTime = time.Now()
}

//server to calculate current uptime
func Uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)
	return now.Sub(startTime).String()
}
