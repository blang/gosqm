package sqm

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Pos represents a byte position in the input text.
type Pos uint

func (p Pos) Position() Pos {
	return p
}

// Positions gets item position in lexer string as column and line.
func (l *lexer) Position(item *item) (col int, line int) {
	text := l.input[:item.pos]
	byteNum := strings.LastIndex(text, "\n")
	startLastLine := 0
	if byteNum == -1 {
		startLastLine = 0
		//col = pos // On first line.
	} else {
		byteNum++ // After the newline.
		startLastLine = byteNum
		//col = pos - byteNum
	}
	col = utf8.RuneCountInString(text[startLastLine:])
	line = 1 + strings.Count(text, "\n")
	return
}

const eof rune = -1

type itemType int

const (
	itemError              itemType = iota // Error occurred. Message in val
	itemEOF                                // EOF if all input is read
	itemInt                                // Number without float point seperator
	itemFloat                              // Float with point seperator
	itemIdentifier                         // Identifier includes classnames and property names
	itemIdentifierArrayDec                 // [] Array declaration
	itemEqual                              // =
	itemSemicolon                          // ; seperates logical blocks
	itemSpace                              // also newlines, tabs merged into one itemSpace
	itemOpenBlock                          // {
	itemCloseBlock                         // };
	itemOpenArray                          // {
	itemCloseArray                         // }
	itemClass                              // class
	itemArraySeperator                     // ,
	itemStringDelim                        // "
	itemString                             // String
)

type item struct {
	typ itemType // Type of this item.
	pos Pos      // Starting position in bytes in the input string.
	val string   // Value of this item.
}

func (i item) String() string {
	switch i.typ {
	case itemError:
		return fmt.Sprintf("Error: %q", i.val)
	}
	return fmt.Sprintf("type: %d val: %q", i.typ, i.val)
}

type stateFn func(*lexer) stateFn
type lexer struct {
	name  string    //used for errors
	input string    //string being scanned
	start Pos       // start position of the item
	pos   Pos       // current position in the input
	width Pos       // width of last rune read
	items chan item // channel of scanned items
}

// Starting state of state machine
var startState stateFn = lexInsideClass

// Initializes a lexer
func makeLexer(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, 10),
	}
	return l
}

// Run the lexer as a state machine until state is nil
func (l *lexer) run() {
	for state := startState; state != nil; {
		state = state(l)
	}
	close(l.items) //TODO: close 'emits' null type, maybe unwanted
}

// Emits an item
func (l *lexer) emit(t itemType) {
	l.items <- item{
		typ: t,
		pos: l.start,
		val: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

// Next rune
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	//TODO: Possible error, use l.width = Pos(w)
	// as seen in http://golang.org/src/pkg/text/template/parse/lex.go
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func (l *lexer) last() (r rune) {
	l.backup()
	return l.next()
}

// Ignores the last rune and moves start to current position
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

//hasContent checks if lexed at least one character
func (l *lexer) hasContent() bool {
	return l.pos > l.start
}

// accept consumes the next rune if it's in the valid set
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) int {
	i := 0
	for l.accept(valid) {
		i++
	}
	return i
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

const digits = "0123456789"
const numberWithSign = "+-" + digits

func isNumber(r rune) bool {
	return (strings.IndexRune(numberWithSign, r) >= 0)
}

const alphaLower = "abcdefghijklmnopqrstuvwxyz"
const alphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func isAlpha(r rune) bool {
	return (strings.IndexRune(alphaLower+alphaUpper, r) >= 0)
}

const space = " \n\r\t"

func isSpace(r rune) bool {
	return (strings.IndexRune(space, r) >= 0)
}

const openBracket = "{"
const closeBracket = "}"

func isOpenBracket(r rune) bool {
	return (strings.IndexRune(openBracket, r) >= 0)
}
func isCloseBracket(r rune) bool {
	return (strings.IndexRune(closeBracket, r) >= 0)
}

func lexOptionalSpace(l *lexer) bool {
	if i := l.acceptRun(space); i > 0 {
		l.emit(itemSpace)
		return true
	}
	return false
}

func lexIdentifier(l *lexer) stateFn {
	lexOptionalSpace(l)
	if !l.accept(alphaLower + alphaUpper) {
		return l.errorf("Identifier does not start with an alpha character")
	}
	l.acceptRun(alphaLower + alphaUpper + digits)
	if l.input[l.start:l.pos] == "class" {
		l.emit(itemClass)
		return lexSpaceBeforeClassIdentifier
	} else {
		l.emit(itemIdentifier)
		if l.peek() == eof {
			return l.errorf("unclosed assignment")
		}
		if l.peek() == '[' {
			return lexAttributeArrayDeclaration
		}
		return lexAssignmentEqual
	}
}

func lexNumber(l *lexer) stateFn {
	isFloat := false
	l.accept("+-")
	l.acceptRun(digits)
	if l.accept(".") {
		isFloat = true
		l.acceptRun(digits)
	}
	if isFloat {
		l.emit(itemFloat)
	} else {
		l.emit(itemInt)
	}
	lexOptionalSpace(l)
	if r := l.next(); r != ';' {
		return l.errorf("Unclosed number assignment")
	}
	l.emit(itemSemicolon)

	return lexInsideClass
}

func lexAttributeArrayDeclaration(l *lexer) stateFn {
	if r := l.next(); r != '[' {
		return l.errorf("Missing open array bracket")
	}
	if r := l.next(); r != ']' {
		return l.errorf("Missing closing array bracket")
	}
	l.emit(itemIdentifierArrayDec)
	return lexArrayAssignmentEqual
}

func lexArrayAssignmentEqual(l *lexer) stateFn {
	lexOptionalSpace(l)
	if r := l.next(); r != '=' {
		return l.errorf("Array assignment missing equals, was %q", r)
	}
	l.emit(itemEqual)
	return lexArrayOpenBracket
}

func lexArrayOpenBracket(l *lexer) stateFn {
	lexOptionalSpace(l)
	if r := l.next(); r != '{' {
		return l.errorf("Missing array open curly bracket")
	}
	l.emit(itemOpenArray)
	return lexInsideArray
}

func lexArrayClose(l *lexer) stateFn {
	lexOptionalSpace(l)
	if r := l.next(); r != '}' {
		return l.errorf("Missing array closing curly bracket")
	}
	if r := l.next(); r != ';' {
		return l.errorf("Missing array closing semicolon")
	}
	l.emit(itemCloseArray)
	return lexInsideClass
}

// lexArrayString lexes a string inside an array
// TODO: Support escape, newline
func lexArrayString(l *lexer) stateFn {
	if !doString(l) {
		return nil
	}
	return lexInsideArray
}

func lexArrayNumber(l *lexer) stateFn {
	lexOptionalSpace(l)
	isFloat := false
	l.accept("+-")
	l.acceptRun(digits)
	if l.accept(".") {
		isFloat = true
		l.acceptRun(digits)
	}
	if !l.hasContent() {
		return l.errorf("Missing number")
	}
	if isFloat {
		l.emit(itemFloat)
	} else {
		l.emit(itemInt)
	}
	return lexInsideArray
}

func lexInsideArray(l *lexer) stateFn {
	lexOptionalSpace(l)
	//Allowed seperator, closing array, string, number
	switch r := l.next(); {
	case r == ',':
		l.emit(itemArraySeperator)
		return lexInsideArray
	case r == '}':
		l.backup()
		return lexArrayClose
	case r == '"':
		l.backup()
		return lexArrayString
	case isNumber(r):
		l.backup()
		return lexArrayNumber
	default:
		return l.errorf("unrecognized character inside array: %#U", r)
	}
}

func lexAssignmentEqual(l *lexer) stateFn {
	lexOptionalSpace(l)
	if r := l.next(); r != '=' {
		return l.errorf("assignment missing equals, was %q", r)
	}
	l.emit(itemEqual)
	return lexAssignmentValue
}

func doString(l *lexer) bool {
	lexOptionalSpace(l)
	if r := l.next(); r != '"' {
		l.errorf("Missing string doublequote")
		return false
	}
	l.emit(itemStringDelim)
	for {
		r := l.next()
		if r != '"' {
			continue
		} else {
			if l.peek() == '"' {
				l.next()
				continue
			}
			if r == eof {
				l.errorf("Unclosed string")
				return false
			}
			l.backup()
			l.emit(itemString)
			break
		}
	}
	if r := l.next(); r != '"' {
		l.errorf("Unclosed string")
		return false
	}
	l.emit(itemStringDelim)
	return true
}

func lexAssignmentString(l *lexer) stateFn {
	if !doString(l) {
		return nil
	}
	if r := l.next(); r != ';' {
		return l.errorf("Unclosed string")
	}
	l.emit(itemSemicolon)
	return lexInsideClass
}

func lexAssignmentValue(l *lexer) stateFn {
	lexOptionalSpace(l)
	switch r := l.next(); {
	case isNumber(r):
		l.backup()
		return lexNumber(l)
	case r == '"':
		l.backup()
		return lexAssignmentString(l)
	case r == eof:
		l.emit(itemEOF)
		return nil
	case r == ';': //TODO: make state func
		l.emit(itemSemicolon)
		return nil
	default:
		return l.errorf("unrecognized character in assignment value: %#U", r)
	}
}

func lexSpaceBeforeClassIdentifier(l *lexer) stateFn {
	l.acceptRun(space)
	if !l.hasContent() {
		return l.errorf("Missing space after class keyword")
	}
	l.emit(itemSpace)
	return lexClassIdentifier
}

func lexClassIdentifier(l *lexer) stateFn {
	if !l.accept(alphaLower + alphaUpper) {
		return l.errorf("Class identifier does not start with an alpha character")
	}
	l.acceptRun(alphaLower + alphaUpper + digits)
	l.emit(itemIdentifier)
	return lexClassOpenBracket
}

func lexClassOpenBracket(l *lexer) stateFn {
	lexOptionalSpace(l)
	if r := l.next(); r != '{' {
		return l.errorf("Missing class opening bracket")
	}
	l.emit(itemOpenBlock)
	return lexInsideClass
}

func lexInsideClassCloseBracket(l *lexer) stateFn {
	lexOptionalSpace(l)
	if !l.accept(closeBracket) {
		return l.errorf("Missing closing bracket")
	}
	if r := l.next(); r != ';' {
		return l.errorf("Missing semicolon after closing bracket")
	}
	l.emit(itemCloseBlock)
	return lexInsideClass
}

func lexInsideClass(l *lexer) stateFn {
	lexOptionalSpace(l)
	// class, attribute, closing curly bracket
	switch r := l.next(); {
	case isAlpha(r):
		l.backup()
		return lexIdentifier
	case isCloseBracket(r):
		l.backup()
		return lexInsideClassCloseBracket
	case r == eof:
		l.emit(itemEOF)
		return nil
	default:
		return l.errorf("unrecognized character inside class: %#U", r)
	}
}
