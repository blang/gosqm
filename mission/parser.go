package mission

import (
	"fmt"
	"github.com/blang/gosqm/sqm"
	"sync"
)

type Parser struct {
	wg *sync.WaitGroup
}

func (p *Parser) Parse(class *sqm.Class) (*MissionFile, error) {
	mf := &MissionFile{}
	mf.class = class

	mf.Intro = &Mission{}
	mf.Mission = &Mission{}
	mf.OutroLoose = &Mission{}
	mf.OutroWin = &Mission{}

	p.wg = &sync.WaitGroup{}

	//set version
	for _, val := range class.Props {
		if val.Name == "version" {
			mf.Version = val.Value
		}
	}
	for _, stage := range class.Classes {
		switch stage.Name {
		case "Intro":
			parseMission(stage, mf.Intro)
		case "Mission":
			parseMission(stage, mf.Mission)
		case "OutroWin":
			parseMission(stage, mf.OutroWin)
		case "OutroLoose":
			parseMission(stage, mf.OutroLoose)

		default:
			return nil, fmt.Errorf("Unrecognized base class %s", stage.Name)
		}
	}
	// p.wg.Wait()
	return mf, nil
}

func parseMission(class *sqm.Class, mission *Mission) {
	parseMissionAddons(class, mission)
	for _, baseClass := range class.Classes {
		switch baseClass.Name {
		case "Intel":
			parseIntel(baseClass, mission)
		case "Groups":
			parseGroups(baseClass, mission)
		case "Markers":
			parseMarkers(baseClass, mission)
		case "Sensors":
			parseSensors(baseClass, mission)
		case "Vehicles":
			parseVehicles(baseClass, mission)
		}

	}
}

func parseMissionAddons(class *sqm.Class, mission *Mission) {
	for _, prop := range class.Arrprops {
		switch prop.Name {
		case "addOns":
			mission.Addons = prop.Values
		case "addOnsAuto":
			mission.AddonsAuto = prop.Values
		}
	}
}

func parseIntel(class *sqm.Class, mission *Mission) {
	intel := &Intel{}
	intel.class = class
	for _, prop := range class.Props {
		switch prop.Name {
		case "resistanceWest":
			intel.ResistanceWest = prop.Value == "1"
		case "startWeather":
			intel.StartWeather = prop.Value
		case "forecastWeather":
			intel.ForecastWeather = prop.Value
		case "year":
			intel.Year = prop.Value
		case "month":
			intel.Month = prop.Value
		case "day":
			intel.Day = prop.Value
		case "hour":
			intel.Hour = prop.Value
		case "minute":
			intel.Minute = prop.Value
		}
	}
	mission.Intel = intel
}

func parseGroups(class *sqm.Class, mission *Mission) {
	for _, groupClass := range class.Classes {
		group := &Group{}

		mission.Groups = append(mission.Groups, group)
		parseGroup(groupClass, group)
	}
}

//TODO: Cross Side grouping possible in editor?
func parseGroup(class *sqm.Class, group *Group) {
	group.class = class
	//parse side
	for _, prop := range class.Props {
		if prop.Name == "side" {
			group.Side = prop.Value
		}
	}
	for _, subclass := range class.Classes {
		switch subclass.Name {
		case "Vehicles":
			parseGroupMembers(subclass, group)
		case "Waypoints":
			parseGroupWaypoints(subclass, group)
		}
	}
}

func parseGroupWaypoints(class *sqm.Class, group *Group) {
	for _, wpClass := range class.Classes {
		wp := &Waypoint{}
		group.Waypoints = append(group.Waypoints, wp)
		parseGroupWaypoint(wpClass, wp)
	}
}

func parseGroupWaypoint(class *sqm.Class, wp *Waypoint) {
	wp.class = class
	for _, prop := range class.Props {
		switch prop.Name {
		case "type":
			wp.Type = prop.Value
		case "showWP":
			wp.ShowWP = prop.Value
		}
	}
	for _, arrprop := range class.Arrprops {
		switch arrprop.Name {
		case "position":
			wp.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		}
	}
	if len(wp.class.Classes) > 0 && wp.class.Classes[0].Name == "Effects" {
		wp.classEffects = wp.class.Classes[0]
	}
}

func parseGroupMembers(class *sqm.Class, group *Group) {
	for _, unitClass := range class.Classes {
		unit := &Unit{}
		group.Units = append(group.Units, unit)
		parseGroupMember(unitClass, unit)
	}
}

func parseGroupMember(class *sqm.Class, unit *Unit) {
	unit.class = class
	for _, prop := range class.Props {
		switch prop.Name {
		case "name":
			unit.Name = prop.Value
		case "vehicle":
			unit.Classname = prop.Value
		case "skill":
			unit.Skill = prop.Value
		case "azimut":
			unit.Direction = prop.Value
		case "special":
			unit.Formation = prop.Value
		case "leader":
			unit.IsLeader = true
		}
	}
	for _, arrprop := range class.Arrprops {
		switch arrprop.Name {
		case "position":
			unit.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		}
	}
}

func parseMarkers(class *sqm.Class, mission *Mission) {
	for _, markerClass := range class.Classes {
		marker := &Marker{}
		mission.Markers = append(mission.Markers, marker)
		parseMarker(markerClass, marker)
	}
}

func parseMarker(c *sqm.Class, marker *Marker) {
	marker.class = c
	for _, prop := range c.Props {
		switch prop.Name {
		case "name":
			marker.Name = prop.Value
		case "text":
			marker.Text = prop.Value
		case "type":
			marker.Type = prop.Value
		case "markerType":
			marker.MarkerType = prop.Value
		case "colorName":
			marker.ColorName = prop.Value
		case "fillName":
			marker.FillName = prop.Value
		case "a":
			marker.Size[0] = prop.Value
		case "b":
			marker.Size[1] = prop.Value
		case "drawBorder":
			marker.DrawBorder = prop.Value == "1"
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			marker.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		}
	}
}

func parseSensors(class *sqm.Class, mission *Mission) {
	for _, sensorClass := range class.Classes {
		sensor := &Sensor{}
		mission.Sensors = append(mission.Sensors, sensor)
		parseSensor(sensorClass, sensor)
	}
}

func parseSensor(c *sqm.Class, sensor *Sensor) {
	sensor.class = c
	for _, prop := range c.Props {
		switch prop.Name {
		case "name":
			sensor.Name = prop.Value
		case "a":
			sensor.Size[0] = prop.Value
		case "b":
			sensor.Size[1] = prop.Value
		case "angle":
			sensor.Angle = prop.Value
		case "rectangular":
			sensor.IsRectangle = prop.Value == "1"
		case "activationBy":
			sensor.ActivationBy = prop.Value
		case "repeating":
			sensor.IsRepeating = prop.Value == "1"
		case "age":
			sensor.Age = prop.Value
		case "expCond":
			sensor.Condition = prop.Value
		case "expActiv":
			sensor.OnActivation = prop.Value
		case "interruptable":
			sensor.IsInterruptible = prop.Value == "1"
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			sensor.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		}
	}
	if len(sensor.class.Classes) > 0 && sensor.class.Classes[0].Name == "Effects" {
		sensor.classEffects = sensor.class.Classes[0]
	}
}

func parseVehicles(class *sqm.Class, mission *Mission) {
	for _, vehClass := range class.Classes {
		veh := &Vehicle{}
		mission.Vehicles = append(mission.Vehicles, veh)
		parseVehicle(vehClass, veh)
	}
}

func parseVehicle(c *sqm.Class, veh *Vehicle) {
	veh.class = c
	for _, prop := range c.Props {
		switch prop.Name {
		case "name":
			veh.Name = prop.Value
		case "azimut":
			veh.Angle = prop.Value
		case "vehicle":
			veh.Classname = prop.Value
		case "side":
			veh.Side = prop.Value
		case "skill":
			veh.Skill = prop.Value
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			veh.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		}
	}
}
