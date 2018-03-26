package controller

import (
	"sync"
	"time"

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
	OnlyHighest   bool
	Channel       chan domain.Journey
	FinishTimeout time.Duration
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
	finishTimeoutTimer := time.NewTimer(c.config.FinishTimeout)
	for {
		select {
		case journey := <-c.finished:
			finishTimeoutTimer.Reset(c.config.FinishTimeout)
			c.onJourneyFinished(journey)
		case <-finishTimeoutTimer.C:
			c.skipJourney()

		}
	}
}

func (c *Controller) skipJourney() {
	// Check if there are elements in the queue or not
	if c.queue.Len() == 0 {
		return
	}
	nextID := c.storedjourneyService.GetNextID()

	// If an element that finishes that was already skiped, it remains in the queue.
	// check if we are skiping a journey that hasn't arrived yet
	if nextID > c.storedjourneyService.GetReceivedHighestID() {
		return
	}

	c.storedjourneyService.SetLast(&domain.Journey{ID: nextID})
	c.metricsService.CounterInc("journeys_skipped")
	log.WithField("ID", nextID).WithField("reason", "timeout").Info("skipped")
	c.checkAndStoreJourney()
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
