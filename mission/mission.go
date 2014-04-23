package mission

type missionFile struct {
  version string
  mission mission
  intro mission
  outroWin mission
  outroLoose mission
}


type mission struct {
	addons []string
	addonsAuto []string
	randomSeed int
	intel *intel
	group
	vehicle
	marker
	trigger


}

type intel struct {
	resistanceWest string
	startWeather string
	forecastWeather string
	year string
	day
	hour
	minute
}

//intel
//groups
//vehicles
//markers
//trigger