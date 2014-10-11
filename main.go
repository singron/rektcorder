package main

import (
	"flag"
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
}

func init() {
	flag.BoolVar(&Config.Debug, "debug", false, "show debug output")
	flag.BoolVar(&Config.Joins, "joins", false, "show joins and leaves")
	flag.StringVar(&Config.Sid, "sid", "", "use sid to login")
	flag.BoolVar(&Config.Record, "record", false, "log messages for playback")
	flag.StringVar(&Config.RecordDir, "record-dir", ".", "dir in which to log messages for playback")
	flag.BoolVar(&Config.Play, "play", false, "playback message logs instead of listening")
	flag.StringVar(&Config.Start, "start", "", "time to start playback")
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

func listen() {
	var wsconfig websocket.Config
	if Config.Sid != "" {
		wsconfig.Header.Set("Cookie", "sid="+Config.Sid)
	}
	ws, err := websocket.Dial("ws://www.destiny.gg:9998/ws", "", "http://www.destiny.gg")
	if err != nil {
		log.Fatal(err)
	}
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
	for {
		var v interface{}
		err = Codec.Receive(ws, &v)
		if err != nil {
			log.Fatal(err)
		}
		switch v := v.(type) {
		case Msg:
			v.Print()
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
		default:
			log.Printf("unknown value %T %v\n", v, v)
		}
	}
}

func playback() {
	p := NewPlayer(Config.RecordDir)
	start, err := time.Parse("2006-01-02T15:04:05", Config.Start)
	if err != nil {
		log.Fatal(err)
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
