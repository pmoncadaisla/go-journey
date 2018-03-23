package controller

import (
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/eventbus"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/storedjourney"
)

// A Journey represents a journey with id and duration time
type Journey struct {
	domain.Journey
	bus                  evbus.Bus
	storedjourneyService storedjourneyservice.Interface
}

type Interface interface {
	Start()
}

func New() *Journey {
	j := &Journey{}
	j.bus = eventbus.Instance()
	j.storedjourneyService = storedjourneyservice.Instance()

	j.Start()

	return j
}

func (j *Journey) Start() {
	j.bus.Publish(eventbus.JOURNEY_STARTED.String(), j)
	go j.Run()
}

func (j *Journey) Run() {
	select {
	case <-time.After(1 * j.Time):
		j.Finish()
	}
}

func (j *Journey) Finish() {
	j.bus.Publish(eventbus.JOURNEY_FINISHED.String(), j)
}

func (j *Journey) Store() {
	j.storedjourneyService.SetLast(&j.Journey)
	j.bus.Publish(eventbus.JOURNEY_STORED.String(), j)

}

func (j *Journey) OnJourneyFinished(fj *domain.Journey) {
	if j.storedjourneyService.GetNextID() == j.GetJourneyID() {
		j.Store()
	}

}
