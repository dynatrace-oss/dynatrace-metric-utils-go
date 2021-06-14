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
// limitations under the License

package normalize_test

import (
	"strings"
	"testing"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/normalize"
)

func TestMetricKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{key: "basecase"},
			want: "basecase",
		},
		{
			name: "valid base case",
			args: args{key: "just.a.normal.key"},
			want: "just.a.normal.key",
		},
		{
			name: "valid leading underscore",
			args: args{key: "_case"},
			want: "_case",
		},
		{
			name: "valid underscore",
			args: args{key: "case_case"},
			want: "case_case",
		},
		{
			name: "valid number",
			args: args{key: "case1"},
			want: "case1",
		},
		{
			name: "invalid leading number",
			args: args{key: "1case"},
			want: "_case",
		},
		{
			name: "invalid multiple leading",
			args: args{key: "!@#case"},
			want: "_case",
		},
		{
			name: "invalid multiple trailing",
			args: args{key: "case!@#"},
			want: "case_",
		},
		{
			name: "valid leading uppercase",
			args: args{key: "Case"},
			want: "Case",
		},
		{
			name: "valid all uppercase",
			args: args{key: "CASE"},
			want: "CASE",
		},
		{
			name: "valid intermittent uppercase",
			args: args{key: "someCase"},
			want: "someCase",
		},
		{
			name: "valid multiple sections",
			args: args{key: "prefix.case"},
			want: "prefix.case",
		},
		{
			name: "valid multiple sections upper",
			args: args{key: "This.Is.Valid"},
			want: "This.Is.Valid",
		},
		{
			name: "invalid multiple sections leading number",
			args: args{key: "0a.b"},
			want: "_a.b",
		},
		{
			name: "valid multiple section leading underscore",
			args: args{key: "_a.b"},
			want: "_a.b",
		},
		{
			name: "valid leading number second section",
			args: args{key: "a.0"},
			want: "a.0",
		},
		{
			name: "valid leading number second section 2",
			args: args{key: "a.0.c"},
			want: "a.0.c",
		},
		{
			name: "valid leading number second section 3",
			args: args{key: "a.0b.c"},
			want: "a.0b.c",
		},
		{
			name: "invalid leading hyphen",
			args: args{key: "-dim"},
			want: "_dim",
		},
		{
			name: "valid trailing hyphen",
			args: args{key: "dim-"},
			want: "dim-",
		},
		{
			name: "valid trailing hyphens",
			args: args{key: "dim---"},
			want: "dim---",
		},
		{
			name:    "invalid empty",
			args:    args{key: ""},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid only number",
			args: args{key: "000"},
			want: "_",
		},
		{
			name: "invalid key first section only number",
			args: args{key: "0.section"},
			want: "_.section",
		},
		{
			name: "invalid leading character",
			args: args{key: "~key"},
			want: "_key",
		},
		{
			name: "invalid leading characters",
			args: args{key: "~0#key"},
			want: "_key",
		},
		{
			name: "invalid intermittent character",
			args: args{key: "some~key"},
			want: "some_key",
		},
		{
			name: "invalid intermittent characters",
			args: args{key: "some#~äkey"},
			want: "some_key",
		},
		{
			name: "invalid two consecutive dots",
			args: args{key: "a..b"},
			want: "a.b",
		},
		{
			name: "invalid five consecutive dots",
			args: args{key: "a.....b"},
			want: "a.b",
		},
		{
			name:    "invalid just a dot",
			args:    args{key: "."},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid leading dot",
			args:    args{key: ".a"},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid trailing dot",
			args: args{key: "a."},
			want: "a",
		},
		{
			name:    "invalid enclosing dots",
			args:    args{key: ".a."},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid consecutive leading underscores",
			args: args{key: "___a"},
			want: "___a",
		},
		{
			name: "valid consecutive trailing underscores",
			args: args{key: "a___"},
			want: "a___",
		},
		{
			name: "invalid trailing invalid chars",
			args: args{key: "a$%@"},
			want: "a_",
		},
		{
			name: "invalid trailing invalid chars groups",
			args: args{key: "a.b$%@.c"},
			want: "a.b_.c",
		},
		{
			name: "valid consecutive enclosed underscores",
			args: args{key: "a___b"},
			want: "a___b",
		},
		{
			name:    "invalid mixture dots underscores",
			args:    args{key: "._._._a_._._."},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid mixture dots underscores 2",
			args: args{key: "_._._.a_._"},
			want: "_._._.a_._",
		},
		{
			name: "invalid empty section",
			args: args{key: "an..empty.section"},
			want: "an.empty.section",
		},
		{
			name: "invalid characters",
			args: args{key: "a,,,b  c=d\\e\\ =,f"},
			want: "a_b_c_d_e_f",
		},
		{
			name: "invalid characters long",
			args: args{key: "a!b\"c#d$e%f&g'h(i)j*k+l,m-n.o/p:q;r<s=t>u?v@w[x]y\\z^0 1_2;3{4|5}6~7"},
			want: "a_b_c_d_e_f_g_h_i_j_k_l_m-n.o_p_q_r_s_t_u_v_w_x_y_z_0_1_2_3_4_5_6_7",
		},
		{
			name: "invalid trailing characters",
			args: args{key: "a.b.+"},
			want: "a.b._",
		},
		{
			name: "valid combined test",
			args: args{key: "metric.key-number-1.001"},
			want: "metric.key-number-1.001",
		},
		{
			name: "valid example 1",
			args: args{key: "MyMetric"},
			want: "MyMetric",
		},
		{
			name: "invalid example 1",
			args: args{key: "0MyMetric"},
			want: "_MyMetric",
		},
		{
			name: "invalid example 2",
			args: args{key: "mÄtric"},
			want: "m_tric",
		},
		{
			name: "invalid example 3",
			args: args{key: "metriÄ"},
			want: "metri_",
		},
		{
			name: "invalid example 4",
			args: args{key: "Ätric"},
			want: "_tric",
		},
		{
			name: "invalid example 5",
			args: args{key: "meträääääÖÖÖc"},
			want: "metr_c",
		},
		{
			name: "invalid truncate key too long",
			args: args{key: strings.Repeat("a", 270)},
			want: strings.Repeat("a", 250),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalize.MetricKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
