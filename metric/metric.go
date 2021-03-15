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
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/serialize"
)

type Metric struct {
	name         string
	prefix       string
	value        metricValue
	dimensions   dimensions.NormalizedDimensionSet
	isDelta      bool
	timestamp    time.Time
	timestampSet bool
}

func joinStrings(key, dim, value, timestamp string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}

	if value == "" {
		return "", errors.New("values cannot be empty")
	}

	var sb strings.Builder

	sb.WriteString(key)
	if dim != "" {
		sb.WriteString(" ")
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
	if m.name == "" && m.prefix == "" {
		return errors.New("metric key and prefix empty, cannot create metric name")
	}

	if m.value == nil {
		return errors.New("metric value not set, cannot create metric")
	}

	return nil
}

func (m Metric) Serialize() (string, error) {
	keyString, err := serialize.MetricName(m.name, m.prefix)
	if err != nil {
		return "", err
	}

	dimString := serialize.NormalizedDimensions(m.dimensions)

	timeString := ""
	if m.timestampSet {
		timeString = strconv.FormatInt(m.timestamp.Unix(), 10)
	}

	return joinStrings(keyString, dimString, m.value.serialize(), timeString)

}

func NewMetric(name string, options ...MetricOption) (*Metric, error) {
	m := &Metric{
		name:         name,
		timestampSet: false,
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

type MetricOption func(m *Metric) error

func checkValueSet(m *Metric) error {
	if m.value != nil {
		return errors.New("cannot set two values on one metric.")
	}
	return nil
}

func WithPrefix(prefix string) MetricOption {
	return func(m *Metric) error {
		m.prefix = prefix
		return nil
	}
}

// if this function is passed multiple times, the last set of dimensions is used.
func WithDimensions(dims dimensions.NormalizedDimensionSet) MetricOption {
	return func(m *Metric) error {
		m.dimensions = dims

		return nil
	}
}

func WithIntCounterValue(val int64) MetricOption {
	return func(m *Metric) error {
		if err := checkValueSet(m); err != nil {
			return err
		}
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}
		m.value = intCounterValue{value: val}

		return nil
	}
}

func WithFloatCounterValue(val float64) MetricOption {
	return func(m *Metric) error {
		if err := checkValueSet(m); err != nil {
			return err
		}
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}
		m.value = floatCounterValue{value: val}

		return nil
	}
}

func WithIntSummaryValue(min, max, sum, count int64) MetricOption {
	return func(m *Metric) error {
		if err := checkValueSet(m); err != nil {
			return err
		}
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		m.value = intSummaryValue{min: min, max: max, sum: sum, count: count}

		return nil
	}
}

func WithFloatSummaryValue(min, max, sum float64, count int64) MetricOption {
	return func(m *Metric) error {
		if err := checkValueSet(m); err != nil {
			return err
		}
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		m.value = floatSummaryValue{min: min, max: max, sum: sum, count: count}

		return nil
	}
}

func WithDelta(isDelta bool) MetricOption {
	return func(m *Metric) error {
		m.isDelta = isDelta
		return nil
	}
}

func WithTimestamp(t time.Time) MetricOption {
	return func(m *Metric) error {
		m.timestamp = t
		m.timestampSet = true
		return nil
	}
}

func WithCurrentTime() MetricOption {
	return WithTimestamp(time.Now())
}
