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

package metric

import (
	"reflect"
	"testing"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
)

func TestMetric_Serialize(t *testing.T) {
	type fields struct {
		name         string
		prefix       string
		value        metricValue
		dimensions   dimensions.DimensionSet
		isDelta      bool
		timestamp    time.Time
		timestampSet bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "valid required items",
			fields: fields{
				name:  "name",
				value: intCounterValue{value: 123, absolute: false},
			},
			want: "name count,123",
		},
		{
			name: "invalid missing name",
			fields: fields{
				value: intCounterValue{value: 123, absolute: false},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid missing value",
			fields: fields{
				name: "name",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid empty name",
			fields: fields{
				name:  "",
				value: intCounterValue{value: 123, absolute: false},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty value uses zero values",
			fields: fields{
				name:  "name",
				value: intCounterValue{},
			},
			want: "name count,0",
		},
		{
			name: "with timestamp",
			fields: fields{
				name:      "name",
				value:     intCounterValue{value: 123, absolute: false},
				timestamp: time.Unix(1615800000, 0),
			},
			want: "name count,123 1615800000",
		},
		{
			name: "with dimensions",
			fields: fields{
				name:       "name",
				value:      intCounterValue{value: 123, absolute: false},
				dimensions: dimensions.CreateDimensionSet(dimensions.NewDimension("key1", "value1"), dimensions.NewDimension("key2", "value2")),
			},
			want: "name,key1=value1,key2=value2 count,123",
		},
		{
			name: "with timestamp and dimensions",
			fields: fields{
				name:       "name",
				value:      intCounterValue{value: 123, absolute: false},
				timestamp:  time.Unix(1615800000, 0),
				dimensions: dimensions.CreateDimensionSet(dimensions.NewDimension("key1", "value1"), dimensions.NewDimension("key2", "value2")),
			},
			want: "name,key1=value1,key2=value2 count,123 1615800000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{
				name:       tt.fields.name,
				prefix:     tt.fields.prefix,
				value:      tt.fields.value,
				dimensions: tt.fields.dimensions,
				timestamp:  tt.fields.timestamp,
			}
			got, err := m.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Metric.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetric(t *testing.T) {
	type args struct {
		name    string
		options []MetricOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Metric
		wantErr bool
	}{
		{
			name:    "no options",
			args:    args{options: []MetricOption{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "just name",
			args:    args{name: "name", options: []MetricOption{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "just value",
			args: args{name: "", options: []MetricOption{
				WithIntCounterValue(3),
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "name and value",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
			}},
			want: &Metric{name: "name", value: intCounterValue{value: 3, absolute: false}},
		},
		{
			name: "with prefix",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
				WithPrefix("prefix"),
			}},
			want: &Metric{name: "name", prefix: "prefix", value: intCounterValue{value: 3, absolute: false}},
		},
		{
			name: "name and monotonic int counter value",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
			}},
			want: &Metric{name: "name", value: intCounterValue{value: 3, absolute: false}},
		},
		{
			name: "name and absolute int counter value",
			args: args{name: "name", options: []MetricOption{
				WithIntAbsoluteCounterValue(3),
			}},
			want: &Metric{name: "name", value: intCounterValue{value: 3, absolute: true}},
		},
		{
			name: "name and monotonic float counter value",
			args: args{name: "name", options: []MetricOption{
				WithFloatCounterValue(3.1415),
			}},
			want: &Metric{name: "name", value: floatCounterValue{value: 3.1415, absolute: false}},
		},
		{
			name: "name and absolute float counter value",
			args: args{name: "name", options: []MetricOption{
				WithFloatAbsoluteCounterValue(3.1415),
			}},
			want: &Metric{name: "name", value: floatCounterValue{value: 3.1415, absolute: true}},
		},
		{
			name: "name and int summary value",
			args: args{name: "name", options: []MetricOption{
				WithIntSummaryValue(0, 10, 25, 7),
			}},
			want: &Metric{name: "name", value: intSummaryValue{min: 0, max: 10, sum: 25, count: 7}},
		},
		{
			name: "name and float summary value",
			args: args{name: "name", options: []MetricOption{
				WithFloatSummaryValue(0.4, 10.87, 25.4, 7),
			}},
			want: &Metric{name: "name", value: floatSummaryValue{min: 0.4, max: 10.87, sum: 25.4, count: 7}},
		},
		{
			name: "invalid monotonic int counter value",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(-3),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid absolute int counter value",
			args: args{name: "name", options: []MetricOption{
				WithIntAbsoluteCounterValue(-3),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid monotonic float counter value",
			args: args{name: "name", options: []MetricOption{
				WithFloatCounterValue(-3.1415),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid absolute float counter value",
			args: args{name: "name", options: []MetricOption{
				WithFloatAbsoluteCounterValue(-3.1415),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid int summary value",
			args: args{name: "name", options: []MetricOption{
				WithIntSummaryValue(0, 10, 25, -7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid int summary value 2",
			args: args{name: "name", options: []MetricOption{
				WithIntSummaryValue(0, 100, 25, 7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid float summary value",
			args: args{name: "name", options: []MetricOption{
				WithFloatSummaryValue(0.4, 10.87, 25.4, -7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid float summary value 2",
			args: args{name: "name", options: []MetricOption{
				WithFloatSummaryValue(0.4, 100.87, 25.4, 7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error on adding two values",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
				WithIntCounterValue(5),
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test with timestamp",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
				WithTimestamp(time.Unix(1615800000, 0)),
			}},
			want: &Metric{name: "name", value: intCounterValue{value: 3, absolute: false}, timestamp: time.Unix(1615800000, 0)},
		},
		{
			name: "test with timestamp",
			args: args{name: "name", options: []MetricOption{
				WithIntCounterValue(3),
				WithDimensions(dimensions.CreateDimensionSet(dimensions.NewDimension("key1", "value1"))),
			}},
			want: &Metric{
				name: "name", value: intCounterValue{value: 3, absolute: false},
				dimensions: dimensions.CreateDimensionSet(dimensions.NewDimension("key1", "value1")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.name, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
