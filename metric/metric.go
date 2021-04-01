// Copyright 2021 Dynatrace LLC
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/serialize"
)

// Metric contains all information needed to create a string representation of the accumulated metric data.
type Metric struct {
	metricKey  string
	prefix     string
	value      metricValue
	dimensions dimensions.NormalizedDimensionList
	timestamp  time.Time
}

// MetricOption represents the function interface used to set options on the metric object.
type MetricOption func(m *Metric) error

func joinStrings(key, dim, value, timestamp string) (string, error) {
	var sb strings.Builder

	sb.WriteString(key)
	if dim != "" {
		sb.WriteString(",")
		sb.WriteString(dim)
	}

	sb.WriteString(" ")
	sb.WriteString(value)

	if timestamp != "" {
		sb.WriteString(" ")
		sb.WriteString(timestamp)
	}

	return sb.String(), nil

}

func (m Metric) ensureRequiredFieldsSet() error {
	if m.metricKey == "" && m.prefix == "" {
		return errors.New("metric key and prefix empty, cannot create metric name")
	}

	if m.value == nil {
		return errors.New("metric value not set, cannot create metric")
	}

	return nil
}

// Serialize creates the string representation of the Metric object.
func (m Metric) Serialize() (string, error) {
	keyString, err := serialize.MetricKey(m.metricKey, m.prefix)
	if err != nil {
		return "", err
	}
	if m.value == nil {
		return "", errors.New("cannot serialize nil value")
	}

	dimString := serialize.Dimensions(m.dimensions)
	valueString := m.value.serialize()
	timeString := serialize.Timestamp(m.timestamp)

	return joinStrings(keyString, dimString, valueString, timeString)
}

// NewMetric creates a new metric with a mandatory name and options. At least one value option must be set.
func NewMetric(name string, options ...MetricOption) (*Metric, error) {
	m := &Metric{
		metricKey: name,
	}

	for _, option := range options {
		err := option(m)
		if err != nil {
			return nil, err
		}
	}

	err := m.ensureRequiredFieldsSet()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func checkValueAlreadySet(m *Metric) error {
	if m.value != nil {
		return errors.New("cannot set two values on one metric.")
	}
	return nil
}

// WithPrefix sets the prefix for Metric creation.
func WithPrefix(prefix string) MetricOption {
	return func(m *Metric) error {
		m.prefix = prefix
		return nil
	}
}

// WithDimensions sets the passed dimension list as the dimensions to export.
// Pass only dimension lists that have been deduplicated (by  Merge)
func WithDimensions(dims dimensions.NormalizedDimensionList) MetricOption {
	return func(m *Metric) error {
		m.dimensions = dims

		return nil
	}
}

func trySetValue(m *Metric, val metricValue) error {
	if err := checkValueAlreadySet(m); err != nil {
		return err
	}
	m.value = val
	return nil
}

// WithIntCounterValueTotal sets a value on the metric that will be formatted as "count,<value>"
func WithIntCounterValueTotal(val int64) MetricOption {
	return func(m *Metric) error {
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}

		return trySetValue(m, intCounterValue{value: val, isDelta: false})
	}
}

// WithIntCounterValueDelta sets a value on the metric that will be formatted as "count,delta=<value>"
func WithIntCounterValueDelta(val int64) MetricOption {
	return func(m *Metric) error {
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}

		return trySetValue(m, intCounterValue{value: val, isDelta: true})
	}
}

// WithFloatCounterValueTotal sets a value on the metric that will be formatted as "count,<value>"
func WithFloatCounterValueTotal(val float64) MetricOption {
	return func(m *Metric) error {
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}

		return trySetValue(m, floatCounterValue{value: val, isDelta: false})
	}
}

// WithFloatCounterValueDelta sets a value on the metric that will be formatted as "count,delta=<value>"
func WithFloatCounterValueDelta(val float64) MetricOption {
	return func(m *Metric) error {
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}

		return trySetValue(m, floatCounterValue{value: val, isDelta: true})
	}
}

// WithIntSummaryValue sets a summary statistic on the metric that will be formatted as
// "gauge,min=<min>,max=<max>,sum=<sum>,count=<count>".
func WithIntSummaryValue(min, max, sum, count int64) MetricOption {
	return func(m *Metric) error {
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		if min > max {
			return fmt.Errorf("min (%d) cannot be greater than max (%d)", min, max)
		}

		return trySetValue(m, intSummaryValue{min: min, max: max, sum: sum, count: count})
	}
}

//  WithFloatSummaryValue sets a summary statistic on the metric that will be formatted as
// "gauge,min=<min>,max=<max>,sum=<sum>,count=<count>".
func WithFloatSummaryValue(min, max, sum float64, count int64) MetricOption {
	return func(m *Metric) error {
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		if min > max {
			return fmt.Errorf("min (%.3f) cannot be greater than max (%.3f)", min, max)
		}

		return trySetValue(m, floatSummaryValue{min: min, max: max, sum: sum, count: count})
	}
}

// WithIntGaugeValue sets a gauge value on the metric that will be formatted as "gauge,<value>"
func WithIntGaugeValue(val int64) MetricOption {
	return func(m *Metric) error {
		return trySetValue(m, intGaugeValue{value: val})
	}
}

// WithFloatGaugeValue sets a gauge value on the metric that will be formatted as "gauge,<value>"
func WithFloatGaugeValue(val float64) MetricOption {
	return func(m *Metric) error {
		return trySetValue(m, floatGaugeValue{value: val})
	}
}

// WithTimestamp sets a specific timestamp for the metric.
func WithTimestamp(t time.Time) MetricOption {
	return func(m *Metric) error {
		m.timestamp = t
		return nil
	}
}

// WithCurrentTime sets the current time as timestamp for the metric.
func WithCurrentTime() MetricOption {
	return WithTimestamp(time.Now())
}
