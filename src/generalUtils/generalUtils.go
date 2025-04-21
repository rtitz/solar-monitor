package generalUtils

import (
	"fmt"
	"strings"
	"time"
)

// Function go get current time and date, return a string for date and a string for time
func GetDateAndTime() (string, string, string) {
	currentTime := time.Now()
	dateNow := currentTime.Format("2006/01/02")
	timeNow := currentTime.Format("15:04:05")
	timezoneNow, _ := currentTime.Zone()
	return dateNow, timeNow, timezoneNow
}

func CreateJSONReaderKeyValue(key, value string) *strings.Reader {
	jsonStr := fmt.Sprintf(`{
        	"%s": "%s"
    	}`, key, value)
	return strings.NewReader(jsonStr)
}
