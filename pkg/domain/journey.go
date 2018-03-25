package domain

import (
	"time"
)

// A Journey represents a journey with id and duration time
type Journey struct {
	ID         int           `json:"journey_id"`
	Time       time.Duration `json:"journey_time"`
	StartTime  time.Time     `json:"-"`
	FinishTime time.Time     `json:"-"`
	StoreTime  time.Time     `json:"-"`
}

// GetJourneyID returns Journeys ID
func (j *Journey) GetJourneyID() int {
	return j.ID
}

// GetJourneyTime returns Journey's Time
func (j *Journey) GetJourneyTime() time.Duration {
	return j.Time
}

// SetStartTime sets start time to given time
func (j *Journey) SetStartTime(t time.Time) {
	j.StartTime = t
}

// SetFinishTime sets finish time to given time
func (j *Journey) SetFinishTime(t time.Time) {
	j.StartTime = t
}

// SetStoreTime sets store time to given time
func (j *Journey) SetStoreTime(t time.Time) {
	j.StartTime = t
}

// SetStartTimeNow sets start time to time.Now()
func (j *Journey) SetStartTimeNow() {
	j.SetStartTime(time.Now())
}

// SetFinishTimeNow sets finish time to time.Now()
func (j *Journey) SetFinishTimeNow() {
	j.SetFinishTime(time.Now())
}

// SetStoreimeNow sets store time to time.Now()
func (j *Journey) SetStoreTimeNow() {
	j.SetStoreTime(time.Now())
}
