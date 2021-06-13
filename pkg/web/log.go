package web

import "log"

const (
	debugLevel = 0
	off        = 5
)

type logger struct {
	mode int
}



func (l *logger) debug(v ...interface{}) {
	if l.mode <= debugLevel {
		log.Println(v...)
	}
}

func (l *logger) debugf(format string, v ...interface{}) {
	if l.mode <= debugLevel {
		log.Printf(format, v...)
	}
}
