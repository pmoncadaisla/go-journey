package main

import (
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

	// Log to stdout
	log.SetWriter(json.New(os.Stdout))

	// Log above Info level
	log.SetLevel(log.LevelInfo)

	// In this channel we will receive the finished journey events
	finished = make(chan domain.Journey)

	// allStored receives signals whenever all pending journeys have been stored
	// adn there are no more journeys to store
	allStored = make(chan bool)

	// Configure Store Controller
	// if OnlyHighest is set to TRUE, then it only stores the journeys
	// with the highest cardinality as proposed
	storeController := storecontroler.Instance(storecontroler.StoreConfig{
		Channel:          finished,
		OnlyHighest:      true,
		AllStoredChannel: allStored,
	})
	storeController.Start()

	// These are singlenton services used by the application
	queueService = queueservice.Instance()
	metricsService = metricsservice.Instance()
	storedJourneyService = storedjourneyservice.Instance()

	// Start HTTP API.
	// Uses Gin Framework
	go webserver()

	// This is just for testing, it adds 10 journes by default
	// First 5 journeys are the 5 proposed in the definition document
	// 5,2,1,4,3
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

	// Intercept SIGINT and SIGTERM signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.WithField("signal", sig).Info("signal received")

	// Wait until all journeys have completed before existing
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
