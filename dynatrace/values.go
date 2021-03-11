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

package dynatrace

import (
	"fmt"
	"strconv"
	"strings"
)

// SerializeIntSummaryValue returns the value part of an metrics ingestion line for the given integers
func SerializeIntSummaryValue(min, max, sum, count int64) string {
	return fmt.Sprintf("gauge,min=%d,max=%d,sum=%d,count=%d", min, max, sum, count)
}

// SerializeIntCountValue transforms the integer given integer into a valid ingestion line value part.
func SerializeIntCountValue(value int64) string {
	return fmt.Sprintf("count,%d", value)
}

// SerializeFloatSummaryValue returns the value part of an metrics ingestion line for the given floats, and an integer count
func SerializeFloatSummaryValue(min, max, sum float64, count int64) string {
	return fmt.Sprintf("gauge,min=%s,max=%s,sum=%s,count=%d", serializeFloat64(min), serializeFloat64(max), serializeFloat64(sum), count)
}

// SerializeFloatCountValue transforms the float given integer into a valid ingestion line value part.
func SerializeFloatCountValue(value float64) string {
	return fmt.Sprintf("count,%s", serializeFloat64(value))
}

func serializeFloat64(n float64) string {
	str := strings.TrimRight(strconv.FormatFloat(n, 'f', 6, 64), "0.")
	if str == "" {
		// if everything was trimmed away, number was 0.000000
		return "0"
	}
	return str
}
