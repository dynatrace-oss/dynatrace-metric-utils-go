package apiconstants

const (
	defaultOneAgentEndpoint = "http://127.0.0.1:14499/metrics/ingest"
	payloadLinesLimit       = 1000
)

// GetDefaultOneAgentEndpoint returns the default OneAgent metrics ingest endpoint.
func GetDefaultOneAgentEndpoint() string {
	return defaultOneAgentEndpoint
}

// GetMetricPayloadLinesLimit returns the maximum number of lines per POST request.
func GetPayloadLinesLimit() int {
	return payloadLinesLimit
}
