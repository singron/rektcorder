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
	l.logTime = logTime(time.Now())
	l.Reset()
	return &l
}

func logTime(n time.Time) time.Time {
	return n.UTC().Truncate(time.Hour)
}

func logFileName(n time.Time) string {
	return fmt.Sprintf("chat-%s.log", n.UTC().Format("2006-01-02-15"))
}

func (l *MessageLog) Reset() error {
	if l.log != nil {
		l.log.Close()
	}
	filename := logFileName(l.logTime)
	log.Printf("Recording to %s\n", filename)
	var err error
	l.log, err = os.OpenFile(filepath.Join(l.dir, filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	l.encoder = json.NewEncoder(l.log)
	return nil
}

func (l *MessageLog) Write(m *Msg) error {
	mtime := logTime(m.Timestamp.Time)
	if !mtime.Equal(l.logTime) {
		log.Printf("Log message time %s != logTime %s, rotating\n", mtime, l.logTime)
		l.logTime = mtime
		err := l.Reset()
		if err != nil {
			return err
		}
	}
	return l.encoder.Encode(m)
}

func (l *MessageLog) Close() error {
	l.encoder = nil
	return l.log.Close()
}
