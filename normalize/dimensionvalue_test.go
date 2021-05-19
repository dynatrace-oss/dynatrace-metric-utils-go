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

package normalize_test

import (
	"strings"
	"testing"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/normalize"
)

func TestDimensionValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid value",
			args: args{value: "value"},
			want: "value",
		},
		{
			name: "valid empty",
			args: args{value: ""},
			want: "",
		},
		{
			name: "valid uppercase",
			args: args{value: "VALUE"},
			want: "VALUE",
		},
		{
			name: "valid colon",
			args: args{value: "a:3"},
			want: "a:3",
		},
		{
			name: "valid value 2",
			args: args{value: "~@#ä"},
			want: "~@#ä",
		},
		{
			name: "escape spaces",
			args: args{value: "a b"},
			want: "a\\ b",
		},
		{
			name: "escape comma",
			args: args{value: "a,b"},
			want: "a\\,b",
		},
		{
			name: "escape equals",
			args: args{value: "a=b"},
			want: "a\\=b",
		},
		{
			name: "escape backslash",
			args: args{value: "a\\b"},
			want: "a\\\\b",
		},
		{
			name: "escape multiple special chars",
			args: args{value: " ,=\\"},
			want: "\\ \\,\\=\\\\",
		},
		{
			name: "escape consecutive special chars",
			args: args{value: "  ,,==\\\\"},
			want: "\\ \\ \\,\\,\\=\\=\\\\\\\\",
		},
		{
			name: "escape quoted multiple special chars",
			args: args{value: `"\ ""`},
			want: `\"\\\ \"\"`,
		},
		{
			name: "escape key-value pair",
			args: args{value: "key=\"value\""},
			want: "key\\=\\\"value\\\"",
		},
		{
			name: "invalid unicode",
			args: args{value: "\u0000a\u0007"},
			want: "a",
		},
		{
			name: "invalid only unicode",
			args: args{value: "\u0000\u0007"},
			want: "",
		},
		{
			name: "invalid unicode space", // \u0001 is a space in unicode
			args: args{value: "a\u0001b"},
			want: "a_b",
		},
		{
			name: "invalid unicode spaces",
			args: args{value: "a\u0001\u0001\u0001b"},
			want: "a_b",
		},
		{
			name: "valid unicode", // 'Ab' in unicode
			args: args{value: "\u0034\u0066"},
			want: "\u0034\u0066",
		},
		{
			name: "valid unicode", //A umlaut, a with ring, O umlaut, U umlaut, all valid.
			args: args{value: "\u0132_\u0133_\u0150_\u0156"},
			want: "\u0132_\u0133_\u0150_\u0156",
		},
		{
			name: "invalid leading unicode NUL",
			args: args{value: "\u0000a"},
			want: "a",
		},
		{
			name: "invalid consecutive leading unicode NUL",
			args: args{value: "\u0000\u0000\u0000a"},
			want: "a",
		},
		{
			name: "invalid trailing unicode NUL",
			args: args{value: "a\u0000"},
			want: "a",
		},
		{
			name: "invalid consecutive trailing unicode NUL",
			args: args{value: "a\u0000\u0000\u0000"},
			want: "a",
		},
		{
			name: "invalid enclosed unicode NUL",
			args: args{value: "a\u0000b"},
			want: "a_b",
		},
		{
			name: "invalid consecutive enclosed unicode NUL",
			args: args{value: "a\u0000\u0000\u0000b"},
			want: "a_b",
		},
		{
			name: "invalid truncate value too long",
			args: args{value: strings.Repeat("a", 270)},
			want: strings.Repeat("a", 250),
		},
		{
			name: "escape sequence not broken apart 1",
			args: args{value: strings.Repeat("a", 249) + "="},
			want: strings.Repeat("a", 249),
		},
		{
			name: "escape sequence not broken apart 2",
			args: args{value: strings.Repeat("a", 248) + "=="},
			want: strings.Repeat("a", 248) + "\\=",
		},
		{
			name: "escape sequence not broken apart 3",
			// 3 trailing backslashes before escaping
			args: args{value: strings.Repeat("a", 247) + "\\\\\\"},
			// 1 escaped trailing backslash
			want: strings.Repeat("a", 247) + "\\\\",
		},
		{
			name: "dimension value of only backslashes",
			args: args{value: strings.Repeat("\\", 270)},
			want: strings.Repeat("\\\\", 125),
		},
		{
			name: "escape too long string",
			args: args{value: strings.Repeat("=", 250)},
			want: strings.Repeat("\\=", 125),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize.DimensionValue(tt.args.value); got != tt.want {
				t.Errorf("DimensionValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
