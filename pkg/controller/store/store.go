package controller

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	queueservice "github.com/pmoncadaisla/go-journey/pkg/service/queue"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/store"
	"github.com/thehivecorporation/log"
)

type Controller struct {
	storedjourneyService storedjourneyservice.Interface
	finished             chan domain.Journey
	queue                queueservice.Interface
	config               StoreConfig
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
		c.config = config
		c.queue = queueservice.Instance()
	})
	return c
}

func (c *Controller) Start() {
	go c.Run()
}

func (c *Controller) Run() {
	for {
		select {
		case journey := <-c.finished:
			c.queue.Push(journey)
			c.OnJourneyFinished(&journey)

		}
	}
}

func (c *Controller) OnJourneyFinished(j *domain.Journey) {
	c.checkAndStoreJourney()
}

func (c *Controller) checkAndStoreJourney() {
	if c.queue.Len() > 0 && c.storedjourneyService.GetNextID() == c.queue.Get().ID {
		//log.WithField("ID", c.queue.Get().ID).Info("c.queue.Get().ID")
		journey := c.queue.Pop()
		c.store(journey)
		c.checkAndStoreJourney()
	}

}

func (c *Controller) store(journey domain.Journey) {
	if c.config.OnlyHighest && c.queue.Len() > 0 && c.queue.Get().ID == (journey.ID+1) {
		journey := c.queue.Pop()
		c.store(journey)
	} else {
		log.WithField("ID: ", journey.ID).WithField("Duration: ", journey.Time).Info("Process Stored")
		c.storedjourneyService.SetLast(&journey)
	}
}
