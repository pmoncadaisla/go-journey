package controller

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
	metricsservice "github.com/pmoncadaisla/go-journey/pkg/service/metrics"
	queueservice "github.com/pmoncadaisla/go-journey/pkg/service/queue"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/store"
	"github.com/pmoncadaisla/go-journey/pkg/storage"
	"github.com/thehivecorporation/log"
)

type Controller struct {
	storedjourneyService storedjourneyservice.Interface
	finished             chan domain.Journey
	queue                queueservice.Interface
	config               StoreConfig
	storagedriver        storage.Interface
	metricsService       metricsservice.Interface
}

type StoreConfig struct {
	OnlyHighest      bool
	Channel          chan domain.Journey
	AllStoredChannel chan bool
	StorageDriver    storage.Interface
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
		c.storagedriver = storage.NewStdoutStorage()
		c.config = config
		c.queue = queueservice.Instance()
	})
	return c
}

// Start just runs in a goroutine
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

// onJourneyFinished is called whenever a journey has finished
func (c *Controller) onJourneyFinished(j domain.Journey) {
	c.queue.Push(j)
	c.metricsService.GaugeInc(metrics.JOURNEYS_FINISHED.String())
	c.checkAndStoreJourney(&j)
}

// checkAndStoreJourney checks if a journey is suposed to be stored
// if it is not going to be stored, it just skips
func (c *Controller) checkAndStoreJourney(j *domain.Journey) {
	if c.queue.Len() > 0 && c.storedjourneyService.GetNextID() == c.queue.Get().ID {
		journey := c.queue.Pop()
		c.store(journey)
		c.metricsService.GaugeDec(metrics.JOURNEYS_FINISHED.String())
		c.checkAndStoreJourney(j)
	}
}

// Store the journey
// It uses a Storage driver
// We are using a stdout storage driver, which just prints using log library
func (c *Controller) store(journey domain.Journey) {
	if c.config.OnlyHighest && c.queue.Len() > 0 && c.queue.Get().ID == (journey.ID+1) {
		journey := c.queue.Pop()
		// Updates pending in batch, ex:
		// if last was 2, and we are storing 5, then 3 and 4 also stop pending
		// Pending = Pending - (Last ID - This ID)
		// We decrease 5-2=3
		pending := c.metricsService.GetGaugeValue(metrics.JOURNEYS_PENDING.String()) - journey.ID - c.metricsService.GetGaugeValue(metrics.JOURNEYS_LAST_STORED_ID.String())
		c.metricsService.GaugeSet(metrics.JOURNEYS_PENDING.String(), pending)
		log.WithField("Decreased count: ", pending).Info("decreased pending")
		c.store(journey)
		journey.SetStoreTimeNow()
	} else {
		c.storagedriver.Store(&journey)
		c.metricsService.CounterInc(metrics.JOURNEYS_STORED.String())
		c.metricsService.GaugeDec(metrics.JOURNEYS_PENDING.String())
		log.WithField("Decreased count: ", 1).Info("decreased pending")
		c.metricsService.GaugeSet(metrics.JOURNEYS_LAST_STORED_ID.String(), journey.ID)
		c.storedjourneyService.SetLast(&journey)
		if c.queue.Len() == 0 {
			go func() {
				c.config.AllStoredChannel <- true
			}()
		}
	}
}
