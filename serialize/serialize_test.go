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

import (
	"testing"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/serialize"
)

func TestMetricName(t *testing.T) {
	type args struct {
		name   string
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid name no prefix",
			args: args{name: "name", prefix: ""},
			want: "name",
		},
		{
			name: "no name valid prefix",
			args: args{name: "", prefix: "prefix"},
			want: "prefix",
		},
		{
			name: "valid name valid prefix",
			args: args{name: "name", prefix: "prefix"},
			want: "prefix.name",
		},
		{
			name:    "no name no prefix",
			args:    args{name: "", prefix: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid name no prefix",
			args:    args{name: "~~~", prefix: ""},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid name valid prefix",
			args: args{name: "~~~", prefix: "prefix"},
			want: "prefix",
		},
		{
			name:    "valid name invalid prefix",
			args:    args{name: "name", prefix: "~~~"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid name invalid prefix",
			args:    args{name: "~~~", prefix: "~~~"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serialize.MetricName(tt.args.name, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizedDimensions(t *testing.T) {
	type args struct {
		dims dimensions.DimensionSet
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no dimensions",
			args: args{dims: dimensions.CreateDimensionSet()},
			want: "",
		},
		{
			name: "one dimension",
			args: args{dims: dimensions.CreateDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
			)},
			want: "dim1=val1",
		},
		{
			name: "two dimensions",
			args: args{dims: dimensions.CreateDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
				dimensions.NewDimension("dim2", "val2"),
			)},
			want: "dim1=val1,dim2=val2",
		},
		{
			name: "five dimensions",
			args: args{dims: dimensions.CreateDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
				dimensions.NewDimension("dim2", "val2"),
				dimensions.NewDimension("dim3", "val3"),
				dimensions.NewDimension("dim4", "val4"),
				dimensions.NewDimension("dim5", "val5"),
			)},
			want: "dim1=val1,dim2=val2,dim3=val3,dim4=val4,dim5=val5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.Dimensions(tt.args.dims); got != tt.want {
				t.Errorf("NormalizedDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid time",
			args: args{t: time.Unix(1615800000, 0)},
			want: "1615800000",
		},
		{
			name: "empty time",
			args: args{ /* using the time.Time zero value if nothing is specified */ },
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.Timestamp(tt.args.t); got != tt.want {
				t.Errorf("Timestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSummaryValue(t *testing.T) {
	type args struct {
		min   int64
		max   int64
		sum   int64
		count int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid test",
			args: args{min: 0, max: 10, sum: 30, count: 7},
			want: "gauge,min=0,max=10,sum=30,count=7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.IntSummaryValue(tt.args.min, tt.args.max, tt.args.sum, tt.args.count); got != tt.want {
				t.Errorf("IntSummaryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntCountValue(t *testing.T) {
	type args struct {
		value    int64
		absolute bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "monotonic counter",
			args: args{value: 300, absolute: false},
			want: "count,300",
		},
		{
			name: "absolute counter",
			args: args{value: 300, absolute: true},
			want: "count,delta=300",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.IntCountValue(tt.args.value, tt.args.absolute); got != tt.want {
				t.Errorf("IntCountValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatSummaryValue(t *testing.T) {
	type args struct {
		min   float64
		max   float64
		sum   float64
		count int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid test",
			args: args{min: 0.3, max: 10.5, sum: 30.7, count: 7},
			want: "gauge,min=0.3,max=10.5,sum=30.7,count=7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.FloatSummaryValue(tt.args.min, tt.args.max, tt.args.sum, tt.args.count); got != tt.want {
				t.Errorf("FloatSummaryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatCountValue(t *testing.T) {
	type args struct {
		value    float64
		absolute bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "monotonic counter",
			args: args{value: 300.456, absolute: false},
			want: "count,300.456",
		},
		{
			name: "absolute counter",
			args: args{value: 300.456, absolute: true},
			want: "count,delta=300.456",
		},
		{
			name: "rounded monotonic counter",
			args: args{value: 300.123456789, absolute: false},
			want: "count,300.123457",
		},
		{
			name: "rounded absolute counter",
			args: args{value: 300.123456789, absolute: true},
			want: "count,delta=300.123457",
		},
		{
			name: "zero value monotonic counter",
			args: args{value: 0.0000, absolute: false},
			want: "count,0",
		},
		{
			name: "rounded absolute counter",
			args: args{value: 0.0000, absolute: true},
			want: "count,delta=0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.FloatCountValue(tt.args.value, tt.args.absolute); got != tt.want {
				t.Errorf("FloatCountValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntGaugeValue(t *testing.T) {
	type args struct {
		value int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid test",
			args: args{value: 3},
			want: "gauge,3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.IntGaugeValue(tt.args.value); got != tt.want {
				t.Errorf("IntGaugeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatGaugeValue(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid test",
			args: args{value: 3},
			want: "gauge,3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.FloatGaugeValue(tt.args.value); got != tt.want {
				t.Errorf("FloatGaugeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
