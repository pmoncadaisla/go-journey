package controller

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
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
	OnlyHighest      bool
	Channel          chan domain.Journey
	AllStoredChannel chan bool
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
	c.metricsService.GaugeInc(metrics.JOURNEYS_FINISHED.String())
	c.checkAndStoreJourney(&j)
}

func (c *Controller) checkAndStoreJourney(j *domain.Journey) {
	if c.queue.Len() > 0 && c.storedjourneyService.GetNextID() == c.queue.Get().ID {
		journey := c.queue.Pop()
		c.store(journey)
		c.metricsService.GaugeDec(metrics.JOURNEYS_FINISHED.String())
		c.checkAndStoreJourney(j)
	}
}

func (c *Controller) store(journey domain.Journey) {
	if c.config.OnlyHighest && c.queue.Len() > 0 && c.queue.Get().ID == (journey.ID+1) {
		journey := c.queue.Pop()
		c.store(journey)
		journey.SetStoreTimeNow()
	} else {
		log.WithField("ID: ", journey.ID).WithField("Duration: ", journey.Time).Info("stored")
		c.metricsService.CounterInc(metrics.JOURNEYS_STORED.String())
		c.metricsService.GaugeDec(metrics.JOURNEYS_PENDING.String())
		c.metricsService.GaugeSet(metrics.JOURNEYS_LAST_STORED_ID.String(), journey.ID)
		c.storedjourneyService.SetLast(&journey)
		if c.queue.Len() == 0 {
			go func() {
				c.config.AllStoredChannel <- true
			}()
		}
	}
}
