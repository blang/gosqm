package sqm

import (
	"io/ioutil"
	"testing"
)

func TestStructure(t *testing.T) {
	class := &Class{Name: "mission"}
	if class.Props != nil {
		t.Errorf("Class has props")
	}
	if len(class.Props) > 0 {
		t.Errorf("Props length is %d", len(class.Props))
	}
	class.Props = append(class.Props, &Property{"test", TNumber, "value"})
	if class.Props == nil {
		t.Errorf("Class has no props")
	}

	if len(class.Props) != 1 {
		t.Errorf("Props length is %d", len(class.Props))
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
						{"version", TNumber, "11"},
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
				{"version", TNumber, "11"},
				{"string", TString, "teststring"},
				{"float1", TNumber, "123.456"},
				{"float2", TNumber, "+123.456"},
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
				{"arr", TNumber, []string{"1", "2", "3"}},
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
						{"arr", TNumber, []string{"1", "2", "3"}},
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
				{"arr", TNumber, []string{"1.2", "2.3", "3.4"}},
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
						{"arr", TNumber, []string{"1.2", "2.3", "3.4"}},
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
		p := MakeParser(test.input)
		c, err := p.Run()
		if err != nil {
			t.Errorf("Parser returned with error %q", err)
			continue
		}
		testClass(t, test.class, c)

	}
}

func TestParseSimple(t *testing.T) {
	p := MakeParser("class testclass { version=11; };")
	c, err := p.Run()
	if err != nil {
		t.Errorf("Parser returned with error %q", err)
	}
	if len(c.Classes) != 1 {
		t.Errorf("Class not processed")
	}
	tc := c.Classes[0]
	if tc.Name != "testclass" {
		t.Errorf("Wrong class name")
	}
	if len(tc.Props) != 1 {
		t.Errorf("Attribute not processed")
	}
	if tc.Arrprops != nil && len(tc.Arrprops) > 0 {
		t.Errorf("Array props are not empty")
	}
	at := tc.Props[0]
	if at.Name != "version" {
		t.Errorf("Prop wrong identifier")
	}
	if at.Typ != TNumber {
		t.Errorf("Type of prop wrong")
	}
	if at.Value != "11" {
		t.Errorf("Wrong prop value")
	}

}

func testClass(t *testing.T, tclass tclass, class *Class) {
	if tclass.name != class.Name {
		t.Errorf("Classname is %s but should be %s", class.Name, tclass.name)
		return
	}
	if (class.Arrprops == nil && len(tclass.arrprops) > 0) || (class.Arrprops != nil && len(tclass.props) != len(class.Props)) {
		t.Errorf("Class %s Prop length is %d but should be %d", tclass.name, len(class.Props), len(tclass.props))
		return
	}
	if (class.Arrprops == nil && len(tclass.arrprops) > 0) || (class.Arrprops != nil && len(tclass.arrprops) != len(class.Arrprops)) {
		t.Errorf("Class %s Prop length is %d but should be %d", tclass.name, len(class.Props), len(tclass.props))
		return
	}
	if (class.Classes == nil && len(tclass.classes) > 0) || (class.Classes != nil && len(tclass.classes) != len(class.Classes)) {
		t.Errorf("Class %s Classes length is %d but should be %d", tclass.name, len(class.Classes), len(tclass.classes))
		return
	}
	for i, tprop := range tclass.props {
		prop := class.Props[i]
		if tprop.Name != prop.Name {
			t.Errorf("Class %s Propname %s but should be %s", tclass.name, prop.Name, tprop.Name)
			return
		}
		if tprop.Typ != prop.Typ {
			t.Errorf("Class %s Proptype %s but should be %s", tclass.name, prop.Typ, tprop.Typ)
			return
		}
		if tprop.Value != prop.Value {
			t.Errorf("Class %s Propvalue %s but should be %s", tclass.name, prop.Value, tprop.Value)
			return
		}
	}
	for i, tprop := range tclass.arrprops {
		prop := class.Arrprops[i]
		if tprop.Name != prop.Name {
			t.Errorf("Class %s ArrPropname %s but should be %s", tclass.name, prop.Name, tprop.Name)
			return
		}
		if tprop.Typ != prop.Typ {
			t.Errorf("Class %s ArrProptype %s but should be %s", tclass.name, prop.Typ, tprop.Typ)
			return
		}
		for j, tval := range tprop.Values {
			val := prop.Values[j]
			if tval != val {
				t.Errorf("Class %s ArrPropvalue %s but should be %s", tclass.name, val, tval)
				return
			}
		}
	}
	for k, tsubc := range tclass.classes {
		subc := class.Classes[k]
		testClass(t, tsubc, subc)
	}

}

func TestParseMissionSQM(t *testing.T) {
	const name = "Mission.sqm parser"
	if testing.Short() {
		t.Skip("Skip mission.sqm in short mode")
		return
	}
	buf, err := ioutil.ReadFile("../testdata/mission.sqm")
	if err != nil {
		t.Errorf("Could not open mission.sqm")
		return
	}
	p := MakeParser(string(buf))
	// c, perr := p.Run()
	_, perr := p.Run()
	if perr != nil {
		t.Errorf("Parser returned with error %q", perr)
	}
	//t.Logf("Class parsed: %q", c)
}

func BenchmarkParseMissionSQM(b *testing.B) {
	buf, err := ioutil.ReadFile("../testdata/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		b.Errorf("Could not open mission.sqm")
		return
	}
	for n := 0; n < b.N; n++ {
		p := MakeParser(bufstr)
		// c, perr := p.Run()
		_, perr := p.Run()
		if perr != nil {
			b.Errorf("Parser returned with error %q", perr)
		}
	}
}

func BenchmarkParseMissionSQMReadfile(b *testing.B) {

	for n := 0; n < b.N; n++ {
		buf, err := ioutil.ReadFile("../testdata/mission.sqm")
		bufstr := string(buf)
		if err != nil {
			b.Errorf("Could not open mission.sqm")
			return
		}
		p := MakeParser(bufstr)
		// c, perr := p.Run()
		_, perr := p.Run()
		if perr != nil {
			b.Errorf("Parser returned with error %q", perr)
		}
	}
}
