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
	"reflect"
	"sort"
	"testing"
)

func TestCreateDimensionSet(t *testing.T) {
	type args struct {
		dims []Dimension
	}
	tests := []struct {
		name string
		args args
		want DimensionSet
	}{
		{
			name: "empty set",
			args: args{
				dims: []Dimension{},
			},
			want: DimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "non-colliding set",
			args: args{
				dims: []Dimension{
					NewDimension("key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key3", "value3"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
				NewDimension("key3", "value3"),
			}},
		},
		{
			name: "colliding set",
			args: args{
				dims: []Dimension{
					NewDimension("key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value3"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
				NewDimension("key1", "value3"),
			}},
		},
		{
			name: "retain order",
			args: args{
				dims: []Dimension{
					NewDimension("key3", "value3"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value1"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{
				NewDimension("key3", "value3"),
				NewDimension("key2", "value2"),
				NewDimension("key1", "value1"),
			}},
		},
		{
			name: "normalized key retained",
			args: args{
				dims: []Dimension{
					NewDimension("~~~key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value3"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
				NewDimension("key1", "value3"),
			}},
		},
		{
			name: "empty on invalid key",
			args: args{
				dims: []Dimension{
					NewDimension("~!@$$", "value1"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "empty on empty key",
			args: args{
				dims: []Dimension{
					NewDimension("", "value1"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "discard invalid key",
			args: args{
				dims: []Dimension{
					NewDimension("key1", "value1"),
					NewDimension("~!@$", "value2"),
					NewDimension("key3", "value3"),
				},
			},
			want: DimensionSet{dimensions: []Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key3", "value3"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDimensionSet(tt.args.dims...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDimensionSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeSets(t *testing.T) {
	type args struct {
		dimensionSets []DimensionSet
	}
	tests := []struct {
		name string
		args args
		want DimensionSet
	}{
		{
			name: "nothing passed",
			args: args{dimensionSets: []DimensionSet{}},
			want: DimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "empty sets",
			args: args{dimensionSets: []DimensionSet{DimensionSet{dimensions: []Dimension{}}, DimensionSet{dimensions: []Dimension{}}}},
			want: DimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "elements in first set",
			args: args{
				dimensionSets: []DimensionSet{
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
					DimensionSet{dimensions: []Dimension{}},
				},
			},
			want: DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in second set",
			args: args{
				dimensionSets: []DimensionSet{
					DimensionSet{dimensions: []Dimension{}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
				},
			},
			want: DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in first and second set",
			args: args{
				dimensionSets: []DimensionSet{
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1")}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim2", "val2")}},
				},
			},
			want: DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in first three sets",
			args: args{
				dimensionSets: []DimensionSet{
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1")}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim2", "val2")}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim3", "val3")}},
				},
			},
			want: DimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2"), NewDimension("dim3", "val3")}},
		},
		{
			name: "elements overwritten",
			args: args{
				dimensionSets: []DimensionSet{
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "default1"), NewDimension("dim2", "default2"), NewDimension("dim3", "default3")}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "label1"), NewDimension("dim2", "label2")}},
					DimensionSet{dimensions: []Dimension{NewDimension("dim1", "overwriting1")}},
				},
			},
			want: DimensionSet{dimensions: []Dimension{NewDimension("dim1", "overwriting1"), NewDimension("dim2", "label2"), NewDimension("dim3", "default3")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// sortedDeepEqual sorts the dimensions by key and then calls deepequal to make sure
			// the order does not matter (since we are using a map to deduplicate)
			if got := MergeSets(tt.args.dimensionSets...); !sortedDeepEqual(got, tt.want) {
				t.Errorf("MergeSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortedDeepEqual(got, want DimensionSet) bool {
	if len(got.dimensions) != len(want.dimensions) {
		return false
	}

	gotDims := got.dimensions
	wantDims := want.dimensions

	sort.SliceStable(gotDims, func(i, j int) bool {
		return gotDims[i].Key < gotDims[j].Key
	})
	sort.SliceStable(wantDims, func(i, j int) bool {
		return wantDims[i].Key < wantDims[j].Key
	})

	return reflect.DeepEqual(gotDims, wantDims)
}
