package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"solar-monitor/generalUtils"
	"solar-monitor/homeassistantApi"
	"solar-monitor/rctApi"
	"solar-monitor/rctUtils"
	"solar-monitor/variables"
	"strconv"
	"strings"
	"time"
)

func main() {
	//lastUsingSolarPowerOnly := false

	fmt.Printf("%s %s (%s/%s)\n\n", variables.AppName, variables.AppVersion, variables.GOOS, variables.GOARCH)
	c, err := rctApi.NewConnection(variables.InverterIp, time.Second*2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	time.Sleep(time.Millisecond * time.Duration(variables.StartDelayInMilliSeconds)) // Start delay; sometimes needed for restart
	go rctUtils.BuildData(c)

	for variables.RctValues[0].TimeNow == "" {
		fmt.Println("Waiting for data...")
		time.Sleep(time.Millisecond * time.Duration(variables.StartDelayInMilliSeconds))
	}
	fmt.Printf("\n =========================================== \n\n")

	// For-loop counter
	var lc int = 0
	variables.ForceUpdateHomeAssistant = true
	for {
		lc++

		if variables.ClearLoopOutput {
			var cmd *exec.Cmd
			switch os := variables.GOOS; os {
			case "darwin":
				cmd = exec.Command("clear") // macOS example, its tested
			case "linux":
				cmd = exec.Command("clear") // Linux example, its tested
			case "windows":
				cmd = exec.Command("cmd", "/c", "cls") // Windows example, its tested
			default:
				cmd = exec.Command("") // default
			}
			cmd.Stdout = os.Stdout
			cmd.Run()
		}

		dateNow := variables.RctValues[0].DateNow
		timeNow := variables.RctValues[0].TimeNow
		timeZone := variables.RctValues[0].TimeZoneNow
		solarGenAPowerW, _ := strconv.ParseFloat(variables.RctValues[0].Output, 32)
		solarGenBPowerW, _ := strconv.ParseFloat(variables.RctValues[1].Output, 32)
		inverterACPowerW, _ := strconv.ParseFloat(variables.RctValues[2].Output, 32)
		totalSolarPowerW, _ := strconv.ParseFloat(variables.RctValues[3].Output, 32)
		totalHouseConsumptionW := variables.TotalHouseComsumptionW
		totalGridPowerW, _ := strconv.ParseFloat(variables.RctValues[4].Output, 32)
		usingSolarPowerOnly := variables.UsingSolarPowerOnly

		jsonString := generalUtils.CreateJsonString(
			dateNow,
			timeNow,
			timeZone,
			solarGenAPowerW,
			solarGenBPowerW,
			inverterACPowerW,
			totalSolarPowerW,
			totalHouseConsumptionW,
			totalGridPowerW,
			usingSolarPowerOnly,
		)

		// Feed HomeAssistant
		if !variables.SolarPowerAvailable { // Force HomeAssistant Update if SolarPower becomes unavailable
			variables.ForceUpdateHomeAssistant = true
		}
		if variables.ForceUpdateHomeAssistant || lc >= variables.HomeAssistantUpdateIntervalInOutputRefreshLoopIterations {
			variables.ForceUpdateHomeAssistant = false
			lc = 0
			if variables.SolarPowerAvailable {
				homeassistantApi.CallHomeAssisatantApiSet(variables.ApiBoolSolarPowerAvailable, generalUtils.CreateJSONReaderKeyValue("state", "on"))
			} else {
				homeassistantApi.CallHomeAssisatantApiSet(variables.ApiBoolSolarPowerAvailable, generalUtils.CreateJSONReaderKeyValue("state", "off"))
			}

			if variables.UsingSolarPowerOnly {
				//if !lastUsingSolarPowerOnly {
				//	lastUsingSolarPowerOnly = true
				homeassistantApi.CallHomeAssisatantApiSet(variables.ApiBoolSolarPowerOnly, generalUtils.CreateJSONReaderKeyValue("state", "on"))
				//}
			} else {
				//if lastUsingSolarPowerOnly {
				//	lastUsingSolarPowerOnly = false
				homeassistantApi.CallHomeAssisatantApiSet(variables.ApiBoolSolarPowerOnly, generalUtils.CreateJSONReaderKeyValue("state", "off"))
				//}
			}

			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntSolarGenAPowerW, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", solarGenAPowerW)))
			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntSolarGenBPowerW, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", solarGenBPowerW)))
			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntSolarGenTotalPowerW, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", inverterACPowerW)))
			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntPVDelivery, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", totalSolarPowerW)))
			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntHouseConsumption, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", totalHouseConsumptionW)))
			homeassistantApi.CallHomeAssisatantApiSet(variables.ApiIntGridFeed, generalUtils.CreateJSONReaderKeyValue("state", fmt.Sprintf("%.0f", totalGridPowerW)))

		}

		// Write to JSON file
		if variables.WriteJsonToFile {
			generalUtils.WriteJsonToFile(jsonString, variables.JsonFile)
		}

		// OUTPUT
		if variables.OutputJson { // JSON output
			//fmt.Printf("%s\n", string(jsonString))
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, jsonString, "", "    "); err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(prettyJSON.String())

		} else { // TEXT output
			//timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
			timeStamp := time.Now().String()
			fmt.Printf("%s %s %s (%s)\n", variables.RctValues[0].DateNow, variables.RctValues[0].TimeNow, variables.RctValues[0].TimeZoneNow, timeStamp)
			for i := range variables.RctValues {
				if i == 4 {
					continue
				} // Skip TotalGridPowerW
				fmt.Printf("%s: %s\n", variables.RctValues[i].Id, variables.RctValues[i].Output)
			}
			fmt.Printf("%s: %s\n", "Total house comsumption [W]", fmt.Sprintf("%.2f", variables.TotalHouseComsumptionW))
			if variables.UsingSolarPowerOnly {
				fmt.Printf("%s: %s\n", "Feed to grid [W]", strings.TrimPrefix(variables.RctValues[4].Output, "-"))
				fmt.Printf("%s: %s\n", "Using solar power only", "True")
			} else {
				fmt.Printf("%s: %s\n", "Consumtion from grid [W]", variables.RctValues[4].Output)
				fmt.Printf("%s: %s\n", "Using solar power only", "False")
			}
		}
		fmt.Println()
		time.Sleep(time.Millisecond * time.Duration(variables.OutputRefreshIntervalInMilliSeconds))
	}

}
