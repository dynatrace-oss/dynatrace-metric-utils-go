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

package serialize_test

// import (
// 	"reflect"
// 	"sort"
// 	"strings"
// 	"testing"
// )

// func TestNewStaticDimensions(t *testing.T) {
// 	type args struct {
// 		tags         []Dimension
// 		oneAgentData []Dimension
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want MetricSerializer
// 	}{
// 		{
// 			name: "valid test",
// 			args: args{
// 				tags:         []Dimension{NewDimension("t1", "tv1")},
// 				oneAgentData: []Dimension{NewDimension("o1", "ov1")},
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{"t1": "tv1"},
// 				overridingDimensions: map[string]string{"o1": "ov1"},
// 			},
// 		},
// 		{
// 			name: "pass nil tags test",
// 			args: args{
// 				tags:         nil,
// 				oneAgentData: []Dimension{NewDimension("o1", "ov1")},
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{},
// 				overridingDimensions: map[string]string{"o1": "ov1"},
// 			},
// 		},
// 		{
// 			name: "pass nil oneAgentData",
// 			args: args{
// 				tags:         []Dimension{NewDimension("t1", "tv1")},
// 				oneAgentData: nil,
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{"t1": "tv1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 		},
// 		{
// 			name: "pass both nil",iserialeze
// 			args: args{
// 				tags:         nil,
// 				oneAgentData: nil,
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{},
// 				overridingDimensions: map[string]string{},
// 			},
// 		},
// 		{
// 			name: "test normalization",
// 			args: args{
// 				tags:         []Dimension{NewDimension("t1", "tv1")},
// 				oneAgentData: []Dimension{NewDimension("~~t1", "ov1")},
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{"t1": "tv1"},
// 				overridingDimensions: map[string]string{"t1": "ov1"},
// 			},
// 		},
// 		{
// 			name: "invalid dimension key",
// 			args: args{
// 				tags:         []Dimension{NewDimension("t1", "tv1")},
// 				oneAgentData: []Dimension{NewDimension("~~~", "ov1")},
// 			},
// 			want: MetricSerializer{
// 				defaultDimensions:    map[string]string{"t1": "tv1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewMetricSerializer(tt.args.tags, tt.args.oneAgentData); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewMetricSerializer() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_joinPrefix(t *testing.T) {
// 	type args struct {
// 		name   string
// 		prefix string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "valid test",
// 			args: args{name: "name", prefix: "prefix"},
// 			want: "prefix.name",
// 		},
// 		{
// 			name: "no prefix",
// 			args: args{name: "name", prefix: ""},
// 			want: "name",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := joinPrefix(tt.args.name, tt.args.prefix); got != tt.want {
// 				t.Errorf("joinPrefix() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_serializeDimensions(t *testing.T) {
// 	type args struct {
// 		dims map[string]string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "passing empty map",
// 			args: args{dims: make(map[string]string)},
// 			want: "",
// 		},
// 		{
// 			name: "passing only one value",
// 			args: args{dims: map[string]string{"dim1": "val1"}},
// 			want: "dim1=val1",
// 		},
// 		{
// 			name: "passing two values",
// 			args: args{dims: map[string]string{"dim1": "val1", "dim2": "val2"}},
// 			want: "dim1=val1,dim2=val2",
// 		},
// 		{
// 			name: "passing more values",
// 			args: args{dims: map[string]string{"dim1": "val1", "dim2": "val2", "dim3": "val3", "dim4": "val4", "dim5": "val5"}},
// 			want: "dim1=val1,dim2=val2,dim3=val3,dim4=val4,dim5=val5",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			want := strings.Split(tt.want, ",")
// 			got := strings.Split(serializeDimensions(tt.args.dims), ",")
// 			sort.Strings(want)
// 			sort.Strings(got)

// 			if !reflect.DeepEqual(got, want) {
// 				t.Errorf("serializeDimensions() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMetricSerializer_SerializeDescriptor(t *testing.T) {
// 	type fields struct {
// 		staticDimensions map[string]string
// 	}
// 	type args struct {
// 		name   string
// 		prefix string
// 		dims   []Dimension
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "base test",
// 			fields: fields{staticDimensions: map[string]string{"sdkey1": "sdval1"}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{NewDimension("dk1", "dv1")}},
// 			want:   "metricPrefix.metricName sdkey1=sdval1,dk1=dv1",
// 		},
// 		{
// 			name:   "empty metrics",
// 			fields: fields{staticDimensions: map[string]string{}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{}},
// 			want:   "metricPrefix.metricName",
// 		},
// 		{
// 			name:   "test with key normalization",
// 			fields: fields{staticDimensions: map[string]string{"sdkey1": "sdval1"}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{NewDimension("~~~dk1", "dv1")}},
// 			want:   "metricPrefix.metricName sdkey1=sdval1,dk1=dv1",
// 		},
// 		{
// 			name:   "test with overwriting static dims",
// 			fields: fields{staticDimensions: map[string]string{"sdkey1": "sdval1"}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{NewDimension("sdkey1", "dv1")}},
// 			want:   "metricPrefix.metricName sdkey1=dv1",
// 		},
// 		{
// 			name:   "test only static dims",
// 			fields: fields{staticDimensions: map[string]string{"sdkey1": "sdval1"}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{}},
// 			want:   "metricPrefix.metricName sdkey1=sdval1",
// 		},
// 		{
// 			name:   "test only dynamic dims",
// 			fields: fields{staticDimensions: map[string]string{}},
// 			args:   args{name: "metricName", prefix: "metricPrefix", dims: []Dimension{NewDimension("dk1", "dv1")}},
// 			want:   "metricPrefix.metricName dk1=dv1",
// 		},
// 		{
// 			name:    "invalid prefix",
// 			fields:  fields{staticDimensions: map[string]string{}},
// 			args:    args{name: "metricName", prefix: "~~~", dims: []Dimension{NewDimension("dk1", "dv1")}},
// 			want:    "",
// 			wantErr: true,
// 		},
// 		{
// 			name:    "invalid name no prefix",
// 			fields:  fields{staticDimensions: map[string]string{}},
// 			args:    args{name: "~~~", prefix: "", dims: []Dimension{NewDimension("dk1", "dv1")}},
// 			want:    "",
// 			wantErr: true,
// 		},
// 		{
// 			name:   "invalid name with valid prefix",
// 			fields: fields{staticDimensions: map[string]string{}},
// 			args:   args{name: "~~~", prefix: "metricPrefix", dims: []Dimension{NewDimension("dk1", "dv1")}},
// 			want:   "metricPrefix dk1=dv1",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m := MetricSerializer{
// 				defaultDimensions: tt.fields.staticDimensions,
// 			}
// 			got, err := m.SerializeDescriptor(tt.args.name, tt.args.prefix, tt.args.dims)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("MetricSerializer.SerializeDescriptor() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("MetricSerializer.SerializeDescriptor() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMetricSerializer_makeUniqueDimensions(t *testing.T) {
// 	type fields struct {
// 		staticDimensions     map[string]string
// 		overridingDimensions map[string]string
// 	}
// 	type args struct {
// 		dims []Dimension
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   map[string]string
// 	}{
// 		{
// 			name: "valid test",
// 			fields: fields{
// 				staticDimensions:     map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{NewDimension("newDim1", "dimVal1")}},
// 			want: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1", "newdim1": "dimVal1"},
// 		},
// 		{
// 			name: "overwrite dimensions",
// 			fields: fields{
// 				staticDimensions:     map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{NewDimension("staticdim1", "dimVal1")}},
// 			want: map[string]string{"staticdim1": "dimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 		},
// 		{
// 			name: "overwrite dimensions with overriting dimensions",
// 			fields: fields{
// 				staticDimensions:     map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 				overridingDimensions: map[string]string{"staticdim1": "Overwritten"},
// 			},
// 			args: args{dims: []Dimension{NewDimension("staticdim1", "dimVal1")}},
// 			want: map[string]string{"staticdim1": "Overwritten", "oneagentdim1": "oneAgentDimVal1"},
// 		},
// 		{
// 			name: "pass nil",
// 			fields: fields{
// 				staticDimensions:     map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: nil},
// 			want: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 		},
// 		{
// 			name: "pass empty slice",
// 			fields: fields{
// 				staticDimensions:     map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{}},
// 			want: map[string]string{"staticdim1": "staticDimVal1", "oneagentdim1": "oneAgentDimVal1"},
// 		},
// 		{
// 			name: "add to empty static dims",
// 			fields: fields{
// 				staticDimensions:     map[string]string{},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{NewDimension("dim1", "value1"), NewDimension("dim2", "value2")}},
// 			want: map[string]string{"dim1": "value1", "dim2": "value2"},
// 		},
// 		{
// 			name: "add same key twice",
// 			fields: fields{
// 				staticDimensions:     map[string]string{},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{NewDimension("dim1", "value1"), NewDimension("~~dim1", "value2")}},
// 			want: map[string]string{"dim1": "value2"},
// 		},
// 		{
// 			name: "invalid dimension key",
// 			fields: fields{
// 				staticDimensions:     map[string]string{},
// 				overridingDimensions: map[string]string{},
// 			},
// 			args: args{dims: []Dimension{NewDimension("~~~", "dimValue")}},
// 			want: map[string]string{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := MetricSerializer{
// 				defaultDimensions:    tt.fields.staticDimensions,
// 				overridingDimensions: tt.fields.overridingDimensions,
// 			}
// 			if got := s.makeUniqueDimensions(tt.args.dims); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MetricSerializer.makeUniqueDimensions() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
