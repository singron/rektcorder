package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Player struct {
	log     *os.File
	logTime time.Time
	decoder *json.Decoder
	dir     string
	offset  time.Duration
}

func NewPlayer(dir string) *Player {
	var p Player
	p.dir = dir
	return &p
}

func (p *Player) Reset() error {
	n := p.Then()
	p.logTime = logTime(n)
	filename := logFileName(n)
	if p.log != nil {
		p.log.Close()
	}
	log.Printf("Playing file %s\n", filename)
	var err error
	p.log, err = os.Open(filepath.Join(p.dir, filename))
	if err != nil {
		return err
	}
	p.decoder = json.NewDecoder(p.log)
	return nil
}

// Playback plays back logs starting at t
func (p *Player) Playback(t time.Time) error {
	log.Printf("Playing back starting at %s\n", t)
	p.offset = time.Now().Sub(t)
	err := p.Reset()
	if err != nil {
		return err
	}
	for {
		var m Msg
		err := p.decoder.Decode(&m)
		if err != nil {
			return err
		}
		d := m.Timestamp.Time.Sub(p.Then())
		debugf("difference is %f seconds", d.Seconds())
		if d > 0 {
			time.Sleep(d)
		}
		m.Print()
	}
}

func (p *Player) Then() time.Time {
	return time.Now().Add(-p.offset)
}

func (p *Player) Close() error {
	p.decoder = nil
	return p.log.Close()
}
