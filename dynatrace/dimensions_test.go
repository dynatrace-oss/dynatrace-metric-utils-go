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
	"reflect"
	"testing"
)

func Test_insertNormalizedDimensions(t *testing.T) {
	type args struct {
		target map[string]string
		dims   []Dimension
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "valid add",
			args: args{
				target: make(map[string]string),
				dims:   []Dimension{NewDimension("dim1", "dv1"), NewDimension("dim2", "dv2")},
			},
			want: map[string]string{"dim1": "dv1", "dim2": "dv2"},
		},
		{
			name: "pass nil dims",
			args: args{
				target: make(map[string]string),
				dims:   nil,
			},
			want: map[string]string{},
		},
		{
			name: "pass nil map",
			args: args{
				target: nil,
				dims:   []Dimension{NewDimension("dim1", "dv1"), NewDimension("dim2", "dv2")},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insertNormalizedDimensions(tt.args.target, tt.args.dims)
			if got := tt.args.target; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insertNormalizedDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStaticDimensions(t *testing.T) {
	type args struct {
		tags         []Dimension
		oneAgentData []Dimension
	}
	tests := []struct {
		name string
		args args
		want StaticDimensions
	}{
		{
			name: "valid test",
			args: args{
				tags:         []Dimension{NewDimension("t1", "tv1")},
				oneAgentData: []Dimension{NewDimension("o1", "ov1")},
			},
			want: StaticDimensions{
				items: map[string]string{"t1": "tv1", "o1": "ov1"},
			},
		},
		{
			name: "overwriting test",
			args: args{
				tags:         []Dimension{NewDimension("t1", "tv1"), NewDimension("t2", "tv2")},
				oneAgentData: []Dimension{NewDimension("t2", "oneagent_overrides")},
			},
			want: StaticDimensions{
				items: map[string]string{"t1": "tv1", "t2": "oneagent_overrides"},
			},
		},
		{
			name: "pass nil tags test",
			args: args{
				tags:         nil,
				oneAgentData: []Dimension{NewDimension("o1", "ov1")},
			},
			want: StaticDimensions{
				items: map[string]string{"o1": "ov1"},
			},
		},
		{
			name: "pass nil oneAgentData",
			args: args{
				tags:         []Dimension{NewDimension("t1", "tv1")},
				oneAgentData: nil,
			},
			want: StaticDimensions{
				items: map[string]string{"t1": "tv1"},
			},
		},
		{
			name: "pass both nil",
			args: args{
				tags:         nil,
				oneAgentData: nil,
			},
			want: StaticDimensions{
				items: map[string]string{},
			},
		},
		{
			name: "test overwrite after normalization",
			args: args{
				tags:         []Dimension{NewDimension("t1", "tv1")},
				oneAgentData: []Dimension{NewDimension("~~t1", "ov1")},
			},
			want: StaticDimensions{
				items: map[string]string{"t1": "ov1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStaticDimensions(tt.args.tags, tt.args.oneAgentData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStaticDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaticDimensions_MakeUniqueDimensions(t *testing.T) {
	type fields struct {
		items map[string]string
	}
	type args struct {
		dims []Dimension
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name:   "valid test",
			fields: fields{items: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"}},
			args:   args{dims: []Dimension{NewDimension("newDim1", "dimVal1")}},
			want:   map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1", "newdim1": "dimVal1"},
		},
		{
			name:   "overwrite dimensions",
			fields: fields{items: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"}},
			args:   args{dims: []Dimension{NewDimension("staticdim1", "dimVal1")}},
			want:   map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
		},
		{
			name:   "pass nil",
			fields: fields{items: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"}},
			args:   args{dims: nil},
			want:   map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
		},
		{
			name:   "pass empty slice",
			fields: fields{items: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"}},
			args:   args{dims: []Dimension{}},
			want:   map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
		},
		{
			name:   "add to empty static dims",
			fields: fields{items: map[string]string{}},
			args:   args{dims: []Dimension{NewDimension("dim1", "value1"), NewDimension("dim2", "value2")}},
			want:   map[string]string{"dim1": "value1", "dim2": "value2"},
		},
		{
			name:   "add same key twice",
			fields: fields{items: map[string]string{}},
			args:   args{dims: []Dimension{NewDimension("dim1", "value1"), NewDimension("~~dim1", "value2")}},
			want:   map[string]string{"dim1": "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := StaticDimensions{
				items: tt.fields.items,
			}
			if got := sd.MakeUniqueDimensions(tt.args.dims); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StaticDimensions.MakeUniqueDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}
