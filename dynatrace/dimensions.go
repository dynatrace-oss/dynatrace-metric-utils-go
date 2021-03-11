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
	"log"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/dynatrace/normalize"
)

type StaticDimensions struct {
	items map[string]string
}

type Dimension struct {
	Key   string
	Value string
}

func NewDimension(key, value string) Dimension {
	return Dimension{Key: key, Value: value}
}

func NewStaticDimensions(tags, oneAgentData []Dimension) StaticDimensions {
	items := make(map[string]string)

	// later calls will overwrite tags.
	insertNormalizedDimensions(items, tags)
	insertNormalizedDimensions(items, oneAgentData)

	return StaticDimensions{items: items}
}

func insertNormalizedDimensions(target map[string]string, dims []Dimension) {
	if dims == nil {
		return
	}

	if target == nil {
		target = make(map[string]string)
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
func (sd StaticDimensions) MakeUniqueDimensions(dims []Dimension) map[string]string {
	items := make(map[string]string)

	// insert the dimensions passed to this function. these will be overwritten by static dimensions
	insertNormalizedDimensions(items, dims)

	for k, v := range sd.items {
		items[k] = v
	}

	return items
}
