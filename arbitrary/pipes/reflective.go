package pipes

import "gitlab.itoodev.com/wrkgrp/spdg"

/* ReflectivePipe is still a relatively normal throuhgput into the internal networking pool
only this time being able to target a little more specifically the destination of the payload
by first reflecting on the incoming type and gathering some more information where possible.
It might also be a dead end and turn out not to be needed, but it feels right. */
type ReflectivePipe struct {
}

func (pipe ReflectivePipe) Want() {
}

func (pipe ReflectivePipe) Have(dat spdg.Type) {
}
