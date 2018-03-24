package domain

import (
	"time"
)

// A Journey represents a journey with id and duration time
type Journey struct {
	ID   int           `json:"journey_id"`
	Time time.Duration `json:"journey_time"`
}

// GetJourneyID returns Journeys ID
func (j *Journey) GetJourneyID() int {
	return j.ID
}

// GetJourneyTime returns Journey's Time
func (j *Journey) GetJourneyTime() time.Duration {
	return j.Time
}
