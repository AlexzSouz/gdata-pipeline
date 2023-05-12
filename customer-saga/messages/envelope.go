package messages

type Envelope[T any] struct {
	TraceId string `json:"trace_id"`
	SpanId  string `json:"span_id"`
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}
