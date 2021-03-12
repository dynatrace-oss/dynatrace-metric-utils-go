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
	values       []metricValue
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

func (m Metric) Serialize() ([]string, error) {
	lines := []string{}
	keyString, err := serialize.MetricName(m.name, m.prefix)
	if err != nil {
		return lines, err
	}

	dimString := serialize.NormalizedDimensions(m.dimensions)

	timeString := ""
	if m.timestampSet {
		timeString = strconv.FormatInt(m.timestamp.Unix(), 10)
	}

	for _, value := range m.values {
		line, err := joinStrings(keyString, dimString, value.serialize(), timeString)
		if err != nil {
			fmt.Println(err)
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func NewMetric(name string, options ...MetricOption) (*Metric, error) {
	m := &Metric{
		name:         name,
		values:       []metricValue{},
		timestampSet: false,
	}

	for _, option := range options {
		err := option(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Metric) AddOption(option MetricOption) error {
	err := option(m)
	if err != nil {
		return err
	}
	return nil
}

type MetricOption func(m *Metric) error

func checkValueSet(m *Metric) error {
	if m.values != nil {
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
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}
		m.values = append(m.values, intCounterValue{value: val})

		return nil
	}
}

func WithFloatCounterValue(val float64) MetricOption {
	return func(m *Metric) error {
		if val < 0 {
			return fmt.Errorf("value must be greater than 0, was %v", val)
		}
		m.values = append(m.values, floatCounterValue{value: val})

		return nil
	}
}

func WithIntSummaryValue(min, max, sum, count int64) MetricOption {
	return func(m *Metric) error {
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		m.values = append(m.values, intSummaryValue{min: min, max: max, sum: sum, count: count})

		return nil
	}
}

func WithFloatSummaryValue(min, max, sum float64, count int64) MetricOption {
	return func(m *Metric) error {
		if count < 0 {
			return fmt.Errorf("count cannot be smaller than 0, was %v", count)
		}
		m.values = append(m.values, floatSummaryValue{min: min, max: max, sum: sum, count: count})

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
