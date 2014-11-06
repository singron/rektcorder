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
	"flag"
	"io"
	"log"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/wsxiaoys/terminal/color"
)

var Config struct {
	Debug     bool
	Joins     bool
	Sid       string
	RecordDir string
	Record    bool
	Play      bool
	Start     string
	Silent    bool
	Timeout   time.Duration
}

func init() {
	flag.BoolVar(&Config.Debug, "debug", false, "show debug output")
	flag.BoolVar(&Config.Joins, "joins", false, "show joins and leaves")
	flag.StringVar(&Config.Sid, "sid", "", "use sid to login")
	flag.BoolVar(&Config.Record, "record", false, "log messages for playback")
	flag.StringVar(&Config.RecordDir, "record-dir", ".", "dir in which to log messages for playback")
	flag.BoolVar(&Config.Play, "play", false, "playback message logs instead of listening")
	flag.StringVar(&Config.Start, "start", "", "time to start playback")
	flag.BoolVar(&Config.Silent, "silent", false, "don't print messages")
	flag.DurationVar(&Config.Timeout, "timeout", time.Duration(30)*time.Second,
		"panic after no messages have been received for this long")
	flag.Parse()
}

func debug(v ...interface{}) {
	if Config.Debug {
		log.Panicln(v...)
	}
}

func debugf(f string, v ...interface{}) {
	if Config.Debug {
		log.Printf(f, v...)
	}
}

func connect() *websocket.Conn {
	var wsconfig websocket.Config
	if Config.Sid != "" {
		wsconfig.Header.Set("Cookie", "sid="+Config.Sid)
	}
	for {
		log.Printf("Connecting...\n")
		ws, err := websocket.Dial("ws://www.destiny.gg:9998/ws", "", "http://www.destiny.gg")
		if err == nil {
			return ws
		}
	}
}

func timeout() {
	panic("Chat timed out")
}

func listen() {
	ws := connect()
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var ml *MessageLog
	if Config.Record {
		ml = NewMessageLog(Config.RecordDir)
		defer ml.Close()
	}
	timer := time.AfterFunc(Config.Timeout, timeout)
	for {
		var v interface{}
		err := Codec.Receive(ws, &v)
		if err == io.EOF {
			ws.Close()
			debugf("EOF, Reconnect\n")
			ws = connect()
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		timer.Reset(Config.Timeout)
		switch v := v.(type) {
		case Msg:
			if !Config.Silent {
				v.Print()
			}
			if Config.Record {
				err = ml.Write(&v)
				if err != nil {
					log.Fatal(err)
				}
			}
		case Ping:
			var p Pong
			p.Data = v.Data
			err := Codec.Send(ws, p)
			if err != nil {
				log.Fatal(err)
			}
		case Names:
		case Join:
			if Config.Joins {
				color.Printf("@{/}* %s joined@{|}\n", v.Nick)
			}
		case Quit:
			if Config.Joins {
				color.Printf("@{/}* %s quit@{|}\n", v.Nick)
			}
		case Broadcast:
			if Config.Joins {
				color.Printf("@{/}* %s@{|}\n", v.Data)
			}
		default:
			log.Printf("unknown value %T %v\n", v, v)
		}
	}
}

func playback() {
	p := NewPlayer(Config.RecordDir)
	start, err := time.Parse("2006-01-02T15:04:05", Config.Start)
	if err != nil {
		debugf("%s isn't a date, trying twitch\n", Config.Start)
		var t Twitch
		v, err := t.Video(Config.Start)
		if err != nil {
			log.Fatal(err)
		}
		start = v.RecordedAt.UTC()
	}
	err = p.Playback(start)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if Config.Play {
		playback()
	} else {
		listen()
	}
}
