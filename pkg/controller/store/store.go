package controller

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/storedjourney"
)

type Controller struct {
	storedjourneyService storedjourneyservice.Interface
}

type Interface interface {
}

var once sync.Once
var c *Controller

// Instance returns eventbus singleton
func Instance() *Controller {
	once.Do(func() {
		c = &Controller{}
		c.storedjourneyService = storedjourneyservice.Instance()
	})
	return c
}

func (c *Controller) OnJourneyFinished(j *domain.Journey) {
	if c.storedjourneyService.GetNextID() == j.GetJourneyID() {

	}

}
