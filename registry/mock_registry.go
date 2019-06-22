package registry

import (
	"sync"
	"fmt"
)

type mockRegistry struct {
	sync.RWMutex
	m map[string]*Service
}

func newMockRegistry() Registry {
	return &mockRegistry{
		m: make(map[string]*Service),
	}
}

func (r *mockRegistry) UnRegister(s *Service) error {
	r.Lock()
	delete(r.m, s.Name)
	r.Unlock()
	return nil
}

func (r *mockRegistry) Register(s *Service) error {
	r.Lock()
	r.m[s.Name] = s
	r.Unlock()
	return nil
}

func (r *mockRegistry) GetService(name string) ([]*Service, error) {
	r.RLock()
	defer r.RUnlock()
	if s, ok := r.m[name]; ok {
		return []*Service{s}, nil
	}
	return nil, fmt.Errorf("service %s not exist", name)
}

func (r *mockRegistry) String() string {
	return "mock"
}
