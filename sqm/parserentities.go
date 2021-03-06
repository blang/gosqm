package sqm

import (
	"fmt"
)

type PropType int

const (
	TString PropType = iota // Arma escaped string
	TNumber                 // Integer or float
)

type Property struct {
	Name  string
	Typ   PropType
	Value string
}

type ArrayProperty struct {
	Name   string
	Typ    PropType
	Values []string
}

type Class struct {
	Name     string
	Props    []*Property
	Arrprops []*ArrayProperty
	Classes  []*Class
	parent   *Class
}

func (p Property) String() string {
	return fmt.Sprintf("%s='%s' (Type: %d)\n", p.Name, p.Value, p.Typ)
}

func (c Class) String() string {
	return fmt.Sprintf("class (name: %s), props: %s, arrprops: %s, classes: %s\n", c.Name, c.Props, c.Arrprops, c.Classes)
}

func (t PropType) String() string {
	switch t {
	case TString:
		return "TString"
	case TNumber:
		return "TNumber"
	}
	return "Unkown"
}
