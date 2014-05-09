package mission

import (
	"fmt"
	"github.com/blang/gosqm/sqm"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func ShouldContainProp(actual interface{}, expected ...interface{}) string {
	switch expected[0].(type) {
	case *sqm.Property:
		props := actual.([]*sqm.Property)
		prop := expected[0].(*sqm.Property)
		for _, p := range props {
			if p.Name == prop.Name {
				if p.Typ != prop.Typ {
					return fmt.Sprintf("Type differs for property %s, expected: %s, actual: %s", p.Name, prop.Typ, p.Typ)
				}
				if p.Value != prop.Value {
					return fmt.Sprintf("Type differs for property %s, expected: %s, actual: %s", p.Name, prop.Value, p.Value)
				}
				return ""
			}
		}
		return fmt.Sprintf("Could not find property %s", prop.Name)
	case *sqm.ArrayProperty:
		props := actual.([]*sqm.ArrayProperty)
		prop := expected[0].(*sqm.ArrayProperty)
		for _, p := range props {
			if p.Name == prop.Name {
				if p.Typ != prop.Typ {
					return fmt.Sprintf("Type differs for property %s, expected: %s, actual: %s", p.Name, prop.Typ, p.Typ)
				}
				if r := ShouldResemble(p.Values, prop.Values); r != "" {
					return r
				}
				return ""
			}
		}
		return fmt.Sprintf("Could not find property %s", prop.Name)
	}
	return "Invalid type"
}

func ShouldContainClassWithName(actual interface{}, expected ...interface{}) string {
	classes := actual.([]*sqm.Class)
	classname := expected[0].(string)
	for _, c := range classes {
		if c.Name == classname {
			return ""
		}
	}
	return fmt.Sprintf("Class with name %s not found", classname)
}

func TestEncodeIntel(t *testing.T) {
	Convey("Given a fresh intel", t, func() {
		intel := &Intel{
			ResistanceWest:  false,
			StartWeather:    "0.3",
			ForecastWeather: "0.8",
			Year:            "2009",
			Month:           "10",
			Day:             "28",
			Hour:            "6",
			Minute:          "5",
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"startFog", sqm.TNumber, "0.1"},
				},
			},
		}
		Convey("When encoding intel", func() {
			class := &sqm.Class{}
			encodeIntel(intel, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"resistanceWest", sqm.TNumber, "0"})
				So(class.Props, ShouldContainProp, &sqm.Property{"startWeather", sqm.TNumber, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"forecastWeather", sqm.TNumber, "0.8"})
				So(class.Props, ShouldContainProp, &sqm.Property{"year", sqm.TNumber, "2009"})
				So(class.Props, ShouldContainProp, &sqm.Property{"month", sqm.TNumber, "10"})
				So(class.Props, ShouldContainProp, &sqm.Property{"day", sqm.TNumber, "28"})
				So(class.Props, ShouldContainProp, &sqm.Property{"hour", sqm.TNumber, "6"})
				So(class.Props, ShouldContainProp, &sqm.Property{"minute", sqm.TNumber, "5"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"startFog", sqm.TNumber, "0.1"})
			})

		})
	})
}

func TestEncodeMissionProps(t *testing.T) {
	Convey("Given fresh mission properties", t, func() {
		mission := &Mission{
			Addons:     []string{"add1", "add2"},
			AddonsAuto: []string{"add3", "add4"},
			RandomSeed: "1749348",
		}
		Convey("When encoding mission properties", func() {
			class := &sqm.Class{}
			encodeMissionProperties(mission, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOns", sqm.TString, []string{"add1", "add2"}})
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, []string{"add3", "add4"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"randomSeed", sqm.TNumber, "1749348"})
			})
		})
	})
}

func TestEncodeVehicle(t *testing.T) {
	Convey("Given fresh unit/vehicle", t, func() {
		veh := &Vehicle{
			Name:                "name",
			Position:            [3]string{"1.0", "2.0", "3.0"},
			Angle:               "0.3",
			Classname:           "classname",
			Skill:               "0.1",
			Special:             "FORM",
			IsLeader:            true,
			Player:              "PLAYER COMMANDER",
			Description:         "Description",
			Presence:            "0.3",
			PresenceCond:        "true",
			Placement:           "20",
			Age:                 "5 MIN",
			Lock:                "UNLOCKED",
			Rank:                "CORPORAL",
			Health:              "0.1",
			Fuel:                "0.2",
			Ammo:                "0.3",
			Init:                "hint a",
			Side:                "WEST",
			ForceHeadlessClient: true,
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"test", sqm.TString, "init"},
				},
			},
		}
		var idCount counter
		Convey("When encoding vehicle", func() {
			class := &sqm.Class{}
			encodeVehicle(veh, class, &idCount)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"id", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"text", sqm.TString, "name"})
				So(class.Props, ShouldContainProp, &sqm.Property{"azimut", sqm.TNumber, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"vehicle", sqm.TString, "classname"})
				So(class.Props, ShouldContainProp, &sqm.Property{"skill", sqm.TNumber, "0.1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"special", sqm.TString, "FORM"})
				So(class.Props, ShouldContainProp, &sqm.Property{"leader", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"player", sqm.TString, "PLAYER COMMANDER"})
				So(class.Props, ShouldContainProp, &sqm.Property{"description", sqm.TString, "Description"})
				So(class.Props, ShouldContainProp, &sqm.Property{"presence", sqm.TNumber, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"presenceCondition", sqm.TString, "true"})
				So(class.Props, ShouldContainProp, &sqm.Property{"placement", sqm.TNumber, "20"})
				So(class.Props, ShouldContainProp, &sqm.Property{"age", sqm.TString, "5 MIN"})
				So(class.Props, ShouldContainProp, &sqm.Property{"lock", sqm.TString, "UNLOCKED"})
				So(class.Props, ShouldContainProp, &sqm.Property{"rank", sqm.TString, "CORPORAL"})
				So(class.Props, ShouldContainProp, &sqm.Property{"health", sqm.TNumber, "0.1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"fuel", sqm.TNumber, "0.2"})
				So(class.Props, ShouldContainProp, &sqm.Property{"ammo", sqm.TNumber, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"init", sqm.TString, "hint a"})
				So(class.Props, ShouldContainProp, &sqm.Property{"side", sqm.TString, "WEST"})
				So(class.Props, ShouldContainProp, &sqm.Property{"forceHeadlessClient", sqm.TNumber, "1"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"test", sqm.TString, "init"})
			})
		})
	})
}

func TestEncodeWaypoint(t *testing.T) {
	Convey("Given fresh waypoint", t, func() {
		wp := &Waypoint{
			Type:     "AND",
			Position: [3]string{"1.0", "2.0", "3.0"},
			ShowWP:   "NEVER",
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"missing", sqm.TString, "missing"},
				},
			},
			Effects: &Effects{
				Sound:       "sound",
				Voice:       "voice",
				SoundEnv:    "soundenv",
				SoundDet:    "sounddet",
				Track:       "track",
				TitleType:   "titletype",
				Title:       "title",
				TitleEffect: "titleeffect",
			},
		}
		Convey("When encoding waypoint", func() {
			class := &sqm.Class{}
			encodeWaypoint(wp, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"showWP", sqm.TString, "NEVER"})
				So(class.Props, ShouldContainProp, &sqm.Property{"type", sqm.TString, "AND"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
			})
			Convey("Effects class should be set", func() {
				So(len(class.Classes), ShouldEqual, 1)
				effclass := class.Classes[0]
				So(effclass.Name, ShouldEqual, "Effects")
				Convey("All effect attributes should be set correctly", func() {
					So(effclass.Props, ShouldContainProp, &sqm.Property{"sound", sqm.TString, "sound"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"voice", sqm.TString, "voice"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"soundEnv", sqm.TString, "soundenv"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"soundDet", sqm.TString, "sounddet"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"track", sqm.TString, "track"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"titleType", sqm.TString, "titletype"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"title", sqm.TString, "title"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"titleEffect", sqm.TString, "titleeffect"})
				})
			})
		})
	})
}

func TestEncodeMarker(t *testing.T) {
	Convey("Given fresh marker", t, func() {
		m := &Marker{
			Name:       "marker",
			Position:   [3]string{"1.0", "2.0", "3.0"},
			Type:       "Empty",
			MarkerType: "ELLIPSE",
			Text:       "text",
			ColorName:  "ColorRed",
			FillName:   "Border",
			DrawBorder: true,
			Size:       [2]string{"100", "200"},
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"missing", sqm.TString, "missing"},
				},
			},
		}
		Convey("When encoding marker", func() {
			class := &sqm.Class{}
			encodeMarker(m, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "marker"})
				So(class.Props, ShouldContainProp, &sqm.Property{"type", sqm.TString, "Empty"})
				So(class.Props, ShouldContainProp, &sqm.Property{"markerType", sqm.TString, "ELLIPSE"})
				So(class.Props, ShouldContainProp, &sqm.Property{"text", sqm.TString, "text"})
				So(class.Props, ShouldContainProp, &sqm.Property{"colorName", sqm.TString, "ColorRed"})
				So(class.Props, ShouldContainProp, &sqm.Property{"fillName", sqm.TString, "Border"})
				So(class.Props, ShouldContainProp, &sqm.Property{"drawBorder", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"a", sqm.TNumber, "100"})
				So(class.Props, ShouldContainProp, &sqm.Property{"b", sqm.TNumber, "200"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
			})
		})
	})
}

func TestEncodeSensor(t *testing.T) {
	Convey("Given a fresh sensor", t, func() {
		s := &Sensor{
			Name:            "sensor",
			Position:        [3]string{"1.0", "2.0", "3.0"},
			Size:            [2]string{"100", "200"},
			Angle:           "12.3",
			IsRectangle:     true,
			ActivationBy:    "ANY",
			ActivationType:  "GUER D",
			TimeoutMin:      "1",
			TimeoutMid:      "2",
			TimeoutMax:      "3",
			Type:            "EAST G",
			IsRepeating:     true,
			Age:             "UNKNOWN",
			Condition:       "isServer",
			OnActivation:    "hint test",
			OnDeactivation:  "hint test2",
			IsInterruptible: true,
			Effects: &Effects{
				Sound:       "sound",
				Voice:       "voice",
				SoundEnv:    "soundenv",
				SoundDet:    "sounddet",
				Track:       "track",
				TitleType:   "titletype",
				Title:       "title",
				TitleEffect: "titleeffect",
			},
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"missing", sqm.TString, "missing"},
				},
			},
		}
		Convey("When encoding sensor", func() {
			class := &sqm.Class{}
			encodeSensor(s, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "sensor"})
				So(class.Props, ShouldContainProp, &sqm.Property{"a", sqm.TNumber, "100"})
				So(class.Props, ShouldContainProp, &sqm.Property{"b", sqm.TNumber, "200"})
				So(class.Props, ShouldContainProp, &sqm.Property{"angle", sqm.TNumber, "12.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"rectangular", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"activationBy", sqm.TString, "ANY"})
				So(class.Props, ShouldContainProp, &sqm.Property{"activationType", sqm.TString, "GUER D"})
				So(class.Props, ShouldContainProp, &sqm.Property{"timeoutMin", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"timeoutMid", sqm.TNumber, "2"})
				So(class.Props, ShouldContainProp, &sqm.Property{"timeoutMax", sqm.TNumber, "3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"type", sqm.TString, "EAST G"})
				So(class.Props, ShouldContainProp, &sqm.Property{"repeating", sqm.TNumber, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"age", sqm.TString, "UNKNOWN"})
				So(class.Props, ShouldContainProp, &sqm.Property{"expCond", sqm.TString, "isServer"})
				So(class.Props, ShouldContainProp, &sqm.Property{"expActiv", sqm.TString, "hint test"})
				So(class.Props, ShouldContainProp, &sqm.Property{"expDesactiv", sqm.TString, "hint test2"})
				So(class.Props, ShouldContainProp, &sqm.Property{"interruptable", sqm.TNumber, "1"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
			})
			Convey("Effects class should be set", func() {
				So(len(class.Classes), ShouldEqual, 1)
				effclass := class.Classes[0]
				So(effclass.Name, ShouldEqual, "Effects")
				Convey("All effect attributes should be set correctly", func() {
					So(effclass.Props, ShouldContainProp, &sqm.Property{"sound", sqm.TString, "sound"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"voice", sqm.TString, "voice"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"soundEnv", sqm.TString, "soundenv"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"soundDet", sqm.TString, "sounddet"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"track", sqm.TString, "track"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"titleType", sqm.TString, "titletype"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"title", sqm.TString, "title"})
					So(effclass.Props, ShouldContainProp, &sqm.Property{"titleEffect", sqm.TString, "titleeffect"})
				})
			})
		})
	})
}

func TestEncodeMission(t *testing.T) {
	Convey("Given a fresh mission", t, func() {
		m := &Mission{
			Addons:     []string{"add1", "add2"},
			AddonsAuto: []string{"add3", "add4"},
			Intel: &Intel{
				ResistanceWest:  false,
				StartWeather:    "0.2",
				ForecastWeather: "0.3",
				Year:            "2009",
				Month:           "10",
				Day:             "5",
				Hour:            "10",
				Minute:          "3",
			},
			Groups: []*Group{
				&Group{
					Side: "WEST",
					Units: []*Vehicle{
						&Vehicle{
							Name:     "unit",
							Position: [3]string{"1.0", "2.0", "3.0"},
						},
					},
					Waypoints: []*Waypoint{
						&Waypoint{
							Position: [3]string{"1.0", "2.0", "3.0"},
							Type:     "AND",
						},
					},
				},
			},
			Vehicles: []*Vehicle{
				&Vehicle{
					Name:      "veh",
					Position:  [3]string{"1.0", "2.0", "3.0"},
					Classname: "classname",
				},
			},
			Markers: []*Marker{
				&Marker{
					Name:     "marker",
					Position: [3]string{"1.0", "2.0", "3.0"},
					Type:     "empty",
				},
			},
			Sensors: []*Sensor{
				&Sensor{
					Name:     "sensor",
					Position: [3]string{"1.0", "2.0", "3.0"},
					Size:     [2]string{"100", "200"},
				},
			},
		}
		Convey("When encoding mission", func() {
			class := &sqm.Class{}
			e := NewEncoder()
			e.encodeMission(m, class)
			Convey("Intel was set", func() {
				So(class.Classes, ShouldContainClassWithName, "Intel")
			})
			Convey("Addons properties was set", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOns", sqm.TString, []string{"add1", "add2"}})
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, []string{"add3", "add4"}})
			})
			Convey("Groups was set", func() {
				So(class.Classes, ShouldContainClassWithName, "Groups")
			})
			Convey("Vehicles was set", func() {
				So(class.Classes, ShouldContainClassWithName, "Vehicles")
			})
			Convey("Markers was set", func() {
				So(class.Classes, ShouldContainClassWithName, "Markers")
			})
			Convey("Sensors was set", func() {
				So(class.Classes, ShouldContainClassWithName, "Sensors")
			})
		})
	})
}
