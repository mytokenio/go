package driver

import (
	"os"
	"io/ioutil"
	"errors"
	"path/filepath"
	"strings"
)

type fileDriver struct {
	path string
}

func NewFileDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &fileDriver{
		path: options.Path,
	}
}

func (d *fileDriver) List() ([]*Value, error) {
	var vals []*Value
	return vals, nil
}

func (d *fileDriver) Get(key string) (*Value, error) {
	path := key
	if d.path != "" {
		path = d.path
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	v := NewValue(key, b)
	v.Format = strings.TrimLeft(filepath.Ext(path), ".")
	return v, nil
}

func (d *fileDriver) Set(value *Value) error {
	return errors.New("TODO")
}

func (d *fileDriver) String() string {
	return "file"
}