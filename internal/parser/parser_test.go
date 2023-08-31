// Copyright (C) 2023 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jaredallard/mangarn/internal/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     parser.Page
	}{
		{
			filename: "[Releaser] A Random Name Vol. 1 Chapter. 01.jpg",
			want: parser.Page{
				Title:              "A Random Name",
				Volume:             1,
				Chapter:            1,
				PageNumber:         ^uint(0),
				AbsolutePageNumber: ^uint(0),
			},
		},
		{
			filename: "0001_A_Series_name_c001_v01_p000_Source_Quality_Release.jpg",
			want: parser.Page{
				Title:              "A Series name",
				Volume:             1,
				Chapter:            1,
				AbsolutePageNumber: 1,
				PageNumber:         0,
			},
		},
		{
			filename: "1009_A_Random_Name_c118_v14_Releaser_HQ_60.jpg",
			want: parser.Page{
				Title:              "A Random Name",
				Volume:             14,
				Chapter:            118,
				PageNumber:         60,
				AbsolutePageNumber: 1009,
			},
		},
	}
	for _, tt := range tests {
		if tt.name == "" {
			tt.name = "Parse(" + tt.filename + ")"
		}

		// set the filename for the want struct to equal the provided filename
		tt.want.FileName = tt.filename

		t.Run(tt.name, func(t *testing.T) {
			res, _ := parser.Parse(tt.filename)
			if diff := cmp.Diff(tt.want, res); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
