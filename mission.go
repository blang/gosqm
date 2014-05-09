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
	Units     []*Vehicle
	class     *sqm.Class
	// Leader    *Unit
}

type Waypoint struct {
	Position         [3]string
	Type             string
	ShowWP           string
	Effects          *Effects
	Synchronizations []string
	class            *sqm.Class
}

type Vehicle struct {
	Name                string
	Position            [3]string
	Angle               string
	Classname           string
	Skill               string
	Special             string
	IsLeader            bool
	Player              string
	Description         string
	Presence            string
	PresenceCond        string
	Placement           string
	Age                 string
	Lock                string
	Rank                string
	Health              string
	Fuel                string
	Ammo                string
	Init                string
	Side                string
	Markers             []string
	ForceHeadlessClient bool
	class               *sqm.Class
}

type Marker struct {
	Name       string
	Position   [3]string
	Angle      string
	Type       string
	MarkerType string
	Text       string
	ColorName  string
	FillName   string
	DrawBorder bool
	Size       [2]string
	class      *sqm.Class
}

type Sensor struct {
	Name             string
	Position         [3]string
	Size             [2]string
	Angle            string
	IsRectangle      bool
	ActivationBy     string
	ActivationType   string
	TimeoutMin       string
	TimeoutMid       string
	TimeoutMax       string
	Type             string
	IsRepeating      bool
	Age              string
	Condition        string
	OnActivation     string
	OnDeactivation   string
	IsInterruptible  bool
	Text             string
	Synchronizations []string
	VehicleID        string
	Effects          *Effects
	class            *sqm.Class
}

type Effects struct {
	Sound       string
	Voice       string
	SoundEnv    string
	SoundDet    string
	Track       string
	TitleType   string
	Title       string
	TitleEffect string
}

type Mission struct {
	Addons     []string
	AddonsAuto []string
	RandomSeed string
	Intel      *Intel
	Groups     []*Group
	Vehicles   []*Vehicle
	Markers    []*Marker
	Sensors    []*Sensor
	class      *sqm.Class
}

type Intel struct {
	ResistanceWest  bool
	StartWeather    string
	ForecastWeather string
	Year            string
	Month           string
	Day             string
	Hour            string
	Minute          string
	class           *sqm.Class
}
