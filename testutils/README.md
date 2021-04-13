# Dynatrace API Test Utilities

This package is for testing purposes and is untested and unmaintained. It is not intended for public or production use.

## Packages

### config

Parses a dtconfig.json file into a go struct with test configurations.

### query

```go
// Query a metric by name and filter
func GetMostRecentValueForMetric(s selector, cfg config.Config) (*json.Number, error)
func GetMetrics(s selector, endpoint, apiToken string) ([]json.Number, error)

// filters
func Eq(key, value string) *filter
func And(f1, f2 *filter) *filter
func Or(f1, f2 *filter) *filter

// selector
func Selector(name string, filter *filter) *selector
```

```go
// query the most recent value of a metric with name
// "metric.name" with dimensions nonce=abc123 AND test=true
s := Selector("metric.name", And(Eq("nonce", "abc123"), Eq("test", "true")))
lastValue := GetMostRecentValueForMetric(s)

// query all values of a metric with name
// "metric.name" with dimensions nonce=abc123 AND test=true
values := GetMostRecentValueForMetric(Selector("metric.name", And(Eq("nonce", "abc123"), Eq("test", "true"))))
```