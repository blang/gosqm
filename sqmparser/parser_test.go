package sqmparser

import (
	"testing"
	"io/ioutil"
)

func TestStructure(t *testing.T) {
	class := &Class{name: "mission"}
	if class.props != nil {
		t.Errorf("Class has props")
	}
	if len(class.props) > 0 {
		t.Errorf("Props length is %d", len(class.props))
	}
	class.props = append(class.props, &Property{"test", TInt, "value"})
	if class.props == nil {
		t.Errorf("Class has no props")
	}

	if len(class.props) != 1 {
		t.Errorf("Props length is %d", len(class.props))
	}
}

type parseTest struct {
	name  string
	input string
	class tclass
}

type tclass struct {
	name     string
	props    []Property
	arrprops []ArrayProperty
	classes  []tclass
}

var parseTests = []parseTest{
	{
		"simple class", "class testclass { version=11; };",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{},
			[]tclass{
				{"testclass",
					[]Property{
						{"version", TInt, "11"},
					},
					[]ArrayProperty{},
					[]tclass{},
				},
			},
		},
	},

	{
		"attributes", "version=11; string=\"teststring\"; float1=123.456; float2=+123.456;",
		tclass{"mission",
			[]Property{
				{"version", TInt, "11"},
				{"string", TString, "teststring"},
				{"float1", TFloat, "123.456"},
				{"float2", TFloat, "+123.456"},
			},
			[]ArrayProperty{},
			[]tclass{},
		},
	},

	{
		"array int attributes", "arr[]={1,2,3};",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{
				{"arr", TInt, []string{"1", "2", "3"}},
			},
			[]tclass{},
		},
	},

	{
		"subclass array int attributes", "class test { arr[]={1,2,3}; };",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{},
			[]tclass{
				{"test",
					[]Property{},
					[]ArrayProperty{
						{"arr", TInt, []string{"1", "2", "3"}},
					},
					[]tclass{},
				},
			},
		},
	},

	{
		"array float attributes", "arr[]={1.2,2.3,3.4};",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{
				{"arr", TFloat, []string{"1.2", "2.3", "3.4"}},
			},
			[]tclass{},
		},
	},

	{
		"subclass array float attributes", "class test { arr[]={1.2,2.3,3.4}; };",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{},
			[]tclass{
				{"test",
					[]Property{},
					[]ArrayProperty{
						{"arr", TFloat, []string{"1.2", "2.3", "3.4"}},
					},
					[]tclass{},
				},
			},
		},
	},

	{
		"array string attributes", "arr[]={\"a\",\"b\",\"c\"};",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{
				{"arr", TString, []string{"a", "b", "c"}},
			},
			[]tclass{},
		},
	},

	{
		"subclass array string attributes", "class test { arr[]={\"a\",\"b\",\"c\"}; };",
		tclass{"mission",
			[]Property{},
			[]ArrayProperty{},
			[]tclass{
				{"test",
					[]Property{},
					[]ArrayProperty{
						{"arr", TString, []string{"a", "b", "c"}},
					},
					[]tclass{},
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		p := makeParser(test.input)
		c, err := p.run()
		if err != nil {
			t.Errorf("Parser returned with error %q", err)
			continue
		}
		testClass(t, test.class, c)

	}
}

func TestParseSimple(t *testing.T) {
	p := makeParser("class testclass { version=11; };")
	c, err := p.run()
	if err != nil {
		t.Errorf("Parser returned with error %q", err)
	}
	if len(c.classes) != 1 {
		t.Errorf("Class not processed")
	}
	tc := c.classes[0]
	if tc.name != "testclass" {
		t.Errorf("Wrong class name")
	}
	if len(tc.props) != 1 {
		t.Errorf("Attribute not processed")
	}
	if tc.arrprops != nil && len(tc.arrprops) > 0 {
		t.Errorf("Array props are not empty")
	}
	at := tc.props[0]
	if at.name != "version" {
		t.Errorf("Prop wrong identifier")
	}
	if at.typ != TInt {
		t.Errorf("Type of prop wrong")
	}
	if at.value != "11" {
		t.Errorf("Wrong prop value")
	}

}

func testClass(t *testing.T, tclass tclass, class *Class) {
	if tclass.name != class.name {
		t.Errorf("Classname is %s but should be %s", class.name, tclass.name)
		return
	}
	if (class.arrprops == nil && len(tclass.arrprops) > 0) || (class.arrprops != nil && len(tclass.props) != len(class.props)) {
		t.Errorf("Class %s Prop length is %d but should be %d", tclass.name, len(class.props), len(tclass.props))
		return
	}
	if (class.arrprops == nil && len(tclass.arrprops) > 0) || (class.arrprops != nil && len(tclass.arrprops) != len(class.arrprops)) {
		t.Errorf("Class %s Prop length is %d but should be %d", tclass.name, len(class.props), len(tclass.props))
		return
	}
	if (class.classes == nil && len(tclass.classes) > 0) || (class.classes != nil && len(tclass.classes) != len(class.classes)) {
		t.Errorf("Class %s Classes length is %d but should be %d", tclass.name, len(class.classes), len(tclass.classes))
		return
	}
	for i, tprop := range tclass.props {
		prop := class.props[i]
		if tprop.name != prop.name {
			t.Errorf("Class %s Propname %s but should be %s", tclass.name, prop.name, tprop.name)
			return
		}
		if tprop.typ != prop.typ {
			t.Errorf("Class %s Proptype %s but should be %s", tclass.name, prop.typ, tprop.typ)
			return
		}
		if tprop.value != prop.value {
			t.Errorf("Class %s Propvalue %s but should be %s", tclass.name, prop.value, tprop.value)
			return
		}
	}
	for i, tprop := range tclass.arrprops {
		prop := class.arrprops[i]
		if tprop.name != prop.name {
			t.Errorf("Class %s ArrPropname %s but should be %s", tclass.name, prop.name, tprop.name)
			return
		}
		if tprop.typ != prop.typ {
			t.Errorf("Class %s ArrProptype %s but should be %s", tclass.name, prop.typ, tprop.typ)
			return
		}
		for j, tval := range tprop.values {
			val := prop.values[j]
			if tval != val {
				t.Errorf("Class %s ArrPropvalue %s but should be %s", tclass.name, val, tval)
				return
			}
		}
	}
	for k, tsubc := range tclass.classes {
		subc := class.classes[k]
		testClass(t, tsubc, subc)
	}

}


func TestParseMissionSQM(t *testing.T) {
	const name = "Mission.sqm parser"
	if testing.Short() {
		t.Skip("Skip mission.sqm in short mode")
		return
	}
	buf, err := ioutil.ReadFile("./mission.sqm")
	if err != nil {
		t.Errorf("Could not open mission.sqm")
		return
	}
	p := makeParser(string(buf))
	// c, perr := p.run()
  _, perr := p.run()
	if perr != nil {
		t.Errorf("Parser returned with error %q", perr)
	}
	//t.Logf("Class parsed: %q", c)
}