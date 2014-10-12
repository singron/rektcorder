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
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/wsxiaoys/terminal/color"
)

var (
	ErrMessageNoHeader      = errors.New("No header in message")
	ErrMessageUnknownHeader = errors.New("Unknown header in message")
)

var Codec = websocket.Codec{marshall, unmarshall}

func unmarshall(msg []byte, payloadType byte, v interface{}) error {
	debugf("receive msg: %s\n", msg)
	i := bytes.IndexByte(msg, byte(' '))
	if i == -1 {
		return ErrMessageNoHeader
	}
	header := string(msg[0:i])
	body := msg[i+1 : len(msg)]
	var err error = ErrMessageNoHeader
	switch header {
	case "MSG":
		var o Msg
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "PING":
		var o Ping
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "PONG":
		var o Pong
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "NAMES":
		var o Names
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "JOIN":
		var o Join
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "QUIT":
		var o Quit
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "MUTE":
		var o Mute
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	case "BROADCAST":
		var o Broadcast
		err = json.Unmarshal(body, &o)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(o))
	default:
		log.Printf("Unknown header %s\n", header)
		log.Printf("    Body: %s\n", body)
		err = ErrMessageUnknownHeader
	}
	return err
}

func marshall(v interface{}) (msg []byte, payloadType byte, err error) {
	switch v.(type) {
	case Pong:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, 0, err
		}
		b = append([]byte("PONG "), b...)
		debugf("send msg: %s\n", string(b))
		return b, websocket.TextFrame, nil
	}
	return nil, 0, errors.New("unknown payload type")
}

type Msg struct {
	Nick      string   `json:"nick"`
	Features  []string `json:"features"`
	Timestamp Time     `json:"timestamp"`
	Data      string   `json:"data"`
}

func (m *Msg) Print() {
	color.Printf("<@{!b}%s@{|}>: %s\n", m.Nick, m.Data)
}

type User struct {
	Nick     string   `json:"nick"`
	Features []string `json:"features"`
}

type Names struct {
	Users []User `json:"users"`
}

type Ping struct {
	Data int64 `json:"data"`
}

type Pong struct {
	Data int64 `json:"data"`
}

type Join struct {
	Nick      string   `json:"nick"`
	Features  []string `json:"features"`
	Timestamp Time     `json:"timestamp"`
}

type Quit struct {
	Nick      string   `json:"nick"`
	Features  []string `json:"features"`
	Timestamp Time     `json:"timestamp"`
}

type Mute struct {
	Nick      string   `json:"nick"` // the nick who muted
	Features  []string `json:"features"`
	Timestamp Time     `json:"timestamp"`
	Data      string   `json:"data"` // the nick to mute
}

type Broadcast struct {
	Timestamp Time   `json:"timestamp"`
	Data      string `json:"data"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i/1000, (i%1000)*1000000)
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.UnixNano()/1000000, 10)), nil
}
