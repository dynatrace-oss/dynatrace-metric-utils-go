# dynatrace-metric-utils-go

Utility for preparing communication with the [Dynatrace Metrics API v2](https://www.dynatrace.com/support/help/dynatrace-api/environment-api/metric-v2/).

## Usage

An example for how to use this library can be found in [example/main.go](example/main.go).
It shows how to create metrics lines that can be sent to a [Dynatrace metrics ingest endpoint](https://www.dynatrace.com/support/help/dynatrace-api/environment-api/metric-v2/post-ingest-metrics/) using an HTTP client library.

### Preparation

The standard workflow consists of creating `NormalizedDimensionList` objects that contain `Dimension`s.
There is no upper limit to how many `NormalizedDimensionLists` can be used, but we found that the following schema covers most use cases.
In this schema we use three lists, all of which are optional.

* The first list of dimensions are default dimensions, which contain information that can be shared between multiple metrics, e.g. the hostname of the machine.
* The second `NormalizedDimensionList` contains dimensions that are specific to a certain metric, e.g. the information whether or not a HTTP request was successful.
* The third list are the Dimensions created by the OneAgent metadata enricher, which is described further below.

The serialization function accepts a merged `NormalizedDimensionList` which can be acquired using the `dimensions.MergeLists` function.
Dimensions in lists passed further right with the same (normalized) dimension keys overwrite dimensions passed in lists further left.

> Note that the MergeLists function must be called every time a new dimension is added to any of the lists!

### Metric line creation

After the creation of the dimensions, the `metric` package allows for the creation of metric lines.
To do so, use the following pattern:

```go
m, err := metric.NewMetric(
  "the_metric_key",
  metric.WithPrefix("prefix"),
  metric.WithIntCounterValueTotal(30),
  metric.WithDimensions(merged),
  metric.WithCurrentTime(),
)
// handle potential errors... 
serialized, err := m.Serialize()
// handle potential errors... 
fmt.Println(serialized)
```

The serialized data point is ready to be sent to a Dynatrace metrics ingest endpoint using an HTTP client library.

#### Metric line creation options

* `WithPrefix`: set a prefix that will be prepended to the metric key.
* `WithDimensions`: sets a `NormalizedDimensionList` for serialization.
  Lists should be de-duplicated and combined before being passed to this function by running them through the `MergeLists` function.
  If only one list is present, `MergeLists` will still do the de-duplication.
* `WithIntCounterValueTotal` / `WithFloatCounterValueTotal`: sets a single value that is serialized as `count,<value>`.
* `WithIntCounterValueDelta` / `WithFloatCounterValueDelta`: sets a single value that is serialized as `count,delta=<value>`.
* `WithIntGaugeValue` / `WithFloatGaugeValue`: sets a single value that is serialized as `gauge,<value>`.
* `WithIntSummaryValue` / `WithFloatSummaryValue`: sets min, max, sum and count values that are serialized as `gauge,min=<min>,max=<max>,sum=<sum>,count=<count>`.
* `WithTimestamp`: sets a specific `time.Time` object on the metric that will be used to create the timestamp on the metric line.
* `WithCurrentTime`: sets the current timestamp to the `Metric` object.

A metric line can be serialized only if it has a valid name (including the optional prefix) and exactly one `Value` attribute set.
Timestamps and dimensions are optional.

### OneAgent Enrichment

Due to how OneAgent metadata is read and how Go reads files, it is at the moment not possible to read metadata on Unix/Linux systems.
OneAgent enrichment for Go therefore only functions on Windows hosts at the moment.
On Unix/Linux hosts, an empty list will be returned without any errors, if it is called nevertheless.
The same applies if no OneAgent is installed on the monitored host.

To acquire a list of OneAgent metadata dimensions, use the following method:

```go
oneAgentDimensions := oneagentenrichment.GetOneAgentMetadata()
```

These dimensions can then be passed to the `MergeLists` function as shown in [the example](example/main.go).

### Common constants

The library also provides constants that might be helpful in the projects consuming this library.

To access the constants, call the respective methods from the `apiconstants` package:

```go
defaultOneAgentEndpoint := apiconstants.GetDefaultOneAgentEndpoint()
```

Currently available constants are:

* the default [local OneAgent metric API](https://www.dynatrace.com/support/help/how-to-use-dynatrace/metrics/metric-ingestion/ingestion-methods/local-api/) endpoint (`GetDefaultOneAgentEndpoint()`)
* the limit for how many metric lines can be ingested in one request (`GetPayloadLinesLimit()`)
