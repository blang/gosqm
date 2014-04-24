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
	Position     [3]string
	Type         string
	ShowWP       string
	class        *sqm.Class
	classEffects *sqm.Class
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
	Name     string
	Position [3]string
	// Angle      string //TODO: Angle?
	Type       string
	IsEllipse  bool
	Text       string
	ColorName  string
	FillName   string
	DrawBorder bool
	Size       [2]string
	class      *sqm.Class
}

type Sensor struct {
	Name            string
	Position        [3]string
	Size            [2]string
	Angle           string
	IsRectangle     bool
	ActivationBy    string
	IsRepeating     bool
	Age             string
	Condition       string
	OnActivation    string
	IsInterruptible bool
	class           *sqm.Class
	classEffects    *sqm.Class
}

type Vehicle struct {
	Name      string
	Position  [3]string
	Angle     string
	Classname string
	Skill     string
	class     *sqm.Class
	Side      string //Always empty?
}

type Mission struct {
	Addons     []string
	AddonsAuto []string
	Intel      *Intel
	Groups     []*Group
	Vehicles   []*Vehicle
	Markers    []*Marker
	Sensors    []*Sensor
	class      *sqm.Class
}

type Intel struct {
	ResistanceWest  string
	StartWeather    string
	ForecastWeather string
	Year            string
	Month           string
	Day             string
	Hour            string
	Minute          string
	class           *sqm.Class
}

//intel
//groups
//vehicles
//markers
//trigger
