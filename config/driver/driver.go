package driver

import (
	"time"
)

type Driver interface {
	List() ([]*Value, error)
	Get(string) (*Value, error)
	Set(*Value) error
	String() string
}

var (
	DefaultConfigFile = "./config.toml"
	DefaultDriver     = NewFileDriver(Path(DefaultConfigFile))
)

type Option func(*Options)

type Options struct {
	Path      string //for file driver
	Host      string //for http driver
	Timeout   time.Duration
	SubDriver Driver //for cache driver
	TTL       time.Duration
}

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Path(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}

func SubDriver(reg Driver) Option {
	return func(o *Options) {
		o.SubDriver = reg
	}
}
