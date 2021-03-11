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

// MetricSerializer provides a way to set up static dimensions (for example OneAgent metadata)
type MetricSerializer struct {
	defaultDimensions    map[string]string
	overridingDimensions map[string]string
}

// Dimension is a KeyValue pair mapping string to string
type Dimension struct {
	Key   string
	Value string
}

// NewDimension creates a new KeyValue pair
func NewDimension(key, value string) Dimension {
	return Dimension{Key: key, Value: value}
}

// NewMetricSerializer creates a new MetricSerializer, which normalizes and then stores
// default and overwriting tags, so they dont have to be re-normalized for each normalize operation
func NewMetricSerializer(dimensions, oneAgentData []Dimension) MetricSerializer {
	statDims := normalizeDimensions(dimensions)
	overridingDimensions := normalizeDimensions(oneAgentData)

	return MetricSerializer{
		defaultDimensions:    statDims,
		overridingDimensions: overridingDimensions,
	}
}

func normalizeDimensions(dims []Dimension) map[string]string {
	items := make(map[string]string)
	if dims == nil {
		return items
	}

	for _, tag := range dims {
		normKey, err := normalize.DimensionKey(tag.Key)
		if err != nil {
			log.Printf("Could not parse '%s' as dimension key. Skipping... (Error: %s)", tag.Key, err)
			continue
		}
		items[normKey] = normalize.DimensionValue(tag.Value)
	}

	return items
}

// makeUniqueDimensions use the static dimensions prepared earlier to create a map of unique keys.
// Dimensions passed to this function will be overwritten by dimensions already stored in static
// dimensions.
func (m MetricSerializer) makeUniqueDimensions(dims []Dimension) map[string]string {
	items := make(map[string]string)
	// static dimensions are added first, these can be overwritten.
	for k, v := range m.defaultDimensions {
		items[k] = v
	}

	// then, the passed dimensions are normalized and added.
	for _, dim := range dims {
		normKey, err := normalize.DimensionKey(dim.Key)
		if err != nil {
			log.Printf("Could not parse '%s' as dimension key. Skipping... (Error: %s)", dim.Key, err)
			continue
		}
		items[normKey] = normalize.DimensionValue(dim.Value)
	}

	// finally, OneAgent dimensions overwrite already existing tags with the same name.
	for k, v := range m.overridingDimensions {
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

// SerializeDescriptor normalized the given name and prefix, and eliminates duplicate dimensions using the
// values stored in MetricSerializer. It returns a string of the concatenated name, prefix and dimensions.
func (m MetricSerializer) SerializeDescriptor(name, prefix string, dims []Dimension) (string, error) {
	metricKey, err := normalize.MetricKey(joinPrefix(name, prefix))
	if err != nil {
		return "", fmt.Errorf("error when normalizing metric key: %s", err)
	}
	dimsString := serializeDimensions(m.makeUniqueDimensions(dims))

	if dimsString != "" {
		return fmt.Sprintf("%s %s", metricKey, dimsString), nil
	}
	return metricKey, nil
}
