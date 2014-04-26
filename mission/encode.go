package mission

import (
	"github.com/blang/gosqm/sqm"
	"strconv"
)

type Encoder struct {
}

func NewEncoder() *Encoder {
	e := &Encoder{}
	return e
}

func (e *Encoder) Encode(missionFile *MissionFile) *sqm.Class {
	return e.encodeMissionFile(missionFile)
}

func (e *Encoder) encodeMissionFile(missionFile *MissionFile) *sqm.Class {
	mainClass := &sqm.Class{
		Name: "mission",
	}
	mainClass.Props = append(mainClass.Props, &sqm.Property{"version", sqm.TInt, missionFile.Version})

	missionClass := &sqm.Class{
		Name: "Mission",
	}
	encodeMission(missionFile.Mission, missionClass)
	mainClass.Classes = append(mainClass.Classes, missionClass)

	introClass := &sqm.Class{
		Name: "Intro",
	}
	encodeMission(missionFile.Intro, introClass)
	mainClass.Classes = append(mainClass.Classes, introClass)

	outroWinClass := &sqm.Class{
		Name: "OutroWin",
	}
	encodeMission(missionFile.OutroWin, outroWinClass)
	mainClass.Classes = append(mainClass.Classes, outroWinClass)

	outroLooseClass := &sqm.Class{
		Name: "OutroLoose",
	}
	encodeMission(missionFile.OutroLoose, outroLooseClass)
	mainClass.Classes = append(mainClass.Classes, outroLooseClass)

	return mainClass
}

func encodeMission(mission *Mission, class *sqm.Class) {
	encodeMissionProperties(mission, class)
	intelClass := &sqm.Class{
		Name: "Intel",
	}
	encodeIntel(mission.Intel, intelClass)
	class.Classes = append(class.Classes, intelClass)

	groupsClass := &sqm.Class{
		Name: "Groups",
	}
	encodeGroups(mission.Groups, groupsClass)
	class.Classes = append(class.Classes, groupsClass)

	markersClass := &sqm.Class{
		Name: "Markers",
	}
	encodeMarkers(mission.Markers, markersClass)
	class.Classes = append(class.Classes, markersClass)

	sensorsClass := &sqm.Class{
		Name: "Sensors",
	}
	encodeSensors(mission.Sensors, sensorsClass)
	class.Classes = append(class.Classes, sensorsClass)

	vehsClass := &sqm.Class{
		Name: "Vehicles",
	}
	encodeVehicles(mission.Vehicles, vehsClass)
	class.Classes = append(class.Classes, vehsClass)

}

func encodeMissionProperties(mission *Mission, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"addOns", sqm.TString, mission.Addons})
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, mission.AddonsAuto})
	if mission.class != nil {
		class.Props = addMissingProps(reg, class.Props, mission.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, mission.class.Arrprops)
	}
}

func encodeIntel(i *Intel, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"resistanceWest", sqm.TInt, i.ResistanceWest})
	class.Props = addProp(reg, class.Props, &sqm.Property{"startWeather", sqm.TFloat, i.StartWeather})
	class.Props = addProp(reg, class.Props, &sqm.Property{"forecastWeather", sqm.TFloat, i.ForecastWeather})
	class.Props = addProp(reg, class.Props, &sqm.Property{"year", sqm.TInt, i.Year})
	class.Props = addProp(reg, class.Props, &sqm.Property{"month", sqm.TInt, i.Month})
	class.Props = addProp(reg, class.Props, &sqm.Property{"day", sqm.TInt, i.Day})
	class.Props = addProp(reg, class.Props, &sqm.Property{"hour", sqm.TInt, i.Hour})
	class.Props = addProp(reg, class.Props, &sqm.Property{"minute", sqm.TInt, i.Minute})

	if i.class != nil {
		class.Props = addMissingProps(reg, class.Props, i.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, i.class.Arrprops)
	}
}

func encodeVehicles(vehs []*Vehicle, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(vehs))})
	for i, v := range vehs {
		vehClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeVehicle(v, vehClass)
		class.Classes = append(class.Classes, vehClass)
	}
}

func encodeVehicle(v *Vehicle, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TFloat, v.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, v.Name})
	class.Props = addProp(reg, class.Props, &sqm.Property{"angle", sqm.TFloat, v.Angle})
	class.Props = addProp(reg, class.Props, &sqm.Property{"vehicle", sqm.TString, v.Classname})
	class.Props = addProp(reg, class.Props, &sqm.Property{"skill", sqm.TFloat, v.Skill})
	class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, v.Side})
	if v.class != nil {
		class.Props = addMissingProps(reg, class.Props, v.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, v.class.Arrprops)
	}
}

func encodeSensors(sensors []*Sensor, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(sensors))})
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
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TFloat, s.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, s.Name})
	class.Props = addProp(reg, class.Props, &sqm.Property{"a", sqm.TFloat, s.Size[0]})
	class.Props = addProp(reg, class.Props, &sqm.Property{"b", sqm.TFloat, s.Size[1]})
	class.Props = addProp(reg, class.Props, &sqm.Property{"angle", sqm.TFloat, s.Angle})
	class.Props = addProp(reg, class.Props, &sqm.Property{"activationBy", sqm.TString, s.ActivationBy})
	if s.IsRectangle {
		class.Props = addProp(reg, class.Props, &sqm.Property{"rectangular", sqm.TInt, "1"})
	}
	if s.IsRepeating {
		class.Props = addProp(reg, class.Props, &sqm.Property{"repeating", sqm.TInt, "1"})
	}
	if s.IsInterruptible {
		class.Props = addProp(reg, class.Props, &sqm.Property{"interruptable", sqm.TInt, "1"})
	}
	class.Props = addProp(reg, class.Props, &sqm.Property{"age", sqm.TString, s.Age})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"expCond", sqm.TString, s.Condition})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"expActiv", sqm.TString, s.OnActivation})
	if s.class != nil {
		class.Props = addMissingProps(reg, class.Props, s.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, s.class.Arrprops)
	}
	if s.classEffects != nil {
		class.Classes = append(class.Classes, s.classEffects)
	}
}

func encodeMarkers(markers []*Marker, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(markers))})
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
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TFloat, m.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, m.Name})
	class.Props = addProp(reg, class.Props, &sqm.Property{"type", sqm.TString, m.Type})
	class.Props = addProp(reg, class.Props, &sqm.Property{"text", sqm.TString, m.Text})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"markerType", sqm.TString, m.MarkerType})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"colorName", sqm.TString, m.ColorName})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"fillName", sqm.TString, m.FillName})
	var drawBorder string
	if m.DrawBorder {
		drawBorder = "1"
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"drawBorder", sqm.TInt, drawBorder})
	class.Props = addProp(reg, class.Props, &sqm.Property{"a", sqm.TFloat, m.Size[0]})
	class.Props = addProp(reg, class.Props, &sqm.Property{"b", sqm.TFloat, m.Size[1]})
	if m.class != nil {
		class.Props = addMissingProps(reg, class.Props, m.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, m.class.Arrprops)
	}

}
func encodeGroups(groups []*Group, class *sqm.Class) {
	class.Props = append(class.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(groups))})
	for i, g := range groups {
		groupClass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeGroup(g, groupClass)
		class.Classes = append(class.Classes, groupClass)
	}
}

func encodeGroup(g *Group, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Props = addProp(reg, class.Props, &sqm.Property{"side", sqm.TString, g.Side})
	if len(g.Units) > 0 {
		groupMemberClass := &sqm.Class{
			Name: "Vehicles",
		}
		groupMemberClass.Props = append(groupMemberClass.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(g.Units))})
		encodeGroupMembers(g.Units, groupMemberClass)
		class.Classes = append(class.Classes, groupMemberClass)
	}

	if len(g.Waypoints) > 0 {
		waypointsClass := &sqm.Class{
			Name: "Waypoints",
		}
		waypointsClass.Props = append(waypointsClass.Props, &sqm.Property{"items", sqm.TInt, strconv.Itoa(len(g.Waypoints))})
		encodeWaypoints(g.Waypoints, waypointsClass)
		class.Classes = append(class.Classes, waypointsClass)
	}

	if g.class != nil {
		class.Props = addMissingProps(reg, class.Props, g.class.Props)
		class.Arrprops = addMissingArrProps(reg, class.Arrprops, g.class.Arrprops)
	}
}

func encodeGroupMembers(units []*Unit, class *sqm.Class) {
	for i, unit := range units {
		unitclass := &sqm.Class{
			Name: "Item" + strconv.Itoa(i),
		}
		encodeUnit(unit, unitclass)
		class.Classes = append(class.Classes, unitclass)
	}
}

func encodeUnit(u *Unit, class *sqm.Class) {
	reg := make(map[string]bool)
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TFloat, u.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"name", sqm.TString, u.Name})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"azimut", sqm.TFloat, u.Direction})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"vehicle", sqm.TString, u.Classname})
	var leader string
	if u.IsLeader {
		leader = "1"
	}
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"leader", sqm.TInt, leader})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"special", sqm.TString, u.Formation})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"skill", sqm.TFloat, u.Skill})

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
	class.Arrprops = addArrProp(reg, class.Arrprops, &sqm.ArrayProperty{"position", sqm.TFloat, w.Position[:]})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"type", sqm.TString, w.Type})
	class.Props = addPropOmitEmpty(reg, class.Props, &sqm.Property{"showWP", sqm.TString, w.ShowWP})
	if w.classEffects != nil {
		class.Classes = append(class.Classes, w.classEffects)
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
