# Metrics

## Usage

```go
//import metrics backend service
import "github.com/mytokenio/go/metrics"

// defer to close metrics
defer metrics.Close()

metrics.Count("key", 1)
metrics.Gauge("key", 1)
metrics.Gauge("key", "string_value")
