package os_test

import (
	"os"
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

//gotests -w  -all go/os/env.go
//gotests -only GetEnvAsIntOrFallback go/os/env.go
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
