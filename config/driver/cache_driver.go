package driver

import (
	"github.com/mytokenio/go/log"
	"sync"
	"time"
)

//for TTL check
type cacheValue struct {
	K         string
	V         *Value
	Timestamp int64
}

type cacheDriver struct {
	sync.RWMutex
	SubDriver Driver
	TTL       time.Duration
	Data      map[string]cacheValue
}

func NewCacheDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	//ttl minimum 5 seconds
	minTTL := time.Second * 5
	if options.TTL > 0 {
		if options.TTL > minTTL {
			minTTL = options.TTL
		} else {
			log.Warnf("minimum ttl 5 seconds")
		}
	}

	if options.SubDriver == nil {
		options.SubDriver = DefaultDriver
	}

	return &cacheDriver{
		TTL:       minTTL,
		SubDriver: options.SubDriver,
		Data:      map[string]cacheValue{},
	}
}

//do not cache for list
func (d *cacheDriver) List() ([]*Value, error) {
	return d.SubDriver.List()
}

func (d *cacheDriver) Get(key string) (*Value, error) {
	if cache := d.cacheGet(key); cache != nil {
		return cache, nil
	}

	v, err := d.SubDriver.Get(key)
	if err != nil {
		return nil, err
	}

	d.cacheSet(v)

	return v, nil
}

func (d *cacheDriver) Set(value *Value) error {
	d.Lock()
	delete(d.Data, value.K)
	d.Unlock()

	return d.SubDriver.Set(value)
}

func (d *cacheDriver) cacheGet(key string) *Value {
	d.RLock()
	v, ok := d.Data[key]
	d.RUnlock()

	if !ok {
		return nil
	}
	//expired ?
	if v.Timestamp+int64(d.TTL.Seconds()) < time.Now().Unix() {
		d.Lock()
		delete(d.Data, key)
		d.Unlock()

		return nil
	}

	return v.V
}

func (d *cacheDriver) cacheSet(value *Value) error {
	d.Lock()
	defer d.Unlock()

	d.Data[value.K] = cacheValue{
		K:         value.K,
		V:         value,
		Timestamp: time.Now().Unix(),
	}
	return nil
}

func (d *cacheDriver) String() string {
	return "cache"
}
