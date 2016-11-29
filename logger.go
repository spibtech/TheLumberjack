package main

import (
	"fmt"
	"os"
	"log"
	"time"
	"sync"
)


// Logger is the actual logging structure
type Logger struct {
	file    *os.File
	tag     string
	logname string
	year    int
	day     int
	month   int
	hour    int
	size    int64
	lock    sync.Mutex
}

func (l *Logger) log(t *time.Time, data string) {
	mins := getMinuteBlock(t.Minute())
	tag := fmt.Sprintf("%04d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), mins)
	l.lock.Lock()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	defer l.lock.Unlock()
	if l.tag == "" || l.tag != tag || l.file == nil {
		gConf.purgeLock.Lock()
		hasLocked := true
		defer func() {
			if hasLocked {
				gConf.purgeLock.Unlock()
			}
		}()
		// reaches limit of number of log files
		filename := fmt.Sprintf("%s%s.%s.log.%04d-%02d-%02d-%02d-%02d", gConf.pathPrefix, l.logname,
			hostname, t.Year(), t.Month(), t.Day(), t.Hour(), mins)
		newfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Error opening log file: %s - %s", filename, err)
			return
		}
		gConf.curfiles++
		gConf.purgeLock.Unlock()
		hasLocked = false

		l.file.Close()
		l.file = newfile
		l.tag = tag
		l.size = 0

	}

	n, _ := l.file.WriteString(data)
	l.size += int64(n)
}
