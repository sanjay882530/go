package test

import (
	"github.com/sirupsen/logrus"
	"github.com/hcnet/go/support/log"
)

var testLogger *log.Entry

func init() {
	testLogger = log.New()
	testLogger.DisableColors()
	testLogger.SetLevel(logrus.DebugLevel)
}
