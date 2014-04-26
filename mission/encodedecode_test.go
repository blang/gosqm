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

func BenchmarkMissionEncode(b *testing.B) {
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
		enc.Encode(missionFile)
	}
}

func BenchmarkMissionDecode(b *testing.B) {
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
	for n := 0; n < b.N; n++ {
		mp := NewParser()
		missionFile, err := mp.Parse(class)
		if err != nil {
			b.Errorf("Can't parse class to missionfile, %q", err)
		}
		if len(missionFile.Mission.Groups) < 5 {
			b.Errorf("Error while parsing mission, %q", missionFile)
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
		var ebuf []byte
		buffer := bytes.NewBuffer(ebuf)
		sqmenc := sqm.NewEncoder(buffer)
		err = sqmenc.Encode(eclass)
		if err != nil {
			b.Errorf("Can't encode class, %q", err)
		}

	}
}

func TestFullEncodeDecode(t *testing.T) {
	buf, err := ioutil.ReadFile("../sqm/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		t.Errorf("Could not open mission.sqm")
		return
	}
	p := sqm.MakeParser(bufstr)
	// c, perr := p.Run()
	class, perr := p.Run()
	if perr != nil {
		t.Errorf("Parser returned with error %q", perr)
	}
	mp := NewParser()
	missionFile, err := mp.Parse(class)
	if err != nil {
		t.Errorf("Can't parse class to missionfile, %q", err)
	}
	enc := NewEncoder()
	eclass := enc.Encode(missionFile)
	var ebuf []byte
	buffer := bytes.NewBuffer(ebuf)
	sqmenc := sqm.NewEncoder(buffer)
	err = sqmenc.Encode(eclass)
	if err != nil {
		t.Errorf("Can't encode class, %q", err)
	}
}
