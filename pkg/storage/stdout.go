package storage

import (
	"sync"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/thehivecorporation/log"
)

type StdoutStorage struct {
}

var once sync.Once
var ss *StdoutStorage

func NewStdoutStorage() *StdoutStorage {
	log.WithField("Driver: ", "STDout").Info("storage driver created")
	once.Do(func() {
		ss = &StdoutStorage{}

	})
	return ss
}

func (ss *StdoutStorage) Store(journey *domain.Journey) {
	log.WithField("ID: ", journey.ID).WithField("Duration: ", journey.Time).Info("stored")
}
