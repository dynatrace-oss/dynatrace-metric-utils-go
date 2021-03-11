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

package mint

import (
	"testing"
)

func TestSerializeIntSummaryValue(t *testing.T) {
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
			args: args{min: 0, max: 10, sum: 20, count: 3},
			want: "gauge,min=0,max=10,sum=20,count=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeIntSummaryValue(tt.args.min, tt.args.max, tt.args.sum, tt.args.count); got != tt.want {
				t.Errorf("IntSummaryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeIntCountValue(t *testing.T) {
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
			args: args{value: 30},
			want: "count,30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeIntCountValue(tt.args.value); got != tt.want {
				t.Errorf("IntCountValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeFloatSummaryValue(t *testing.T) {
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
			args: args{min: 0.5, max: 7.3, sum: 12.7, count: 3},
			want: "gauge,min=0.5,max=7.3,sum=12.7,count=3",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeFloatSummaryValue(tt.args.min, tt.args.max, tt.args.sum, tt.args.count); got != tt.want {
				t.Errorf("FloatSummaryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeFloatCountValue(t *testing.T) {
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
			args: args{value: 30.885},
			want: "count,30.885",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeFloatCountValue(tt.args.value); got != tt.want {
				t.Errorf("FloatCountValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeFloat64(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid float test",
			args: args{n: 3.141592},
			want: "3.141592",
		},
		{
			name: "rounded float test",
			args: args{n: 3.14159265359},
			want: "3.141593",
		},
		{
			name: "trim test",
			args: args{n: 2.500000},
			want: "2.5",
		},
		{
			name: "pass in zero test",
			args: args{n: 0.},
			want: "0",
		},
		{
			name: "truncate zero test",
			args: args{n: 0.0000},
			want: "0",
		},
		{
			name: "negative test",
			args: args{n: -10.24},
			want: "-10.24",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serializeFloat64(tt.args.n); got != tt.want {
				t.Errorf("serializeFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
