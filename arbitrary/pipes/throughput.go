package pipes

import "gitlab.itoodev.com/wrkgrp/spdg"

/* ThroughputPipe is an arbitary opening to the internal network of the application and
is only concerned with piping the data through with minimal impact on transfer speeds.
Its primary use-case is relaying or rerouting traffic, should there be a need for that.
For instance if an object is recovering and there is another object available to take
over a workload. */
type ThroughputPipe struct {
}

func (pipe ThroughputPipe) Want() {

}

func (pipe ThroughputPipe) Have(dat spdg.Type) {

}
