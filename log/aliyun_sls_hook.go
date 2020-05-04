package log

import (
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"
	"os"
)

const LOGSTORE string = "golog"

type AliyunSLSHook struct {
	producerInstance *producer.Producer
	logstore string
}

func NewAliyunSLSHook(endpoint string, key string, secret string, debug bool) *AliyunSLSHook {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = endpoint
	producerConfig.AccessKeyID = key
	producerConfig.AccessKeySecret = secret
	if !debug {
		producerConfig.AllowLogLevel = "warn"
	}

	producerInstance := producer.InitProducer(producerConfig)
	producerInstance.Start()

	logstore := os.Getenv(envLogStore)
	if logstore == "" {
		logstore = LOGSTORE
	}
	return &AliyunSLSHook{producerInstance, logstore}
}

func (hook *AliyunSLSHook) Fire(entry *logrus.Entry) error {
	_type := entry.Data[typeField]
	var typ string
	if _type == nil {
		typ = "default"
	} else {
		typ = _type.(string)
	}

	delete(entry.Data, typeField)

	const PROJECT string = "mytoken-open-api-log"

	go func() {
		// GenerateLog  is producer's function for generating SLS format logs
		// GenerateLog has low performance, and native Log interface is the best choice for high performance.
		extra_out, _ := json.Marshal(extra)
		data_out, _ := json.Marshal(entry.Data)
		log := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{
			"Level":    entry.Level.String(),
			"Time":     entry.Time.UTC().Format(TimeFormat),
			"UniqueID": uniqueId,
			"Host":     hostname,
			"Extra":    string(extra_out),
			"Message":  entry.Message,
			"Context":  string(data_out)})
		err := hook.producerInstance.SendLog(PROJECT, hook.logstore, typ, hostname, log)
		if err != nil {
			// TODO
		}
	}()

	return nil
}

func (hook *AliyunSLSHook) Levels() []logrus.Level {
	//without debug level
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}
