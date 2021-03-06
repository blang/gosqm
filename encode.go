package gosqm

import (
	"github.com/blang/gosqm/sqm"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
)

type Encoder struct {
	wg *sync.WaitGroup
	w  io.Writer
}

type counter int32

func (c *counter) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}
func (c *counter) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func NewClassEncoder() *Encoder {
	e := &Encoder{
		wg: &sync.WaitGroup{},
	}
	return e
}

func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		wg: &sync.WaitGroup{},
		w:  w,
	}
	return e
}

func (e *Encoder) Encode(missionFile *MissionFile) error {
	class := e.EncodeToClass(missionFile)
	sqmenc := sqm.NewEncoder(e.w)
	return sqmenc.Encode(class)
}

func (e *Encoder) EncodeToClass(missionFile *MissionFile) *sqm.Class {
	e.wg = &sync.WaitGroup{}
	c := e.encodeMissionFile(missionFile)
	e.wg.Wait()
	return c
}

func (e *Encoder) encodeMissionFile(missionFile *MissionFile) *sqm.Class {
	mainClass := &sqm.Class{
		Name: "mission",
	}
	mainClass.Props = append(mainClass.Props, &sqm.Property{"version", sqm.TNumber, missionFile.Version})

	missionClass := &sqm.Class{
		Name: "Mission",
	}
	e.wg.Add(1)
	go func() {
		e.encodeMission(missionFile.Mission, missionClass)
		e.wg.Done()
	}()
	mainClass.Classes = append(mainClass.Classes, missionClass)
	introClass := &sqm.Class{
		Name: "Intro",
	}
	e.wg.Add(1)
	go func() {
		e.encodeMission(missionFile.Intro, introClass)
		e.wg.Done()
	}()

	mainClass.Classes = append(mainClass.Classes, introClass)

	outroWinClass := &sqm.Class{
		Name: "OutroWin",
	}
	e.wg.Add(1)
	go func() {
		e.encodeMission(missionFile.OutroWin, outroWinClass)
		e.wg.Done()
	}()
	mainClass.Classes = append(mainClass.Classes, outroWinClass)

	outroLooseClass := &sqm.Class{
		Name: "OutroLoose",
	}
	e.wg.Add(1)
	go func() {
		e.encodeMission(missionFile.OutroLoose, outroLooseClass)
		e.wg.Done()
	}()
	mainClass.Classes = append(mainClass.Classes, outroLooseClass)

	return mainClass
}

func (e *Encoder) encodeMission(mission *Mission, class *sqm.Class) {
	var counter counter = -1
	encodeMissionProperties(mission, class)
	intelClass := &sqm.Class{
		Name: "Intel",
	}
	e.wg.Add(1)
	go func() {
		encodeIntel(mission.Intel, intelClass)
		e.wg.Done()
	}()
	class.Classes = append(class.Classes, intelClass)

	if len(mission.Groups) > 0 {
		groupsClass := &sqm.Class{
			Name: "Groups",
		}
		e.wg.Add(1)
		go func() {
			e.encodeGroups(mission.Groups, groupsClass, &counter)
			e.wg.Done()
		}()
		class.Classes = append(class.Classes, groupsClass)
	}

	if len(mission.Markers) > 0 {
		markersClass := &sqm.Class{
			Name: "Markers",
		}
		e.wg.Add(1)
		go func() {
			encodeMarkers(mission.Markers, markersClass)
			e.wg.Done()
		}()
		class.Classes = append(class.Classes, markersClass)
	}

	if len(mission.Sensors) > 0 {
		sensorsClass := &sqm.Class{
			Name: "Sensors",
		}
		e.wg.Add(1)
		go func() {
			encodeSensors(mission.Sensors, sensorsClass)
			e.wg.Done()
		}()
		class.Classes = append(class.Classes, sensorsClass)
	}
	if len(mission.Vehicles) > 0 {
		vehsClass := &sqm.Class{
			Name: "Vehicles",
		}
		e.wg.Add(1)
		go func() {
			e.encodeVehicles(mission.Vehicles, vehsClass, &counter)
			e.wg.Done()
		}()
		class.Classes = append(class.Classes, vehsClass)
	}
}

func encodeMissionProperties(mission *Mission, class *sqm.Class) {
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"addOns", sqm.TString, mission.Addons})
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, mission.AddonsAuto})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"randomSeed", sqm.TNumber, mission.RandomSeed})
}

func encodeIntel(i *Intel, class *sqm.Class) {
	var resistanceWest string
	if i.ResistanceWest {
		resistanceWest = "1"
	} else {
		resistanceWest = "0"
	}
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"resistanceWest", sqm.TNumber, resistanceWest})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"startWeather", sqm.TNumber, i.StartWeather})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"forecastWeather", sqm.TNumber, i.ForecastWeather})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"year", sqm.TNumber, i.Year})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"month", sqm.TNumber, i.Month})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"day", sqm.TNumber, i.Day})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"hour", sqm.TNumber, i.Hour})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"minute", sqm.TNumber, i.Minute})
}

func (e *Encoder) encodeVehicles(vehs []*Vehicle, class *sqm.Class, counter *counter) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(vehs))})
	for i, v := range vehs {
		vehClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}

		encodeVehicle(v, vehClass, counter)
		class.Classes = append(class.Classes, vehClass)
	}
}

func encodeSensors(sensors []*Sensor, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(sensors))})
	for i, s := range sensors {
		sensorClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeSensor(s, sensorClass)
		class.Classes = append(class.Classes, sensorClass)
	}
}

func encodeSensor(s *Sensor, class *sqm.Class) {
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, s.Position[:]})
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"synchronizations", sqm.TNumber, s.Synchronizations[:]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"name", sqm.TString, s.Name})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"a", sqm.TNumber, s.Size[0]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"b", sqm.TNumber, s.Size[1]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"angle", sqm.TNumber, s.Angle})
	class.Props = addProp(class.Props, &sqm.Property{"activationBy", sqm.TString, s.ActivationBy})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"activationType", sqm.TString, s.ActivationType})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"timeoutMin", sqm.TNumber, s.TimeoutMin})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"timeoutMid", sqm.TNumber, s.TimeoutMid})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"timeoutMax", sqm.TNumber, s.TimeoutMax})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"type", sqm.TString, s.Type})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"text", sqm.TString, s.Text})
	if s.IsRectangle {
		class.Props = addProp(class.Props, &sqm.Property{"rectangular", sqm.TNumber, "1"})
	}
	if s.IsRepeating {
		class.Props = addProp(class.Props, &sqm.Property{"repeating", sqm.TNumber, "1"})
	}
	if s.IsInterruptible {
		class.Props = addProp(class.Props, &sqm.Property{"interruptable", sqm.TNumber, "1"})
	}
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"age", sqm.TString, s.Age})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"expCond", sqm.TString, s.Condition})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"expActiv", sqm.TString, s.OnActivation})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"expDesactiv", sqm.TString, s.OnDeactivation})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"idVehicle", sqm.TNumber, s.VehicleID})
	if s.Effects != nil {
		effClass := &sqm.Class{
			Name: "Effects",
		}
		class.Classes = append(class.Classes, effClass)
		encodeEffects(s.Effects, effClass)
	}
}

func encodeEffects(e *Effects, class *sqm.Class) {
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"sound", sqm.TString, e.Sound})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"voice", sqm.TString, e.Voice})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"soundDet", sqm.TString, e.SoundDet})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"soundEnv", sqm.TString, e.SoundEnv})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"title", sqm.TString, e.Title})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"titleEffect", sqm.TString, e.TitleEffect})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"titleType", sqm.TString, e.TitleType})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"track", sqm.TString, e.Track})
}

func encodeMarkers(markers []*Marker, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(markers))})
	for i, m := range markers {
		markerClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeMarker(m, markerClass)
		class.Classes = append(class.Classes, markerClass)
	}
}

func encodeMarker(m *Marker, class *sqm.Class) {
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, m.Position[:]})
	class.Props = addProp(class.Props, &sqm.Property{"name", sqm.TString, m.Name})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"angle", sqm.TNumber, m.Angle})
	class.Props = addProp(class.Props, &sqm.Property{"type", sqm.TString, m.Type})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"text", sqm.TString, m.Text})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"markerType", sqm.TString, m.MarkerType})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"colorName", sqm.TString, m.ColorName})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"fillName", sqm.TString, m.FillName})
	var drawBorder string
	if m.DrawBorder {
		drawBorder = "1"
	}
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"drawBorder", sqm.TNumber, drawBorder})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"a", sqm.TNumber, m.Size[0]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"b", sqm.TNumber, m.Size[1]})
}
func (e *Encoder) encodeGroups(groups []*Group, class *sqm.Class, counter *counter) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(groups))})
	for i, g := range groups {
		groupClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		e.wg.Add(1)
		go func(g *Group, groupClass *sqm.Class) {
			e.encodeGroup(g, groupClass, counter)
			e.wg.Done()
		}(g, groupClass)

		class.Classes = append(class.Classes, groupClass)
	}
}

func (e *Encoder) encodeGroup(g *Group, class *sqm.Class, counter *counter) {
	class.Props = addProp(class.Props, &sqm.Property{"side", sqm.TString, g.Side})
	if len(g.Units) > 0 {
		groupMemberClass := &sqm.Class{
			Name: "Vehicles",
		}
		groupMemberClass.Props = append(groupMemberClass.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(g.Units))})
		e.encodeGroupMembers(g.Units, groupMemberClass, counter)
		class.Classes = append(class.Classes, groupMemberClass)
	}

	if len(g.Waypoints) > 0 {
		waypointsClass := &sqm.Class{
			Name: "Waypoints",
		}
		waypointsClass.Props = append(waypointsClass.Props, &sqm.Property{"items", sqm.TNumber, strconv.Itoa(len(g.Waypoints))})
		encodeWaypoints(g.Waypoints, waypointsClass)
		class.Classes = append(class.Classes, waypointsClass)
	}
}

func (e *Encoder) encodeGroupMembers(units []*Vehicle, class *sqm.Class, counter *counter) {
	for i, unit := range units {
		unitclass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}

		encodeVehicle(unit, unitclass, counter)
		class.Classes = append(class.Classes, unitclass)
	}
}

func encodeVehicle(v *Vehicle, class *sqm.Class, counter *counter) {
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"id", sqm.TNumber, strconv.Itoa(int(counter.inc()))})
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, v.Position[:]})
	class.Arrprops = addArrPropOmitEmpty(class.Arrprops, &sqm.ArrayProperty{"markers", sqm.TString, v.Markers[:]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"text", sqm.TString, v.Name})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"azimut", sqm.TNumber, v.Angle})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"vehicle", sqm.TString, v.Classname})
	var leader string
	if v.IsLeader {
		leader = "1"
	}
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"leader", sqm.TNumber, leader})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"special", sqm.TString, v.Special})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"skill", sqm.TNumber, v.Skill})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"player", sqm.TString, v.Player})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"description", sqm.TString, v.Description})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"presence", sqm.TNumber, v.Presence})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"presenceCondition", sqm.TString, v.PresenceCond})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"placement", sqm.TNumber, v.Placement})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"age", sqm.TString, v.Age})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"lock", sqm.TString, v.Lock})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"rank", sqm.TString, v.Rank})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"health", sqm.TNumber, v.Health})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"fuel", sqm.TNumber, v.Fuel})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"ammo", sqm.TNumber, v.Ammo})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"init", sqm.TString, v.Init})
	if v.Side == "" {
		class.Props = addProp(class.Props, &sqm.Property{"side", sqm.TString, "EMPTY"})
	} else {
		class.Props = addProp(class.Props, &sqm.Property{"side", sqm.TString, v.Side})
	}

	if v.ForceHeadlessClient {
		class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"forceHeadlessClient", sqm.TNumber, "1"})
	}
}

func encodeWaypoints(waypoints []*Waypoint, class *sqm.Class) {
	for i, waypoint := range waypoints {
		waypointclass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeWaypoint(waypoint, waypointclass)
		class.Classes = append(class.Classes, waypointclass)
	}
}

func encodeWaypoint(w *Waypoint, class *sqm.Class) {
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, w.Position[:]})
	class.Arrprops = addArrProp(class.Arrprops, &sqm.ArrayProperty{"synchronizations", sqm.TNumber, w.Synchronizations[:]})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"type", sqm.TString, w.Type})
	class.Props = addPropOmitEmpty(class.Props, &sqm.Property{"showWP", sqm.TString, w.ShowWP})
	if w.Effects != nil {
		effClass := &sqm.Class{
			Name: "Effects",
		}
		class.Classes = append(class.Classes, effClass)
		encodeEffects(w.Effects, effClass)
	}
}

func addArrProp(props []*sqm.ArrayProperty, prop *sqm.ArrayProperty) []*sqm.ArrayProperty {
	return append(props, prop)
}

func addArrPropOmitEmpty(props []*sqm.ArrayProperty, prop *sqm.ArrayProperty) []*sqm.ArrayProperty {
	if len(prop.Values) != 0 {
		return append(props, prop)
	}
	return props
}

func addProp(props []*sqm.Property, prop *sqm.Property) []*sqm.Property {
	return append(props, prop)
}
func addPropOmitEmpty(props []*sqm.Property, prop *sqm.Property) []*sqm.Property {
	if prop.Value != "" {
		return append(props, prop)
	}
	return props
}
