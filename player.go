//rektcorder is a recorder and player for destiny.gg chat logs
//Copyright (C) 2013 Eric Culp
//
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU Affero General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU Affero General Public License for more details.
//
//You should have received a copy of the GNU Affero General Public License
//along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
