// Copyright 2022 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIstioRelease(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr bool
		release IstioRelease
	}{
		{
			input: "1.11.1",
			release: IstioRelease{
				Major: 11,
				Minor: 1,
			},
		},
		{
			input: "1.4.0",
			release: IstioRelease{
				Major: 4,
				Minor: 0,
			},
		},
		{
			input:   "foo",
			wantErr: true,
		},
		{
			input:   "a.b.c",
			wantErr: true,
		},
		{
			input:   "2.1.1",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			e, r := IstioReleaseFromString(tc.input)
			if tc.wantErr {
				require.Error(t, e)
			} else {
				require.Equal(t, tc.release, r)
			}
		})
	}
}
