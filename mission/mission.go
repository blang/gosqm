package mission

import (
	sqm "github.com/blang/gosqm/sqmparser"
)

type MissionFile struct {
	Version    string
	Mission    Mission
	Intro      Mission
	OutroWin   Mission
	OutroLoose Mission
}

type Mission struct {
	class sqm.Class
	Group
	// Vehicle
	// Marker
	// Trigger
}

type Intel struct {
	ResistanceWest  string
	StartWeather    string
	ForecastWeather string
	Year            string
	Day
	Hour
	Minute
}

//intel
//groups
//vehicles
//markers
//trigger
