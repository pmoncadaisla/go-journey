package controller

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	metricsservice "github.com/pmoncadaisla/go-journey/pkg/service/metrics"
	queueservice "github.com/pmoncadaisla/go-journey/pkg/service/queue"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/store"
	"github.com/thehivecorporation/log"
)

type Controller struct {
	storedjourneyService storedjourneyservice.Interface
	finished             chan domain.Journey
	queue                queueservice.Interface
	config               StoreConfig
	metricsService       metricsservice.Interface
}

type StoreConfig struct {
	OnlyHighest bool
	Channel     chan domain.Journey
}

var once sync.Once
var c *Controller

// Instance returns eventbus singleton
func Instance(config StoreConfig) *Controller {
	once.Do(func() {
		c = &Controller{}
		c.storedjourneyService = storedjourneyservice.Instance()
		c.finished = config.Channel
		c.metricsService = metricsservice.Instance()
		c.config = config
		c.queue = queueservice.Instance()
	})
	return c
}

func (c *Controller) Start() {
	go c.run()
}

func (c *Controller) run() {
	for {
		select {
		case journey := <-c.finished:
			c.onJourneyFinished(journey)
		}
	}
}

func (c *Controller) onJourneyFinished(j domain.Journey) {
	c.queue.Push(j)
	c.metricsService.GaugeInc("journeys_finished")
	c.checkAndStoreJourney()
}

func (c *Controller) checkAndStoreJourney() {
	if c.queue.Len() > 0 && c.storedjourneyService.GetNextID() == c.queue.Get().ID {
		//log.WithField("ID", c.queue.Get().ID).Info("c.queue.Get().ID")
		journey := c.queue.Pop()
		c.store(journey)
		c.metricsService.GaugeDec("journeys_finished")
		c.checkAndStoreJourney()
	}
}

func (c *Controller) store(journey domain.Journey) {
	if c.config.OnlyHighest && c.queue.Len() > 0 && c.queue.Get().ID == (journey.ID+1) {
		journey := c.queue.Pop()
		c.store(journey)
		journey.SetStoreTimeNow()
	} else {
		log.WithField("ID: ", journey.ID).WithField("Duration: ", journey.Time).Info("Process Stored")
		c.metricsService.CounterInc("journeys_stored")
		c.storedjourneyService.SetLast(&journey)
	}
}
