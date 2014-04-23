package mission

import (
	"fmt"
	sqm "github.com/blang/gosqm/sqmparser"
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
		case "Mission":
			parseMission(stage, mf.Mission)
		case "Intro":
			parseMission(stage, mf.Intro)
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
	for _, base := range class.Classes {
		switch base.Name {
		case "Groups":
			parseGroups(base, mission)
		}
	}
}

func parseGroups(class *sqm.Class, Mission *Mission) {
	for _, groupClass := range class.Classes {
		group := &Group{}
		group.class = groupClass
		Mission.Groups = append(Mission.Groups, group)
		parseGroup(groupClass, group)
	}
}

//TODO: Cross Side grouping possible in editor?
func parseGroup(class *sqm.Class, group *Group) {
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
		}
	}
}

func parseGroupMembers(class *sqm.Class, group *Group) {
	for _, unitClass := range class.Classes {
		unit := &Unit{}
		unit.class = unitClass
		group.Units = append(group.Units, unit)
		parseGroupMember(class, unit)
	}
}

func parseGroupMember(class *sqm.Class, unit *Unit) {
	for _, prop := range class.Props {
		switch prop.Name {
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
