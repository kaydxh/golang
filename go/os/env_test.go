/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package os_test

import (
	"os"
	"slices"
	"testing"

	"github.com/google/uuid"
	os_ "github.com/kaydxh/golang/go/os"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsStringOrFallback(t *testing.T) {
	const expected = "foo"

	assert := assert.New(t)
	key := "FLOCKER_SET_VAR"
	os.Setenv(key, expected)

	assert.Equal(expected, os_.GetEnvAsStringOrFallback(key, "~"+expected))

	key = "FLOCKER_UNSET_VAR"
	assert.Equal(expected, os_.GetEnvAsStringOrFallback(key, expected))
	assert.NotEqual(expected, os_.GetEnvAsStringOrFallback(key, "~"+expected))
}

// gotests -w  -all go/os/env.go
// gotests -only GetEnvAsIntOrFallback go/os/env.go
func TestUUID(t *testing.T) {
	id := uuid.New().String()
	t.Logf("uuid: %v ", id)
}

func TestGetEnvAsIntOrFallback(t *testing.T) {
	type args struct {
		key          string
		defaultValue int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				key:          "FLOCKER_SET_VAR_INT",
				defaultValue: 0,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := os_.GetEnvAsIntOrFallback(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnvAsIntOrFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEnvAsIntOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsFloat64OrFallback(t *testing.T) {
	type args struct {
		key          string
		defaultValue float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				key:          "FLOCKER_SET_VAR_FLOAT",
				defaultValue: 1.0,
			},
			want:    1.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := os_.GetEnvAsFloat64OrFallback(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnvAsFloat64OrFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEnvAsFloat64OrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillEnvironmentPreservesNonEmptyBase(t *testing.T) {
	base := []string{"FROM_PROCESS=process", "EMPTY=", "CODEX_HOME=/internal"}
	got, changed := os_.FillEnvironment(base, map[string]string{
		"FROM_PROCESS": "stored",
		"EMPTY":        "filled",
		"CODEX_HOME":   "/stored",
		"NEW_KEY":      "new",
		"EMPTY_NAME":   "",
	})
	if !changed {
		t.Fatal("expected environment change")
	}
	for _, want := range []string{
		"FROM_PROCESS=process", "EMPTY=filled", "CODEX_HOME=/internal", "NEW_KEY=new",
	} {
		if !slices.Contains(got, want) {
			t.Fatalf("missing %q in %#v", want, got)
		}
	}
	if slices.Contains(got, "FROM_PROCESS=stored") || slices.Contains(got, "CODEX_HOME=/stored") {
		t.Fatalf("stored value overrode base: %#v", got)
	}
	if !slices.Equal(base, []string{"FROM_PROCESS=process", "EMPTY=", "CODEX_HOME=/internal"}) {
		t.Fatalf("input environment changed: %#v", base)
	}
}

func TestSetEnvironmentReplacesAllDuplicates(t *testing.T) {
	got := os_.SetEnvironment([]string{"PORT=3000", "X=1", "PORT=4000"}, "PORT", "5000")
	if !slices.Equal(got, []string{"X=1", "PORT=5000"}) {
		t.Fatalf("got %#v", got)
	}
}
