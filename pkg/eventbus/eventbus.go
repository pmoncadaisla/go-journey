package eventbus

import (
	"sync"

	evbus "github.com/asaskevich/EventBus"
)

var once sync.Once
var bus evbus.Bus

// Instance returns eventbus singleton
func Instance() evbus.Bus {
	once.Do(func() {
		bus = evbus.New()
	})
	return bus
}
