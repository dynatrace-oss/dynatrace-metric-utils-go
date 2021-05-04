package apiconstants

const (
	defaultOneAgentEndpoint = "http://127.0.0.1:14499/metrics/ingest"
	lineLengthLimit         = 2000
	metricPayloadLinesLimit = 1000
)

// GetDefaultOneAgentEndpoint returns the default OneAgent metrics ingest endpoint.
func GetDefaultOneAgentEndpoint() string {
	return defaultOneAgentEndpoint
}

// GetLineLengthLimit returns the maximum number of characters allowed per metric line.
func GetLineLengthLimit() int {
	return lineLengthLimit
}

// GetMetricPayloadLinesLimit returns the maximum number of lines per POST request.
func GetMetricPayloadLinesLimit() int {
	return metricPayloadLinesLimit
}
