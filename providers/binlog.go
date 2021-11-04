package providers

import (
	"context"

	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/prometheus/common/log"
	"gitlab.itoodev.com/wrkgrp/binlog/arbitrary"
	"gitlab.itoodev.com/wrkgrp/binlog/guards"
	"gitlab.itoodev.com/wrkgrp/binlog/tester"
	"gitlab.itoodev.com/wrkgrp/binlog/transport"
	"gitlab.itoodev.com/wrkgrp/binlog/utils"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type Binlog struct {
	Q    transport.Queue
	ctx  context.Context
	cfg  replication.BinlogSyncerConfig
	errs []error
}

func NewBinlog(ctx context.Context) Data {
	return Binlog{
		ctx: ctx,
		cfg: replication.BinlogSyncerConfig{
			ServerID: 100,
			Flavor:   "mysql",
			Host:     "192.168.178.129",
			Port:     3306,
			User:     "root",
			Password: "QTkq81vE",
		},
	}
}

func (provider Binlog) RequestChannel(dat spdg.Type) {
	state, reason := dat.State()
	machine := arbitrary.NewMachine(arbitrary.STATE)

	machine.Walker(
		state, reason,
		func(spdg.Status, spdg.Reason) {
			// TODO: Inject status and reason into dat.
		},
		func(spdg.Status, spdg.Reason) {
			provider.unit(&dat)
		},
	)
}

func (provider Binlog) unit(dat *spdg.Type) {
	syncer := replication.NewBinlogSyncer(provider.cfg)
	provider.pusher(dat, syncer)
}

func (provider Binlog) pusher(dat *spdg.Type, raw *replication.BinlogSyncer) {
	out, err := utils.Encode(raw)
	provider.recover(err)

	*dat = utils.Wrap(provider.ctx, out)
}

func (provider Binlog) addError(err error) {
	provider.errs = append(provider.errs, err)
}

func (provider Binlog) logErrors() {
	for _, err := range provider.errs {
		log.Errorln(err)
	}
}

func (provider Binlog) recover(err error) {
	if err != nil {
		log.Errorln("attempting to recover from error state")

		// Perform basic error reporting duties.
		provider.addError(err)
		provider.logErrors()

		// Add a rescue guard to this process, which provides guaranteed up state.
		// Generally considered a bad idea, unless you have a reason.
		// Do not get stuck!
		guards.Rescue(provider.dontPanic)
	}
}

/* dontPanic is called when recover fails, or if any critical fatal error or panic is thrown.
This is done via the Rescue guard, which captures these states and accepts a recue method as
a function. This technique makes it impossible for this code to crash, however know how to
deal with that situation before using something like this! Do not get stuck in infinite
recursion, not good for memory. */
func (provider Binlog) dontPanic() {
	// Nest a self-similar guard and cause recursion. Crash/Retry/ad. inifi.
	guards.Rescue(provider.dontPanic)

	log.Errorln("attemmpting to recover from critical state")

	log.Debugln("dumping accessible objects")
	log.Debugln(provider)

	log.Debugln("reinitializing object")
	replacement := NewBinlog(context.Background())

	// Run any type of testing framework on the replacement for deep inspection.
	log.Debugln("testing replacement")
	log.Debugln(replacement)
	tester := tester.Errors{
		Inspect: replacement,
	}

	if ok := tester.Asserts(); !ok {
		provider.dontPanic()
	}

	// Encode to the last known data type.
	log.Debugln("encoding object")
	encoded, _ := utils.Encode(replacement)

	// Wrap into an SPDG.
	log.Debugln("wrapping object")
	wrapped := utils.Wrap(provider.ctx, encoded)

	// Open up an arbitrary pipe and send the wrapper. No longer your problem.
	// If somewhere down the line the context is cancelled this will get cleaned up.
	// If the context happens to have a timeout associated, it will make everything
	// that knows and cares go back to retrying. But don't get stuck!
	log.Debugln("replacing object")
	pipe := arbitrary.NewPipe(arbitrary.REFLECTIVE)
	pipe.Have(wrapped)
}

func (provider Binlog) RequestReport() []error {
	return provider.errs
}
