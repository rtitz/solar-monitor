package rctUtils

import (
	"errors"
	"fmt"
	"solar-monitor/generalUtils"
	"solar-monitor/variables"
	"strconv"
	"strings"
	"time"

	"github.com/mlnoga/rct"
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

func GetValueRct(connection rct.Connection, id rct.Identifier, query string) (string, error) {
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

func BuildData(c *rct.Connection) {
	for {
		// Get time and date
		dateNow, timeNow, timezoneNow := generalUtils.GetDateAndTime()
		//fmt.Printf("%s %s %s\n", dateNow, timeNow, timezoneNow)

		// Get values
		for i := range variables.RctValues {
			for {
				output, err := GetValueRct(*c, variables.RctValues[i].Id, variables.RctValues[i].Query)
				if err == nil {
					variables.RctValues[i].Output = output
					variables.RctValues[i].DateNow = dateNow
					variables.RctValues[i].TimeNow = timeNow
					variables.RctValues[i].TimeZoneNow = timezoneNow
					break
				} else {
					fmt.Printf(" !! RETRY: %s\n", "FAILED TO GET DATA")
					time.Sleep(time.Millisecond * time.Duration(variables.RetryIntervalInMilliSeconds))
				}
			}
		}
		var RealPowerW, TotalGridPowerW float64
		RealPowerW, _ = strconv.ParseFloat(variables.RctValues[3].Output, 32)
		TotalGridPowerW, _ = strconv.ParseFloat(variables.RctValues[4].Output, 32)
		variables.TotalHouseComsumptionW = RealPowerW + TotalGridPowerW
		if TotalGridPowerW < 0 {
			variables.UsingSolarPowerOnly = true
		} else {
			variables.UsingSolarPowerOnly = false
		}
		time.Sleep(time.Millisecond * time.Duration(variables.FetchIntervalInMilliSeconds))
	}
}
