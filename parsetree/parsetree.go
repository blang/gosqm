package main

import (
	"bytes"
	"flag"
	"fmt"
	sqm "github.com/blang/gosqm/sqmparser"
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
	if classname != "" {
		fmt.Printf("%s%s: %s %s %d\n", indent(level), class.Name, classname, side, countAttributes(class))
	} else {
		fmt.Printf("%s%s: %d\n", indent(level), class.Name, countAttributes(class))
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

func indent(level int) string {
	var buffer bytes.Buffer

	for i := 0; i < level*(*indentSpaces); i++ {
		buffer.WriteString(" ")
	}

	return buffer.String()
}
