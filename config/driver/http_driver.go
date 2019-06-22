package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mytokenio/go/log"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	EnvConfigServer = "CONFIG_SERVER"
	Env             = "ENV"
	ServiceName     = "SERVICE_NAME"
	JobID           = "JOB_ID"
	CodeSuccess     = 0
)

type httpDriver struct {
	Host       string
	HttpClient *http.Client
	sync.Mutex
}

type Request struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Comment   string `json:"comment"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

type Response struct {
	Code      int             `json:"code"`
	Msg       string          `json:"msg"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// {"code":-1,"msg":"error","data":{"error":"sql: no rows in result set","reason":"Query data error"},"timestsamp":1536223015}
type DataError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type DataConfig struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Comment   string `json:"comment"`
	UpdatedBy string `json:"updated_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c DataConfig) toMetadata() map[string]string {
	return map[string]string{
		"comment":    c.Comment,
		"updated_by": c.UpdatedBy,
		"updated_at": c.UpdatedAt,
		"created_at": c.CreatedAt,
	}
}

func NewHttpDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	timeout := time.Second * 3
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	if options.Host == "" {
		options.Host = os.Getenv(EnvConfigServer)
	}

	return &httpDriver{
		Host:       options.Host,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (d *httpDriver) List() ([]*Value, error) {
	var vals []*Value

	uri := "/v1/config/item"
	b, err := d.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	cs := &[]*DataConfig{}
	err = json.Unmarshal(b, cs)
	if err != nil {
		log.Errorf("json unmarshal error %s", err)
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}

	for _, c := range *cs {
		v := NewValue(c.Key, []byte(c.Value))
		v.Metadata = c.toMetadata()
		vals = append(vals, v)
	}

	return vals, nil
}

func (d *httpDriver) Get(key string) (*Value, error) {
	uri := fmt.Sprintf("/v1/config/item/%s", key)
	b, err := d.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	c := &DataConfig{}
	err = json.Unmarshal(b, c)
	if err != nil {
		log.Errorf("json unmarshal error %s", err)
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}

	v := NewValue(c.Key, []byte(c.Value))
	return v, nil
}

func (d *httpDriver) Set(value *Value) error {
	uri := "/v1/config/item"
	method := "POST" //create
	req := Request{
		Key:   value.K,
		Value: value.String(),
	}

	//check for update
	existValue, err := d.Get(value.K)
	if existValue != nil {
		uri = fmt.Sprintf("/v1/config/item/%s", value.K)
		method = "PATCH"
		req.UpdatedBy = "sdk"
	} else {
		req.CreatedBy = "sdk"
	}

	reqBytes, _ := json.Marshal(req)
	_, err = d.request(method, uri, reqBytes)
	if err != nil {
		return fmt.Errorf("post failed %s", err)
	}
	return nil
}

func (d *httpDriver) request(method string, path string, data []byte) (json.RawMessage, error) {
	if d.Host == "" {
		return nil, errors.New("config server host empty")
	}

	url := d.Host + path
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := d.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d, body %s", resp.StatusCode, respBody)
	}

	rsp := &Response{}
	err = json.Unmarshal(respBody, rsp)
	if err != nil {
		return nil, fmt.Errorf("response json unmarshal error %s, body %s", err, respBody)
	}
	if rsp.Code != CodeSuccess {
		dataErr := &DataError{}
		json.Unmarshal(rsp.Data, dataErr)
		return nil, fmt.Errorf("response msg: %s, url: %s, body: %s, error: %s, reason: %s", rsp.Msg, url, body, dataErr.Error, dataErr.Reason)
	}

	return rsp.Data, err
}

func (d *httpDriver) String() string {
	return "http"
}
