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

package main

import (
	"fmt"
	"log"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
	"github.com/dynatrace-oss/dynatrace-metric-utils-go/oneagentenrichment"
)

func main() {
	// default and one agent dimensions will usually not change during in one exporter, so they can be
	// normalized once and reused.
	// these are the default values, they will be overwritten
	defaultDimensions := dimensions.CreateDimensionSet(
		dimensions.NewDimension("default1", "value1"),
		dimensions.NewDimension("dim1", "default1"),
		dimensions.NewDimension("dim2", "default2"),
	)

	// these are the oneAgent values. They will overwrite all other dimensions if they have the same key.
	// oneAgentDimensions := dimensions.FromOneAgentMetadata()
	oneAgentDimensions := oneagentenrichment.GetOneAgentMetadata()

	// these are labels, usually set by an instrument and therefore created in each exporter iteration.
	// it might also be possible to cache these (if the user is sure that they are the same on each export)
	labels := dimensions.CreateDimensionSet(
		dimensions.NewDimension("someLabel", "labelVal"),
		dimensions.NewDimension("dim1", "label1"),
	)

	merged := dimensions.MergeSets(defaultDimensions, labels, oneAgentDimensions)

	m, err := metric.NewMetric(
		"name",
		metric.WithPrefix("prefix"),
		metric.WithIntAbsoluteCounterValue(30),
		metric.WithDimensions(merged),
		metric.WithCurrentTime(),
	)
	if err != nil {
		log.Fatal(err)
	}

	serialized, err := m.Serialize()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(serialized)
}
