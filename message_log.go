package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type MessageLog struct {
	dir     string
	log     *os.File
	logTime time.Time
	encoder *json.Encoder
}

func NewMessageLog(messageDir string) *MessageLog {
	var l MessageLog
	l.dir = messageDir
	l.Reset()
	return &l
}

func (l *MessageLog) Reset() error {
	n := time.Now()
	l.logTime = time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), 0, 0, 0, time.Local)
	filename := fmt.Sprintf("chat-%s.log", time.Now().Format("2006-01-02-15"))
	if l.log != nil {
		l.log.Close()
	}
	var err error
	log.Printf("Recording to %s\n", filename)
	l.log, err = os.OpenFile(filepath.Join(l.dir, filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	l.encoder = json.NewEncoder(l.log)
	return nil
}

func (l *MessageLog) Write(m *Msg) error {
	hours := time.Now().Sub(l.logTime).Hours()
	if hours >= 1 {
		log.Printf("Log is %f hours behind, rotating\n", hours)
		err := l.Reset()
		if err != nil {
			return err
		}
	}
	return l.encoder.Encode(m)
}

func (l *MessageLog) Close() error {
	return l.log.Close()
}
