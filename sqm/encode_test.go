package sqm

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func newBufEncoder() (*Encoder, *bytes.Buffer) {
	var e *Encoder
	var b []byte
	buf := bytes.NewBuffer(b)
	e = NewEncoder(buf)
	return e, buf

}

func TestEncode(t *testing.T) {
	var e *Encoder
	var buf *bytes.Buffer
	Convey("Given a fresh encoder", t, func() {
		e, buf = newBufEncoder()
		Convey("Given a Array Property of type String with 3 members", func() {
			arrProp := &ArrayProperty{
				Name:   "ArrayProp",
				Typ:    TString,
				Values: []string{"first", "second", "third"},
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeArrProperty(arrProp, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  ArrayProp[]={"first","second","third"};`+LINEBREAK)
				})
			})
		})

		Convey("Given a Array Property of type Integer with 3 members", func() {
			arrProp := &ArrayProperty{
				Name:   "ArrayProp",
				Typ:    TNumber,
				Values: []string{"1", "2", "3"},
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeArrProperty(arrProp, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  ArrayProp[]={1,2,3};`+LINEBREAK)
				})
			})
		})

		Convey("Given a Array Property of type Float with 3 members", func() {
			arrProp := &ArrayProperty{
				Name:   "ArrayProp",
				Typ:    TNumber,
				Values: []string{"1.123", "2.123", "3.123"},
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeArrProperty(arrProp, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  ArrayProp[]={1.123,2.123,3.123};`+LINEBREAK)
				})
			})
		})

		Convey("Given a Array Property of type String with special name 'addOns' with 3 members", func() {
			arrProp := &ArrayProperty{
				Name:   "addOns",
				Typ:    TString,
				Values: []string{"first", "second", "third"},
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeArrProperty(arrProp, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`  addOns[]=`+LINEBREAK+
							`  {`+LINEBREAK+
							`    "first",`+LINEBREAK+
							`    "second",`+LINEBREAK+
							`    "third"`+LINEBREAK+
							`  };`+LINEBREAK)
				})
			})
		})

		Convey("Given a Array Property of type String with special name 'addOnsAuto' with 3 members", func() {
			arrProp := &ArrayProperty{
				Name:   "addOnsAuto",
				Typ:    TString,
				Values: []string{"first", "second", "third"},
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeArrProperty(arrProp, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`  addOnsAuto[]=`+LINEBREAK+
							`  {`+LINEBREAK+
							`    "first",`+LINEBREAK+
							`    "second",`+LINEBREAK+
							`    "third"`+LINEBREAK+
							`  };`+LINEBREAK)
				})
			})
		})

		Convey("Given a String Property", func() {
			prop := &Property{
				Name:  "key",
				Typ:   TString,
				Value: "value",
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeProperty(prop, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  key="value";`+LINEBREAK)
				})
			})
		})

		Convey("Given an Integer Property", func() {
			prop := &Property{
				Name:  "key",
				Typ:   TNumber,
				Value: "123",
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeProperty(prop, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  key=123;`+LINEBREAK)
				})
			})
		})

		Convey("Given an Float Property", func() {
			prop := &Property{
				Name:  "key",
				Typ:   TNumber,
				Value: "123.456",
			}
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeProperty(prop, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  key=123.456;`+LINEBREAK)
				})
			})
		})

		Convey("Given an emtpy Class", func() {
			c := &Class{
				Name: "myclass",
			}
			Convey("Encoding it with level 0", func() {
				Convey("Should write correct string", func() {
					err := e.encodeClass(c, 0)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `class myclass`+LINEBREAK+`{`+LINEBREAK+`};`+LINEBREAK)
				})
			})
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeClass(c, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, `  class myclass`+LINEBREAK+`  {`+LINEBREAK+`  };`+LINEBREAK)
				})
			})
		})

		Convey("Given a Class with an empty subclass", func() {
			c := &Class{
				Name: "myclass",
				Classes: []*Class{
					&Class{
						Name: "subclass",
					},
				},
			}
			Convey("Encoding it with level 0", func() {
				Convey("Should write correct string", func() {
					err := e.encodeClass(c, 0)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`class myclass`+LINEBREAK+
							`{`+LINEBREAK+
							`  class subclass`+LINEBREAK+
							`  {`+LINEBREAK+
							`  };`+LINEBREAK+
							`};`+LINEBREAK)
				})
			})
			Convey("Encoding it with level 1", func() {
				Convey("Should write correct string", func() {
					err := e.encodeClass(c, 1)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`  class myclass`+LINEBREAK+
							`  {`+LINEBREAK+
							`    class subclass`+LINEBREAK+
							`    {`+LINEBREAK+
							`    };`+LINEBREAK+
							`  };`+LINEBREAK)
				})
			})
		})

		Convey("Given a full Class", func() {
			c := &Class{
				Name: "mainclass",
				Props: []*Property{
					&Property{
						Name:  "version",
						Typ:   TNumber,
						Value: "1",
					},
				},
				Classes: []*Class{
					&Class{
						Name: "myclass",
						Arrprops: []*ArrayProperty{
							&ArrayProperty{
								Name:   "addOnsAuto",
								Typ:    TString,
								Values: []string{"first", "second", "third"},
							},
						},
						Props: []*Property{
							&Property{
								Name:  "key1",
								Typ:   TString,
								Value: "value1",
							},
							&Property{
								Name:  "key2",
								Typ:   TString,
								Value: "value2",
							},
						},
						Classes: []*Class{
							&Class{
								Name: "subclass",
								Props: []*Property{
									&Property{
										Name:  "key1",
										Typ:   TString,
										Value: "value1",
									},
								},
							},
						},
					},
				},
			}
			Convey("Encoding it with level 0", func() {
				Convey("Should write correct string", func() {
					err := e.encodeClass(c, 0)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`class mainclass`+LINEBREAK+
							`{`+LINEBREAK+
							`  version=1;`+LINEBREAK+
							`  class myclass`+LINEBREAK+
							`  {`+LINEBREAK+
							`    addOnsAuto[]=`+LINEBREAK+
							`    {`+LINEBREAK+
							`      "first",`+LINEBREAK+
							`      "second",`+LINEBREAK+
							`      "third"`+LINEBREAK+
							`    };`+LINEBREAK+
							`    key1="value1";`+LINEBREAK+
							`    key2="value2";`+LINEBREAK+
							`    class subclass`+LINEBREAK+
							`    {`+LINEBREAK+
							`      key1="value1";`+LINEBREAK+
							`    };`+LINEBREAK+
							`  };`+LINEBREAK+
							`};`+LINEBREAK)
				})
			})

			Convey("Encoding it as main class", func() {
				Convey("Should write correct string", func() {
					err := e.Encode(c)
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual,
						`version=1;`+LINEBREAK+
							`class myclass`+LINEBREAK+
							`{`+LINEBREAK+
							`  addOnsAuto[]=`+LINEBREAK+
							`  {`+LINEBREAK+
							`    "first",`+LINEBREAK+
							`    "second",`+LINEBREAK+
							`    "third"`+LINEBREAK+
							`  };`+LINEBREAK+
							`  key1="value1";`+LINEBREAK+
							`  key2="value2";`+LINEBREAK+
							`  class subclass`+LINEBREAK+
							`  {`+LINEBREAK+
							`    key1="value1";`+LINEBREAK+
							`  };`+LINEBREAK+
							`};`+LINEBREAK)
				})
			})
		})
	})
	Convey("Indent should return correct amount of spaces", t, func() {
		So(indent(0), ShouldEqual, "")
		So(indent(1), ShouldEqual, "  ")
		So(indent(2), ShouldEqual, "    ")
	})
}
