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
	"testing"
)

func TestMergeSets(t *testing.T) {
	type args struct {
		dimensions            NormalizedDimensionSet
		overwritingDimensions []NormalizedDimensionSet
	}
	tests := []struct {
		name string
		args args
		want NormalizedDimensionSet
	}{
		{
			name: "empty sets",
			args: args{dimensions: NormalizedDimensionSet{dimensions: []Dimension{}}, overwritingDimensions: []NormalizedDimensionSet{}},
			want: NormalizedDimensionSet{dimensions: []Dimension{}},
		},
		{
			name: "elements in first set",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")},
				},
				overwritingDimensions: []NormalizedDimensionSet{},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in second set",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in first and second set",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim1", "val1")},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim2", "val2")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements in first three sets",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim1", "val1")},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim2", "val2")}},
					{dimensions: []Dimension{NewDimension("dim3", "val3")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim1", "val1"), NewDimension("dim2", "val2"), NewDimension("dim3", "val3")}},
		},
		{
			name: "elements stay ordered",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim3", "val3")},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim2", "val2")}},
					{dimensions: []Dimension{NewDimension("dim1", "val1")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim3", "val3"), NewDimension("dim2", "val2"), NewDimension("dim1", "val1")}},
		},
		{
			name: "elements overwritten",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim1", "val1")},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim2", "val2")}},
					{dimensions: []Dimension{NewDimension("dim1", "val3")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim1", "val3"), NewDimension("dim2", "val2")}},
		},
		{
			name: "elements overwritten keep order",
			args: args{
				dimensions: NormalizedDimensionSet{
					dimensions: []Dimension{NewDimension("dim4", "val1")},
				},
				overwritingDimensions: []NormalizedDimensionSet{
					{dimensions: []Dimension{NewDimension("dim2", "val2")}},
					{dimensions: []Dimension{NewDimension("dim1", "val1")}},
					{dimensions: []Dimension{NewDimension("dim4", "val4")}},
				},
			},
			want: NormalizedDimensionSet{dimensions: []Dimension{NewDimension("dim4", "val4"), NewDimension("dim2", "val2"), NewDimension("dim1", "val1")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeSets(tt.args.dimensions, tt.args.overwritingDimensions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeSet(t *testing.T) {
	type args struct {
		dimset DimensionSet
	}
	tests := []struct {
		name string
		args args
		want NormalizedDimensionSet
	}{
		{
			name: "empty set",
			args: args{
				dimset: NewDimensionSet(),
			},
			want: NormalizedDimensionSet{[]Dimension{}},
		},
		{
			name: "non-colliding set",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key3", "value3"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
				NewDimension("key3", "value3"),
			}},
		},
		{
			name: "colliding set",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value3"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
			}},
		},
		{
			name: "retain order",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("key3", "value3"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value1"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{
				NewDimension("key3", "value3"),
				NewDimension("key2", "value2"),
				NewDimension("key1", "value1"),
			}},
		},
		{
			name: "normalized key retained",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("~~~key1", "value1"),
					NewDimension("key2", "value2"),
					NewDimension("key1", "value3"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key2", "value2"),
			}},
		},
		{
			name: "empty on invalid key",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("~!@$$", "value1"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{}},
		},
		{
			name: "empty on empty key",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("", "value1"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{}},
		},
		{
			name: "discard invalid key",
			args: args{
				dimset: NewDimensionSet(
					NewDimension("key1", "value1"),
					NewDimension("~!@$", "value2"),
					NewDimension("key3", "value3"),
				),
			},
			want: NormalizedDimensionSet{[]Dimension{
				NewDimension("key1", "value1"),
				NewDimension("key3", "value3"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeSet(tt.args.dimset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
