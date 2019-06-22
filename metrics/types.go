package metrics

import (
	"sync"

	"github.com/Shopify/sarama"
)

type ReportStatePkg struct {
	JobID       int64  `json:"job_id"`
	ServiceName string `json:"service_name"`
	Status      int    `json:"status"`
	EnvType     int    `json:"env_type"`
	StartTime   int64  `json:"start_time"`
	StopTime    int64  `json:"stop_time"`
	HeartTime   int64  `json:"heart_time"`
	ExitCode    int    `json:"exit_code"`
	Host        string `json:"host"`
	ProcessID   int    `json:"process_id"`
	Memory      int    `json:"memory"`
	Load        int    `json:"load"`
	NetIn       int64  `json:"net_in"`
	NetOut      int64  `json:"net_out"`
	Extend      string `json:"extend"`
}

type ReportAlarmPkg struct {
	JobID       int64  `json:"job_id"`
	ServiceName string `json:"service_name"`
	Content     string `json:"content"`
	HeartTime   int64  `json:"heart_time"`
}

type serviceInfo struct {
	jobID       int64
	serviceName string
	envType     int
	host        string
	processID   int
}

type kafkaInfo struct {
	mutex                  sync.Mutex
	isInitialized          bool
	producer               sarama.AsyncProducer
	brokers                []string
	reportStateTopic       string
	reportAlarmTopic       string
	chanStateProducerValue chan string
	chanAlarmProducerValue chan string
}
