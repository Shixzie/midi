package midi

import (
	"reflect"
	"testing"
)

func TestMpqnFromBpm(t *testing.T) {
	type args struct {
		bpm Timing
	}
	tests := []struct {
		name string
		args args
		want Timing
	}{
		{
			"mpqn to bpm",
			args{200},
			300000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MpqnFromBpm(tt.args.bpm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MpqnFromBpm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBpmFromMpqn(t *testing.T) {
	type args struct {
		mpqn Timing
	}
	tests := []struct {
		name string
		args args
		want Timing
	}{
		{
			"bpm to mpqn",
			args{300000},
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BpmFromMpqn(tt.args.mpqn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BpmFromMpqn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateTickTime(t *testing.T) {
	type args struct {
		ticks int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"ticks to time",
			args{128},
			[]byte{128, 129, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateTickTime(tt.args.ticks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TranslateTickTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
