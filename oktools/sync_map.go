package oktools

import "sync"

var ()

type (
	SyncMap struct {
		Mutex sync.Mutex
		Cache map[string]interface{}
	}

	Aa struct {
	}
)

func (s *SyncMap) initMap() {
	if s.Cache == nil {
		s.Cache = map[string]interface{}{}
	}
}

func (s *SyncMap) Get(key string) (interface{}, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.initMap()

	v, ok := s.Cache[key]
	return v, ok
}

func (s *SyncMap) GetString(key string) (string, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.initMap()

	v, ok := s.Cache[key]
	if !ok {
		return "", false
	} else {
		return v.(string), ok
	}
}

func (s *SyncMap) GetBool(key string) (bool, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.initMap()

	v, ok := s.Cache[key]
	if !ok {
		return false, false
	} else {
		return v.(bool), true
	}
}

func (s *SyncMap) Put(key string, value interface{}) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.initMap()

	s.Cache[key] = value
}

func (s *SyncMap) Size() int {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	return len(s.Cache)
}
