package guards

import "github.com/prometheus/common/log"

func Rescue(handler func()) func() {
	return func() {
		recoverer(handler)
	}
}

func recoverer(handler func()) {
	if r := recover(); r != nil {
		errorhandler(handler)
	}
}

func errorhandler(handler func()) {
	if handler == nil {
		log.Errorln("rescued without handler")
	} else {
		log.Errorln("rescued with handler")
		handler()
	}
}
