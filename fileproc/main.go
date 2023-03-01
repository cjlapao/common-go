package fileproc

import (
	log "github.com/cjlapao/common-go-logger"
)

var logger = log.Get()

// Process Process a file replacing the variables in it
func Process(content []byte, variables ...Variable) []byte {
	logger.Debug("Starting to process content to replace variables")
	return ReplaceAll(content, variables...)
}
