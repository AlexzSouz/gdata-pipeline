package aggregates

type Aggregate struct {
	Id               int `json:"id"`
	UncommitedEvents []interface{}
}

// Raise domain events in a Event-Sourcing architecture
func (a *Aggregate) Raise(event interface{}) {
	if a.UncommitedEvents == nil || len(a.UncommitedEvents) == 0 {
		a.UncommitedEvents = make([]interface{}, 0)
	}

	// TODO : Implement custom Marshal/Unmarshal for serializing event type
	a.UncommitedEvents = append(a.UncommitedEvents, event)
}
