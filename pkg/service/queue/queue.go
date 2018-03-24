package service

import (
	"container/heap"
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/util"
)

var once sync.Once
var s *Service

type Service struct {
	sync.RWMutex
	lastJourney *domain.Journey
	queue       util.PriorityQueue
}

type Interface interface {
	Push(j domain.Journey)
	Pop() domain.Journey
	Len() int
}

// New create a Service.
func Instance() *Service {
	once.Do(func() {
		s = &Service{}
		s.queue = make(util.PriorityQueue, 0)
		heap.Init(&s.queue)

	})
	return s
}

func (s *Service) Len() int {
	return s.queue.Len()
}

func (s *Service) Push(j domain.Journey) {
	s.Lock()
	defer s.Unlock()
	item := &util.Item{Value: j, Priority: j.ID}
	heap.Push(&s.queue, item)
}

func (s *Service) Pop() domain.Journey {
	s.Lock()
	defer s.Unlock()
	item := heap.Pop(&s.queue).(*util.Item)
	return item.Value.(domain.Journey)
}
