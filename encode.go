package mission

import (
	"github.com/blang/gosqm/sqm"
	"strconv"
	"sync"
	"sync/atomic"
)

type Encoder struct {
	wg *sync.WaitGroup
}

type counter int32

func (c *counter) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}
func (c *counter) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func NewEncoder() *Encoder {
	e := &Encoder{
		wg: &sync.WaitGroup{},
	}
	return e
}

func (e *Encoder) Encode(missionFile *MissionFile) *sqm.Class {
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

	groupsClass := &sqm.Class{
		Name: "Groups",
	}
	e.wg.Add(1)
	go func() {
		e.encodeGroups(mission.Groups, groupsClass, &counter)
		e.wg.Done()
	}()
	class.Classes = append(class.Classes, groupsClass)

	markersClass := &sqm.Class{
		Name: "Markers",
	}
	e.wg.Add(1)
	go func() {
		encodeMarkers(mission.Markers, markersClass)
		e.wg.Done()
	}()
	class.Classes = append(class.Classes, markersClass)

	sensorsClass := &sqm.Class{
		Name: "Sensors",
	}
	e.wg.Add(1)
	go func() {
		encodeSensors(mission.Sensors, sensorsClass)
		e.wg.Done()
	}()
	class.Classes = append(class.Classes, sensorsClass)

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

func encodeMissionProperties(mission *Mission, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"addOns", sqm.TString, mission.Addons})
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, mission.AddonsAuto})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"randomSeed", sqm.TNumber, mission.RandomSeed})
}

func encodeIntel(i *Intel, class *sqm.Class) {
	reg := make(map[string]bool)
	var resistanceWest string
	if i.ResistanceWest {
		resistanceWest = "1"
	} else {
		resistanceWest = "0"
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"resistanceWest", sqm.TNumber, resistanceWest})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"startWeather", sqm.TNumber, i.StartWeather})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"forecastWeather", sqm.TNumber, i.ForecastWeather})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"year", sqm.TNumber, i.Year})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"month", sqm.TNumber, i.Month})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"day", sqm.TNumber, i.Day})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"hour", sqm.TNumber, i.Hour})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"minute", sqm.TNumber, i.Minute})

	if i.class != nil {
		class.Props = addMissingProps(reg, class.Props, i.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, i.class.Arrprops)
	}
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

func encodeVehicle(v *Vehicle, class *sqm.Class, counter *counter) {
	reg := make(map[string]bool)
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"id", sqm.TNumber, strconv.Itoa(int(counter.inc()))})
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, v.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, v.Name})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"angle", sqm.TNumber, v.Angle})
	class.Props = addProp(reg, class.Props, &sqm.Property{"vehicle", sqm.TString, v.Classname})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"skill", sqm.TNumber, v.Skill})
	if v.Side == "" {
		class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, "EMPTY"})
	} else {
		class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, v.Side})
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"presence", sqm.TNumber, v.Presence})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"presenceCondition", sqm.TString, v.PresenceCond})
	if v.class != nil {
		class.Props = addMissingProps(reg, class.Props, v.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, v.class.Arrprops)
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
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, s.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, s.Name})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"a", sqm.TNumber, s.Size[0]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"b", sqm.TNumber, s.Size[1]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"angle", sqm.TNumber, s.Angle})
	class.Props = addProp(reg, class.Props, &sqm.Property{"activationBy", sqm.TString, s.ActivationBy})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"activationType", sqm.TString, s.ActivationType})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"timeoutMin", sqm.TNumber, s.TimeoutMin})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"timeoutMid", sqm.TNumber, s.TimeoutMid})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"timeoutMax", sqm.TNumber, s.TimeoutMax})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"type", sqm.TString, s.Type})
	if s.IsRectangle {
		class.Props = addProp(reg, class.Props, &sqm.Property{"rectangular", sqm.TNumber, "1"})
	}
	if s.IsRepeating {
		class.Props = addProp(reg, class.Props, &sqm.Property{"repeating", sqm.TNumber, "1"})
	}
	if s.IsInterruptible {
		class.Props = addProp(reg, class.Props, &sqm.Property{"interruptable", sqm.TNumber, "1"})
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"age", sqm.TString, s.Age})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"expCond", sqm.TString, s.Condition})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"expActiv", sqm.TString, s.OnActivation})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"expDesactiv", sqm.TString, s.OnDeactivation})
	if s.class != nil {
		class.Props = addMissingProps(reg, class.Props, s.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, s.class.Arrprops)
	}
	if s.Effects != nil {
		effClass := &sqm.Class{
			Name: "Effects",
		}
		class.Classes = append(class.Classes, effClass)
		encodeEffects(s.Effects, effClass)
	}
}

func encodeEffects(e *Effects, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"sound", sqm.TString, e.Sound})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"voice", sqm.TString, e.Voice})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"soundDet", sqm.TString, e.SoundDet})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"soundEnv", sqm.TString, e.SoundEnv})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"title", sqm.TString, e.Title})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"titleEffect", sqm.TString, e.TitleEffect})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"titleType", sqm.TString, e.TitleType})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"track", sqm.TString, e.Track})
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
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, m.Position[:]})
	class.Props = addProp(reg, class.Props, &sqm.Property{"name", sqm.TString, m.Name})
	class.Props = addProp(reg, class.Props, &sqm.Property{"type", sqm.TString, m.Type})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"text", sqm.TString, m.Text})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"markerType", sqm.TString, m.MarkerType})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"colorName", sqm.TString, m.ColorName})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"fillName", sqm.TString, m.FillName})
	var drawBorder string
	if m.DrawBorder {
		drawBorder = "1"
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"drawBorder", sqm.TNumber, drawBorder})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"a", sqm.TNumber, m.Size[0]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"b", sqm.TNumber, m.Size[1]})
	if m.class != nil {
		class.Props = addMissingProps(reg, class.Props, m.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, m.class.Arrprops)
	}

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
	reg := make(map[string]bool)
	class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, g.Side})
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

	if g.class != nil {
		class.Props = addMissingProps(reg, class.Props, g.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, g.class.Arrprops)
	}
}

func (e *Encoder) encodeGroupMembers(units []*Unit, class *sqm.Class, counter *counter) {
	for i, unit := range units {
		unitclass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}

		encodeUnit(unit, unitclass, counter)
		class.Classes = append(class.Classes, unitclass)
	}
}

func encodeUnit(u *Unit, class *sqm.Class, counter *counter) {
	reg := make(map[string]bool)
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"id", sqm.TNumber, strconv.Itoa(int(counter.inc()))})
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, u.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"text", sqm.TString, u.Name})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"azimut", sqm.TNumber, u.Direction})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"vehicle", sqm.TString, u.Classname})
	var leader string
	if u.IsLeader {
		leader = "1"
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"leader", sqm.TNumber, leader})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"special", sqm.TString, u.Special})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"skill", sqm.TNumber, u.Skill})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"player", sqm.TString, u.Player})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"description", sqm.TString, u.Description})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"presence", sqm.TNumber, u.Presence})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"presenceCondition", sqm.TString, u.PresenceCond})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"placement", sqm.TNumber, u.Placement})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"age", sqm.TString, u.Age})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"lock", sqm.TString, u.Lock})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"rank", sqm.TString, u.Rank})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"health", sqm.TNumber, u.Health})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"fuel", sqm.TNumber, u.Fuel})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"ammo", sqm.TNumber, u.Ammo})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"init", sqm.TString, u.Init})
	class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, u.Side})

	if u.class != nil {
		class.Props = addMissingProps(reg, class.Props, u.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, u.class.Arrprops)
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
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TNumber, w.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"type", sqm.TString, w.Type})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"showWP", sqm.TString, w.ShowWP})
	if w.Effects != nil {
		effClass := &sqm.Class{
			Name: "Effects",
		}
		class.Classes = append(class.Classes, effClass)
		encodeEffects(w.Effects, effClass)
	}
	if w.class != nil {
		class.Props = addMissingProps(reg, class.Props, w.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, w.class.Arrprops)
	}
}

func addMissingProps(register map[string]bool, existingProps []*sqm.Property, additional []*sqm.Property) []*sqm.Property {
	for _, p := range additional {
		if !register[p.Name] {
			existingProps = addProp(register, existingProps, p)
		}
	}
	return existingProps
}

func addMissingArrProps(register map[string]bool, existingProps []*sqm.ArrayProperty, additional []*sqm.ArrayProperty) []*sqm.ArrayProperty {
	for _, p := range additional {
		if !register[p.Name] {
			existingProps = addArrProp(register, existingProps, p)
		}
	}
	return existingProps
}

func addArrProp(register map[string]bool, props []*sqm.ArrayProperty, prop *sqm.ArrayProperty) []*sqm.ArrayProperty {
	register[prop.Name] = true
	return append(props, prop)
}

func addArrPropOmitEmpty(register map[string]bool, props []*sqm.ArrayProperty, prop *sqm.ArrayProperty) []*sqm.ArrayProperty {
	register[prop.Name] = true
	if len(prop.Values) != 0 {
		return append(props, prop)
	}
	return props
}

func addProp(register map[string]bool, props []*sqm.Property, prop *sqm.Property) []*sqm.Property {
	register[prop.Name] = true
	return append(props, prop)
}
func addPropOmitEmpty(register map[string]bool, props []*sqm.Property, prop *sqm.Property) []*sqm.Property {
	register[prop.Name] = true
	if prop.Value != "" {
		return append(props, prop)
	}
	return props
}
