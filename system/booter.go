package system

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/siddontang/go-log/log"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type Booter interface {
	Add(Booter)
	Kick() bool
	Inspect() spdg.Type
}

type SystemBooter struct {
	booters  []Booter
	logs     spdg.Type
	contexts []context.Context
	cancells []context.CancelFunc
}

/* SystemBooter is the top-level object that starts the program. It is a self-nesting
type so after this it is booters all the way down. */
func NewSystemBooter() SystemBooter {
	return SystemBooter{
		contexts: make([]context.Context, 0),
		cancells: make([]context.CancelFunc, 0),
	}
}

func (booter SystemBooter) Failsafe() {
	sigint := make(chan os.Signal, 1)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGTRAP)

	go func() {
		booter.sighandle(sigint, sigs, done)
	}()
}

func (booter SystemBooter) dispose(done chan bool) {
	for _, cancel := range booter.cancells {
		cancel()
	}

	done <- true
}

/* Add appends a sub booter to the object. */
func (booter SystemBooter) Add(btr Booter) {
	booter.booters = append(booter.booters, btr)
}

/* Kick cycles the booter and propegates downwards. */
func (booter SystemBooter) Kick() bool {
	out := false

	for _, booter := range booter.booters {
		out = booter.Kick()
	}

	return out
}

/* Inspect the state of the booter, including historical. All SPDG types are WORMs
(Write Once Read Many) so no data can ever be lost. */
func (booter SystemBooter) Inspect() spdg.Type {
	return booter.logs
}

func (booter SystemBooter) sighandle(sigint chan os.Signal, sigs chan os.Signal, done chan bool) {
	select {
	case <-sigint:
		log.Warnln("user terminate, attempting clean exit")
		booter.dispose(done)
	case <-sigs:
		log.Errorln("system terminate, attempting clean restart")
		restart := NewSystemBooter()

		// Leaking goroutines for fun and profit, but it actually makes sense here.
		go func() {
			if ok := restart.Kick(); !ok {
				panic(restart.Inspect())
			}
		}()

		booter.dispose(done)
	case <-done:
		os.Exit(1)
	default:
	}
}
