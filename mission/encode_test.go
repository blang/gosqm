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
					&sqm.Property{"startFog", sqm.TFloat, "0.1"},
				},
			},
		}
		Convey("When encoding intel", func() {
			class := &sqm.Class{}
			encodeIntel(intel, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"resistanceWest", sqm.TInt, "0"})
				So(class.Props, ShouldContainProp, &sqm.Property{"startWeather", sqm.TFloat, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"forecastWeather", sqm.TFloat, "0.8"})
				So(class.Props, ShouldContainProp, &sqm.Property{"year", sqm.TInt, "2009"})
				So(class.Props, ShouldContainProp, &sqm.Property{"month", sqm.TInt, "10"})
				So(class.Props, ShouldContainProp, &sqm.Property{"day", sqm.TInt, "28"})
				So(class.Props, ShouldContainProp, &sqm.Property{"hour", sqm.TInt, "6"})
				So(class.Props, ShouldContainProp, &sqm.Property{"minute", sqm.TInt, "5"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"startFog", sqm.TFloat, "0.1"})
			})

		})
	})
}

func TestEncodeMissionProps(t *testing.T) {
	Convey("Given fresh mission properties", t, func() {
		mission := &Mission{
			Addons:     []string{"add1", "add2"},
			AddonsAuto: []string{"add3", "add4"},
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"randomSeed", sqm.TInt, "1749348"},
				},
			},
		}
		Convey("When encoding mission properties", func() {
			class := &sqm.Class{}
			encodeMissionProperties(mission, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOns", sqm.TString, []string{"add1", "add2"}})
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"addOnsAuto", sqm.TString, []string{"add3", "add4"}})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"randomSeed", sqm.TInt, "1749348"})
			})
		})
	})
}

func TestEncodeUnit(t *testing.T) {
	Convey("Given fresh unit", t, func() {
		unit := &Unit{
			Name:      "name",
			Position:  [3]string{"1.0", "2.0", "3.0"},
			Direction: "0.3",
			Classname: "classname",
			Skill:     "0.1",
			Formation: "FORM",
			IsLeader:  true,
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"init", sqm.TString, "init"},
				},
			},
		}
		Convey("When encoding unit", func() {
			class := &sqm.Class{}
			encodeUnit(unit, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "name"})
				So(class.Props, ShouldContainProp, &sqm.Property{"azimut", sqm.TFloat, "0.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"vehicle", sqm.TString, "classname"})
				So(class.Props, ShouldContainProp, &sqm.Property{"skill", sqm.TFloat, "0.1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"special", sqm.TString, "FORM"})
				So(class.Props, ShouldContainProp, &sqm.Property{"leader", sqm.TInt, "1"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"init", sqm.TString, "init"})
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
			classEffects: &sqm.Class{
				Name: "Effects",
			},
		}
		Convey("When encoding waypoint", func() {
			class := &sqm.Class{}
			encodeWaypoint(wp, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"showWP", sqm.TString, "NEVER"})
				So(class.Props, ShouldContainProp, &sqm.Property{"type", sqm.TString, "AND"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
			})
			Convey("Effects class was set", func() {
				So(len(class.Classes), ShouldBeGreaterThan, 0)
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
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "marker"})
				So(class.Props, ShouldContainProp, &sqm.Property{"type", sqm.TString, "Empty"})
				So(class.Props, ShouldContainProp, &sqm.Property{"markerType", sqm.TString, "ELLIPSE"})
				So(class.Props, ShouldContainProp, &sqm.Property{"text", sqm.TString, "text"})
				So(class.Props, ShouldContainProp, &sqm.Property{"colorName", sqm.TString, "ColorRed"})
				So(class.Props, ShouldContainProp, &sqm.Property{"fillName", sqm.TString, "Border"})
				So(class.Props, ShouldContainProp, &sqm.Property{"drawBorder", sqm.TInt, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"a", sqm.TFloat, "100"})
				So(class.Props, ShouldContainProp, &sqm.Property{"b", sqm.TFloat, "200"})
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
			IsRepeating:     true,
			Age:             "UNKNOWN",
			Condition:       "isServer",
			OnActivation:    "hint test",
			IsInterruptible: true,
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"missing", sqm.TString, "missing"},
				},
			},
			classEffects: &sqm.Class{
				Name: "Effects",
			},
		}
		Convey("When encoding sensor", func() {
			class := &sqm.Class{}
			encodeSensor(s, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "sensor"})
				So(class.Props, ShouldContainProp, &sqm.Property{"a", sqm.TFloat, "100"})
				So(class.Props, ShouldContainProp, &sqm.Property{"b", sqm.TFloat, "200"})
				So(class.Props, ShouldContainProp, &sqm.Property{"angle", sqm.TFloat, "12.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"rectangular", sqm.TInt, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"activationBy", sqm.TString, "ANY"})
				So(class.Props, ShouldContainProp, &sqm.Property{"repeating", sqm.TInt, "1"})
				So(class.Props, ShouldContainProp, &sqm.Property{"age", sqm.TString, "UNKNOWN"})
				So(class.Props, ShouldContainProp, &sqm.Property{"expCond", sqm.TString, "isServer"})
				So(class.Props, ShouldContainProp, &sqm.Property{"expActiv", sqm.TString, "hint test"})
				So(class.Props, ShouldContainProp, &sqm.Property{"interruptable", sqm.TInt, "1"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
			})
			Convey("Effects class should be set", func() {
				So(len(class.Classes), ShouldEqual, 1)
			})
		})
	})
}

func TestEncodeVehicle(t *testing.T) {
	Convey("Given a fresh vehicle", t, func() {
		v := &Vehicle{
			Name:      "vehicle",
			Position:  [3]string{"1.0", "2.0", "3.0"},
			Angle:     "12.3",
			Classname: "classname",
			Skill:     "0.2",
			Side:      "EMPTY",
			class: &sqm.Class{
				Props: []*sqm.Property{
					&sqm.Property{"missing", sqm.TString, "missing"},
				},
			},
		}
		Convey("When encoding vehicle", func() {
			class := &sqm.Class{}
			encodeVehicle(v, class)
			Convey("Class properties should be set correctly", func() {
				So(class.Arrprops, ShouldContainProp, &sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}})
				So(class.Props, ShouldContainProp, &sqm.Property{"name", sqm.TString, "vehicle"})
				So(class.Props, ShouldContainProp, &sqm.Property{"angle", sqm.TFloat, "12.3"})
				So(class.Props, ShouldContainProp, &sqm.Property{"vehicle", sqm.TString, "classname"})
				So(class.Props, ShouldContainProp, &sqm.Property{"skill", sqm.TFloat, "0.2"})
				So(class.Props, ShouldContainProp, &sqm.Property{"side", sqm.TString, "EMPTY"})
			})
			Convey("Missing properties should be taken from parent class", func() {
				So(class.Props, ShouldContainProp, &sqm.Property{"missing", sqm.TString, "missing"})
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
					Units: []*Unit{
						&Unit{
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
			encodeMission(m, class)
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
