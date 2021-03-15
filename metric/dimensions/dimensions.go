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
			keys[k] = struct{}{}
		}
	}

	return newNormalizedDimensionSet(normalizedDims)
}

// todo
func FromOneAgentMetadata() NormalizedDimensionSet {
	return NormalizeSet(NewDimensionSet())
}

// MergeSets combines one or more NormalizedDimensionSets into one. Dimensions in sets passed further to the right but containing the
// same keys as sets further to the left will overwrite the values. The resulting set contains no duplicate keys. If duplicate
// keys appear in different sets, the value of the resulting set will be the one from the last set passed to this function and
// containing the key. The order of keys is retained in the sense that keys seen first will also appear in the output set first,
// e. g. keys passed in the leftmost set will appear before keys in the rightmost set, while for values the opposite is true.
func MergeSets(dimensions NormalizedDimensionSet, overwritingDimensions ...NormalizedDimensionSet) NormalizedDimensionSet {
	allDimensions := append([]NormalizedDimensionSet{dimensions}, overwritingDimensions...)

	approxElements := 0
	for _, dims := range allDimensions {
		approxElements += len(dims.dimensions)
	}

	// when using the map type, the order of elements is no longer guaranteed.
	uniqueDimensions := make(map[string]string, approxElements)
	keyOrder := []string{}

	// iterate NormalizedDimensionSets
	for _, dims := range allDimensions {
		// iterate dimensions within the sets.
		for _, dim := range dims.dimensions {
			if _, ok := uniqueDimensions[dim.Key]; !ok {
				// key does not yet exist, so we remember the first time it showed up
				// to keep the keys in order.
				keyOrder = append(keyOrder, dim.Key)
			}
			// overwrite the dimension value with the last occurrence for each key
			uniqueDimensions[dim.Key] = dim.Value
		}
	}

	// create an ordered set of dimensions
	orderedDimensions := make([]Dimension, len(keyOrder))
	for i, key := range keyOrder {
		orderedDimensions[i] = NewDimension(key, uniqueDimensions[key])
	}

	// all dimensions must be normalized in order to enter this function, so we can be sure that the result is also normalized.
	return newNormalizedDimensionSet(orderedDimensions)
}
