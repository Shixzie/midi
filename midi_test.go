package midi

import (
	"reflect"
	"testing"
)

func TestCodes_String(t *testing.T) {
	tests := []struct {
		name string
		c    Codes
		want string
	}{
		{
			"codes to string",
			Codes{0x5, 0x2, 0xaf},
			"0502af",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Codes.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCodes(t *testing.T) {
	type args struct {
		str        string
		finalBytes int
	}
	tests := []struct {
		name string
		args args
		want Codes
	}{
		{
			"string to codes",
			args{"0502af", 0},
			Codes{0x5, 0x2, 0xaf},
		},
		{
			"string to codes final bytes",
			args{"0502af", 3},
			Codes{0x5, 0x2, 0xaf},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCodes(tt.args.str, tt.args.finalBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCodes() = %#v, want %v", got, tt.want)
			}
		})
	}
}
