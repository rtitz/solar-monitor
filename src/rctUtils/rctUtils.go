package rctUtils

import (
	"errors"
	"fmt"
	"solar-monitor/generalUtils"
	"solar-monitor/rctApi"
	"solar-monitor/variables"
	"strconv"
	"strings"
	"time"
)

/*func getValueCalc(query string) (string, error) {
	var vStr string
	switch query {
	case "float32":
		var v float32 = 1 + 1
		vStr = fmt.Sprintf("%.2f", v)
		if strings.Contains(strings.ToLower(vStr), "invalid") {
			return vStr, errors.New("FAILED TO GET DATA")
		}
	}
	return vStr, nil
}*/

func formatStringTo2Digits(s string) string {
	// Parse string to float
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ""
	}
	// Format to 2 decimal places
	return fmt.Sprintf("%.2f", value)
}

func GetValueRct(connection *rctApi.Connection, id rctApi.Identifier, query string) (string, error) {
	var vStr string
	switch query {
	case "float32":
		v, err := connection.QueryFloat32(id)
		if err != nil {
			//fmt.Println(err)
			return "", err
		}
		vStr = fmt.Sprintf("%.2f", v)
		if strings.Contains(strings.ToLower(vStr), "invalid") {
			return vStr, errors.New("FAILED TO GET DATA")
		}
	}
	return vStr, nil
}

func BuildData(c *rctApi.Connection) {
	var failedToGetDataCounter = 0
	for {
		// Get time and date
		dateNow, timeNow, timezoneNow, _, _ := generalUtils.GetDateAndTime()
		//fmt.Printf("%s %s %s\n", dateNow, timeNow, timezoneNow)

		// Get values
		for i := range variables.RctValues {
			for {
				output, err := GetValueRct(c, variables.RctValues[i].Id, variables.RctValues[i].Query)
				if err == nil {
					failedToGetDataCounter = 0
					variables.RctValues[i].Output = formatStringTo2Digits(output)
					variables.RctValues[i].DateNow = dateNow
					variables.RctValues[i].TimeNow = timeNow
					variables.RctValues[i].TimeZoneNow = timezoneNow
					break
				} else {
					failedToGetDataCounter++
					if failedToGetDataCounter >= variables.MaxRetriesToFetchDataUntilInvalid {
						fmt.Printf(" !! RETRY: %s\n", "FAILED TO GET DATA")
						variables.RctValues[i].Output = "0"
						variables.ForceUpdateHomeAssistant = true
					}
					time.Sleep(time.Millisecond * time.Duration(variables.RetryIntervalInMilliSeconds))
				}
			}
		}

		var RealPowerW, TotalGridPowerW float64
		RealPowerW, _ = strconv.ParseFloat(variables.RctValues[3].Output, 32)
		TotalGridPowerW, _ = strconv.ParseFloat(variables.RctValues[4].Output, 32)
		variables.TotalHouseComsumptionW = RealPowerW + TotalGridPowerW
		if RealPowerW > 0 {
			variables.SolarPowerAvailable = true
		} else {
			variables.SolarPowerAvailable = false
		}
		if TotalGridPowerW < 0 {
			variables.UsingSolarPowerOnly = true
		} else {
			variables.UsingSolarPowerOnly = false
		}
		time.Sleep(time.Millisecond * time.Duration(variables.FetchIntervalInMilliSeconds))
	}
}
