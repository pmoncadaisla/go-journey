package metrics

// MetricName ... type JORNEYS_STORED, JOURNEYS_FINISHED, JOURNEYS_STARTED
type Metric int

// journeys_stored
// journeys_finished
// journeys_started
const (
	// JORNEYS_STORED ...
	JOURNEYS_STORED Metric = 1 + iota
	// JOURNEYS_FINISHED ...
	JOURNEYS_FINISHED
	// JOURNEYS_STARTED ...
	JOURNEYS_STARTED
	// JOURNEYS_RUNNING ...
	JOURNEYS_RUNNING
	// JOURNEYS_PENDING ...
	JOURNEYS_PENDING
	// JOURNEYS_STARTED_ALLTIME ...
	JOURNEYS_STARTED_ALLTIME
	// JOURNEYS_FINISHED_ALLTIME ...
	JOURNEYS_FINISHED_ALLTIME
	// JOURNEYS_HIGHEST_RECEIVED_ID ...
	JOURNEYS_HIGHEST_RECEIVED_ID
	// JOURNEYS_LAST_STORED_ID ...
	JOURNEYS_LAST_STORED_ID
	// JOURNEYS_RECEIVED ...
	JOURNEYS_RECEIVED
)

var metric = [...]string{
	"journeys_stored",
	"journeys_finished",
	"journeys_started",
	"journeys_running",
	"journeys_pending",
	"journeys_started_alltime",
	"journeys_finished_alltime",
	"journeys_highest_received_id",
	"journeys_last_stored_id",
	"journeys_received",
}

func (n Metric) String() string {
	return metric[n-1]
}
