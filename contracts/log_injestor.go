package contracts

import "time"

type LogEntry struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ResourceID string    `json:"resourceId"`
	Timestamp  time.Time `json:"timestamp"`
	TraceID    string    `json:"traceId"`
	SpanID     string    `json:"spanId"`
	Commit     string    `json:"commit"`
	Metadata   Metadata  `json:"metadata"`
}

type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}
