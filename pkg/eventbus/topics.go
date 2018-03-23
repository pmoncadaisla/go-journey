package eventbus

// ChannelName ... type JOURNEY_RECEIVED, JOURNEY_STARTED, JOURNEY_FINISHED, JOURNEY_STORED
type TopicName int

const (
	// JOURNEY_RECEIVED ...
	JOURNEY_RECEIVED TopicName = 1 + iota
	// JOURNEY_STARTED ...
	JOURNEY_STARTED
	// JOURNEY_FINISHED ...
	JOURNEY_FINISHED
	// JOURNEY_STORED ...
	JOURNEY_STORED
)

var topicName = [...]string{
	"journey:received",
	"journey:started",
	"journey:finished",
	"journey:stored",
}

func (n TopicName) String() string {
	return topicName[n-1]
}
