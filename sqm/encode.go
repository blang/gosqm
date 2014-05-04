package sqm

import (
	"bytes"
	"io"
)

const LINEBREAK = "\r\n"
const INDENT = 2

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{}
	e.w = w
	return e
}

func (e *Encoder) writeString(s string) error {
	_, err := e.w.Write([]byte(s))
	return err
}

func (e *Encoder) Encode(class *Class) error {
	return e.encodeMainClass(class, 0)
}

func (e *Encoder) encodeClass(class *Class, level int) error {
	err := e.writeString(indent(level) + "class " + class.Name + LINEBREAK + indent(level) + "{" + LINEBREAK)
	if err != nil {
		return err
	}
	err = e.encodeSubElements(class, level+1)
	if err != nil {
		return err
	}

	err = e.writeString(indent(level) + "};" + LINEBREAK)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeMainClass(class *Class, level int) error {
	err := e.encodeSubElements(class, level)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeSubElements(class *Class, level int) error {
	//encode arr properties
	for _, arrProp := range class.Arrprops {
		err := e.encodeArrProperty(arrProp, level)
		if err != nil {
			return err
		}
	}

	//encode properties
	for _, prop := range class.Props {
		err := e.encodeProperty(prop, level)
		if err != nil {
			return err
		}
	}

	//encode subclasses
	for _, subclass := range class.Classes {
		err := e.encodeClass(subclass, level)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeProperty(p *Property, level int) error {
	var err error
	switch p.Typ {
	case TString:
		err = e.writeString(indent(level) + p.Name + "=\"" + p.Value + "\";" + LINEBREAK)
	case TInt:
		fallthrough
	case TFloat:
		err = e.writeString(indent(level) + p.Name + "=" + p.Value + ";" + LINEBREAK)
	}
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeArrProperty(arrProp *ArrayProperty, level int) error {
	if arrProp.Name == "addOns" || arrProp.Name == "addOnsAuto" {
		return e.encodeAddonsArrProperty(arrProp, level)
	} else {
		return e.encodeNormalArrProperty(arrProp, level)
	}
}

func (e *Encoder) encodeAddonsArrProperty(arrProp *ArrayProperty, level int) error {
	err := e.writeString(indent(level) + arrProp.Name + "[]=" + LINEBREAK + indent(level) + "{" + LINEBREAK)
	if err != nil {
		return err
	}
	for i, val := range arrProp.Values {
		if i > 0 {
			err = e.writeString("," + LINEBREAK)
			if err != nil {
				return err
			}

		}

		err = e.writeString(indent(level+1) + "\"" + val + "\"")
		if err != nil {
			return err
		}

	}
	err = e.writeString(LINEBREAK + indent(level) + "};" + LINEBREAK)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeNormalArrProperty(arrProp *ArrayProperty, level int) error {
	err := e.writeString(indent(level) + arrProp.Name + "[]={")
	if err != nil {
		return err
	}
	for i, val := range arrProp.Values {
		if i > 0 {
			err = e.writeString(",")
			if err != nil {
				return err
			}

		}
		if arrProp.Typ == TString {
			err = e.writeString("\"" + val + "\"")
			if err != nil {
				return err
			}
		} else {
			err := e.writeString(val)
			if err != nil {
				return err
			}
		}
	}
	err = e.writeString("};" + LINEBREAK)
	if err != nil {
		return err
	}
	return nil
}

const indentCacheMax = 50

var indentCache [indentCacheMax]*string

func indent(level int) string {
	switch {
	case level == 0:
		return ""
	case level <= 50:
		if indentCache[level-1] == nil {
			var buffer bytes.Buffer

			for i := 0; i < level*INDENT; i++ {
				buffer.WriteString(" ")
			}
			str := buffer.String()
			indentCache[level-1] = &str
			return str
		} else {
			return *(indentCache[level-1])
		}

	case level > 50:
		var buffer bytes.Buffer

		for i := 0; i < level*INDENT; i++ {
			buffer.WriteString(" ")
		}
		return buffer.String()
	}
	return "" //Invalid level input
}
