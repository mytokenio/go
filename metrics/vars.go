package metrics

import (
	"sync"
	"time"
)

var (
	mutex             sync.Mutex
	countMap          map[string]int64
	gaugeIntMap       map[string]int64
	gaugeStrMap       map[string]string
	globalServiceInfo serviceInfo
	globalKafka       kafkaInfo
	exitChan          chan struct{}
)

var (
	cronInterval = 10 * time.Second // 定时任务的间隔时间
)

var isDefaultKey = map[string]bool{
	"job_id":       true,
	"service_name": true,
	"status":       true,
	"env_type":     true,
	"start_time":   true,
	"stop_time":    true,
	"heart_time":   true,
	"exit_code":    true,
	"host":         true,
	"process_id":   true,
	"memory":       true,
	"load":         true,
	"net_in":       true,
	"net_out":      true,
}

// ---------------------------------------------------------------------------------------------------------------------

const (
	STATUS_ERROR     = 1 // 状态异常
	STATUS_OK        = 0 // 状态正常
	EXIT_CODE_UNEXIT = 0 // 服务未退出
	EXIT_CODE_OK     = 1 // 服务正常退出
	EXIT_CODE_ERROR  = 2 // 服务异常退出
	EXIT_CODE_KILL   = 3 // 服务被杀
)

const (
	ENV_DEV  = "dev"
	ENV_BETA = "beta"
	ENV_PRO  = "product"

	ENV_TYPE_DEV  = 0
	ENV_TYPE_BETA = 1
	ENV_TYPE_PRO  = 4

	ENV_JOB_ID       = "JOB_ID"
	ENV_SERVICE_NAME = "SERVICE_NAME"
	ENV_ENV_TYPE     = "ENV"

	DEF_SERVICE_NAME = "undefined"
)

const (
	default_map_caps           = 20
	default_producer_msg_caps  = 20
	default_roport_state_topic = "state_monitor_center"
	default_roport_alarm_topic = "alarm_monitor_center"
	dev_default_kafka_brokers  = "dev.kafka.mytoken.org:9092"
	beta_default_kafka_brokers = "172.17.1.63:9092"
	pro_default_kafka_brokers  = "172.16.0.131:9092,172.16.0.132:9092,172.16.0.133:9092"
)

const (
	state_producter_msg_key = "state_monitor_center"
	alarm_producter_msg_key = "alarm_monitor_center"
)
