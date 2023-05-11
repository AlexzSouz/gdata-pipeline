package events

type DomainEvent struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
}
