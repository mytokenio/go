package driver

import (
	"fmt"
	"github.com/mytokenio/go/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
)

var mockResponse = map[string]map[string]string {
	"GET": {
		"/v1/config/item/test.key" : `{"code":0,"msg":"success","data":{"comment":"","created_at":"2018-09-05T18:00:11+0800","created_by":"sdk","deleted_at":"1970-01-01T00:00:00+0800","deleted_by":"","id":4,"key":"test.key","state":0,"updated_at":"1970-01-01T00:00:00+0800","updated_by":"","value":"test value"},"timestsamp":1536235332}`,
		"/v1/config/item" : `{"code":0,"msg":"success","data":[{"comment":"","created_at":"2018-09-05T16:53:38+0800","created_by":"sdk","deleted_at":"1970-01-01T00:00:00+0800","deleted_by":"","id":1,"key":"mt.service.user","state":0,"updated_at":"2018-09-06T20:03:46+0800","updated_by":"sdk","value":"dddd"},{"comment":"","created_at":"2018-09-05T18:00:11+0800","created_by":"sdk","deleted_at":"1970-01-01T00:00:00+0800","deleted_by":"","id":4,"key":"user","state":0,"updated_at":"1970-01-01T00:00:00+0800","updated_by":"","value":"dafda"},{"comment":"","created_at":"2018-09-06T17:15:57+0800","created_by":"sdk","deleted_at":"1970-01-01T00:00:00+0800","deleted_by":"","id":8,"key":"mt.service.config-manager","state":0,"updated_at":"2018-09-06T18:18:09+0800","updated_by":"sdk","value":"#just for demo test\r\nTitle = \"MyToken Config Center\""}],"timestsamp":1536235537}`,
		"/v1/config/item/not_exist" : `{"code":-1,"msg":"error","data":{"error":"sql: no rows in result set","reason":"Query data error"},"timestsamp":1536235976}`,
		"/v1/config/item/error_json" : `error`,
	},
	"POST": {
		"/v1/config/item": `{"code":0,"msg":"success","data":{"key":"test.key"},"timestsamp":1536235426}`,
	},
	"PATCH": {
		"/v1/config/item/test.key": `{"code":0,"msg":"success","data":{"key":"test.key"},"timestsamp":1536235426}`,
	},
}

func assert(t *testing.T, actual interface{}, expect interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName = strings.Replace(fileName, wd, "", 1)
	if actual != expect {
		t.Errorf("expect %v, got %v at (%v:%v)", expect, actual, fileName, line)
	}
}

func testNewHttpDriver(host string) Driver {
	return NewHttpDriver(Host(host))
}

func TestHttp(t *testing.T) {
	var server *httptest.Server

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rsp, ok := mockResponse[r.Method][r.RequestURI]; ok {
			log.Infof("mock http resp %s %s", r.RequestURI, rsp)
			fmt.Fprint(w, rsp)
		} else {
			fmt.Fprint(w, `error`)
		}
	}))
	defer server.Close()
	//os.Setenv("CONFIG_SERVER", "http://"+server.URL)

	var v *Value
	var err error

	log.Infof("server url: %s", server.URL)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", server.URL+"/v1/item", nil)
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Infof("body %s", respBody)

	//get
	d := testNewHttpDriver(server.URL)
	v, err = d.Get("test.key")
	assert(t, err, nil)
	assert(t, v.String(), "test value")

	//get error
	v, err = d.Get("not_exist")
	if !strings.Contains(err.Error(), "Query data error") {
		t.Errorf("get not_exist failed %s", err.Error())
	}

	//get error json
	v, err = d.Get("error_json")
	if !strings.Contains(err.Error(), "json unmarshal error") {
		t.Errorf("get error_json failed %s", err.Error())
	}

	//create
	err = d.Set(NewValue("test.key", []byte("xxx")))
	assert(t, err, nil)

	//update
	err = d.Set(NewValue("test.key", []byte("xxx")))
	assert(t, err, nil)

	//list
	ls, err := d.List()
	assert(t, err, nil)
	for _, v = range ls {
		if len(v.String()) == 0 {
			t.Errorf("list value empty %s", v.K)
		}
	}
}

