package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thehivecorporation/log"
	"github.com/thehivecorporation/log/writers/json"

	storecontroler "github.com/pmoncadaisla/go-journey/pkg/controller/store"
	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/journey"
	queueservice "github.com/pmoncadaisla/go-journey/pkg/service/queue"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetWriter(json.New(os.Stdout))
	log.SetLevel(log.LevelInfo)
	log.Info("Started")

	finished := make(chan domain.Journey)

	storeController := storecontroler.Instance(storecontroler.StoreConfig{
		Channel:     finished,
		OnlyHighest: false,
	})
	storeController.Start()

	wait := make(chan bool)

	queueservice.Instance()

	journeys := []domain.Journey{
		domain.Journey{ID: 5, Time: time.Millisecond * 1},
		domain.Journey{ID: 2, Time: time.Millisecond * 2},
		domain.Journey{ID: 1, Time: time.Millisecond * 3},
		domain.Journey{ID: 4, Time: time.Millisecond * 4},
		domain.Journey{ID: 3, Time: time.Second * 2},
		domain.Journey{ID: 6, Time: time.Millisecond * 6},
		domain.Journey{ID: 9, Time: time.Millisecond * 7},
		domain.Journey{ID: 10, Time: time.Millisecond * 8},
		domain.Journey{ID: 8, Time: time.Millisecond * 9},
		domain.Journey{ID: 7, Time: time.Millisecond * 10},
	}

	for _, j := range journeys {
		journey.Receive(j.ID, j.Time, finished)
	}

	// var i int
	// for i = 9; i > 0; i-- {
	// 	j := domain.Journey{ID: rand.Intn(1000), Time: time.Second * time.Duration(rand.Int63n(10)+1)}
	// 	queue.Push(j)
	// }

	// for queue.Len() > 0 {
	// 	element := queue.Pop()
	// 	log.Info(element)
	// }

	// for i = 1; i < 10; i++ {
	// 	journey.New(i, time.Second*time.Duration(rand.Int63n(10)+1), finished)
	// }

	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

	<-wait

}
