package transport

import "gitlab.itoodev.com/wrkgrp/spdg"

type Queue interface {
	Cycle()
	Add([]chan spdg.Type)
}

type SystemQueue struct {
	io  [][]chan spdg.Type
	msg [][][]spdg.Type
}

func NewSystemQueue() Queue {
	return SystemQueue{}
}

func (q SystemQueue) Add(listener []chan spdg.Type) {
	q.io[1] = append(q.io[1], listener[1])
	q.io[0] = append(q.io[0], listener[0])
}

func (q SystemQueue) Cycle() {
	for i, qio := range q.io {
		q.cycleio(i, qio)
	}
}

func (q SystemQueue) cycleio(i int, qio []chan spdg.Type) {
	for j, io := range qio {
		// Pop from the msg stack at the corresponding index for this channel.
		var pop spdg.Type
		pop, q.msg[i][j] = q.msg[i][j][0], q.msg[i][j][1:]

		// Send the message to the listener.
		io <- pop
	}
}
