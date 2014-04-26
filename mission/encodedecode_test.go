package mission

import (
	"bytes"
	"github.com/blang/gosqm/sqm"
	"io/ioutil"
	"testing"
)

func BenchmarkFullDecode(b *testing.B) {
	buf, err := ioutil.ReadFile("../sqm/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		b.Errorf("Could not open mission.sqm")
		return
	}
	for n := 0; n < b.N; n++ {
		p := sqm.MakeParser(bufstr)
		// c, perr := p.Run()
		class, perr := p.Run()
		if perr != nil {
			b.Errorf("Parser returned with error %q", perr)
		}
		mp := NewParser()
		_, err := mp.Parse(class)
		if err != nil {
			b.Errorf("Can't parse class to missionfile, %q", err)
		}

	}
}

func BenchmarkFullEncode(b *testing.B) {
	buf, err := ioutil.ReadFile("../sqm/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		b.Errorf("Could not open mission.sqm")
		return
	}
	p := sqm.MakeParser(bufstr)
	// c, perr := p.Run()
	class, perr := p.Run()
	if perr != nil {
		b.Errorf("Parser returned with error %q", perr)
	}
	mp := NewParser()
	missionFile, err := mp.Parse(class)
	if err != nil {
		b.Errorf("Can't parse class to missionfile, %q", err)
	}
	for n := 0; n < b.N; n++ {
		enc := NewEncoder()
		eclass := enc.Encode(missionFile)
		var buf []byte
		buffer := bytes.NewBuffer(buf)
		sqmenc := sqm.NewEncoder(buffer)
		err := sqmenc.Encode(eclass)
		if err != nil {
			b.Errorf("Can't encode class, %q", err)
		}

	}
}

func BenchmarkFullEncodeDecode(b *testing.B) {
	buf, err := ioutil.ReadFile("../sqm/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		b.Errorf("Could not open mission.sqm")
		return
	}
	for n := 0; n < b.N; n++ {
		p := sqm.MakeParser(bufstr)
		// c, perr := p.Run()
		class, perr := p.Run()
		if perr != nil {
			b.Errorf("Parser returned with error %q", perr)
		}
		mp := NewParser()
		missionFile, err := mp.Parse(class)
		if err != nil {
			b.Errorf("Can't parse class to missionfile, %q", err)
		}
		enc := NewEncoder()
		eclass := enc.Encode(missionFile)
		var buf []byte
		buffer := bytes.NewBuffer(buf)
		sqmenc := sqm.NewEncoder(buffer)
		err = sqmenc.Encode(eclass)
		if err != nil {
			b.Errorf("Can't encode class, %q", err)
		}

	}
}
