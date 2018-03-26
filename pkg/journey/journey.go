package journey

import (
	"time"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	metricsservice "github.com/pmoncadaisla/go-journey/pkg/service/metrics"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/store"
	"github.com/thehivecorporation/log"
)

// A Journey represents a journey with id and duration time
type Journey struct {
	domain.Journey
	finished             chan domain.Journey
	storedjourneyService storedjourneyservice.Interface
	metricsService       metricsservice.Interface
}

type Interface interface {
	Start()
}

func Receive(id int, time time.Duration, finished chan domain.Journey) {
	j := &Journey{Journey: domain.Journey{ID: id, Time: time}}
	j.finished = finished
	j.metricsService = metricsservice.Instance()
	j.storedjourneyService = storedjourneyservice.Instance()
	j.Start()
}

func (j *Journey) Start() {
	j.SetStartTimeNow()
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("started")
	j.metricsService.CounterInc("journeys_started_alltime")
	go j.Run()
}

func (j *Journey) Run() {
	time.Sleep(1 * j.Time)
	j.Finish()
}

func (j *Journey) Finish() {
	j.SetFinishTimeNow()
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("finished")
	j.metricsService.CounterInc("journeys_finished_alltime")
	go func() {
		j.finished <- j.Journey
	}()

}
