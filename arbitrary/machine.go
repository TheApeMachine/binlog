package arbitrary

import "gitlab.itoodev.com/wrkgrp/spdg"

type MachineType uint

const (
	STATE MachineType = iota
)

/* Machine is an arbitrary type for any function, or function of nested functions, and
can operate in two modes. First as a walker, where it takes in multiple state and reason
scenarios which is automatically maps over handler functions. Second as a runner where
it will recursively loop through the pipeline, automatically handing over input and
return values. The runner looks deceptively simple, but is very powerful. */
type Machine interface {
	Walker(spdg.Status, spdg.Reason, func(spdg.Status, spdg.Reason), func(spdg.Status, spdg.Reason))
	Runner(spdg.Type) spdg.Type
}

type StateMachine struct {
}

func (machine StateMachine) Walker(
	state spdg.Status, reason spdg.Reason,
	unhappy func(spdg.Status, spdg.Reason),
	happy func(spdg.Status, spdg.Reason),
) {
	switch state {
	case spdg.NO, spdg.ERR:
		unhappy(state, reason)
	case spdg.OK:
		happy(state, reason)
	}
}

/* Runner looks simple, but is powerful. Why? Think what you can do if any mutation,
or any logic circuit, puts out the same type as came in, no matter what? That leaves
you only with infinite options. Just decode the SPDG, use reflection in the worst
case to get the real type back, and use it for any arbitrary workload.
Then encode out to SPDG again and return. */
func (machine StateMachine) Runner(t spdg.Type) spdg.Type {
	return t
}

func NewMachine(t MachineType) Machine {
	switch t {
	case STATE:
		return StateMachine{}
	}

	return nil
}
