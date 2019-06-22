package web

import (
	"net/http"
	"sync"
	"os"
	"os/signal"
	"syscall"
	"net"
	"github.com/mytokenio/go/registry"
	"github.com/mytokenio/go/log"
	"time"
	"strconv"
)

type Service interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Run() error
}

func NewService(opts ...Option) Service {
	return newService(opts...)
}

type service struct {
	sync.Mutex
	opts    Options
	mux     *http.ServeMux
	regSrv *registry.Service
	running bool
	exit    chan chan error
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)
	s := &service{
		opts: options,
		mux:  http.NewServeMux(),
	}
	s.regSrv = s.toRegistryService()
	return s
}

func (s *service) toRegistryService() *registry.Service {
	host, port, _ := net.SplitHostPort(s.opts.Address)
	if s.opts.Advertise != "" {
		host, port, _ = net.SplitHostPort(s.opts.Advertise)
	}

	portInt, _ := strconv.Atoi(port)
	return &registry.Service{
		Name:    s.opts.Name,
		Version: s.opts.Version,
		Nodes: []registry.Node{{
			Name: host,
			Host: host,
			Port: portInt,
		}},
	}
}

func (s *service) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	if err := s.register(); err != nil {
		return err
	}

	// start reg loop
	ex := make(chan bool)
	go s.run(ex)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case sig := <-sc:
		log.Infof("Received signal %s\n", sig)
	}

	close(ex)

	if err := s.unregister(); err != nil {
		return err
	}

	return s.stop()
}

func (s *service) listen(network, addr string) (net.Listener, error) {
	var l net.Listener
	var err error

	l, err = net.Listen(network, addr)

	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *service) run(exit chan bool) {
	t := time.NewTicker(time.Minute)

	for {
		select {
		case <-t.C:
			s.register()
		case <-exit:
			t.Stop()
			return
		}
	}
}

func (s *service) register() error {
	return registry.Register(s.regSrv)
}

func (s *service) unregister() error {
	return registry.UnRegister(s.regSrv)
}

func (s *service) start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	l, err := s.listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.opts.Address = l.Addr().String()

	var handler http.Handler
	if s.opts.Handler != nil {
		handler = s.opts.Handler
	} else {
		handler = s.mux
	}

	httpSrv := &http.Server{
		Handler: handler,
	}

	go httpSrv.Serve(l)

	s.exit = make(chan chan error, 1)
	s.running = true

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	log.Infof("listening on %v", l.Addr().String())
	return nil
}

func (s *service) stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false

	log.Infof("stopping")
	return <-ch
}
