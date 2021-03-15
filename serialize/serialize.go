// Copyright 2021 Dynatrace LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serialize

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/normalize"
)

func joinPrefix(name, prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s.%s", prefix, name)
	}
	return name
}

func MetricName(name, prefix string) (string, error) {
	return normalize.MetricKey(joinPrefix(name, prefix))
}

func formatDimensions(dims []dimensions.Dimension) string {
	var sb strings.Builder
	firstIter := true

	for _, dim := range dims {
		if firstIter {
			firstIter = false
		} else {
			sb.WriteString(",")
		}

		sb.WriteString(fmt.Sprintf("%s=%s", dim.Key, dim.Value))
	}

	return sb.String()
}

func NormalizedDimensions(dims dimensions.NormalizedDimensionSet) string {
	return dims.Format(formatDimensions)
}

// IntSummaryValue returns the value part of an metrics ingestion line for the given integers
func IntSummaryValue(min, max, sum, count int64) string {
	return fmt.Sprintf("gauge,min=%d,max=%d,sum=%d,count=%d", min, max, sum, count)
}

// IntCountValue transforms the integer given integer into a valid ingestion line value part.
func IntCountValue(value int64, absolute bool) string {
	if absolute {
		return fmt.Sprintf("count,delta=%d", value)
	}
	return fmt.Sprintf("count,%d", value)
}

// FloatSummaryValue returns the value part of an metrics ingestion line for the given floats, and an integer count
func FloatSummaryValue(min, max, sum float64, count int64) string {
	return fmt.Sprintf("gauge,min=%s,max=%s,sum=%s,count=%d", serializeFloat64(min), serializeFloat64(max), serializeFloat64(sum), count)
}

// FloatCountValue transforms the float given integer into a valid ingestion line value part.
func FloatCountValue(value float64, absolute bool) string {
	if absolute {
		return fmt.Sprintf("count,delta=%s", serializeFloat64(value))
	}
	return fmt.Sprintf("count,%s", serializeFloat64(value))
}

func IntGaugeValue(value int64) string {
	return fmt.Sprintf("gauge,%d", value)
}

func FloatGaugeValue(value float64) string {
	return fmt.Sprintf("gauge,%s", serializeFloat64(value))
}

func serializeFloat64(n float64) string {
	str := strings.TrimRight(strconv.FormatFloat(n, 'f', 6, 64), "0.")
	if str == "" {
		// if everything was trimmed away, number was 0.000000
		return "0"
	}
	return str
}

// Timestamp retruns the current timestamp as Unix time or an empty string if time has not been set.
func Timestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.Unix(), 10)
}
