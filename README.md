# dynatrace-metric-utils-go

Utility for preparing communication with the [Dynatrace v2 metrics API](https://www.dynatrace.com/support/help/dynatrace-api/environment-api/metric-v2/).

## Usage

An example for how to use this library can be found in [example/main.go](example/main.go).
It shows how to create metrics lines that can be sent to any [Dynatrace metrics ingest endpoint](https://www.dynatrace.com/support/help/dynatrace-api/environment-api/metric-v2/post-ingest-metrics/).

### Preparation

The standard workflow consists of creating `NormalizedDimensionList` objects that contain `Dimension`s.
There is no upper limit to how many `NormalizedDimensionLists` can be used, but we found that the following schema covers most use cases.
In this schema we use three lists, all of which are optional.

* The first list of dimensions are default dimensions, which contain information that can be shared between multiple metrics, e.g. the hostname of the machine.
* The second `NormalizedDimensionList` contains dimensions that are specific to a certain metric, e.g. the information whether or not a HTTP request was successful.
* The third list are the Dimensions created by the OneAgent metadata enricher.

The serialization function accepts a merged `NormalizedDimensionList` which can be  acquired using the `dimensions.MergeLists` function.
Dimensions in lists passed further right with the same (normalized) dimension keys  overwrite dimensions passed in lists further left.

> Note that the merge function must be called every time a new dimension is added to any of the lists!

### Metric line creation

After the creation of the dimensions, the `metric` package allows for the creation of metric lines.
To do so, simply use the following pattern:

```go
m, err := metric.NewMetric(
 "name",
  metric.WithPrefix("prefix"),
  metric.WithIntAbsoluteCounterValue(30),
  metric.WithDimensions(merged),
  metric.WithCurrentTime(),
)
// handle potential errors... 
serialized, err := m.Serialize()
// handle potential errors... 
fmt.Println(serialized)
```

The created metric line is ready to be sent to the Dynatrace metrics endpoint!

### OneAgent Enrichment

Due to how OneAgent metadata is read and how Go reads files, it is at the moment not possible to read metadata on Linux systems.
OneAgent enrichment for Go therefore only functions on Windows hosts at the moment.
To acquire a list of OneAgent metadata dimensions, simply call the following method:

```go
oneAgentDimensions := oneagentenrichment.GetOneAgentMetadata()
```

These dimensions can then simply be passed to the `MergeLists` function as shown in [the example](example/main.go).
