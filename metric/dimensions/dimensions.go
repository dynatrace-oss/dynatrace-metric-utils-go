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

package dimensions

import (
	"log"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/normalize"
)

type Dimension struct {
	Key   string
	Value string
}

type DimensionSet struct {
	dimensions []Dimension
}

// NormalizedDimensionSet contains normalized keys and values and no duplicate keys.
type NormalizedDimensionSet struct {
	dimensions []Dimension
}

// pass a function that transforms a slice of dimensions to a string.
// That way, the code for the actual serialization can be stored in the
// serialization package without exporting the dimensions in normalized dimensions
// which in turn restricts manipulation of already normalized values.
func (nds NormalizedDimensionSet) Format(formatter func([]Dimension) string) string {
	return formatter(nds.dimensions)
}

func newNormalizedDimensionSet(dims []Dimension) NormalizedDimensionSet {
	return NormalizedDimensionSet{dimensions: dims}
}

func NewDimension(key, val string) Dimension {
	return Dimension{Key: key, Value: val}
}

func NewDimensionSet(dims ...Dimension) DimensionSet {
	return DimensionSet{dims}
}

// NormalizeSet normalizes all keys and values in the passed DimensionSet.
// If keys collide (after normalization of each of the keys), the first one is retained.
// The order of the passed dimensions is retained (except for removed duplicates).
func NormalizeSet(dimset DimensionSet) NormalizedDimensionSet {
	// this is basically a set, but golang does not offer a set type
	keys := make(map[string]struct{}, len(dimset.dimensions))
	normalizedDims := []Dimension{}

	for _, dim := range dimset.dimensions {
		k, err := normalize.DimensionKey(dim.Key)
		if err != nil {
			log.Printf("normalization for '%s' returned invalid key. Skipping...", dim.Key)
			continue
		}

		// check if the normalized key has been seen yet and execute if it has not.
		// this means, that the first appearance of a key in the dimensions will be added.
		if _, ok := keys[k]; !ok {
			normalizedDims = append(normalizedDims, NewDimension(k, normalize.DimensionValue(dim.Value)))
		}
	}

	return newNormalizedDimensionSet(normalizedDims)
}

func FromOneAgentMetadata() NormalizedDimensionSet {
	return NormalizeSet(NewDimensionSet())
}

func MergeSets(dimensions NormalizedDimensionSet, overwritingDimensions ...NormalizedDimensionSet) NormalizedDimensionSet {
	approxElements := len(dimensions.dimensions)

	for _, dims := range overwritingDimensions {
		approxElements += len(dims.dimensions)
	}

	// when using the map type, the order of elements is no longer guaranteed.
	uniqueDimensions := make(map[string]string, approxElements)
	keyOrder := []string{}

	for _, dim := range dimensions.dimensions {
		if _, ok := uniqueDimensions[dim.Key]; !ok {
			// key does not yet exist
			keyOrder = append(keyOrder, dim.Key)
		}
		uniqueDimensions[dim.Key] = dim.Value
	}

	for _, dims := range overwritingDimensions {
		for _, dim := range dims.dimensions {
			if _, ok := uniqueDimensions[dim.Key]; !ok {
				keyOrder = append(keyOrder, dim.Key)
			}
			uniqueDimensions[dim.Key] = dim.Value
		}
	}

	orderedDimensions := make([]Dimension, len(uniqueDimensions))
	for i, key := range keyOrder {
		orderedDimensions[i] = NewDimension(key, uniqueDimensions[key])
	}

	return newNormalizedDimensionSet(orderedDimensions)
}
