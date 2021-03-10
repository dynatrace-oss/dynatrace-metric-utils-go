# dynatrace-metric-utils-go

Utility for interacting with the Dynatrace Metrics v2 API

## Normalization

The `normalization` package contains functions to normalize metric keys, dimension keys, and dimension values. It can be used by calling the functions, for example:

```go
normalize.MetricKey("some_metric_key)
```

<!-- 
These functions are also used by the `serialize` package:

## Serialization -->
