package journey

import (
	"time"

	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
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

var mutex sync.RWMutex

type Interface interface {
	Start()
}

func Receive(id int, time time.Duration, finished chan domain.Journey) {
	mutex.Lock()
	j := &Journey{Journey: domain.Journey{ID: id, Time: time}}
	j.metricsService = metricsservice.Instance()
	j.metricsService.CounterInc(metrics.JOURNEYS_RECEIVED.String())
	mutex.Unlock()
	j.finished = finished
	j.storedjourneyService = storedjourneyservice.Instance()
	j.Start()
}

func (j *Journey) Start() {
	j.SetStartTimeNow()
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("started")
	j.metricsService.CounterInc(metrics.JOURNEYS_STARTED_ALLTIME.String())
	go j.Run()
}

func (j *Journey) Run() {
	j.metricsService.GaugeInc(metrics.JOURNEYS_RUNNING.String())
	j.metricsService.GaugeInc(metrics.JOURNEYS_PENDING.String())
	time.Sleep(1 * j.Time)
	j.metricsService.GaugeDec(metrics.JOURNEYS_RUNNING.String())
	j.Finish()
}

func (j *Journey) Finish() {
	j.SetFinishTimeNow()
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("finished")
	j.metricsService.CounterInc(metrics.JOURNEYS_FINISHED_ALLTIME.String())
	go func() {
		j.finished <- j.Journey
	}()

}
