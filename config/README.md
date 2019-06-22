## Config

### Usage

import package

```
import (
    "github.com/mytokenio/go/config"
)
```

#### Service Config

custom struct 

```
type MyConfig struct {...}
```

bind config like this:
```
mc := &MyConfig{}
c := config.GetConfig()
c.BindJSON(mc)
// or
c.BindTOML(mc)
```

### Watch Change

default watch interval 5 seconds,
```
c.Watch(func() error {
    err := c.BindTOML(mc)
    if err != nil {
        log.Errorf("config bind error %s", err)
        return err
    }

    log.Infof("service config changed %v", mc)
    return nil
})
```

you can pass second parameter to control interval

`c.Watch(callback, 10 * time.Second)`


### File Driver

the default config driver, default file name `config.toml`
```
c := config.NewConfig()
c.BindTOML(...)
// or
c.Watch(...)
```


### Http Driver

if env `CONFIG_SERVER` not empty, the default config driver would be http driver

```
c := config.NewConfig(config.Service("your-service-name"))
```

or use shortcut:

```
c := config.NewHttpConfig("your-service-name")
```

### UI (for http driver)

moved to [mytokenio/config-manager](https://github.com/mytokenio/config-manager)

### Other

TODO


