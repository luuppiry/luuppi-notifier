package formatters

type EventsFormatter struct{}

func (f *EventsFormatter) Format(s []byte) []byte {
	return []byte{}
}
func NewEventsFormatter() *EventsFormatter {
	return &EventsFormatter{}
}
