package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mytokenio/go/config/driver"
)

func init() {
	var host string

	switch strings.ToLower(os.Getenv(driver.Env)) {
	case ENV_BETA:
		host = BETA_HOST_MONITOR_CENTER
	case ENV_PRO:
		host = PRO_HOST_MONITOR_CENTER
	default:
		host = DEV_HOST_MONITOR_CENTER
	}

	jobId, _ = strconv.ParseInt(os.Getenv(driver.JobID), 10, 64)

	currentConfig = NewConfig(
		Service(os.Getenv(driver.ServiceName)),
		TTL(60*time.Second),
		Driver(
			driver.NewHttpDriver(
				driver.Host(host),
				driver.Timeout(3*time.Second),
			),
		),
	)
}

func Init(assort int, service, host, path string) error {
	if service == "" {
		service = os.Getenv(driver.ServiceName)
	}
	if path == "" {
		path = driver.DefaultConfigFile
	}
	switch assort {
	case CFG_FROM_MONITOR_CENTER:
		currentConfig = NewConfig(
			Service(service),
			TTL(60*time.Second),
			Driver(
				driver.NewHttpDriver(
					driver.Host(host),
					driver.Timeout(3*time.Second),
				),
			),
		)
	case CFG_FROM_LOCAL_FILE:
		currentConfig = NewFileConfig(path)
	default:
		return errors.New("params error")
	}

	return nil
}
