package service

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
)

var once sync.Once
var s *Service

type Service struct {
	sync.RWMutex
	lastJourney *domain.Journey
}

type Interface interface {
	GetLast() *domain.Journey
	SetLast(j *domain.Journey)
	GetNextID() int64
}

// New create a Service.
func Instance() *Service {
	once.Do(func() {
		s = &Service{}
		s.lastJourney = &domain.Journey{ID: 0}

	})
	return s
}

func (s *Service) GetLast() *domain.Journey {
	s.RLock()
	defer s.RUnlock()
	return s.lastJourney
}

func (s *Service) SetLast(j *domain.Journey) {
	s.Lock()
	defer s.Unlock()
	s.lastJourney = j

}

func (s *Service) GetNextID() int64 {
	return s.GetLast().ID + 1
}
