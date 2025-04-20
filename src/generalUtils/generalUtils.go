package generalUtils

import "time"

// Function go get current time and date, return a string for date and a string for time
func GetDateAndTime() (string, string, string) {
	currentTime := time.Now()
	dateNow := currentTime.Format("2006/01/02")
	timeNow := currentTime.Format("15:04:05")
	timezoneNow, _ := currentTime.Zone()
	return dateNow, timeNow, timezoneNow
}
