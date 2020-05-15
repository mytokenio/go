package log

import (
	jsoniter "github.com/json-iterator/go"
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
	defMaxRolls    = 1
	envLogToFile   = "LOG_TO_FILE"          // write log to file
	envLogServer   = "LOG_SERVER"           // log server
	envLogEndpoint = "LOG_ENDPOINT"         // Aliyun SLS endpoint
	envLogStore = "LOG_STORE"         // Aliyun SLS logstore
	envLogKey      = "ALIYUN_ACCESS_KEY"    // Aliyun access key
	envLogSecret   = "ALIYUN_ACCESS_SECRET" // Aliyun access secret
	envJobID       = "JOB_ID"               // job id
	envServiceName = "SERVICE_NAME"         // service name
	envEnv         = "ENV"                  // env
	envDev         = "dev"                  // dev
	envBeta        = "beta"                 // beta
	envTest        = "test"                 // test
	envPro         = "product"              // product
)
