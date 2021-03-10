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
	"testing"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/normalize"
)

func TestDimensionKey(t *testing.T) {
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
			name: "valid case",
			args: args{key: "dim"},
			want: "dim",
		},
		{
			name: "valid number",
			args: args{key: "dim1"},
			want: "dim1",
		},
		{
			name: "valid leading underscore",
			args: args{key: "_dim"},
			want: "_dim",
		},
		{
			name: "invalid leading uppercase",
			args: args{key: "Dim"},
			want: "dim",
		},
		{
			name: "invalid internal uppercase",
			args: args{key: "dIm"},
			want: "dim",
		},
		{
			name: "invalid trailing uppercase",
			args: args{key: "diM"},
			want: "dim",
		},
		{
			name: "invalid all uppercase",
			args: args{key: "DIM"},
			want: "dim",
		},
		{
			name: "valid dimension colon",
			args: args{key: "dim:dim"},
			want: "dim:dim",
		},
		{
			name: "valid dimension underscore",
			args: args{key: "dim_dim"},
			want: "dim_dim",
		},
		{
			name: "valid dimension hyphen",
			args: args{key: "dim-dim"},
			want: "dim-dim",
		},
		{
			name: "invalid leading hyphen",
			args: args{key: "-dim"},
			want: "dim",
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
			name: "invalid leading multiple",
			args: args{key: "~0#dim"},
			want: "dim",
		},
		{
			name: "invalid leading multiple hyphens",
			args: args{key: "---dim"},
			want: "dim",
		},
		{
			name: "invalid leading colon",
			args: args{key: ":dim"},
			want: "dim",
		},
		{
			name:    "invalid chars",
			args:    args{key: "~@#ä"},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid trailing chars",
			args: args{key: "aaa~@#ä"},
			want: "aaa",
		},
		{
			name: "valid trailing underscores",
			args: args{key: "aaa___"},
			want: "aaa___",
		},
		{
			name:    "invalid only numbers",
			args:    args{key: "000"},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid compound key",
			args: args{key: "dim1.value1"},
			want: "dim1.value1",
		},
		{
			name: "invalid compound leading number",
			args: args{key: "dim.0dim"},
			want: "dim.dim",
		},
		{
			name: "invalid compound only number",
			args: args{key: "dim.000"},
			want: "dim",
		},
		{
			name: "invalid compound leading invalid char",
			args: args{key: "dim.~val"},
			want: "dim.val",
		},
		{
			name: "invalid compound trailing invalid char",
			args: args{key: "dim.val~~"},
			want: "dim.val",
		},
		{
			name: "invalid compound only invalid char",
			args: args{key: "dim.~~~"},
			want: "dim",
		},
		{
			name: "valid compound leading underscore",
			args: args{key: "dim._val"},
			want: "dim._val",
		},
		{
			name: "valid compound only underscore",
			args: args{key: "dim.___"},
			want: "dim.___",
		},
		{
			name: "valid compound long",
			args: args{key: "dim.dim.dim.dim"},
			want: "dim.dim.dim.dim",
		},
		{
			name: "invalid two dots",
			args: args{key: "a..b"},
			want: "a.b",
		},
		{
			name: "invalid five dots",
			args: args{key: "a.....b"},
			want: "a.b",
		},
		{
			name: "invalid leading dot",
			args: args{key: ".a"},
			want: "a",
		},
		{
			name: "valid colon in compound",
			args: args{key: "a.b:c.d"},
			want: "a.b:c.d",
		},
		{
			name: "invalid trailing dot",
			args: args{key: "a."},
			want: "a",
		},
		{
			name:    "invalid just a dot",
			args:    args{key: "."},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid trailing dots",
			args: args{key: "a..."},
			want: "a",
		},
		{
			name: "invalid enclosing dots",
			args: args{key: ".a."},
			want: "a",
		},
		{
			name: "invalid leading whitespace",
			args: args{key: "   a"},
			want: "a",
		},
		{
			name: "invalid trailing whitespace",
			args: args{key: "a   "},
			want: "a",
		},
		{
			name: "invalid internal whitespace",
			args: args{key: "a b"},
			want: "a_b",
		},
		{
			name: "invalid internal whitespace",
			args: args{key: "a    b"},
			want: "a_b",
		},
		{
			name:    "invalid empty",
			args:    args{key: ""},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid combined key",
			args: args{key: "dim.val:count.val001"},
			want: "dim.val:count.val001",
		},
		{
			name: "invalid characters",
			args: args{key: "a,,,b  c=d\\e\\ =,f"},
			want: "a_b_c_d_e_f",
		},
		{
			name: "invalid characters long",
			args: args{key: "a!b\"c#d$e%f&g'h(i)j*k+l,m-n.o/p:q;r<s=t>u?v@w[x]y\\z^0 1_2;3{4|5}6~7"},
			want: "a_b_c_d_e_f_g_h_i_j_k_l_m-n.o_p:q_r_s_t_u_v_w_x_y_z_0_1_2_3_4_5_6_7",
		},
		{
			name: "invalid example 1",
			args: args{key: "Tag"},
			want: "tag",
		},
		{
			name: "invalid example 2",
			args: args{key: "0Tag"},
			want: "tag",
		},
		{
			name: "invalid example 3",
			args: args{key: "tÄg"},
			want: "t_g",
		},
		{
			name: "invalid example 4",
			args: args{key: "mytäääg"},
			want: "myt_g",
		},
		{
			name: "invalid example 5",
			args: args{key: "ääätag"},
			want: "tag",
		},
		{
			name: "invalid example 6",
			args: args{key: "ä_ätag"},
			want: "__tag",
		},
		{
			name: "invalid example 7",
			args: args{key: "Bla___"},
			want: "bla___",
		},
		{
			name: "invalid truncate key too long",
			args: args{key: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			want: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalize.DimensionKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DimensionKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DimensionKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
