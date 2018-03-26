package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thehivecorporation/log"
	"github.com/thehivecorporation/log/writers/json"

	storecontroler "github.com/pmoncadaisla/go-journey/pkg/controller/store"
	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/journey"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
	metricsservice "github.com/pmoncadaisla/go-journey/pkg/service/metrics"
	queueservice "github.com/pmoncadaisla/go-journey/pkg/service/queue"
	storedjourneyservice "github.com/pmoncadaisla/go-journey/pkg/service/store"
)

var finished chan domain.Journey
var allStored chan bool
var storedJourneyService storedjourneyservice.Interface
var queueService queueservice.Interface
var metricsService metricsservice.Interface

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	log.SetWriter(json.New(os.Stdout))
	log.SetLevel(log.LevelInfo)
	log.Info("Started")

	finished = make(chan domain.Journey)
	allStored = make(chan bool)

	storeController := storecontroler.Instance(storecontroler.StoreConfig{
		Channel:          finished,
		OnlyHighest:      false,
		AllStoredChannel: allStored,
	})
	storeController.Start()
	queueService = queueservice.Instance()
	metricsService = metricsservice.Instance()
	storedJourneyService = storedjourneyservice.Instance()

	go webserver()

	journeys := []domain.Journey{
		domain.Journey{ID: 5, Time: time.Millisecond * 1},
		domain.Journey{ID: 2, Time: time.Millisecond * 2},
		domain.Journey{ID: 1, Time: time.Millisecond * 3},
		domain.Journey{ID: 4, Time: time.Millisecond * 4},
		domain.Journey{ID: 3, Time: time.Second * 7},
		domain.Journey{ID: 6, Time: time.Millisecond * 6},
		domain.Journey{ID: 9, Time: time.Millisecond * 7},
		domain.Journey{ID: 10, Time: time.Millisecond * 8},
		domain.Journey{ID: 8, Time: time.Millisecond * 9},
		domain.Journey{ID: 7, Time: time.Millisecond * 10},
	}

	for _, j := range journeys {
		journey.Receive(j.ID, j.Time, finished)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.WithField("signal", sig).Info("signal received")
	graceFullShutDownWait()

}

func graceFullShutDownWait() {
	log.Info("Waiting for all journeys to finish...")
	<-allStored
	if metricsService.GetGaugeValue(metrics.JOURNEYS_PENDING.String()) > 0 {
		graceFullShutDownWait()
	}
	log.Info("Endping program")

}
