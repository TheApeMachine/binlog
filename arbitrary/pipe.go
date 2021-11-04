package arbitrary

/* Maybe this package flies in the face of all that is holy, maybe not. Maybe you shouldn't
worry about it. Nothing is out of scope as long as it fits the context of the overall design.
Mine is about not having to worry about complexity, and the rule is if it's not easy, it's not
finished. I want to have arbitrary dumpable objects I can deposit objects in and not care how
they get to the right place. Of course this makes much more sense if you use a singular type
across your whole pipeline, which is what this project is really about. */

import (
	"gitlab.itoodev.com/wrkgrp/binlog/arbitrary/pipes"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type PipeType uint

const (
	THROUGHPUT PipeType = iota
	REFLECTIVE
)

type Pipe interface {
	Want()
	Have(spdg.Type)
}

func NewPipe(t PipeType) Pipe {
	switch t {
	case THROUGHPUT:
		return pipes.ThroughputPipe{}
	case REFLECTIVE:
		return pipes.ReflectivePipe{}
	}

	return nil
}
