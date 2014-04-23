package mission

import (
	sqm "github.com/blang/gosqm/sqmparser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseGroupMember(t *testing.T) {
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
