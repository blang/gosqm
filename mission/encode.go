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
