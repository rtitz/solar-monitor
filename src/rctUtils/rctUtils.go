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

func formatStringTo2Digits(s string) string {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.2f", value)
}

func GetValueRct(connection *rctApi.Connection, id rctApi.Identifier, query string) (string, error) {
	if query == "float32" {
		v, err := connection.QueryFloat32(id)
		if err != nil {
			return "", err
		}
		vStr := fmt.Sprintf("%.2f", v)
		if strings.Contains(strings.ToLower(vStr), "invalid") {
			return vStr, errors.New("FAILED TO GET DATA")
		}
		return vStr, nil
	}
	return "", errors.New("unsupported query type")
}

func BuildData(c *rctApi.Connection) {
	var failedToGetDataCounter int
	for {
		dateNow, timeNow, timezoneNow, _, _ := generalUtils.GetDateAndTime()

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
				}

				failedToGetDataCounter++
				if failedToGetDataCounter >= variables.MaxRetriesToFetchDataUntilInvalid {
					fmt.Printf(" !! RETRY: %s\n", "FAILED TO GET DATA")
					variables.RctValues[i].Output = "0"
					variables.ForceUpdateHomeAssistant = true
				}
				time.Sleep(time.Millisecond * time.Duration(variables.RetryIntervalInMilliSeconds))
			}
		}

		RealPowerW, _ := strconv.ParseFloat(variables.RctValues[3].Output, 32)
		TotalGridPowerW, _ := strconv.ParseFloat(variables.RctValues[4].Output, 32)
		variables.TotalHouseComsumptionW = RealPowerW + TotalGridPowerW
		variables.SolarPowerAvailable = RealPowerW > 0
		variables.UsingSolarPowerOnly = TotalGridPowerW < 0

		time.Sleep(time.Millisecond * time.Duration(variables.FetchIntervalInMilliSeconds))
	}
}
