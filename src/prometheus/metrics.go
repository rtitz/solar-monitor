package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SolarGenAPowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "solar_gen_a_power_watts",
		Help: "Solar generator A power in watts",
	})

	SolarGenBPowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "solar_gen_b_power_watts",
		Help: "Solar generator B power in watts",
	})

	InverterACPowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "inverter_ac_power_watts",
		Help: "Inverter AC power in watts",
	})

	TotalSolarPowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_solar_power_watts",
		Help: "Total solar power in watts",
	})

	TotalHouseConsumptionW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_house_consumption_watts",
		Help: "Total house consumption in watts",
	})

	TotalGridPowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_grid_power_watts",
		Help: "Total grid power in watts (positive = consuming, negative = feeding)",
	})

	UsingSolarPowerOnly = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "using_solar_power_only",
		Help: "Whether using solar power only (1 = true, 0 = false)",
	})

	SolarPowerAvailable = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "solar_power_available",
		Help: "Whether solar power is available (1 = true, 0 = false)",
	})
)

func UpdateMetrics(solarGenA, solarGenB, inverterAC, totalSolar, totalHouse, totalGrid float64, usingSolarOnly, solarAvailable bool) {
	SolarGenAPowerW.Set(solarGenA)
	SolarGenBPowerW.Set(solarGenB)
	InverterACPowerW.Set(inverterAC)
	TotalSolarPowerW.Set(totalSolar)
	TotalHouseConsumptionW.Set(totalHouse)
	TotalGridPowerW.Set(totalGrid)
	
	if usingSolarOnly {
		UsingSolarPowerOnly.Set(1)
	} else {
		UsingSolarPowerOnly.Set(0)
	}
	
	if solarAvailable {
		SolarPowerAvailable.Set(1)
	} else {
		SolarPowerAvailable.Set(0)
	}
}