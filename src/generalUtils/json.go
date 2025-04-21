package generalUtils

import (
	"encoding/json"
	"fmt"
	"solar-monitor/variables"
)

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

	// Convert to JSON string
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return jsonString
}
