package types

type LogMetadata struct {
	LogEntryID int64  `json:"logEntryId,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
}

type LogEntry struct {
	Level      string            `json:"level,omitempty"`
	Message    string            `json:"message,omitempty"`
	ResourceID string            `json:"resourceId,omitempty"`
	Timestamp  string            `json:"timestamp,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
	SpanID     string            `json:"spanId,omitempty"`
	Commit     string            `json:"commit,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
