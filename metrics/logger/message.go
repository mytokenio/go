package logger

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mytokenio/go/log"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Message struct {
	Index string                 `json:"@index"`
	Type  string                 `json:"@type"`
	Time  string                 `json:"datetime"`
	Data  map[string]interface{} `json:"data"`
}

func newCounterMessage(name string, n float64, fields map[string]string) []byte {
	msg := Message{
		Index: fmt.Sprintf("metrics-counter-%s", name),
		Type:  name,
		Time:  time.Now().UTC().Format(TimeFormat),
		Data: map[string]interface{}{
			"value": n,
		},
	}
	if len(fields) > 0 {
		msg.Data["fields"] = fields
	}
	b, _ := json.Marshal(msg)
	log.Infof("metrics counter: %s", b)
	return b
}

func newGaugeMessage(name string, n float64, fields map[string]string) []byte {
	msg := Message{
		Index: fmt.Sprintf("metrics-gauge-%s", name),
		Type:  name,
		Time:  time.Now().UTC().Format(TimeFormat),
		Data: map[string]interface{}{
			"value": n,
		},
	}
	if len(fields) > 0 {
		msg.Data["fields"] = fields
	}
	b, _ := json.Marshal(msg)
	log.Infof("metrics gauge: %s", b)
	return b
}
