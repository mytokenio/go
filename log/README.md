## Log


based on logrus, support sync to MyToken log server

#### Config

sdk read configs from environment

`LOG_LEVEL` log level, `info` `error` `warn` `fatal` `panic` `debug` (default)

`LOG_SERVER` log udp server address, format `ip:port`

### Usage

import package

```
import (
    "github.com/mytokenio/go/log"
    "github.com/sirupsen/logrus" //if need log type or context data (fields)
)
```

usage same as logrus, print to stdout, like:

```
log.Infof("test format info log %s", xxx)
```

support log type and context data:

```
log.Type("your_type").Infof("test format info log %s", xxx)
```

type and context data (fields):
```
fields := logrus.Fields{
    "k1": "v1",
    "k2": "v2",
}
log.Type("test").WithFields(fields).Info("log with fields and type")
```

or fields as pairs, above code equal to:
```
log.With("k1", "v1", "k2", "v2", "type", "test").Info("log with fields and type")
```

log type would be a part of elasticsearch index name, format `golog-{type}-{date}`