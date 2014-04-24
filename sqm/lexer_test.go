package sqm

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestLexerBasics(t *testing.T) {
	var l *lexer
	Convey("Given a fresh lexer with \"input\"", t, func() {
		l = makeLexer("lexer", strings.NewReader("input"))
		Convey("It reads next rune and backup", func() {
			l.next()
			l.backup()
		})
		Convey("It reads till eof", func() {
			for i := 0; i < 5; i++ {
				l.next()
			}
			last := l.next()
			So(last, ShouldEqual, eof)
		})
		Convey("It has valid next runes", func() {
			So(l.next(), ShouldEqual, 'i')
			So(l.next(), ShouldEqual, 'n')
			So(l.next(), ShouldEqual, 'p')
			So(l.next(), ShouldEqual, 'u')
			So(l.next(), ShouldEqual, 't')
			So(l.next(), ShouldEqual, eof)
		})
		Convey("It should support last()", func() {
			l.next()
			So(l.last(), ShouldEqual, 'i')
		})
		Convey("It should support content()", func() {
			l.next()
			l.next()
			So(l.content(), ShouldEqual, "in")
		})
		Convey("It should emit", func() {
			l.next()
			l.next()
			l.emit(itemIdentifier)
			i := l.nextItem()
			So(i.val, ShouldEqual, "in")
			So(l.content(), ShouldEqual, "")
			So(l.next(), ShouldEqual, 'p')
			l.next()
			So(l.content(), ShouldEqual, "pu")

		})
		Convey("It should backup()", func() {
			l.next()
			l.next()
			l.backup()
			So(l.next(), ShouldEqual, 'n')
		})
		Convey("It should ignore()", func() {
			l.next()
			l.next()
			l.ignore()
			l.next()
			So(l.content(), ShouldEqual, "p")
		})
		Convey("It should peek()", func() {
			So(l.peek(), ShouldEqual, 'i')
			So(l.next(), ShouldEqual, 'i')
			So(l.next(), ShouldEqual, 'n')
		})
	})
	Convey("Given a fresh lexer with a long string", t, func() {
		l = makeLexer("lexer", strings.NewReader("inputinputinputinputinputinputinput"))
		So(l.next(), ShouldEqual, 'i')
		So(l.next(), ShouldEqual, 'n')
		l.backup()
		l.emit(itemIdentifier)
		i := l.nextItem()
		So(i.val, ShouldEqual, "i")

	})
}

// // func TestItemType(t *testing.T) {
// // 	var x itemType = itemError
// // 	var y itemType = itemDot
// // 	if x == y {
// // 		t.Errorf("Both types are the same")
// // 	}
// // }

// func ExampleItemString() {
// 	var i item = item{typ: itemEqual, val: "="}
// 	fmt.Println(i)
// 	// Output:
// 	// type: 6 val: "="
// }

// func TestMakeLexer(t *testing.T) {
// 	const name, input = "lexer", "input"
// 	l := makeLexer(name, input)
// 	if l.input != input {
// 		t.Errorf("Lexer input should be %q but was %q", input, l.input)
// 	}

// 	if l.items == nil {
// 		t.Errorf("Lexer returned invalid channel")
// 	}
// }

// func TestLexerNext(t *testing.T) {
// 	const name, input = "lexer", "a"
// 	const testRune rune = 'a'
// 	const start, pos, width = 0, 1, 1
// 	l := makeLexer(name, input)
// 	rune := l.next()
// 	if rune != testRune {
// 		t.Errorf("Next returned wrong rune %q", rune)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerIgnore(t *testing.T) {
// 	const name, input = "lexer", "ab"
// 	const start, pos, width = 1, 1, 1
// 	l := makeLexer(name, input)
// 	l.next()
// 	l.ignore()
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerBackup(t *testing.T) {
// 	const name, input = "lexer", "ab"
// 	const start, pos, width = 0, 0, 1
// 	l := makeLexer(name, input)
// 	l.next()
// 	l.backup()
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerPeek(t *testing.T) {
// 	const name, input = "lexer", "abc"
// 	const start, pos, width = 0, 1, 1
// 	l := makeLexer(name, input)
// 	l.next()
// 	l.peek()
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerLast(t *testing.T) {
// 	const name, input = "lexer", "a"
// 	const testRune rune = 'a'
// 	const start, pos, width = 0, 1, 1
// 	l := makeLexer(name, input)
// 	rune := l.next()
// 	runeLast := l.last()
// 	if rune != testRune {
// 		t.Errorf("Next returned wrong rune %q", rune)
// 	}
// 	if rune != runeLast {
// 		t.Errorf("Last returned wrong rune %q", runeLast)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerAccept(t *testing.T) {
// 	const name, input, accept = "lexer", "abc", "b"
// 	const start, pos, width = 0, 2, 1
// 	l := makeLexer(name, input)
// 	l.next()
// 	if !l.accept(accept) {
// 		t.Errorf("Does not accept %q", accept)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerAcceptFail(t *testing.T) {
// 	const name, input, accept = "lexer", "abc", "c"
// 	const start, pos, width = 0, 1, 1
// 	l := makeLexer(name, input)
// 	l.next()
// 	if l.accept(accept) {
// 		t.Errorf("Does accept %q but should not", accept)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// 	if l.width != width {
// 		t.Errorf("Width wrong %q", l.width)
// 	}
// }

// func TestLexerAcceptRun(t *testing.T) {
// 	const name, input, accept = "lexer", "abc", "ab"
// 	const acceptCount = 2
// 	const start, pos, width = 0, 2, 1
// 	l := makeLexer(name, input)
// 	if times := l.acceptRun(accept); times != acceptCount {
// 		t.Errorf("Does not accept %q times but %q times", acceptCount, times)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// }

// func TestLexerAcceptRunFail(t *testing.T) {
// 	const name, input, accept = "lexer", "abc", "bc"
// 	const acceptCount = 0
// 	const start, pos, width = 0, 0, 0
// 	l := makeLexer(name, input)
// 	if times := l.acceptRun(accept); times != acceptCount {
// 		t.Errorf("Does not accept %q times but %q times", acceptCount, times)
// 	}
// 	if l.start != start {
// 		t.Errorf("Start wrong %q", l.start)
// 	}
// 	if l.pos != pos {
// 		t.Errorf("Pos wrong %q", l.pos)
// 	}
// }

// // func TestLexNumberInt(t *testing.T) {
// // 	const name, input, value = "lexer", "123", "123"
// // 	const typ itemType = itemInt

// // 	l := makeLexer(name, input)
// // 	go l.run()
// // 	item := <-l.items
// // 	if item.typ != typ {
// // 		t.Errorf("Item type expected %q but was %q", typ, item.typ)
// // 	}
// // 	if item.val != value {
// // 		t.Errorf("Item value expected %q but was %q", value, item.val)
// // 	}
// // 	if isClosed := <-l.items; isClosed.typ != itemEOF {
// // 		t.Errorf("Channel was not closed, further input pending?")
// // 	}
// // }

// // func TestLexNumberFloat(t *testing.T) {
// // 	const name, input, value = "lexer", "123.456", "123.456"
// // 	const typ itemType = itemFloat

// // 	l := makeLexer(name, input)
// // 	go l.run()
// // 	item := <-l.items
// // 	if item.typ != typ {
// // 		t.Errorf("Item type expected %q but was %q", typ, item.typ)
// // 	}
// // 	if item.val != value {
// // 		t.Errorf("Item value expected %q but was %q", value, item.val)
// // 	}
// // 	if isClosed := <-l.items; isClosed.typ != itemEOF {
// // 		t.Errorf("Channel was not closed, further input pending?")
// // 	}
// // }

// type lexTest struct {
// 	name  string
// 	input string
// 	items []item
// }

// func collect(t *lexTest) (items []item) {
// 	l := makeLexer(t.name, t.input)
// 	go l.run()
// 	for {
// 		item := l.nextItem()
// 		items = append(items, item)
// 		if item.typ == itemEOF || item.typ == itemError {
// 			break
// 		}
// 	}
// 	return
// }

// func equals(items, expItems []item, checkPos bool) bool {
// 	if len(items) != len(expItems) {
// 		return false
// 	}
// 	for k, testItem := range items {
// 		expItem := expItems[k]
// 		if testItem.typ != expItem.typ {
// 			return false
// 		}
// 		if testItem.val != expItem.val {
// 			return false
// 		}
// 		if checkPos && testItem.pos != expItem.pos {
// 			return false
// 		}
// 	}
// 	return true
// }

// var (
// 	tEOF = item{itemEOF, 0, ""}
// )
// var lexTests = []lexTest{
// 	{"attribute number", "version=12;", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemInt, 7, "12"},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"attribute float", "version=123.456;", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemFloat, 7, "123.456"},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"attribute string", "version=\"test\";", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 7, "test"},
// 		{itemStringDelim, 0, "\""},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"attribute string escaped", "version=\"test=\"\"value\"\";\";", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 7, "test=\"\"value\"\";"},
// 		{itemStringDelim, 0, "\""},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"attribute string double escaped", "version=\"ret2=[\"\"ret=[\"\"\"\"val\"\"\"\"] call fnc;\"\"] call fnc;\";", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 7, "ret2=[\"\"ret=[\"\"\"\"val\"\"\"\"] call fnc;\"\"] call fnc;"},
// 		{itemStringDelim, 0, "\""},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"attribute string triple escaped", "version=\"ret2=[\"\"ret=[\"\"\"\"val\"\"\"\"] call fnc;\"\"] call fnc;\";", []item{
// 		{itemIdentifier, 0, "version"},
// 		{itemEqual, 6, "="},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 7, "ret2=[\"\"ret=[\"\"\"\"val\"\"\"\"] call fnc;\"\"] call fnc;"},
// 		{itemStringDelim, 0, "\""},
// 		{itemSemicolon, 0, ";"},
// 		tEOF,
// 	}},

// 	{"array string", "array[]={\"test1\",\"test2\"};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 0, "test1"},
// 		{itemStringDelim, 0, "\""},
// 		{itemArraySeperator, 0, ","},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 0, "test2"},
// 		{itemStringDelim, 0, "\""},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"array integer", "array[]={123,456};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemInt, 0, "123"},
// 		{itemArraySeperator, 0, ","},
// 		{itemInt, 0, "456"},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"array float", "array[]={123.456,456.789};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemFloat, 0, "123.456"},
// 		{itemArraySeperator, 0, ","},
// 		{itemFloat, 0, "456.789"},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"array empty", "array[]={};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"array single string", "array[]={\"test\"};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 0, "test"},
// 		{itemStringDelim, 0, "\""},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"array string multiline", "array[]={\n\"test\"\n};", []item{
// 		{itemIdentifier, 0, "array"},
// 		{itemIdentifierArrayDec, 0, "[]"},
// 		{itemEqual, 0, "="},
// 		{itemOpenArray, 0, "{"},
// 		{itemSpace, 0, "\n"},
// 		{itemStringDelim, 0, "\""},
// 		{itemString, 0, "test"},
// 		{itemStringDelim, 0, "\""},
// 		{itemSpace, 0, "\n"},
// 		{itemCloseArray, 0, "};"},
// 		tEOF,
// 	}},

// 	{"class empty", "class ident {};", []item{
// 		{itemClass, 0, "class"},
// 		{itemSpace, 0, " "},
// 		{itemIdentifier, 0, "ident"},
// 		{itemSpace, 0, " "},
// 		{itemOpenBlock, 0, "{"},
// 		{itemCloseBlock, 0, "};"},
// 		tEOF,
// 	}},

// 	{"class attribute", "class ident {\nunits=3;\n};", []item{
// 		{itemClass, 0, "class"},
// 		{itemSpace, 0, " "},
// 		{itemIdentifier, 0, "ident"},
// 		{itemSpace, 0, " "},
// 		{itemOpenBlock, 0, "{"},
// 		{itemSpace, 0, "\n"},
// 		{itemIdentifier, 0, "units"},
// 		{itemEqual, 0, "="},
// 		{itemInt, 0, "3"},
// 		{itemSemicolon, 0, ";"},
// 		{itemSpace, 0, "\n"},
// 		{itemCloseBlock, 0, "};"},
// 		tEOF,
// 	}},
// }

// func TestLexer(t *testing.T) {
// 	for _, test := range lexTests {
// 		items := collect(&test)
// 		if !equals(test.items, items, false) {
// 			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
// 		}
// 	}
// }

// func TestMissionSQM(t *testing.T) {
// 	const name = "Mission.sqm parser"
// 	if testing.Short() {
// 		t.Skip("Skip mission.sqm in short mode")
// 		return
// 	}
// 	buf, err := ioutil.ReadFile("./mission.sqm")
// 	if err != nil {
// 		t.Errorf("Could not open mission.sqm")
// 		return
// 	}
// 	input := string(buf)
// 	l := makeLexer(name, input)

// 	go l.run()
// 	i := 0
// 	for {
// 		item := l.nextItem()
// 		i++
// 		if item.typ == itemEOF {
// 			t.Logf("Successfully imported %d items from file", i)
// 			return
// 		} else if item.typ == itemError {
// 			t.Errorf("Got error %q", item)
// 			return
// 		}
// 	}
// }
