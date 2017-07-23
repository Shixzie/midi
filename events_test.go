package midi

import (
	"reflect"
	"testing"
)

func TestNewEvent(t *testing.T) {
	type args struct {
		time    []byte
		_type   EventType
		channel int
		param1  byte
		param2  byte
	}
	tests := []struct {
		name    string
		args    args
		want    *NormalEvent
		wantErr bool
	}{
		{
			"invalid channel",
			args{nil, EventNoteOn, -1, 0, 0},
			nil,
			true,
		},
		{
			"invalid event type",
			args{nil, 0, 0, 0, 0},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEvent(tt.args.time, tt.args._type, tt.args.channel, tt.args.param1, tt.args.param2)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalEvent_SetTime(t *testing.T) {
	type args struct {
		ticks int
	}
	tests := []struct {
		name string
		e    *NormalEvent
		args args
	}{
		{
			"set time",
			&NormalEvent{},
			args{128},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.SetTime(tt.args.ticks)
		})
	}
}

func TestNormalEvent_Bytes(t *testing.T) {
	e := func(t EventType) *NormalEvent {
		return &NormalEvent{
			time:    []byte{0},
			_type:   t,
			channel: 0,
			param1:  60,
			param2:  90,
		}
	}
	tests := []struct {
		name string
		e    *NormalEvent
		want Codes
	}{
		{
			"note-on to bytes",
			e(EventNoteOn),
			Codes{0x0, 0x90, 0x3c, 0x5a},
		},
		{
			"0-type to bytes",
			e(0),
			Codes{0x0, 0x0, 0x3c, 0x5a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalEvent.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetaEvent(t *testing.T) {
	type args struct {
		time  []byte
		_type EventType
		data  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *MetaEvent
		wantErr bool
	}{
		{
			"invalid data type",
			args{nil, 0, nil},
			nil,
			true,
		},
		{
			"invalid event type",
			args{nil, 0xFF, ""},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetaEvent(tt.args.time, tt.args._type, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetaEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetaEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetaEvent_SetTime(t *testing.T) {
	type args struct {
		ticks int
	}
	tests := []struct {
		name string
		e    *MetaEvent
		args args
	}{
		{
			"set time meta",
			&MetaEvent{},
			args{128},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.SetTime(tt.args.ticks)
		})
	}
}

func TestMetaEvent_Bytes(t *testing.T) {
	e := func(d interface{}) *MetaEvent {
		return &MetaEvent{
			data:  d,
			_type: EventSequence,
			time:  []byte{0},
		}
	}
	tests := []struct {
		name string
		e    *MetaEvent
		want Codes
	}{
		{
			"string data",
			e("2"),
			Codes{0x0, 0xFF, 0x7, 0x1, 0x32},
		},
		{
			"timing data",
			e(Timing(2)),
			Codes{0x0, 0xFF, 0x7, 0x1, 0x2},
		},
		{
			"bytes data",
			e([]byte{0x2}),
			Codes{0x0, 0xFF, 0x7, 0x1, 0x2},
		},
		{
			"unknow data",
			e(1),
			Codes{0x0, 0xFF, 0x7, 0x0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetaEvent.Bytes() = %#v, want %v", got, tt.want)
			}
		})
	}
}
