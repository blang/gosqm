package mission

import (
	"github.com/blang/gosqm/sqm"
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
