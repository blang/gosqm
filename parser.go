package mission

import (
	"fmt"
	"github.com/blang/gosqm/sqm"
	"sync"
)

type Context string

func (c Context) String() string {
	return string(c)
}

const (
	ContextUnit            Context = "Unit"
	ContextMission                 = "Mission"
	ContextMissionFile             = "MissionFile"
	ContextIntel                   = "Intel"
	ContextGroup                   = "Group"
	ContextWaypoint                = "Waypoint"
	ContextWaypointEffects         = "WaypointEffects"
	ContextMarker                  = "Marker"
	ContextSensor                  = "Sensor"
	ContextSensorEffects           = "SensorEffects"
	ContextVehicle                 = "Vehicle"
)

type UnkownPropertyError struct {
	ParentClass   *sqm.Class
	Property      *sqm.Property
	ArrayProperty *sqm.ArrayProperty
	Context       Context
}

func (e *UnkownPropertyError) Error() string {
	var propName string
	if e.Property != nil {
		propName = e.Property.Name
	} else if e.ArrayProperty != nil {
		propName = e.ArrayProperty.Name
	}
	if propName != "" && e.ParentClass != nil {
		return "Unknown property " + propName + " in class " + e.ParentClass.Name + " in context " + e.Context.String()
	} else {
		return "Unkown property"
	}
}

type UnkownClassError struct {
	ParentClass *sqm.Class
	Class       *sqm.Class
	Context     Context
}

func (e *UnkownClassError) Error() string {
	if e.ParentClass != nil && e.Class != nil {
		return "Unknown class " + e.Class.Name + " in class " + e.ParentClass.Name + " in context " + e.Context.String()
	} else {
		return "Unkown property"
	}
}

type Parser struct {
	wg     *sync.WaitGroup
	errors []error
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(class *sqm.Class) (*MissionFile, error) {
	p.errors = nil
	if class == nil {
		return nil, fmt.Errorf("Class was nil")
	}

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
			p.parseMission(stage, mf.Intro)
		case "Mission":
			p.parseMission(stage, mf.Mission)
		case "OutroWin":
			p.parseMission(stage, mf.OutroWin)
		case "OutroLoose":
			p.parseMission(stage, mf.OutroLoose)

		default:
			p.saveError(&UnkownClassError{
				ParentClass: class,
				Class:       stage,
				Context:     ContextMissionFile,
			})
		}
	}
	// p.wg.Wait()
	return mf, nil
}

func (p *Parser) Warnings() []error {
	return p.errors
}

func (p *Parser) saveError(e error) {
	p.errors = append(p.errors, e)
}

func (p *Parser) parseMission(class *sqm.Class, mission *Mission) {
	p.parseMissionProps(class, mission)
	for _, baseClass := range class.Classes {
		switch baseClass.Name {
		case "Intel":
			p.parseIntel(baseClass, mission)
		case "Groups":
			p.parseGroups(baseClass, mission)
		case "Markers":
			p.parseMarkers(baseClass, mission)
		case "Sensors":
			p.parseSensors(baseClass, mission)
		case "Vehicles":
			p.parseVehicles(baseClass, mission)
		default:
			p.saveError(&UnkownClassError{
				ParentClass: class,
				Class:       baseClass,
				Context:     ContextMission,
			})
		}

	}
}

func (p *Parser) parseMissionProps(class *sqm.Class, mission *Mission) {
	for _, prop := range class.Arrprops {
		switch prop.Name {
		case "addOns":
			mission.Addons = prop.Values
		case "addOnsAuto":
			mission.AddonsAuto = prop.Values
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   class,
				ArrayProperty: prop,
				Context:       ContextMission,
			})
		}
	}
	for _, prop := range class.Props {
		switch prop.Name {
		case "randomSeed":
			mission.RandomSeed = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: class,
				Property:    prop,
				Context:     ContextMission,
			})
		}
	}
}

func (p *Parser) parseIntel(class *sqm.Class, mission *Mission) {
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
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: class,
				Property:    prop,
				Context:     ContextIntel,
			})
		}
	}
	mission.Intel = intel
}

func (p *Parser) parseGroups(class *sqm.Class, mission *Mission) {
	for _, groupClass := range class.Classes {
		group := &Group{}

		mission.Groups = append(mission.Groups, group)
		p.parseGroup(groupClass, group)
	}
}

//TODO: Cross Side grouping possible in editor?
func (p *Parser) parseGroup(class *sqm.Class, group *Group) {
	group.class = class
	//parse side
	for _, prop := range class.Props {
		switch prop.Name {
		case "side":
			group.Side = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: class,
				Property:    prop,
				Context:     ContextGroup,
			})
		}

	}
	for _, subclass := range class.Classes {
		switch subclass.Name {
		case "Vehicles":
			p.parseGroupMembers(subclass, group)
		case "Waypoints":
			p.parseGroupWaypoints(subclass, group)
		default:
			p.saveError(&UnkownClassError{
				ParentClass: class,
				Class:       subclass,
				Context:     ContextGroup,
			})
		}
	}
}

func (p *Parser) parseGroupWaypoints(class *sqm.Class, group *Group) {
	for _, wpClass := range class.Classes {
		wp := &Waypoint{}
		group.Waypoints = append(group.Waypoints, wp)
		p.parseGroupWaypoint(wpClass, wp)
	}
}

func (p *Parser) parseGroupWaypoint(class *sqm.Class, wp *Waypoint) {
	wp.class = class
	for _, prop := range class.Props {
		switch prop.Name {
		case "type":
			wp.Type = prop.Value
		case "showWP":
			wp.ShowWP = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: class,
				Property:    prop,
				Context:     ContextWaypoint,
			})
		}
	}
	for _, arrprop := range class.Arrprops {
		switch arrprop.Name {
		case "position":
			wp.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   class,
				ArrayProperty: arrprop,
				Context:       ContextWaypoint,
			})
		}
	}
	if len(class.Classes) > 0 {
		for _, subclass := range class.Classes {
			switch subclass.Name {
			case "Effects":
				effects := &Effects{}
				wp.Effects = effects
				p.parseEffects(subclass, effects)
			default:
				p.saveError(&UnkownClassError{
					ParentClass: class,
					Class:       subclass,
					Context:     ContextWaypointEffects,
				})
			}
		}
	}
}

func (p *Parser) parseGroupMembers(class *sqm.Class, group *Group) {
	for _, unitClass := range class.Classes {
		unit := &Unit{}
		group.Units = append(group.Units, unit)
		p.parseGroupMember(unitClass, unit)
	}
}

func (p *Parser) parseGroupMember(class *sqm.Class, unit *Unit) {
	unit.class = class
	for _, prop := range class.Props {
		switch prop.Name {
		case "id":
			// autogenerated
		case "text":
			unit.Name = prop.Value
		case "vehicle":
			unit.Classname = prop.Value
		case "skill":
			unit.Skill = prop.Value
		case "azimut":
			unit.Direction = prop.Value
		case "special":
			unit.Special = prop.Value
		case "leader":
			unit.IsLeader = prop.Value == "1"
		case "player":
			unit.Player = prop.Value
		case "description":
			unit.Description = prop.Value
		case "presence":
			unit.Presence = prop.Value
		case "presenceCondition":
			unit.PresenceCond = prop.Value
		case "placement":
			unit.Placement = prop.Value
		case "age":
			unit.Age = prop.Value
		case "lock":
			unit.Lock = prop.Value
		case "rank":
			unit.Rank = prop.Value
		case "health":
			unit.Health = prop.Value
		case "fuel":
			unit.Fuel = prop.Value
		case "ammo":
			unit.Ammo = prop.Value
		case "init":
			unit.Init = prop.Value
		case "side":
			unit.Side = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: class,
				Property:    prop,
				Context:     ContextUnit,
			})
		}
	}
	for _, arrprop := range class.Arrprops {
		switch arrprop.Name {
		case "position":
			unit.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   class,
				ArrayProperty: arrprop,
				Context:       ContextUnit,
			})
		}
	}
}

func (p *Parser) parseMarkers(class *sqm.Class, mission *Mission) {
	for _, markerClass := range class.Classes {
		marker := &Marker{}
		mission.Markers = append(mission.Markers, marker)
		p.parseMarker(markerClass, marker)
	}
}

func (p *Parser) parseMarker(c *sqm.Class, marker *Marker) {
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
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: c,
				Property:    prop,
				Context:     ContextMarker,
			})
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			marker.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   c,
				ArrayProperty: arrprop,
				Context:       ContextMarker,
			})
		}
	}
}

func (p *Parser) parseSensors(class *sqm.Class, mission *Mission) {
	for _, sensorClass := range class.Classes {
		sensor := &Sensor{}
		mission.Sensors = append(mission.Sensors, sensor)
		p.parseSensor(sensorClass, sensor)
	}
}

func (p *Parser) parseSensor(c *sqm.Class, sensor *Sensor) {
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
		case "activationType":
			sensor.ActivationType = prop.Value
		case "timeoutMin":
			sensor.TimeoutMin = prop.Value
		case "timeoutMid":
			sensor.TimeoutMid = prop.Value
		case "timeoutMax":
			sensor.TimeoutMax = prop.Value
		case "type":
			sensor.Type = prop.Value
		case "repeating":
			sensor.IsRepeating = prop.Value == "1"
		case "age":
			sensor.Age = prop.Value
		case "expCond":
			sensor.Condition = prop.Value
		case "expActiv":
			sensor.OnActivation = prop.Value
		case "expDesactiv":
			sensor.OnDeactivation = prop.Value
		case "interruptable":
			sensor.IsInterruptible = prop.Value == "1"
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: c,
				Property:    prop,
				Context:     ContextSensor,
			})
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			sensor.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   c,
				ArrayProperty: arrprop,
				Context:       ContextSensor,
			})
		}
	}
	if len(c.Classes) > 0 {
		for _, subclass := range c.Classes {
			switch subclass.Name {
			case "Effects":
				effects := &Effects{}
				sensor.Effects = effects
				p.parseEffects(c.Classes[0], effects)
			default:
				p.saveError(&UnkownClassError{
					ParentClass: c,
					Class:       subclass,
					Context:     ContextSensor,
				})
			}
		}
	}
}

func (p *Parser) parseEffects(c *sqm.Class, effects *Effects) {
	for _, prop := range c.Props {
		switch prop.Name {
		case "sound":
			effects.Sound = prop.Value
		case "voice":
			effects.Voice = prop.Value
		case "soundEnv":
			effects.SoundEnv = prop.Value
		case "soundDet":
			effects.SoundDet = prop.Value
		case "title":
			effects.Title = prop.Value
		case "titleType":
			effects.TitleType = prop.Value
		case "titleEffect":
			effects.TitleEffect = prop.Value
		case "track":
			effects.Track = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: c,
				Property:    prop,
				Context:     ContextSensorEffects,
			})
		}
	}
}

func (p *Parser) parseVehicles(class *sqm.Class, mission *Mission) {
	for _, vehClass := range class.Classes {
		veh := &Vehicle{}
		mission.Vehicles = append(mission.Vehicles, veh)
		p.parseVehicle(vehClass, veh)
	}
}

func (p *Parser) parseVehicle(c *sqm.Class, veh *Vehicle) {
	veh.class = c
	for _, prop := range c.Props {
		switch prop.Name {
		case "id":
			// autogenerated
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
		case "presence":
			veh.Presence = prop.Value
		case "presenceCondition":
			veh.PresenceCond = prop.Value
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass: c,
				Property:    prop,
				Context:     ContextVehicle,
			})
		}
	}
	for _, arrprop := range c.Arrprops {
		switch arrprop.Name {
		case "position":
			veh.Position = [3]string{arrprop.Values[0], arrprop.Values[1], arrprop.Values[2]}
		default:
			p.saveError(&UnkownPropertyError{
				ParentClass:   c,
				ArrayProperty: arrprop,
				Context:       ContextVehicle,
			})
		}
	}
}
