package projectors

import (
	"context"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"gitlab.itoodev.com/wrkgrp/binlog/utils"
	"gitlab.itoodev.com/wrkgrp/binlog/viewers"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type Binlog struct {
	ctx context.Context
}

func NewBinlog(ctx context.Context) Binlog {
	return Binlog{
		ctx: ctx,
	}
}

func (projector Binlog) View(viewer viewers.Channel) {
	var decoded *replication.BinlogSyncer
	utils.Decode(decoded)

	// Yes, we just transported an encrypted pointer over a queue and are about
	// to use what it is pointing to. Encryption is a first-class citizen.
	streamer, _ := decoded.StartSyncGTID(&mysql.MysqlGTIDSet{})
	viewer.RequestChannel(projector.generate(streamer))
}

func (projector Binlog) generate(streamer *replication.BinlogStreamer) chan spdg.Type {
	out := make(chan spdg.Type)

	go func() {
		defer close(out)

		for {
			out <- projector.cycle(streamer)
		}
	}()

	return out
}

func (projector Binlog) cycle(streamer *replication.BinlogStreamer) spdg.Type {
	// Upgrade the context with a timeout.
	ctx, cancel := context.WithTimeout(projector.ctx, 2*time.Second)
	ev, err := streamer.GetEvent(ctx)
	cancel()

	if err == context.DeadlineExceeded {
		// Meet timeout.
		return nil
	}

	encoded, _ := utils.Encode(ev)
	return utils.Wrap(projector.ctx, encoded)
}
