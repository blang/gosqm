package gosqm

import (
	"bytes"
	"github.com/blang/gosqm/sqm"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkFullDecode(b *testing.B) {
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
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
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
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
		enc := NewClassEncoder()
		eclass := enc.EncodeToClass(missionFile)
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
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
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
		enc := NewClassEncoder()
		enc.EncodeToClass(missionFile)
	}
}

func BenchmarkMissionDecode(b *testing.B) {
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
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
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
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
		enc := NewClassEncoder()
		eclass := enc.EncodeToClass(missionFile)
		var ebuf []byte
		buffer := bytes.NewBuffer(ebuf)
		sqmenc := sqm.NewEncoder(buffer)
		err = sqmenc.Encode(eclass)
		if err != nil {
			b.Errorf("Can't encode class, %q", err)
		}

	}
}

func TestFullEncodeDecodeSingleSteps(t *testing.T) {
	buf, err := ioutil.ReadFile("testdata/mission.sqm")
	bufstr := string(buf)
	if err != nil {
		t.Fatalf("Could not open mission.sqm")
		return
	}
	p := sqm.MakeParser(bufstr)
	// c, perr := p.Run()
	class, perr := p.Run()
	if perr != nil {
		t.Fatalf("Parser returned with error %q", perr)
		return
	}
	mp := NewParser()
	missionFile, err := mp.Parse(class)
	if err != nil {
		t.Fatalf("Can't parse class to missionfile, %q", err)
		return
	}
	for _, perr := range mp.Warnings() {
		t.Logf("Warning: " + perr.Error() + "\n")
	}
	enc := NewClassEncoder()
	eclass := enc.EncodeToClass(missionFile)
	var ebuf []byte
	buffer := bytes.NewBuffer(ebuf)
	sqmenc := sqm.NewEncoder(buffer)
	err = sqmenc.Encode(eclass)
	if err != nil {
		t.Fatalf("Can't encode class, %q", err)
		return
	}
	// ioutil.WriteFile("mission.out.sqm", buffer.Bytes(), 0666)
}

func TestFullEncodeDecode(t *testing.T) {
	f, err := os.Open("testdata/mission.sqm")
	defer f.Close()
	if err != nil {
		t.Fatal("Could not open testdata")
		return
	}
	dec := NewDecoder(f)
	missionFile, err := dec.Decode()
	if err != nil {
		t.Fatalf("Decode error: %q", err.Error())
		return
	}

	if !(len(missionFile.Mission.Groups) > 0) {
		t.Error("Could not read groups")
	}

	var ebuf []byte
	buffer := bytes.NewBuffer(ebuf)
	enc := NewEncoder(buffer)
	err = enc.Encode(missionFile)
	if err != nil {
		t.Fatalf("Can't encode class, %q", err)
		return
	}
	ioutil.WriteFile("mission.out.sqm", buffer.Bytes(), 0666)
}
