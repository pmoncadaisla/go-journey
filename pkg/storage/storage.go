package storage

import "github.com/pmoncadaisla/go-journey/pkg/domain"

// Interface must be implemented by all storage drivers
type Interface interface {
	Store(*domain.Journey)
}
