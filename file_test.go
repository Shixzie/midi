package midi

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewFile(t *testing.T) {
	type args struct {
		ticks  int
		tracks []*Track
	}
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
	}{
		{
			"default ticks",
			args{0, nil},
			&File{
				ticks: DefaultTicks,
			},
			false,
		},
		{
			"invalid ticks",
			args{-1, nil},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFile(tt.args.ticks, tt.args.tracks...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_AddTrack(t *testing.T) {
	type args struct {
		t *Track
	}
	tests := []struct {
		name string
		f    *File
		args args
	}{
		{
			"nil track",
			&File{},
			args{nil},
		},
		{
			"non-nil track",
			&File{},
			args{NewTrack()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.AddTrack(tt.args.t)
		})
	}
}

func TestFile_Bytes(t *testing.T) {
	f, _ := NewFile(DefaultTicks,
		NewTrack().
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
			NoteOff(3, Note("c2"), nil, 0),
	)
	f2, _ := NewFile(DefaultTicks, NewTrack(), NewTrack())
	tests := []struct {
		name string
		f    *File
		want Codes
	}{
		{
			"file to bytes",
			f,
			Codes{
				0x4d, 0x54, 0x68, 0x64, 0x0, 0x0, 0x0, 0x6,
				0x0, 0x0, 0x0, 0x80, 0x4d, 0x54, 0x72, 0x6b,
				0x0, 0x0, 0x0, 0x29, 0x0, 0x90, 0x3c, 0x5a,
				0x0, 0x90, 0x48, 0x5a, 0x0, 0x90, 0x30, 0x5a,
				0x80, 0x4, 0x80, 0x48, 0x5a, 0x0, 0x80, 0x30,
				0x5a, 0x0, 0x80, 0x24, 0x5a, 0x0, 0xff, 0x2f, 0x0,
			},
		},
		{
			"file to bytes",
			f2,
			Codes{
				0x4d, 0x54, 0x68, 0x64, 0x0, 0x0, 0x0, 0x6,
				0x0, 0x1, 0x0, 0x80, 0x4d, 0x54, 0x72, 0x6b,
				0x0, 0x0, 0x0, 0x4, 0x0, 0xff, 0x2f, 0x0,
				0x4d, 0x54, 0x72, 0x6b, 0x0, 0x0, 0x0, 0x4,
				0x0, 0xff, 0x2f, 0x0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Save(t *testing.T) {
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

	type fields struct {
		ticks  int
		tracks []*Track
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"file save 1",
			fields{DefaultTicks, []*Track{tr}},
			args{"test1"},
			false,
		},
		{
			"file save 2",
			fields{DefaultTicks, []*Track{tr}},
			args{"test2.mid"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				ticks:  tt.fields.ticks,
				tracks: tt.fields.tracks,
			}
			if err := f.Save(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("File.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_Encode(t *testing.T) {
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
	f, _ := NewFile(DefaultTicks, tr)
	type fields struct {
		ticks  int
		tracks []*Track
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantW   string
		wantErr bool
	}{
		{
			"file.encode",
			fields{DefaultTicks, []*Track{tr}},
			49,
			string(f.Bytes()),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				ticks:  tt.fields.ticks,
				tracks: tt.fields.tracks,
			}
			w := &bytes.Buffer{}
			got, err := f.Encode(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("File.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("File.Encode() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("File.Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEncode(t *testing.T) {
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
	f, _ := NewFile(DefaultTicks, tr)
	type args struct {
		f *File
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantW   string
		wantErr bool
	}{
		{
			"encode",
			args{f},
			49,
			string(f.Bytes()),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := Encode(tt.args.f, w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
