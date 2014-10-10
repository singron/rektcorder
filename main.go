package main

import (
	"flag"
	"log"

	"code.google.com/p/go.net/websocket"
	"github.com/wsxiaoys/terminal/color"
)

var Config struct {
	Debug     bool
	Joins     bool
	Sid       string
	RecordDir string
	Record    bool
}

func init() {
	flag.BoolVar(&Config.Debug, "debug", false, "show debug output")
	flag.BoolVar(&Config.Joins, "joins", false, "show joins and leaves")
	flag.StringVar(&Config.Sid, "sid", "", "use sid to login")
	flag.BoolVar(&Config.Record, "record", false, "log messages for playback")
	flag.StringVar(&Config.RecordDir, "record-dir", ".", "dir in which to log messages for playback")
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

func main() {
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
			color.Printf("<@{!b}%s@{|}>: %s\n", v.Nick, v.Data)
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
