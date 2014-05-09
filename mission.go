package mission

type MissionFile struct {
	Version    string
	Mission    *Mission
	Intro      *Mission
	OutroWin   *Mission
	OutroLoose *Mission
}

type Group struct {
	Side      string
	Waypoints []*Waypoint
	Units     []*Vehicle
}

type Waypoint struct {
	Position         [3]string
	Type             string
	ShowWP           string
	Effects          *Effects
	Synchronizations []string
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
}
