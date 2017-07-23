package midi

import "testing"

func TestPitch_pitchie(t *testing.T) {
	tests := []struct {
		name string
		p    Pitch
		want byte
	}{
		{
			"Pitch pitchie",
			Pitch(0),
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.pitchie(); got != tt.want {
				t.Errorf("Pitch.pitchie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNote_pitchie(t *testing.T) {
	tests := []struct {
		name string
		n    Note
		want byte
	}{
		{
			"Note pitchie",
			"",
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.pitchie(); got != tt.want {
				t.Errorf("Note.pitchie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsurePitch(t *testing.T) {
	type args struct {
		p Pitchier
	}
	tests := []struct {
		name    string
		args    args
		want    Pitch
		wantErr bool
	}{
		{
			"nil pitchier",
			args{nil},
			0,
			true,
		},
		{
			"pitch to pitch",
			args{Pitch(60)},
			60,
			false,
		},
		{
			"note to pitch",
			args{Note("c4")},
			60,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnsurePitch(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsurePitch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnsurePitch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPitchFromNote(t *testing.T) {
	type args struct {
		n Note
	}
	tests := []struct {
		name    string
		args    args
		want    Pitch
		wantErr bool
	}{
		{
			"invalid note",
			args{""},
			-1,
			true,
		},
		{
			"note to pitch",
			args{"c4"},
			60,
			false,
		},
		{
			"note sharp to pitch",
			args{"c#4"},
			61,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PitchFromNote(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("PitchFromNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PitchFromNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoteFromPitch(t *testing.T) {
	type args struct {
		p               Pitch
		returnFlattened bool
	}
	tests := []struct {
		name string
		args args
		want Note
	}{
		{
			"octave 1+",
			args{60, false},
			"c4",
		},
		{
			"sharp",
			args{61, false},
			"c#4",
		},
		{
			"flatten",
			args{61, true},
			"db4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoteFromPitch(tt.args.p, tt.args.returnFlattened); got != tt.want {
				t.Errorf("NoteFromPitch() = %v, want %v", got, tt.want)
			}
		})
	}
}
