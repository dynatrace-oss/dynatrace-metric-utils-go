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

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/dynatrace"
)

func main() {
	// these are usually created in the dynatrace exporter
	dims := []dynatrace.Dimension{dynatrace.NewDimension("key1", "value1"), dynatrace.NewDimension("key2", "value2")}
	oneAgentData := []dynatrace.Dimension{dynatrace.NewDimension("key1", "oneagentValue")}

	serializer := dynatrace.NewMetricSerializer(dims, oneAgentData)

	// dimensions could be created for each new instrument or context, and can be used to map behavior (e. g. failed requests...)
	otherDims := []dynatrace.Dimension{dynatrace.NewDimension("mykey1", "myvalue1"), dynatrace.NewDimension("key1", "myvalue1")}
	descriptor, _ := serializer.SerializeDescriptor("name", "prefix", otherDims)

	// this is a crude example of whatever data format the surrounding program stores data,
	sumVal := struct{ min, max, sum, count int64 }{1, 5, 10, 12}
	value := dynatrace.SerializeIntSummaryValue(sumVal.min, sumVal.max, sumVal.sum, sumVal.count)

	// prefix.name key1=oneagentValue,key2=value2,mykey1=myvalue1 gauge,min=1,max=5,sum=10,count=12
	fmt.Println(descriptor, value)
}
