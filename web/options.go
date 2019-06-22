package web

import (
	"net/http"
	"github.com/mytokenio/go/registry"
	"google.golang.org/grpc/metadata"
)

type Option func(o *Options)

type Options struct {
	Name      string
	Version   string
	Metadata  metadata.MD
	Address   string
	Advertise string
	Handler   http.Handler
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:    "web",
		Address: ":0",
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func Metadata(md metadata.MD) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

func Handler(h http.Handler) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		registry.DefaultRegistry = r
	}
}
