package service

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
	metricsservice "github.com/pmoncadaisla/go-journey/pkg/service/metrics"
)

var once sync.Once
var s *Service

type Service struct {
	sync.RWMutex
	lastJourney            *domain.Journey
	highestReceivedJourney *domain.Journey
	metricsService         metricsservice.Interface
}

type Interface interface {
	GetLast() *domain.Journey
	SetLast(j *domain.Journey)
	GetNextID() int
	GetReceivedHighestID() int
	Receive(j *domain.Journey)
}

// New create a Service.
func Instance() *Service {
	once.Do(func() {
		s = &Service{}
		s.lastJourney = &domain.Journey{ID: 0}
		s.metricsService = metricsservice.Instance()
		s.highestReceivedJourney = &domain.Journey{ID: 0}

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
	return

}

func (s *Service) Receive(j *domain.Journey) {
	s.Lock()
	defer s.Unlock()
	s.lastJourney = j
	if j.ID > s.highestReceivedJourney.ID {
		s.highestReceivedJourney = j
		s.metricsService.CounterInc(metrics.JOURNEYS_HIGHEST_RECEIVED_ID.String())
	}
	return
}

func (s *Service) GetNextID() int {
	return s.GetLast().ID + 1
}

func (s *Service) GetReceivedHighestID() int {
	return s.highestReceivedJourney.ID
}
