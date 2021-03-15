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
		dims dimensions.NormalizedDimensionSet
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no dimensions",
			args: args{dims: dimensions.NormalizeSet(dimensions.NewDimensionSet())},
			want: "",
		},
		{
			name: "one dimension",
			args: args{dims: dimensions.NormalizeSet(dimensions.NewDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
			))},
			want: "dim1=val1",
		},
		{
			name: "two dimensions",
			args: args{dims: dimensions.NormalizeSet(dimensions.NewDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
				dimensions.NewDimension("dim2", "val2"),
			))},
			want: "dim1=val1,dim2=val2",
		},
		{
			name: "five dimensions",
			args: args{dims: dimensions.NormalizeSet(dimensions.NewDimensionSet(
				dimensions.NewDimension("dim1", "val1"),
				dimensions.NewDimension("dim2", "val2"),
				dimensions.NewDimension("dim3", "val3"),
				dimensions.NewDimension("dim4", "val4"),
				dimensions.NewDimension("dim5", "val5"),
			))},
			want: "dim1=val1,dim2=val2,dim3=val3,dim4=val4,dim5=val5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serialize.NormalizedDimensions(tt.args.dims); got != tt.want {
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
