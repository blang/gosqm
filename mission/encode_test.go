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

func TestEncodeIntel(t *testing.T) {
	Convey("Given a fresh intel", t, func() {
		intel := &Intel{
			ResistanceWest:  "0",
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
