package viewers

import (
	"bytes"
	"context"

	"gitlab.itoodev.com/wrkgrp/binlog/utils"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type Channel interface {
	Watch()
	RequestChannel(chan spdg.Type)
}

type Binlog struct {
	in  chan spdg.Type
	ctx context.Context
}

func NewBinlog(ctx context.Context) Channel {
	return Binlog{
		in:  make(chan spdg.Type),
		ctx: ctx,
	}
}

func (channel Binlog) Watch() {
	switch channel.in {
	case nil:
		channel.recover(spdg.NO)
	}

	channel.cycle()
}

func (channel Binlog) RequestChannel(ch chan spdg.Type) {
	channel.in = ch
}

func (channel Binlog) cycle() {
	for frame := range channel.in {
		t, v, l := utils.Unwrap(channel.ctx)
		frame.Peek(v)
		utils.Pushwrap(t, v, l, bytes.Buffer{})
	}
}

func (channel Binlog) recover(status spdg.Status) {

}
