## MyToken Go SDK

## Log

log based on `logrus`, support sync to mytoken log server, default write log to /data/logs/

```
go get github.com/mytokenio/go/log
```

```go
import (
    "github.com/mytokenio/go/log"
)

log.Info("xxx log")
log.WithField("kkk", "vvv", "custom_type").Info("log with type & kv data")
```

[more detail](https://github.com/mytokenio/go/tree/master/log)

## Config


```
go get github.com/mytokenio/go/config
```

```go
import (
    "github.com/mytokenio/go/config"
)

mc := &MyConfig{}
c := config.GetConfig()

// bind to struct
c.BindTOML(mc)

// or, watch change
c.Watch(func() error {
    err := c.BindTOML(mc)
    // TODO
    return nil
})
```

[more detail](https://github.com/mytokenio/go/tree/master/config)

## Metrics

```go
// import metrics backend service
import "github.com/mytokenio/go/metrics"

// defer to close metrics
defer metrics.Close()

// count/gauge usage
metrics.Count("key", 1)
metrics.Gauge("key", 1)
metrics.Gauge("key", "string_value")

```

[more detail](https://github.com/mytokenio/go/tree/master/metrics)

## TODO

service registry/broker/health

rpc server/client/protocol

rate limiter, circuit breaker, hytrix

tracing, zipkin or jaeger

cli

queue

demo project

...


