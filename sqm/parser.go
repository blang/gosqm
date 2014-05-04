package sqm

import (
	"fmt"
)

type Parser struct {
	input    string
	class    *Class //current class
	lexer    *lexer
	buff     *itemBuffer
	err      error
	propBuff *propBuffer
}

// A place for the currently processing property
type propBuffer struct {
	prop    *Property
	arrprop *ArrayProperty
}

type itemBuffer struct {
	ch      chan item
	prev    *item
	current *item
	ahead   *item
}

func makeItemBuffer(ch chan item) *itemBuffer {
	return &itemBuffer{ch: ch}
}

func (b *itemBuffer) next() *item {
	b.prev = b.current
	if b.ahead != nil {
		b.current = b.ahead
		b.ahead = nil
		return b.current
	}
	c := <-b.ch
	b.current = &c
	return b.current
}

func (b *itemBuffer) curr() *item {
	return b.current
}

func (b *itemBuffer) lookAhead() *item {
	if b.ahead == nil {
		c := <-b.ch
		b.ahead = &c
	}
	return b.ahead
}

func (b *itemBuffer) lookBack() *item {
	return b.prev
}

func MakeParser(input string) *Parser {
	l := makeLexer("sqm", input)
	class := &Class{Name: "mission"}
	parser := &Parser{
		input: input,
		class: class,
		lexer: l,
		buff:  makeItemBuffer(l.items),
	}
	return parser
}

// func (p *parser) next() item {
// 	return <-p.lexer.items
// }

type parserError struct {
	s string
}

func (ep *parserError) String() string {
	return ep.s
}

func (ep *parserError) Error() string {
	return ep.s
}

//TODO: Make better parser errors by using parserError fields
func (p *Parser) makeParserError(s string) *parserError {
	col, line := p.lexer.Position(p.buff.curr())
	err := fmt.Sprintf("Input:%d:%d: %s, read: %q", line, col, s, p.buff.curr())
	return &parserError{s: err}
}

func (p *Parser) ignoreSpace() {
	if space := p.buff.lookAhead(); space.typ == itemSpace {
		p.buff.next() //ignore space
	}
}

type pstateFn func(*Parser) (pstateFn, *parserError)

var pstartState pstateFn = parseInsideClass

// parseClassOpen parses a class definition till the open bracket
// and adds the class to the stack
func parseClassOpen(p *Parser) (pstateFn, *parserError) {
	var className string
	if p.buff.next().typ != itemClass {
		return nil, p.makeParserError("Missing class definition")
	}

	p.ignoreSpace()

	if classNameItem := p.buff.next(); classNameItem.typ != itemIdentifier {
		return nil, p.makeParserError("Missiong class indentifier")
	} else {
		className = classNameItem.val
	}

	p.ignoreSpace()

	if oblock := p.buff.next(); oblock.typ != itemOpenBlock {
		return nil, p.makeParserError("Missing { after class definition")
	}
	newClass := &Class{Name: className, parent: p.class}
	p.class = newClass
	return parseInsideClass, nil
}

func parseClassClose(p *Parser) (pstateFn, *parserError) {
	if p.buff.next().typ != itemCloseBlock {
		return nil, p.makeParserError("Unclosed class")
	}
	if p.class.parent == nil { //cant close base class
		return nil, p.makeParserError("Closing base class not allowed, unclosed class")
	}

	p.class.parent.Classes = append(p.class.parent.Classes, p.class)
	p.class = p.class.parent
	return parseInsideClass, nil
}

func parseProperty(p *Parser) (pstateFn, *parserError) {
	var name string
	ident := p.buff.next()
	if ident.typ != itemIdentifier {
		return nil, p.makeParserError("Expected identifier")
	}
	name = ident.val
	p.ignoreSpace()

	val := p.buff.next()
	switch val.typ {

	case itemEqual: //string or number
		prop := &Property{Name: name}
		p.propBuff = &propBuffer{prop: prop}
		return parsePropertyValue, nil

	case itemIdentifierArrayDec: //array
		p.ignoreSpace()
		if n := p.buff.next(); n.typ != itemEqual {
			return nil, p.makeParserError("Expected equal sign for array property")
		}
		prop := &ArrayProperty{Name: name}
		p.propBuff = &propBuffer{arrprop: prop}
		return parseArrayPropertyValue, nil

	default:
		return nil, p.makeParserError("Unexpected token in assignment")
	}

}

func parseArrayPropertyValue(p *Parser) (pstateFn, *parserError) {
	p.ignoreSpace()
	if n := p.buff.next(); n.typ != itemOpenArray {
		return nil, p.makeParserError("Expected open curly bracket for array property")
	}
	p.ignoreSpace()
	switch p.buff.lookAhead().typ {
	case itemStringDelim:
		p.propBuff.arrprop.Typ = TString
		err := parseArrayPropertyStringValues(p)
		if err != nil {
			return nil, err
		}
		if n := p.buff.next(); n.typ != itemCloseArray {
			return nil, p.makeParserError("Expected closing array after array string value")
		}
		p.class.Arrprops = append(p.class.Arrprops, p.propBuff.arrprop)
		p.propBuff.arrprop = nil
		return parseInsideClass, nil
	case itemInt:
		p.propBuff.arrprop.Typ = TInt
		err := parseArrayPropertyIntValues(p)
		if err != nil {
			return nil, err
		}
		if n := p.buff.next(); n.typ != itemCloseArray {
			return nil, p.makeParserError("Expected closing array after array int value")
		}
		p.class.Arrprops = append(p.class.Arrprops, p.propBuff.arrprop)
		p.propBuff.arrprop = nil
		return parseInsideClass, nil
	case itemFloat:
		p.propBuff.arrprop.Typ = TFloat
		err := parseArrayPropertyFloatValues(p)
		if err != nil {
			return nil, err
		}
		if n := p.buff.next(); n.typ != itemCloseArray {
			return nil, p.makeParserError("Expected closing array after array float value")
		}
		p.class.Arrprops = append(p.class.Arrprops, p.propBuff.arrprop)
		p.propBuff.arrprop = nil
		return parseInsideClass, nil
	default:
		return nil, p.makeParserError("Unexpected token in array assignment")
	}

	return nil, p.makeParserError("Arrays not implemented yet")
}

func parseArrayPropertyStringValues(p *Parser) *parserError {
	p.ignoreSpace()
	if t := p.buff.next(); t.typ != itemStringDelim {
		return p.makeParserError("Expected doublequote for array string value")
	}
	if t := p.buff.next(); t.typ != itemString {
		return p.makeParserError("Expected string for array string value")
	} else {
		p.propBuff.arrprop.Values = append(p.propBuff.arrprop.Values, t.val)
	}
	if t := p.buff.next(); t.typ != itemStringDelim {
		return p.makeParserError("Expected doublequote for array string value")
	}
	p.ignoreSpace()
	switch p.buff.lookAhead().typ {
	case itemArraySeperator:
		p.buff.next()
		return parseArrayPropertyStringValues(p)
	case itemCloseArray:
		return nil
	default:
		return p.makeParserError("Unexpected token in array string value")
	}
}

func parseArrayPropertyIntValues(p *Parser) *parserError {
	p.ignoreSpace()
	if t := p.buff.next(); t.typ != itemInt {
		return p.makeParserError("Expected int for array int value")
	} else {
		p.propBuff.arrprop.Values = append(p.propBuff.arrprop.Values, t.val)
	}
	p.ignoreSpace()
	switch p.buff.lookAhead().typ {
	case itemArraySeperator:
		p.buff.next()
		return parseArrayPropertyIntValues(p)
	case itemCloseArray:
		return nil
	default:
		return p.makeParserError("Unexpected token in array int value")
	}
}

func parseArrayPropertyFloatValues(p *Parser) *parserError {
	p.ignoreSpace()
	if t := p.buff.next(); t.typ != itemFloat {
		return p.makeParserError("Expected int for array float value")
	} else {
		p.propBuff.arrprop.Values = append(p.propBuff.arrprop.Values, t.val)
	}
	p.ignoreSpace()
	switch p.buff.lookAhead().typ {
	case itemArraySeperator:
		p.buff.next()
		return parseArrayPropertyFloatValues(p)
	case itemCloseArray:
		return nil
	default:
		return p.makeParserError("Unexpected token in array int value")
	}
}

func parsePropertyValue(p *Parser) (pstateFn, *parserError) {
	switch p.buff.lookAhead().typ {
	case itemStringDelim:
		p.propBuff.prop.Typ = TString
		p.buff.next()
		if v := p.buff.next(); v.typ != itemString {
			return nil, p.makeParserError("Expected string after string delimiter")
		} else {
			p.propBuff.prop.Value = v.val
		}
		if v := p.buff.next(); v.typ != itemStringDelim {
			return nil, p.makeParserError("Expected stringdelimiter after string")
		}
		p.ignoreSpace()
		if v := p.buff.next(); v.typ != itemSemicolon {
			return nil, p.makeParserError("Unclosed string assignment")
		}
		p.class.Props = append(p.class.Props, p.propBuff.prop)
		p.propBuff.prop = nil
		return parseInsideClass, nil
	case itemFloat:
		p.propBuff.prop.Typ = TFloat
		v := p.buff.next()
		p.propBuff.prop.Value = v.val
		p.ignoreSpace()
		if v := p.buff.next(); v.typ != itemSemicolon {
			return nil, p.makeParserError("Unclosed float assignment")
		}
		p.class.Props = append(p.class.Props, p.propBuff.prop)
		p.propBuff.prop = nil
		return parseInsideClass, nil
	case itemInt:
		p.propBuff.prop.Typ = TInt
		v := p.buff.next()
		p.propBuff.prop.Value = v.val
		p.ignoreSpace()
		if v := p.buff.next(); v.typ != itemSemicolon {
			return nil, p.makeParserError("Unclosed number assignment")
		}
		p.class.Props = append(p.class.Props, p.propBuff.prop)
		p.propBuff.prop = nil
		return parseInsideClass, nil
	default:
		return nil, p.makeParserError("Unexpected token in property value assigment")

	}
}

func parseInsideClass(p *Parser) (pstateFn, *parserError) {
	p.ignoreSpace()
	i := p.buff.lookAhead()
	switch i.typ {
	case itemError:
		return nil, p.makeParserError("Error while parsing inside class: " + i.val)
	case itemEOF:
		if p.class.parent != nil {
			return nil, p.makeParserError("Closing base class not allowed, unclosed class")
		}
		return nil, nil
	case itemClass:
		return parseClassOpen, nil
	case itemSpace:
		return parseInsideClass, nil
	case itemCloseBlock:
		return parseClassClose, nil
	case itemIdentifier:
		return parseProperty, nil
	default:
		return nil, p.makeParserError(fmt.Sprintf("Unrecognized item %q", i))
	}

	if i.typ == itemEOF || i.typ == itemError {
		return nil, p.makeParserError("Reached EOF or Error")
	}
	return parseInsideClass, nil
}

func (p *Parser) Run() (*Class, error) {
	l := p.lexer
	go l.run()
	var err *parserError

	for state := pstartState; state != nil; {
		state, err = state(p)
		if err != nil {
			return nil, err
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return p.class, nil
}
