package generalUtils

import (
	"encoding/json"
	"fmt"
	"solar-monitor/variables"
)

// Pre-allocate a buffer for JSON marshaling to reduce allocations
var jsonBuffer = make([]byte, 0, 512)

func CreateJsonString(Date, Time, TimeZone string, SolarGenAPowerW, SolarGenBPowerW, InverterACPowerW, TotalSolarPowerW, TotalHouseComsumptionW, TotalGridPowerW float64, UsingSolarPowerOnly bool) []byte {
	data := variables.JsonData{
		Date:                   Date,
		Time:                   Time,
		TimeZone:               TimeZone,
		SolarGenAPowerW:        SolarGenAPowerW,
		SolarGenBPowerW:        SolarGenBPowerW,
		InverterACPowerW:       InverterACPowerW,
		TotalSolarPowerW:       TotalSolarPowerW,
		TotalHouseComsumptionW: TotalHouseComsumptionW,
		TotalGridPowerW:        TotalGridPowerW,
		UsingSolarPowerOnly:    UsingSolarPowerOnly,
	}

	// Reset buffer and reuse it
	jsonBuffer = jsonBuffer[:0]
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return []byte("{}")  // Return empty JSON object instead of nil
	}
	return jsonString
}
