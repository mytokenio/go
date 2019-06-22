package main

import (
	"github.com/mytokenio/go/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.WithField("kk", "vv", "xxx").Infof("test debug log %s", "ddd")
	log.Infof("test info log %s", "ddd")
	log.Warnf("test warn log %s", "ddd")

	log.RefreshUniqueId()

	log.Errorf("test error log %s", "ddd")
	log.WithFields(logrus.Fields{"aa":"bb", "cc": logrus.Fields{"dd":"aaaaaaa"}}, "test").Infof("test obj")

	log.Type("test").Info("fdsafdas")

	fields := logrus.Fields{
		"aa": "bb",
		"cc": "dd",
	}
	log.Type("test").WithFields(fields).Info("log with fields and type")
}
