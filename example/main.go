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
	dims := []dynatrace.Dimension{dynatrace.NewDimension("key1", "value1"), dynatrace.NewDimension("key2", "value2")}
	oneAgentData := []dynatrace.Dimension{dynatrace.NewDimension("key1", "oneagentValue")}

	serializer := dynatrace.NewMetricSerializer(dims, oneAgentData)

	otherDims := []dynatrace.Dimension{dynatrace.NewDimension("mykey1", "myvalue1"), dynatrace.NewDimension("key1", "myvalue1")}
	title, _ := serializer.SerializeDescriptor("name", "prefix", otherDims)

	fmt.Println(title)

}
