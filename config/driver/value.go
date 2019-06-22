package driver

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"errors"
	"strings"
	"gopkg.in/yaml.v2"
)

type Value struct {
	K         string
	V         []byte
	Format    string //json, yml, yaml, toml, ...
	CheckSum  string
	Metadata  map[string]string
}

func NewValue(k string, v []byte) *Value {
	value := &Value{
		K: k,
		V: v,
	}
	value.CheckSum = value.genCheckSum()
	return value
}

func (v *Value) genCheckSum() string {
	s := sha512.New512_256()
	s.Write([]byte(v.K))
	s.Write(v.V)
	return hex.EncodeToString(s.Sum(nil))
}

func (v *Value) Bytes() []byte {
	return v.V
}

func (v *Value) String() string {
	return string(v.V)
}

func (v *Value) Bind(obj interface{}) error {
	if v.Format == "" {
		return errors.New("format error")
	}

	switch strings.ToLower(v.Format) {
	case "json":
		return v.BindJSON(obj)
	case "toml":
		return v.BindTOML(obj)
	case "yml", "yaml":
		return v.BindYAML(obj)
	default:
		return errors.New(v.Format + " format not supported")
	}
}

func (v *Value) BindJSON(obj interface{}) error {
	e := json.Unmarshal(v.Bytes(), obj)
	if e != nil {
		return fmt.Errorf("json unmarshal error %s", e)
	}
	return nil
}

func (v *Value) BindTOML(obj interface{}) error {
	e := toml.Unmarshal(v.Bytes(), obj)
	if e != nil {
		return fmt.Errorf("toml unmarshal error %s", e)
	}
	return nil
}

func (v *Value) BindYAML(obj interface{}) error {
	e := yaml.Unmarshal(v.Bytes(), obj)
	if e != nil {
		return fmt.Errorf("yaml unmarshal error %s", e)
	}
	return nil
}