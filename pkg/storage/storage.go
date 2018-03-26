package storage

import "github.com/pmoncadaisla/go-journey/pkg/domain"

type Interface interface {
	Store(*domain.Journey)
}
