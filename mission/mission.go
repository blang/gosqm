package mission

import (
	"github.com/blang/gosqm/sqm"
)

type MissionFile struct {
	Version    string
	Mission    *Mission
	Intro      *Mission
	OutroWin   *Mission
	OutroLoose *Mission
	class      *sqm.Class
}

type Group struct {
	Side      string
	Waypoints []*Waypoint
	Units     []*Unit
	class     *sqm.Class
	// Leader    *Unit
}

type Waypoint struct {
}

type Unit struct {
	Position  [3]string
	Direction string
	Classname string
	Skill     string
	Formation string
	IsLeader  bool
	class     *sqm.Class
}

type Mission struct {
	Groups []*Group
	class  *sqm.Class
	// Vehicle
	// Marker
	// Trigger
}

type Intel struct {
	ResistanceWest  string
	StartWeather    string
	ForecastWeather string
	Year            string
	Day             string
	Hour            string
	Minute          string
}

//intel
//groups
//vehicles
//markers
//trigger
