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
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
)

func TestMetric_Serialize(t *testing.T) {
	type fields struct {
		name         string
		prefix       string
		value        metricValue
		dimensions   dimensions.NormalizedDimensionList
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
				value: intCounterValue{value: 123, isDelta: false},
			},
			want: "name count,123",
		},
		{
			name: "invalid missing name",
			fields: fields{
				value: intCounterValue{value: 123, isDelta: false},
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
				value: intCounterValue{value: 123, isDelta: false},
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
				value:     intCounterValue{value: 123, isDelta: false},
				timestamp: time.Unix(1615800000, 123000000),
			},
			want: "name count,123 1615800000123",
		},
		{
			name: "with dimensions",
			fields: fields{
				name:       "name",
				value:      intCounterValue{value: 123, isDelta: false},
				dimensions: dimensions.NewNormalizedDimensionList(dimensions.NewDimension("key1", "value1"), dimensions.NewDimension("key2", "value2")),
			},
			want: "name,key1=value1,key2=value2 count,123",
		},
		{
			name: "with timestamp and dimensions",
			fields: fields{
				name:       "name",
				value:      intCounterValue{value: 123, isDelta: false},
				timestamp:  time.Unix(1615800000, 123000000),
				dimensions: dimensions.NewNormalizedDimensionList(dimensions.NewDimension("key1", "value1"), dimensions.NewDimension("key2", "value2")),
			},
			want: "name,key1=value1,key2=value2 count,123 1615800000123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{
				metricKey:  tt.fields.name,
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
		metricKey string
		options   []MetricOption
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
			args:    args{metricKey: "name", options: []MetricOption{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "just value",
			args: args{metricKey: "", options: []MetricOption{
				WithIntCounterValueTotal(3),
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "name and value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
			}},
			want: &Metric{metricKey: "name", value: intCounterValue{value: 3, isDelta: false}},
		},
		{
			name: "with prefix",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
				WithPrefix("prefix"),
			}},
			want: &Metric{metricKey: "name", prefix: "prefix", value: intCounterValue{value: 3, isDelta: false}},
		},
		{
			name: "name and monotonic int counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
			}},
			want: &Metric{metricKey: "name", value: intCounterValue{value: 3, isDelta: false}},
		},
		{
			name: "name and absolute int counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueDelta(3),
			}},
			want: &Metric{metricKey: "name", value: intCounterValue{value: 3, isDelta: true}},
		},
		{
			name: "name and monotonic float counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueTotal(3.1415),
			}},
			want: &Metric{metricKey: "name", value: floatCounterValue{value: 3.1415, isDelta: false}},
		},
		{
			name: "name and monotonic float counter value NaN",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueTotal(math.NaN()),
			}},
			wantErr: true,
		},
		{
			name: "name and monotonic float counter value negative infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueTotal(math.Inf(-1)),
			}},
			wantErr: true,
		},
		{
			name: "name and monotonic float counter value infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueTotal(math.Inf(1)),
			}},
			wantErr: true,
		},
		{
			name: "name and absolute float counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueDelta(3.1415),
			}},
			want: &Metric{metricKey: "name", value: floatCounterValue{value: 3.1415, isDelta: true}},
		},
		{
			name: "name and absolute float counter value NaN",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueDelta(math.NaN()),
			}},
			wantErr: true,
		},
		{
			name: "name and absolute float counter value negative infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueDelta(math.Inf(-1)),
			}},
			wantErr: true,
		},
		{
			name: "name and absolute float counter value positive infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueDelta(math.Inf(1)),
			}},
			wantErr: true,
		},
		{
			name: "name and int summary value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntSummaryValue(0, 10, 25, 7),
			}},
			want: &Metric{metricKey: "name", value: intSummaryValue{min: 0, max: 10, sum: 25, count: 7}},
		},
		{
			name: "name and float summary value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(0.4, 10.87, 25.4, 7),
			}},
			want: &Metric{metricKey: "name", value: floatSummaryValue{min: 0.4, max: 10.87, sum: 25.4, count: 7}},
		},
		{
			name: "name and int gauge value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntGaugeValue(7),
			}},
			want: &Metric{metricKey: "name", value: intGaugeValue{value: 7}},
		},
		{
			name: "name and float gauge value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatGaugeValue(7.34),
			}},
			want: &Metric{metricKey: "name", value: floatGaugeValue{value: 7.34}},
		},
		{
			name: "name and float gauge value NaN",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatGaugeValue(math.NaN()),
			}},
			wantErr: true,
		},
		{
			name: "name and float gauge value negative infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatGaugeValue(math.Inf(-1)),
			}},
			wantErr: true,
		},
		{
			name: "name and float gauge value positive infinity",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatGaugeValue(math.Inf(1)),
			}},
			wantErr: true,
		},
		{
			name: "invalid monotonic int counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(-3),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid absolute int counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueDelta(-3),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid monotonic float counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueTotal(-3.1415),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid absolute float counter value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatCounterValueDelta(-3.1415),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid int summary value",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntSummaryValue(0, 10, 25, -7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid int summary value 2",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntSummaryValue(10, 2, 25, 7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid float summary value",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(0.4, 10.87, 25.4, -7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid float summary value 2",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(10.3, 1.87, 25.4, 7),
			}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid float summary value NaN 1",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(math.NaN(), 1.87, 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value NaN 2",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, math.NaN(), 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value NaN 3",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, 2.34, math.NaN(), 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value negative infinity 1",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(math.Inf(-1), 1.87, 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value negative infinity 2",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, math.Inf(-1), 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value negative infinity 3",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, 2.34, math.Inf(-1), 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value positive infinity 1",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(math.Inf(1), 1.87, 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value positive infinity 2",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, math.Inf(1), 25.4, 7),
			}},
			wantErr: true,
		},
		{
			name: "invalid float summary value positive infinity 3",
			args: args{metricKey: "name", options: []MetricOption{
				WithFloatSummaryValue(1.87, 2.34, math.Inf(1), 7),
			}},
			wantErr: true,
		},
		{
			name: "error on adding two values",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
				WithIntCounterValueTotal(5),
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test with timestamp",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
				WithTimestamp(time.Unix(1615800000, 0)),
			}},
			want: &Metric{metricKey: "name", value: intCounterValue{value: 3, isDelta: false}, timestamp: time.Unix(1615800000, 0)},
		},
		{
			name: "test with timestamp",
			args: args{metricKey: "name", options: []MetricOption{
				WithIntCounterValueTotal(3),
				WithDimensions(dimensions.NewNormalizedDimensionList(dimensions.NewDimension("key1", "value1"))),
			}},
			want: &Metric{
				metricKey: "name", value: intCounterValue{value: 3, isDelta: false},
				dimensions: dimensions.NewNormalizedDimensionList(dimensions.NewDimension("key1", "value1")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.metricKey, tt.args.options...)
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

func TestWithCurrentTime(t *testing.T) {
	m, err := NewMetric(
		"name",
		WithIntGaugeValue(1),
		WithCurrentTime(),
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	serialized, err := m.Serialize()
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedLengthDummy := "name gauge,1 1617294350925"
	expectedStart := "name gauge,1 "

	if !strings.HasPrefix(serialized, expectedStart) {
		t.Errorf("metric does not start with " + expectedStart)
	}

	if len(serialized) < len(expectedLengthDummy) {
		t.Errorf("serialized metric is too short")
	}
}
