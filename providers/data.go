package providers

import (
	"gitlab.itoodev.com/wrkgrp/binlog/tester"
	"gitlab.itoodev.com/wrkgrp/spdg"
)

type Data interface {
	tester.ErrorReporter
	RequestChannel(spdg.Type)
}
