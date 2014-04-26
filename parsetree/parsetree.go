package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/blang/gosqm/sqm"
	"io/ioutil"
)

var (
	missionFile  = flag.String("mission", "./mission.sqm", "mission.sqm")
	indentSpaces = flag.Int("indent", 2, "Indent in spaces")
)

func init() {
	flag.Parse()
}

func main() {
	buf, err := ioutil.ReadFile(*missionFile)
	if err != nil {
		fmt.Printf("Could not open file: %s", err.Error())
		return
	}
	p := sqm.MakeParser(string(buf))
	// c, perr := p.run()
	class, perr := p.Run()
	if perr != nil {
		fmt.Printf("Parser returned with error %q", perr)
	}

	printTreePart(class, 0)
}

func printTreePart(class *sqm.Class, level int) {
	classname, side := parseVehicle(class)
	id, found := parseId(class)
	var name string
	if found {
		name = class.Name + " (" + id + ")"
	} else {
		name = class.Name
	}
	if classname != "" {
		fmt.Printf("%s%s: %s %s %d\n", indent(level), name, classname, side, countAttributes(class))
	} else {
		fmt.Printf("%s%s: %d\n", indent(level), name, countAttributes(class))
	}
	for _, subclass := range class.Classes {
		printTreePart(subclass, level+1)
	}
}

func countAttributes(class *sqm.Class) int {
	i := 0
	i += len(class.Arrprops)
	i += len(class.Props)
	return i
}

func parseVehicle(class *sqm.Class) (classname string, side string) {
	for _, prop := range class.Props {
		switch prop.Name {
		case "vehicle":
			classname = prop.Value

		case "side":
			side = prop.Value
		}
	}
	return
}

func parseId(class *sqm.Class) (string, bool) {
	for _, prop := range class.Props {
		switch prop.Name {
		case "id":
			return prop.Value, true
		}
	}
	return "", false
}

func indent(level int) string {
	var buffer bytes.Buffer

	for i := 0; i < level*(*indentSpaces); i++ {
		buffer.WriteString(" ")
	}

	return buffer.String()
}
