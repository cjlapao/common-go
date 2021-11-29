package database

import (
	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/log"
)

var ctx = executionctx.GetContext()
var logger = log.Get()

type DatabaseCtx struct {
	DatabaseName     string
	ConnectionString string
}
