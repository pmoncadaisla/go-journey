package journey

import (
	"time"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/thehivecorporation/log"
)

// A Journey represents a journey with id and duration time
type Journey struct {
	domain.Journey
	finished chan domain.Journey
}

type Interface interface {
	Start()
}

func New(id int, time time.Duration, finished chan domain.Journey) {
	j := &Journey{Journey: domain.Journey{ID: id, Time: time}}
	j.finished = finished
	j.Start()
}

func (j *Journey) Start() {
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("started")
	go j.Run()
}

func (j *Journey) Run() {
	time.Sleep(1 * j.Time)
	j.Finish()
}

func (j *Journey) Finish() {
	log.WithField("id", j.GetJourneyID()).WithField("time", j.GetJourneyTime()).Info("finished")
	go func() {
		j.finished <- j.Journey
	}()

}
