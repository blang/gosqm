package sqm

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const eof rune = -1
const (
	BUFREADER_SIZE   = 16
	BUFFER_INIT_SIZE = 16
)

type itemType int

const (
	itemEOF itemType = iota //make it the null type for closed channels
	itemError
	itemInt
	itemFloat
	itemIdentifier
	itemIdentifierArrayDec //[]
	itemEqual              // =
	itemSemicolon
	itemSpace          //also newlines
	itemOpenBlock      //{
	itemCloseBlock     //};
	itemOpenArray      //{
	itemCloseArray     //}
	itemClass          //class
	itemArraySeperator //,
	itemStringDelim    //"
	itemString
)

//TODO: item might need position information for later debugging
type item struct {
	typ itemType
	pos int //starting position in bytes in the input string
	// TODO: Should be character instead of byte?
	val string
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
	name    string //used for errors
	rd      *bufio.Reader
	buf     []rune
	start   int // start position of the item
	pos     int // current position in the input
	lastpos int
	width   int       // width of last rune read
	items   chan item // channel of scanned items
}

// var startState stateFn = func(l *lexer) stateFn {
// 	return func(l2 *lexer) stateFn {
// 		return nil
// 	}
// }

var startState stateFn = lexInsideClass

func makeLexer(name string, rd io.Reader) *lexer {
	bufrd := bufio.NewReaderSize(rd, BUFREADER_SIZE)

	l := &lexer{
		name:  name,
		items: make(chan item, 10),
		buf:   make([]rune, 0, BUFFER_INIT_SIZE),
		rd:    bufrd,
	}
	return l
}

// func (l *lexer) lex() (*lexer, chan item) {
// 	go l.run()
// 	return l, l.items
// }

// Run the lexer as a state machine until state is nil
func (l *lexer) run() {
	for state := startState; state != nil; {
		state = state(l)
	}
	close(l.items) //TODO: close 'emits' null type, maybe unwanted
}

func (l *lexer) emit(t itemType) {
	//TODO: transmit correct token positions
	l.items <- item{t, l.start, string(l.buf[l.start:l.pos])}

	//reset positions to reuse buffer
	if l.pos < l.lastpos { //there are read bytes left
		var newbuf []rune
		if (l.lastpos - l.pos) > BUFFER_INIT_SIZE {
			newbuf = make([]rune, 0, l.lastpos-l.pos)
		} else {
			newbuf = make([]rune, 0, BUFFER_INIT_SIZE)
		}
		for i := l.pos; i < l.lastpos; i++ {
			newbuf = append(newbuf, l.buf[i])
		}
		l.buf = newbuf
		l.lastpos = l.lastpos - l.pos
	} else {
		l.lastpos = 0
	}

	l.start = 0
	l.pos = 0

}

func (l *lexer) next() (r rune) {
	// log.Println("next()")
	// defer log.Printf("After next: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d, buf: %s\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf), l.buf)
	//Rune already in buffer
	if l.pos+1 <= l.lastpos {
		// log.Printf("Rune already read: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d, buf: %s\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf), l.buf)
		//rune already read
		l.pos += 1
		return l.buf[l.pos-1]

	} else { // Rune not in buffer
		r, s, err := l.rd.ReadRune()
		if err != nil {
			// log.Printf("error occurred: %s\n", err.Error())
			r = eof
			// if err == io.EOF {
			// 	return eof
			// } else {
			// 	//TODO: Proper error handling
			// 	return eof
			// }
		}
		if s == 0 {
			// log.Printf("No input read")
			r = eof
			// return eof //TODO: Better error handling
		}
		// log.Printf("Read rune from reader %c: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d\n", r, l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf))
		if l.pos+1 <= len(l.buf) {
			// log.Printf("Rune has space: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf))
			//Use existing buffer space by overwriting

			l.buf[l.pos] = r
			l.pos += 1
			l.lastpos = l.pos
			return r
		}
		//Extend buffer
		// log.Printf("Rune needs space, extend buffer: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf))
		l.pos = l.pos + 1
		l.lastpos = l.pos
		l.buf = append(l.buf, r)
		return r
	}
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
	if l.pos <= l.start {
		panic("Can't backup before start") //TODO: Remove debug panic
	}
	// log.Printf("Backup before: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf))
	l.pos -= 1
	// log.Printf("Backup: start: %d, pos: %d, lastpos: %d, len(buf): %d, cap(buf): %d\n", l.start, l.pos, l.lastpos, len(l.buf), cap(l.buf))
}

func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *lexer) content() string {
	return string(l.buf[l.start:l.pos])
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

const space = " \n\r"

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
	if l.content() == "class" {
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
		return l.errorf("unrecognized character in assignment value: %#U", r)
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
		return l.errorf("unrecognized character in assignment value: %#U", r)
	}
}
