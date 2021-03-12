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
	"strings"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/dimensions"
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

// func serializeNormalizedDimensions(dims dimensions.NormalizedDimensionSet) string {
// 	var sb strings.Builder

// 	firstIteration := true
// 	for _, dim := range dims.Dimensions {
// 		if !firstIteration {
// 			sb.WriteString(",")
// 		} else {
// 			firstIteration = false
// 		}
// 		sb.WriteString(fmt.Sprintf("%s=%s", dim.Key, dim.Value))
// 	}

// 	return sb.String()
// }

// // SerializeDescriptor normalized the given name and prefix, and eliminates duplicate dimensions using the
// // values stored in MetricSerializer. It returns a string of the concatenated name, prefix and dimensions.
// func Descriptor(name, prefix string, dims dimensions.NormalizedDimensionSet) (string, error) {
// 	metricKey, err := normalize.MetricKey(joinPrefix(name, prefix))
// 	if err != nil {
// 		return "", fmt.Errorf("error when normalizing metric key: %s", err)
// 	}
// 	dimsString := serializeNormalizedDimensions(dims)

// 	if dimsString != "" {
// 		return fmt.Sprintf("%s %s", metricKey, dimsString), nil
// 	}
// 	return metricKey, nil
// }
