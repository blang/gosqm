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
	Name      string
	Position  [3]string
	Direction string
	Classname string
	Skill     string
	Formation string
	IsLeader  bool
	class     *sqm.Class
}

type Marker struct {
	Name       string
	Position   [3]string
	Type       string
	IsEllipse  bool
	Text       string
	ColorName  string
	FillName   string
	DrawBorder bool
	Size       [2]string
	class      *sqm.Class
}

type Mission struct {
	Groups []*Group
	class  *sqm.Class
	// Vehicle
	Markers []*Marker
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
