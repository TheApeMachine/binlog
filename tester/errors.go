package tester

import (
	"errors"
	"log"
)

type Tester interface {
	Asserts() bool
}

type ErrorReporter interface {
	RequestReport() []error
}

type Errors struct {
	Inspect ErrorReporter
}

func (test Errors) Asserts() bool {
	return test.valuetest()
}

func (test Errors) valuetest() bool {
	go func() {
		if len(test.Inspect.RequestReport()) != 1 || test.Inspect.RequestReport()[0] != errors.New("test") {
			log.Fatalln("it should collect errors")
		}
	}()

	return true
}
