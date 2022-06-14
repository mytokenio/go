package metrics

import (
	"time"
)

// reference to https://github.com/rcrowley/go-metrics
// internal/kv reference to gokit

var (
	BatchInterval = time.Second * 5
)

// abstraction for metrics backend
type Metrics interface {
	//backend name
	String() string
	Counter(id string) Counter
	//Gauge(id string) Gauge
	Close() error
}

type Counter interface {
	Incr(delta int64)
	Decr(delta int64)
	Value() int64
	With(pair ...string) Counter
}

//type Gauge interface {
//	Set(value int64)
//	Value() int64
//	With(pair ...string) Gauge
//}

// 兼容性处理
func Count(param ... interface{}) {
	return
}
func Alarm(param ... interface{}) {
	return
}
func Close() {
	return
}
func Gauge(param ... interface{}) {
	return
}

