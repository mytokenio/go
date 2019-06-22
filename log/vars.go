package log

import (
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var (
	log      *logrus.Logger
	json     jsoniter.API
	uniqueId string
	extra    map[string]interface{}
)

// ---------------------------------------------------------------------------------------------------------------------

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)

const (
	typeField  = "log_type"
	defLogPath = "/data/logs/"
	defLogFile = "default.log"
)

const (
	defMaxRolls    = 3
	envLogToFile   = "LOG_TO_FILE"  // write log to file
	envLogServer   = "LOG_SERVER"   // log server
	envJobID       = "JOB_ID"       // job id
	envServiceName = "SERVICE_NAME" // service name
	envEnv         = "ENV"          // env
	envDev         = "dev"          // dev
	envBeta        = "beta"         // beta
	envTest        = "test"         // test
	envPro         = "product"      // product
)
