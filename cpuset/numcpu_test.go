/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cpuset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNumCPU(t *testing.T) {
	orig := onlineCPUsPath
	t.Cleanup(func() { onlineCPUsPath = orig })

	dir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    int
		wantErr bool
	}{
		{name: "single", content: "0\n", want: 1},
		{name: "range", content: "0-3\n", want: 4},
		{name: "list and ranges", content: "0-1,4-7\n", want: 6},
		{name: "no trailing newline", content: "0-15", want: 16},
		{name: "malformed", content: "nope\n", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, tc.name)
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatal(err)
			}
			onlineCPUsPath = path
			got, err := NumCPU()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("NumCPU() = %d, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("NumCPU() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("NumCPU() = %d, want %d", got, tc.want)
			}
		})
	}

	t.Run("missing file", func(t *testing.T) {
		onlineCPUsPath = filepath.Join(dir, "absent")
		if _, err := NumCPU(); err == nil {
			t.Error("NumCPU() with missing file: want error, got nil")
		}
	})
}
