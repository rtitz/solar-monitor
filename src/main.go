package main

import (
	"fmt"
	"os"
	"os/exec"
	"solar-monitor/rctApi"
	"solar-monitor/rctUtils"
	"solar-monitor/variables"
	"strings"
	"time"

	"github.com/mlnoga/rct"
)

func main() {
	fmt.Printf("%s %s (%s/%s)\n\n", variables.AppName, variables.AppVersion, variables.GOOS, variables.GOARCH)
	rctApi.TestFunction()
	c, err := rct.NewConnection(variables.InverterIp, time.Second*2)
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

	for {
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

		fmt.Println()
		time.Sleep(time.Millisecond * time.Duration(variables.OutputRefreshIntervalInMilliSeconds))
	}

}
