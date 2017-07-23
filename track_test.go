package midi

import (
	"reflect"
	"testing"
)

func TestNewTrack(t *testing.T) {
	tests := []struct {
		name string
		want *Track
	}{
		{
			"new track",
			&Track{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_AddEvent(t *testing.T) {
	type args struct {
		e Event
	}
	ne, _ := NewEvent(nil, EventNoteOff, 0, 1, 0)
	me, _ := NewMetaEvent(nil, EventSequence, Timing(0))
	tests := []struct {
		name    string
		t       *Track
		args    args
		wantErr bool
	}{
		{
			"nil event",
			nil,
			args{nil},
			true,
		},
		{
			"normal event",
			NewTrack(),
			args{ne},
			false,
		},
		{
			"meta event",
			NewTrack(),
			args{me},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.AddEvent(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("Track.AddEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTrack_NoteOn(t *testing.T) {
	type args struct {
		channel  int
		p        Pitchier
		time     []byte
		velocity int
	}
	tr := NewTrack()
	p, _ := PitchFromNote("c5")
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"invalid channel",
			nil,
			args{-1, Note("c4"), nil, 0},
			nil,
		},
		{
			"invalid pitchier",
			nil,
			args{0, Note("c"), nil, 0},
			nil,
		},
		{
			"note-on",
			tr,
			args{0, p, nil, 0},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.NoteOn(tt.args.channel, tt.args.p, tt.args.time, tt.args.velocity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.NoteOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_NoteOff(t *testing.T) {
	type args struct {
		channel  int
		p        Pitchier
		time     []byte
		velocity int
	}
	tr := NewTrack()
	p, _ := PitchFromNote("c5")
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"invalid channel",
			nil,
			args{-1, Note("c4"), nil, 0},
			nil,
		},
		{
			"invalid pitchier",
			nil,
			args{0, Note("c"), nil, 0},
			nil,
		},
		{
			"note-off",
			tr,
			args{0, p, nil, 0},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.NoteOff(tt.args.channel, tt.args.p, tt.args.time, tt.args.velocity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.NoteOff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Note(t *testing.T) {
	type args struct {
		channel  int
		p        Pitchier
		dur      int
		time     []byte
		velocity int
	}
	tr := NewTrack()
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"0 dur",
			tr,
			args{0, Pitch(60), 0, nil, 0},
			tr,
		},
		{
			"1+ dur",
			tr,
			args{0, Pitch(60), 1, nil, 0},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Note(tt.args.channel, tt.args.p, tt.args.dur, tt.args.time, tt.args.velocity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.Note() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Chord(t *testing.T) {
	type args struct {
		channel  int
		chord    []Pitchier
		dur      int
		velocity int
	}
	tr := NewTrack()
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"chord",
			tr,
			args{
				0,
				[]Pitchier{
					Note("c4"),
					Note("c#4"),
					Pitch(60),
				},
				3,
				0,
			},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Chord(tt.args.channel, tt.args.chord, tt.args.dur, tt.args.velocity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.Chord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Instrument(t *testing.T) {
	type args struct {
		channel    int
		instrument byte
		time       []byte
	}
	tr := NewTrack()
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"invalid channel",
			nil,
			args{-1, 0, nil},
			nil,
		},
		{
			"instrument event",
			tr,
			args{0, 0, nil},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Instrument(tt.args.channel, tt.args.instrument, tt.args.time); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.Instrument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Tempo(t *testing.T) {
	type args struct {
		bpm  Timing
		time []byte
	}
	tr := NewTrack()
	tests := []struct {
		name string
		t    *Track
		args args
		want *Track
	}{
		{
			"set tempo event",
			tr,
			args{200, nil},
			tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Tempo(tt.args.bpm, tt.args.time); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.Tempo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Bytes(t *testing.T) {
	tr := NewTrack().
		NoteOn(0, Note("c4"), nil, 0).
		Chord(
			2,
			[]Pitchier{
				Note("c5"),
				Note("c3"),
			},
			4,
			0,
		).
		NoteOff(3, Note("c2"), nil, 0)
	tests := []struct {
		name string
		t    *Track
		want Codes
	}{
		{
			"track to codes",
			tr,
			Codes{
				0x4d, 0x54, 0x72, 0x6b, 0x0, 0x0, 0x0, 0x29,
				0x0, 0x90, 0x3c, 0x5a, 0x0, 0x90, 0x48, 0x5a,
				0x0, 0x90, 0x30, 0x5a, 0x80, 0x4, 0x80, 0x48,
				0x5a, 0x0, 0x80, 0x30, 0x5a, 0x0, 0x80, 0x24,
				0x5a, 0x0, 0xff, 0x2f, 0x0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Track.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
