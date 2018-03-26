package service

import (
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var once sync.Once
var s *Service

type Service struct {
	Counters map[string]*Counter
	Gauges   map[string]*Gauge
}

type Interface interface {
	CounterInc(counter string)
	GaugeInc(gauge string)
	GaugeDec(gauge string)
}

type Counter struct {
	sync.RWMutex
	Name  string
	Help  string
	Value int
	Vec   *prometheus.CounterVec
}

type Gauge struct {
	sync.RWMutex
	Name  string
	Help  string
	Value int
	Vec   *prometheus.GaugeVec
}

// Instance creates a Singleton Service.
func Instance() *Service {
	once.Do(func() {
		s = &Service{}
		s.Counters = make(map[string]*Counter)
		s.Gauges = make(map[string]*Gauge)

	})
	return s
}

func (s *Service) NewCounter(name, help string) {

	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		[]string{"hostname"},
	)
	c := &Counter{
		Name: name,
		Help: help,
		Vec:  vec,
	}

	prometheus.MustRegister(vec)

	s.Counters[name] = c
}

func (s *Service) NewGauge(name, help string) {

	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		[]string{"hostname"},
	)
	c := &Gauge{
		Name: name,
		Help: help,
		Vec:  vec,
	}

	prometheus.MustRegister(vec)

	s.Gauges[name] = c
}

func (s *Service) getLabels() prometheus.Labels {
	hostname, _ := os.Hostname()
	return prometheus.Labels{"hostname": hostname}
}

func (s *Service) getCounter(counter string) *Counter {
	if _, ok := s.Counters[counter]; !ok {
		s.NewCounter(counter, counter)
	}
	return s.Counters[counter]
}

func (s *Service) CounterInc(counter string) {
	c := s.getCounter(counter)
	c.Lock()
	defer c.Unlock()
	c.Value++
	c.Vec.With(s.getLabels()).Inc()
}

func (s *Service) getGauge(gauge string) *Gauge {
	if _, ok := s.Gauges[gauge]; !ok {
		s.NewGauge(gauge, gauge)
	}
	return s.Gauges[gauge]
}

func (s *Service) GetGaugeValue(gauge string) int {
	return s.getGauge(gauge).Value
}

func (s *Service) GaugeInc(gauge string) {
	g := s.getGauge(gauge)
	g.Lock()
	defer g.Unlock()
	g.Value++
	g.Vec.With(s.getLabels()).Inc()
}

func (s *Service) GaugeDec(gauge string) {
	g := s.getGauge(gauge)
	g.Lock()
	defer g.Unlock()
	if g.Value > 0 {
		g.Value--
		g.Vec.With(s.getLabels()).Dec()
	}

}
