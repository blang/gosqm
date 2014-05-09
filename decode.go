package gosqm

import (
	"github.com/blang/gosqm/sqm"
	"io"
	"io/ioutil"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (d *Decoder) Decode() (*MissionFile, error) {
	b, err := ioutil.ReadAll(d.r)
	if err != nil {
		return nil, err
	}
	bufstr := string(b)
	p := sqm.MakeParser(bufstr)
	class, perr := p.Run()
	if perr != nil {
		return nil, perr
	}
	mp := NewParser()
	return mp.Parse(class)
}
