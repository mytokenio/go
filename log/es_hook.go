package log

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

var hostname, _ = os.Hostname()

type EsLogHook struct {
	Server string
	Conn   net.Conn
}

type Message struct {
	Index    string  `json:"@index"`
	Type     string  `json:"@type"`
	Level    string  `json:"level"`
	Time     string  `json:"datetime"`
	UniqueID string  `json:"unique_id"`
	Info     logInfo `json:"info"`
}

type logInfo struct {
	Host    string                 `json:"host"`
	Extra   map[string]interface{} `json:"extra"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context"`
}

func NewEsLogHook(server string) *EsLogHook {
	conn, err := net.Dial("udp", server)
	if err != nil {
		log.Errorf("failed dial log server %v", server)
		os.Exit(1)
	}
	return &EsLogHook{server, conn}
}

func (hook *EsLogHook) Fire(entry *logrus.Entry) error {
	_type := entry.Data[typeField]
	var typ string
	if _type == nil {
		typ = "default"
	} else {
		typ = _type.(string)
	}

	delete(entry.Data, typeField)

	msg := Message{
		Index:    fmt.Sprintf("golog-%s-%s", typ, entry.Time.UTC().Format(DateFormat)),
		Type:     typ,
		Level:    entry.Level.String(),
		Time:     entry.Time.UTC().Format(TimeFormat),
		UniqueID: uniqueId,
		Info: logInfo{
			Host:    hostname,
			Extra:   extra,
			Message: entry.Message,
			Context: entry.Data,
		},
	}
	b, _ := json.Marshal(msg)
	_, err := hook.Conn.Write(b)
	if err != nil {
		log.Errorf("log server write error: %s", err)

		//retry TODO
	}
	return nil
}

func (hook *EsLogHook) Levels() []logrus.Level {
	//without debug level
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}
