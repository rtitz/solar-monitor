package variables

import (
	"runtime"

	"github.com/mlnoga/rct"
)

// App name and version
var AppName string
var AppVersion string

var InverterIp string = "10.0.0.17"
var FetchIntervalInMilliSeconds int = 1000
var RetryIntervalInMilliSeconds int = 200
var StartDelayInMilliSeconds int = 2000

var OutputRefreshIntervalInMilliSeconds int = 2000
var ClearLoopOutput bool = true

var TotalHouseComsumptionW float64 = 0
var UsingSolarPowerOnly bool = false

// RCT DATA
type RctData struct {
	Id          rct.Identifier
	Query       string
	Output      string
	DateNow     string
	TimeNow     string
	TimeZoneNow string
}

var RctValues = [5]RctData{
	{Id: rct.SolarGenAPowerW, Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
	{Id: rct.SolarGenBPowerW, Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
	{Id: rct.InverterACPowerW, Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
	{Id: rct.RealPowerW, Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
	{Id: rct.TotalGridPowerW, Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
}

type CalcData struct {
	Id          string
	Query       string
	Output      string
	DateNow     string
	TimeNow     string
	TimeZoneNow string
}

var CalcValues = [1]CalcData{
	{Id: "TotalHouse", Query: "float32", Output: "", DateNow: "", TimeNow: "", TimeZoneNow: ""},
}

// Other
var GOOS string
var GOARCH string

func init() {
	// App name and version
	AppName = "Solar-Monitor"
	AppVersion = "1.0.0"

	GOOS = runtime.GOOS
	GOARCH = runtime.GOARCH
}
