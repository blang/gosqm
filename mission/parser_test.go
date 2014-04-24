package mission

import (
	"github.com/blang/gosqm/sqm"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseGroups(t *testing.T) {
	Convey("Given a valid groups class with subclasses", t, func() {
		unitclass := &sqm.Class{
			Props: []*sqm.Property{
				&sqm.Property{"azimut", sqm.TFloat, "12.3"},
				&sqm.Property{"vehicle", sqm.TString, "classname"},
				&sqm.Property{"leader", sqm.TInt, "1"},
				&sqm.Property{"special", sqm.TString, "FORM"},
				&sqm.Property{"skill", sqm.TFloat, "0.60000002"},
			},

			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}},
			},
		}
		groupvehiclesclass := &sqm.Class{
			Name: "Vehicles",

			Classes: []*sqm.Class{
				unitclass,
			},
		}
		groupclass := &sqm.Class{
			Name: "Item0",
			Props: []*sqm.Property{
				&sqm.Property{"side", sqm.TString, "WEST"},
			},
			Classes: []*sqm.Class{
				groupvehiclesclass,
			},
		}
		groupsclass := &sqm.Class{
			Name: "Groups",
			Classes: []*sqm.Class{
				groupclass,
			},
		}
		Convey("When parse group member", func() {
			unit := &Unit{}
			parseGroupMember(unitclass, unit)
			Convey("parsed unit should have all attributes", func() {
				So(unit.Classname, ShouldEqual, "classname")
				So(unit.Direction, ShouldEqual, "12.3")
				So(unit.Formation, ShouldEqual, "FORM")
				So(unit.IsLeader, ShouldBeTrue)
				So(unit.Skill, ShouldEqual, "0.60000002")
				So(unit.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
			})
			Convey("Pointer to class was set", func() {
				So(unit.class, ShouldPointTo, unitclass)
			})
		})
		Convey("When parse groups", func() {
			mission := &Mission{}
			parseGroups(groupsclass, mission)
			Convey("Mission has one group", func() {
				So(len(mission.Groups), ShouldEqual, 1)
			})
			Convey("Group should have one member", func() {
				So(len(mission.Groups[0].Units), ShouldEqual, 1)
			})
			Convey("Group should have right side", func() {
				So(mission.Groups[0].Side, ShouldEqual, "WEST")
			})
		})
	})
}

func TestParseMarkers(t *testing.T) {

	Convey("Given a set of valid marker classes", t, func() {
		markerClass := &sqm.Class{
			Name: "Item0",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"name", sqm.TString, "m1"},
				&sqm.Property{"markerType", sqm.TString, "ELLIPSE"},
				&sqm.Property{"type", sqm.TString, "Empty"},
				&sqm.Property{"colorName", sqm.TString, "ColorRed"},
				&sqm.Property{"fillName", sqm.TString, "Border"},
				&sqm.Property{"a", sqm.TInt, "1000"},
				&sqm.Property{"b", sqm.TInt, "2000"},
				&sqm.Property{"drawBorder", sqm.TInt, "1"},
			},
		}
		markersClass := &sqm.Class{
			Name: "Markers",
			Classes: []*sqm.Class{
				markerClass,
			},
		}
		Convey("When parse markers", func() {
			mission := &Mission{}
			parseMarkers(markersClass, mission)
			Convey("Mission has one marker", func() {
				So(len(mission.Markers), ShouldEqual, 1)
			})
			Convey("Marker has type", func() {
				So(mission.Markers[0].Type, ShouldEqual, "Empty")
			})
		})
		Convey("When parse single marker", func() {
			m := &Marker{}
			parseMarker(markerClass, m)
			Convey("All properties are correct", func() {
				So(m.Name, ShouldEqual, "m1")
				So(m.IsEllipse, ShouldBeTrue)
				So(m.Type, ShouldEqual, "Empty")
				So(m.ColorName, ShouldEqual, "ColorRed")
				So(m.FillName, ShouldEqual, "Border")
				So(m.Size, ShouldResemble, [2]string{"1000", "2000"})
				So(m.DrawBorder, ShouldBeTrue)
				So(m.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
			})
			Convey("Pointer to class was set", func() {
				So(m.class, ShouldPointTo, markerClass)
			})
		})
	})
}

func TestParseSensors(t *testing.T) {
	Convey("Given a valid sensor class", t, func() {
		effectsClass := &sqm.Class{
			Name: "Effects",
		}
		sensorClass := &sqm.Class{
			Name: "Item0",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TFloat, []string{"1.0", "2.0", "3.0"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"name", sqm.TString, "s1"},
				&sqm.Property{"a", sqm.TFloat, "1000"},
				&sqm.Property{"b", sqm.TFloat, "2000"},
				&sqm.Property{"angle", sqm.TFloat, "38.8545"},
				&sqm.Property{"rectangular", sqm.TInt, "1"},
				&sqm.Property{"repeating", sqm.TInt, "1"},
				&sqm.Property{"interruptable", sqm.TInt, "1"},
				&sqm.Property{"age", sqm.TString, "UNKNOWN"},
				&sqm.Property{"activationBy", sqm.TString, "ANY"},
				&sqm.Property{"expCond", sqm.TString, "isServer"},
				&sqm.Property{"expActiv", sqm.TString, "hint a1"},
			},
			Classes: []*sqm.Class{effectsClass},
		}
		sensorsClass := &sqm.Class{
			Name: "Sensors",
			Classes: []*sqm.Class{
				sensorClass,
			},
		}

		Convey("When parse sensors", func() {
			mission := &Mission{}
			parseSensors(sensorsClass, mission)
			Convey("Mission has one sensor", func() {
				So(len(mission.Sensors), ShouldEqual, 1)
			})
			Convey("Marker has ActivationBy", func() {
				So(mission.Sensors[0].ActivationBy, ShouldEqual, "ANY")
			})
		})
		Convey("When parse single sensor", func() {
			s := &Sensor{}
			parseSensor(sensorClass, s)
			Convey("All properties are correct", func() {
				So(s.Name, ShouldEqual, "s1")
				So(s.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
				So(s.Size, ShouldResemble, [2]string{"1000", "2000"})
				So(s.Angle, ShouldEqual, "38.8545")
				So(s.IsRectangle, ShouldBeTrue)
				So(s.IsRepeating, ShouldBeTrue)
				So(s.IsInterruptible, ShouldBeTrue)
				So(s.Age, ShouldEqual, "UNKNOWN")
				So(s.ActivationBy, ShouldEqual, "ANY")
				So(s.Condition, ShouldEqual, "isServer")
				So(s.OnActivation, ShouldEqual, "hint a1")
			})
			Convey("Pointer to class was set", func() {
				So(s.class, ShouldPointTo, sensorClass)
			})
			Convey("Pointer to effects class was set", func() {
				So(s.classEffects, ShouldPointTo, effectsClass)
			})
		})
	})
}
