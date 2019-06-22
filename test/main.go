package main

//import (
//	"encoding/json"
//	"github.com/mytokenio/go/config"
//	"github.com/mytokenio/go/config/driver"
//	"github.com/mytokenio/go/log"
//	"github.com/mytokenio/go/metrics/logger"
//	"github.com/mytokenio/go/registry"
//	"google.golang.org/grpc/metadata"
//	"time"
//)
//
//func main() {
//	//testFileConfig()
//	//testHttpConfig()
//	testMetrics()
//}
//
//func testMetrics() {
//	m := logger.New("test", "127.0.0.1:12333")
//	defer m.Close()
//
//	c := m.Counter("counter")
//	g := m.Gauge("test-gauge")
//
//	go func() {
//		for i := 0; i < 1000; i++ {
//			c.Incr(10)
//			log.Infof("counter val %d", c.Value())
//			g.Set(int64(i))
//			log.Infof("gauge val %d", g.Value())
//			time.Sleep(time.Second)
//		}
//	}()
//
//	time.Sleep(time.Minute)
//}
//
//
//type MyConfig struct {
//	API string `toml:"api"`
//	DB struct {
//		Host     string `toml:"host"`
//		User     string `toml:"user"`
//		Password string `toml:"password"`
//		Name     string `toml:"name"`
//	} `toml:"db"`
//	LogServers []string `toml:"log_servers"`
//}
//
//func testFileConfig() {
//	mc := &MyConfig{}
//
//	c := config.NewConfig()
//	c.Watch(func() error {
//		err := c.BindTOML(mc)
//		if err != nil {
//			log.Errorf("config bind error %s", err)
//			return err
//		}
//
//		log.Infof("service config changed %v", mc)
//		return nil
//	})
//
//	//c := config.NewConfig(
//	//	config.TTL(10 * time.Second),
//	//	config.Driver(
//	//		driver.NewFileDriver(driver.Path("config.toml")),
//	//	),
//	//)
//	//c.Watch(func() error {
//	//	err := c.BindTOML(mc)
//	//	if err != nil {
//	//		log.Errorf("config bind error %s", err)
//	//		return err
//	//	}
//	//
//	//	log.Infof("service config changed %v", mc)
//	//	return nil
//	//})
//
//	log.Infof("MyConfig %v", mc)
//}
//
//func testHttpConfig() {
//	myConfigJson := `
//	{
//		"api": "http://api.mytokenapi.com",
//		"db": {
//			"host": "localhost",
//			"user": "root",
//			"password": "",
//			"name": "mytoken"
//		},
//		"log_servers": ["127.0.0.1:12333", "127.0.0.1:12334"]
//	}
//	`
//	mc := &MyConfig{}
//	c := config.NewConfig(
//		config.Service("user"),
//		config.TTL(time.Second*10), //cache ttl
//		config.Driver(
//			driver.NewHttpDriver(
//				driver.Host("http://127.0.0.1:8083"),
//				driver.Timeout(time.Second*3),
//			),
//		),
//	)
//
//	value := driver.NewValue("mt.service."+c.Service, []byte(myConfigJson))
//	c.Driver.Set(value)
//	c.BindTOML(mc)
//}
//
//func testService() {
//	s := registry.Service{
//		Name:     "test",
//		Metadata: metadata.Pairs("kk", "vv", "aa", "bb", "aa", "cc"),
//		Nodes: []registry.Node{
//			{"test", "127.0.0.1", 12345},
//		},
//	}
//	b, e := json.Marshal(s)
//	if e != nil {
//		log.Errorf("error %s", e)
//	} else {
//		log.Infof("service %s", b)
//	}
//}
