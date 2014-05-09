package mission

import (
	"github.com/blang/gosqm/sqm"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseMission(t *testing.T) {
	Convey("Given a fresh mission class", t, func() {
		p := NewParser()
		missionclass := &sqm.Class{
			Name: "Mission",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"addOns", sqm.TString, []string{"addon1", "addon2", "addon3"}},
				&sqm.ArrayProperty{"addOnsAuto", sqm.TString, []string{"addon4", "addon5", "addon6"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"randomSeed", sqm.TNumber, "13617784"},
			},
		}
		Convey("When parse addons", func() {
			m := &Mission{}
			p.parseMissionProps(missionclass, m)
			Convey("All properties are correct", func() {
				So(m.Addons, ShouldResemble, []string{"addon1", "addon2", "addon3"})
				So(m.AddonsAuto, ShouldResemble, []string{"addon4", "addon5", "addon6"})
				So(m.RandomSeed, ShouldEqual, "13617784")
			})
		})
	})
}

func TestParseIntel(t *testing.T) {
	Convey("Given a valid intel class", t, func() {
		p := NewParser()
		intelclass := &sqm.Class{
			Name: "Intel",
			Props: []*sqm.Property{
				&sqm.Property{"resistanceWest", sqm.TNumber, "1"},
				&sqm.Property{"startWeather", sqm.TNumber, "0.3"},
				&sqm.Property{"forecastWeather", sqm.TNumber, "0.8"},
				&sqm.Property{"year", sqm.TNumber, "2009"},
				&sqm.Property{"month", sqm.TNumber, "10"},
				&sqm.Property{"day", sqm.TNumber, "28"},
				&sqm.Property{"hour", sqm.TNumber, "6"},
				&sqm.Property{"minute", sqm.TNumber, "5"},
			},
		}
		Convey("When parse intel", func() {
			mission := &Mission{}
			p.parseIntel(intelclass, mission)
			i := mission.Intel
			Convey("All properties are correct", func() {
				So(i.ResistanceWest, ShouldBeTrue)
				So(i.StartWeather, ShouldEqual, "0.3")
				So(i.ForecastWeather, ShouldEqual, "0.8")
				So(i.Year, ShouldEqual, "2009")
				So(i.Month, ShouldEqual, "10")
				So(i.Day, ShouldEqual, "28")
				So(i.Hour, ShouldEqual, "6")
				So(i.Minute, ShouldEqual, "5")
			})
			Convey("Pointer to class was set", func() {
				So(i.class, ShouldPointTo, intelclass)
			})
		})
	})
}

func TestParseGroups(t *testing.T) {
	Convey("Given a valid groups class with subclasses", t, func() {
		p := NewParser()
		unitclass := &sqm.Class{
			Name: "Item0",
			Props: []*sqm.Property{
				&sqm.Property{"text", sqm.TNumber, "name"},
				&sqm.Property{"azimut", sqm.TNumber, "12.3"},
				&sqm.Property{"vehicle", sqm.TString, "classname"},
				&sqm.Property{"leader", sqm.TNumber, "1"},
				&sqm.Property{"special", sqm.TString, "FORM"},
				&sqm.Property{"skill", sqm.TNumber, "0.60000002"},
				&sqm.Property{"player", sqm.TString, "PLAYER COMMANDER"},
				&sqm.Property{"description", sqm.TString, "Description"},
				&sqm.Property{"presence", sqm.TNumber, "0.3"},
				&sqm.Property{"presenceCondition", sqm.TString, "true"},
				&sqm.Property{"placement", sqm.TNumber, "20"},
				&sqm.Property{"age", sqm.TString, "5 MIN"},
				&sqm.Property{"lock", sqm.TString, "UNLOCKED"},
				&sqm.Property{"rank", sqm.TString, "CORPORAL"},
				&sqm.Property{"health", sqm.TNumber, "0.1"},
				&sqm.Property{"fuel", sqm.TNumber, "0.2"},
				&sqm.Property{"ammo", sqm.TNumber, "0.3"},
				&sqm.Property{"init", sqm.TString, "hint a"},
				&sqm.Property{"side", sqm.TString, "WEST"},
			},

			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}},
			},
		}
		effectsClass := &sqm.Class{
			Name: "Effects",
			Props: []*sqm.Property{
				&sqm.Property{"sound", sqm.TString, "sound"},
				&sqm.Property{"voice", sqm.TString, "voice"},
				&sqm.Property{"soundEnv", sqm.TString, "soundenv"},
				&sqm.Property{"soundDet", sqm.TString, "sounddet"},
				&sqm.Property{"track", sqm.TString, "track"},
				&sqm.Property{"titleType", sqm.TString, "titletype"},
				&sqm.Property{"title", sqm.TString, "title"},
				&sqm.Property{"titleEffect", sqm.TString, "titleeffect"},
			},
		}
		waypointclass := &sqm.Class{
			Name: "Item0",
			Props: []*sqm.Property{
				&sqm.Property{"type", sqm.TString, "AND"},
				&sqm.Property{"showWP", sqm.TString, "NEVER"},
			},
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}},
			},
			Classes: []*sqm.Class{effectsClass},
		}

		groupwaypointsclass := &sqm.Class{
			Name: "Waypoints",

			Classes: []*sqm.Class{
				waypointclass,
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
				groupwaypointsclass,
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
			p.parseGroupMember(unitclass, unit)
			Convey("parsed unit should have all attributes", func() {
				So(unit.Name, ShouldEqual, "name")
				So(unit.Classname, ShouldEqual, "classname")
				So(unit.Direction, ShouldEqual, "12.3")
				So(unit.Special, ShouldEqual, "FORM")
				So(unit.IsLeader, ShouldBeTrue)
				So(unit.Skill, ShouldEqual, "0.60000002")
				So(unit.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
				So(unit.Player, ShouldEqual, "PLAYER COMMANDER")
				So(unit.Description, ShouldEqual, "Description")
				So(unit.Presence, ShouldEqual, "0.3")
				So(unit.PresenceCond, ShouldEqual, "true")
				So(unit.Placement, ShouldEqual, "20")
				So(unit.Age, ShouldEqual, "5 MIN")
				So(unit.Lock, ShouldEqual, "UNLOCKED")
				So(unit.Rank, ShouldEqual, "CORPORAL")
				So(unit.Health, ShouldEqual, "0.1")
				So(unit.Fuel, ShouldEqual, "0.2")
				So(unit.Ammo, ShouldEqual, "0.3")
				So(unit.Init, ShouldEqual, "hint a")
				So(unit.Side, ShouldEqual, "WEST")
			})
			Convey("Pointer to class was set", func() {
				So(unit.class, ShouldPointTo, unitclass)
			})
		})
		Convey("When parse groups", func() {
			mission := &Mission{}
			p.parseGroups(groupsclass, mission)
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
		Convey("When parse group", func() {
			group := &Group{}
			p.parseGroup(groupclass, group)

			Convey("Group has one waypoint", func() {
				So(len(group.Waypoints), ShouldEqual, 1)
			})
		})
		Convey("When parse waypoint", func() {
			wp := &Waypoint{}
			p.parseGroupWaypoint(waypointclass, wp)
			Convey("All properties are correct", func() {
				So(wp.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
				So(wp.Type, ShouldEqual, "AND")
				So(wp.ShowWP, ShouldEqual, "NEVER")
			})
			Convey("Pointer to class was set", func() {
				So(wp.class, ShouldPointTo, waypointclass)
			})
			Convey("Effects was set", func() {
				So(wp.Effects, ShouldNotBeNil)
				eff := wp.Effects
				Convey("All effect properties should be set", func() {
					So(eff.Sound, ShouldEqual, "sound")
					So(eff.Voice, ShouldEqual, "voice")
					So(eff.SoundDet, ShouldEqual, "sounddet")
					So(eff.SoundEnv, ShouldEqual, "soundenv")
					So(eff.Title, ShouldEqual, "title")
					So(eff.TitleEffect, ShouldEqual, "titleeffect")
					So(eff.TitleType, ShouldEqual, "titletype")
					So(eff.Track, ShouldEqual, "track")
				})
			})
		})
	})
}

func TestParseMarkers(t *testing.T) {

	Convey("Given a set of valid marker classes", t, func() {
		p := NewParser()
		markerClass := &sqm.Class{
			Name: "Item0",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"name", sqm.TString, "m1"},
				&sqm.Property{"markerType", sqm.TString, "ELLIPSE"},
				&sqm.Property{"type", sqm.TString, "Empty"},
				&sqm.Property{"colorName", sqm.TString, "ColorRed"},
				&sqm.Property{"fillName", sqm.TString, "Border"},
				&sqm.Property{"a", sqm.TNumber, "1000"},
				&sqm.Property{"b", sqm.TNumber, "2000"},
				&sqm.Property{"drawBorder", sqm.TNumber, "1"},
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
			p.parseMarkers(markersClass, mission)
			Convey("Mission has one marker", func() {
				So(len(mission.Markers), ShouldEqual, 1)
			})
			Convey("Marker has type", func() {
				So(mission.Markers[0].Type, ShouldEqual, "Empty")
			})
		})
		Convey("When parse single marker", func() {
			m := &Marker{}
			p.parseMarker(markerClass, m)
			Convey("All properties are correct", func() {
				So(m.Name, ShouldEqual, "m1")
				So(m.MarkerType, ShouldEqual, "ELLIPSE")
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
		p := NewParser()
		effectsClass := &sqm.Class{
			Name: "Effects",
			Props: []*sqm.Property{
				&sqm.Property{"sound", sqm.TString, "sound"},
				&sqm.Property{"voice", sqm.TString, "voice"},
				&sqm.Property{"soundEnv", sqm.TString, "soundenv"},
				&sqm.Property{"soundDet", sqm.TString, "sounddet"},
				&sqm.Property{"track", sqm.TString, "track"},
				&sqm.Property{"titleType", sqm.TString, "titletype"},
				&sqm.Property{"title", sqm.TString, "title"},
				&sqm.Property{"titleEffect", sqm.TString, "titleeffect"},
			},
		}
		sensorClass := &sqm.Class{
			Name: "Item0",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"name", sqm.TString, "s1"},
				&sqm.Property{"a", sqm.TNumber, "1000"},
				&sqm.Property{"b", sqm.TNumber, "2000"},
				&sqm.Property{"angle", sqm.TNumber, "38.8545"},
				&sqm.Property{"rectangular", sqm.TNumber, "1"},
				&sqm.Property{"repeating", sqm.TNumber, "1"},
				&sqm.Property{"interruptable", sqm.TNumber, "1"},
				&sqm.Property{"age", sqm.TString, "UNKNOWN"},
				&sqm.Property{"activationBy", sqm.TString, "ANY"},
				&sqm.Property{"activationType", sqm.TString, "GUER D"},
				&sqm.Property{"timeoutMin", sqm.TNumber, "1"},
				&sqm.Property{"timeoutMid", sqm.TNumber, "2"},
				&sqm.Property{"timeoutMax", sqm.TNumber, "3"},
				&sqm.Property{"type", sqm.TString, "EAST G"},
				&sqm.Property{"expCond", sqm.TString, "isServer"},
				&sqm.Property{"expActiv", sqm.TString, "hint a1"},
				&sqm.Property{"expDesactiv", sqm.TString, "hint a2"},
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
			p.parseSensors(sensorsClass, mission)
			Convey("Mission has one sensor", func() {
				So(len(mission.Sensors), ShouldEqual, 1)
			})
			Convey("Marker has ActivationBy", func() {
				So(mission.Sensors[0].ActivationBy, ShouldEqual, "ANY")
			})
		})
		Convey("When parse single sensor", func() {
			s := &Sensor{}
			p.parseSensor(sensorClass, s)
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
				So(s.ActivationType, ShouldEqual, "GUER D")
				So(s.TimeoutMin, ShouldEqual, "1")
				So(s.TimeoutMid, ShouldEqual, "2")
				So(s.TimeoutMax, ShouldEqual, "3")
				So(s.Type, ShouldEqual, "EAST G")
				So(s.Condition, ShouldEqual, "isServer")
				So(s.OnActivation, ShouldEqual, "hint a1")
				So(s.OnDeactivation, ShouldEqual, "hint a2")
			})
			Convey("Pointer to class was set", func() {
				So(s.class, ShouldPointTo, sensorClass)
			})
			Convey("Effects should be set", func() {
				So(s.Effects, ShouldNotBeNil)
				eff := s.Effects
				Convey("All effect properties should be set", func() {
					So(eff.Sound, ShouldEqual, "sound")
					So(eff.Voice, ShouldEqual, "voice")
					So(eff.SoundDet, ShouldEqual, "sounddet")
					So(eff.SoundEnv, ShouldEqual, "soundenv")
					So(eff.Title, ShouldEqual, "title")
					So(eff.TitleEffect, ShouldEqual, "titleeffect")
					So(eff.TitleType, ShouldEqual, "titletype")
					So(eff.Track, ShouldEqual, "track")
				})
			})
		})
	})
}

func TestParseVehicles(t *testing.T) {
	Convey("Given a valid vehicle class", t, func() {
		p := NewParser()
		vehClass := &sqm.Class{
			Name: "Item0",
			Arrprops: []*sqm.ArrayProperty{
				&sqm.ArrayProperty{"position", sqm.TNumber, []string{"1.0", "2.0", "3.0"}},
			},
			Props: []*sqm.Property{
				&sqm.Property{"name", sqm.TString, "s1"},
				&sqm.Property{"azimut", sqm.TNumber, "30.2"},
				&sqm.Property{"side", sqm.TString, "EMPTY"},
				&sqm.Property{"vehicle", sqm.TString, "HeliH"},
				&sqm.Property{"skill", sqm.TNumber, "0.6"},
				&sqm.Property{"presence", sqm.TNumber, "0.3"},
				&sqm.Property{"presenceCondition", sqm.TString, "true"},
				&sqm.Property{"special", sqm.TString, "NONE"},
			},
		}
		vehsClass := &sqm.Class{
			Name: "Vehicles",
			Classes: []*sqm.Class{
				vehClass,
			},
		}

		Convey("When parse sensors", func() {
			mission := &Mission{}
			p.parseVehicles(vehsClass, mission)
			Convey("Mission has one vehicle", func() {
				So(len(mission.Vehicles), ShouldEqual, 1)
			})
			Convey("Vehicle has classname", func() {
				So(mission.Vehicles[0].Classname, ShouldEqual, "HeliH")
			})
		})
		Convey("When parse single vehicle", func() {
			v := &Vehicle{}
			p.parseVehicle(vehClass, v)
			Convey("All properties are correct", func() {
				So(v.Name, ShouldEqual, "s1")
				So(v.Position, ShouldResemble, [3]string{"1.0", "2.0", "3.0"})
				So(v.Angle, ShouldEqual, "30.2")
				So(v.Side, ShouldEqual, "EMPTY")
				So(v.Classname, ShouldEqual, "HeliH")
				So(v.Skill, ShouldEqual, "0.6")
				So(v.Presence, ShouldEqual, "0.3")
				So(v.PresenceCond, ShouldEqual, "true")
				So(v.Special, ShouldEqual, "NONE")
			})
			Convey("Pointer to class was set", func() {
				So(v.class, ShouldPointTo, vehClass)
			})
		})
	})
}
