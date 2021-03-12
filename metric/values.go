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

import "github.com/dynatrace-oss/dynatrace-metric-utils-go/serialize"

type metricValue interface {
	serialize() string
}

type intCounterValue struct {
	value int64
}

func (i intCounterValue) serialize() string {
	return serialize.SerializeIntCountValue(i.value)
}

type floatCounterValue struct {
	value float64
}

func (f floatCounterValue) serialize() string {
	return serialize.SerializeFloatCountValue(f.value)
}

type intSummaryValue struct {
	min, max, sum, count int64
}

func (i intSummaryValue) serialize() string {
	return serialize.SerializeIntSummaryValue(i.min, i.max, i.sum, i.count)
}

type floatSummaryValue struct {
	min, max, sum float64
	count         int64
}

func (f floatSummaryValue) serialize() string {
	return serialize.SerializeFloatSummaryValue(f.min, f.max, f.sum, f.count)
}
