# Metrics

metrics interface and implementation based on mytoken logger server

## Interface

```go
// abstraction for metrics backend
type Metrics interface {
    //backend name
    String() string
    Counter(id string) Counter
    Gauge(id string) Gauge
    Close() error
}

type Counter interface {
    Incr(delta int64)
    Decr(delta int64)
    Value() int64
    With(pair ...string) Counter
}

type Gauge interface {
    Set(value int64)
    Value() int64
    With(pair ...string) Gauge
}

```

## Usage

```go
//import metrics backend service
import "github.com/mytokenio/go/metrics/logger"

//init with metrics namespace and logger server address
m := logger.New("test", "127.0.0.1:12333")
defer m.Close()

//create counter/gauge instance
c := m.Counter("counter")
g := m.Gauge("test-gauge")

//call counter/gauge
c.Incr(10)
log.Infof("counter value %d", c.Value())

g.Set(1234)
log.Infof("gauge value %d", g.Value())

// with kv pair
c.With("k1": "v1", "k2", "v2", ...).Incr(123)

g.With("k1": "v1", "k2", "v2", ...).Set(123)
```

## Sync Metrics

Default

- [MyToken Logger Server](https://github.com/mytokenio/go/tree/master/metrics/logger)

TODO

- Prometheus

You can custom backend service by implementing `Metrics` interface, to sync metrics data to anywhere.