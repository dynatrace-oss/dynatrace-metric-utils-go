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
	"log"
	"strings"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/dynatrace/normalize"
)

type MetricSerializer struct {
	staticDimensions map[string]string
}

type Dimension struct {
	Key   string
	Value string
}

func NewDimension(key, value string) Dimension {
	return Dimension{Key: key, Value: value}
}

func NewMetricSerializer(tags, oneAgentData []Dimension) MetricSerializer {
	items := make(map[string]string)

	// later calls will overwrite tags. We can ignore errors here,
	// since the map is initialized in the line above.
	insertNormalizedDimensions(items, tags)
	insertNormalizedDimensions(items, oneAgentData)

	return MetricSerializer{staticDimensions: items}
}

func insertNormalizedDimensions(target map[string]string, dims []Dimension) {
	if dims == nil || target == nil {
		return
	}

	for _, tag := range dims {
		normKey, err := normalize.DimensionKey(tag.Key)
		if err != nil {
			log.Printf("Could not parse '%s' as dimension key. Skipping... (Error: %s)", tag.Key, err)
			continue
		}
		target[normKey] = normalize.DimensionValue(tag.Value)
	}
}

//MakeUniqueDimensions use the static dimensions prepared earlier to create a map of unique keys.
// Dimensions passed to this function will be overwritten by dimensions already stored in static
// dimensions.
func (s MetricSerializer) makeUniqueDimensions(dims []Dimension) map[string]string {
	items := make(map[string]string)

	// insert the dimensions passed to this function. these will be overwritten by static dimensions
	insertNormalizedDimensions(items, dims)

	// add static dimensions
	for k, v := range s.staticDimensions {
		items[k] = v
	}

	return items
}

func joinPrefix(name, prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s.%s", prefix, name)
	}
	return name
}

func serializeDimensions(dims map[string]string) string {
	var sb strings.Builder

	firstIteration := true
	for k, v := range dims {
		if !firstIteration {
			sb.WriteString(",")
		} else {
			firstIteration = false
		}
		sb.WriteString(fmt.Sprintf("%s=%s", k, v))
	}

	return sb.String()
}

func (m MetricSerializer) SerializeDescriptor(name, prefix string, dims []Dimension) (string, error) {
	metricKey, err := normalize.MetricKey(joinPrefix(name, prefix))
	if err != nil {
		return "", fmt.Errorf("error when normalizing metric key: %s", err)
	}
	dimsString := serializeDimensions(m.makeUniqueDimensions(dims))

	return fmt.Sprintf("%s %s", metricKey, dimsString), nil
}
