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
package oneagentenrichment

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
)

func Test_readIndirectionFile(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{reader: strings.NewReader("metadata_file.txt")},
			want: "metadata_file.txt",
		},
		{
			name: "empty file",
			args: args{reader: strings.NewReader("")},
			want: "",
		},
		{
			name: "whitespace",
			args: args{reader: strings.NewReader("\t \t metadata_file.txt\t \t\n")},
			want: "metadata_file.txt",
		},
		{
			name:    "pass nil reader",
			args:    args{reader: nil},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readIndirectionFile(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("readIndirectionFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readIndirectionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readMetadataFile(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{reader: strings.NewReader("key1=value1\nkey2=value2")},
			want: []string{"key1=value1", "key2=value2"},
		},
		{
			name: "metadata file empty",
			args: args{reader: strings.NewReader("")},
			want: []string{},
		},
		{
			name: "ignore whitespace",
			args: args{reader: strings.NewReader("\t \tkey1=value1\t \t\n\t \tkey2=value2\t \t")},
			want: []string{"key1=value1", "key2=value2"},
		},
		{
			name: "empty lines ignored",
			args: args{reader: strings.NewReader("\n\t \t\n")},
			want: []string{},
		},
		{
			name:    "pass nil reader",
			args:    args{reader: nil},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readMetadataFile(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("readMetadataFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readMetadataFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readOneAgentMetadata(t *testing.T) {
	type args struct {
		indirectionBasename string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{indirectionBasename: "testdata/indirection.properties"},
			want: []string{"key1=value1", "key2=value2", "key3=value3"},
		},
		{
			name: "metadata file empty",
			args: args{indirectionBasename: "testdata/indirection_target_empty.properties"},
			want: []string{},
		},
		{
			name:    "indirection file empty",
			args:    args{indirectionBasename: "testdata/indirection_empty.properties"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "indirection file does not exist",
			args:    args{indirectionBasename: "testdata/indirection_file_that_does_not_exist.properties"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "indirection target does not exist",
			args:    args{indirectionBasename: "testdata/indirection_target_nonexistent.properties"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readOneAgentMetadata(tt.args.indirectionBasename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readOneAgentMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readOneAgentMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOneAgentMetadataEnricher_parseOneAgentMetadata(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want []dimensions.Dimension
	}{
		{
			name: "valid case",
			args: args{[]string{"key1=value1", "key2=value2", "key3=value3"}},
			want: []dimensions.Dimension{
				dimensions.NewDimension("key1", "value1"),
				dimensions.NewDimension("key2", "value2"),
				dimensions.NewDimension("key3", "value3"),
			},
		},
		{
			name: "pass empty list",
			args: args{[]string{}},
			want: []dimensions.Dimension{},
		},
		{
			name: "pass invalid strings",
			args: args{[]string{
				"=0x5c14d9a68d569861",
				"otherKey=",
				"",
				"=",
				"===",
			}},
			want: []dimensions.Dimension{},
		},
		{
			name: "pass mixed strings",
			args: args{[]string{
				"invalid1",
				"key1=value1",
				"=invalid",
				"key2=value2",
				"===",
			}},
			want: []dimensions.Dimension{
				dimensions.NewDimension("key1", "value1"),
				dimensions.NewDimension("key2", "value2"),
			},
		},
		{
			name: "valid tailing equal signs",
			args: args{[]string{"key1=value1=="}},
			want: []dimensions.Dimension{dimensions.NewDimension("key1", "value1==")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := parseOneAgentMetadata(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OneAgentMetadataEnricher.parseOneAgentMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_asDimensionSet(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want dimensions.DimensionSet
	}{
		{
			name: "empty set",
			args: args{lines: []string{}},
			want: dimensions.CreateDimensionSet(),
		},
		{
			name: "one element",
			args: args{lines: []string{"key1=value1"}},
			want: dimensions.CreateDimensionSet(dimensions.NewDimension("key1", "value1")),
		},
		{
			name: "multiple elements",
			args: args{lines: []string{"key1=value1", "key2=value2", "key3=value3"}},
			want: dimensions.CreateDimensionSet(
				dimensions.NewDimension("key1", "value1"),
				dimensions.NewDimension("key2", "value2"),
				dimensions.NewDimension("key3", "value3"),
			),
		},
		{
			name: "duplicate keys",
			args: args{lines: []string{"key1=value1", "key2=value2", "key1=value3"}},
			want: dimensions.CreateDimensionSet(
				dimensions.NewDimension("key1", "value1"),
				dimensions.NewDimension("key2", "value2"),
				dimensions.NewDimension("key1", "value3"),
			),
		},
		{
			name: "invalid keys are not formatted",
			args: args{lines: []string{"key1====", "~~#=value2", "=value3"}},
			want: dimensions.CreateDimensionSet(
				dimensions.NewDimension("key1", "==="),
				dimensions.NewDimension("~~#", "value2"),
				// =value3 is discarded since it cannot be split
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := asDimensionSet(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("asDimensionSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
