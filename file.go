package midi

import (
	"errors"
	"os"
	"strings"
)

const (
	// DefaultTicks for a midi file
	DefaultTicks = 128
	// FileHdrChunkID is the file magic cookie
	FileHdrChunkID = "MThd"
)

var (
	// FileHdrChunkSize is the header length for SMF
	FileHdrChunkSize = Codes{0x00, 0x00, 0x00, 0x06}
	// FileHdrType0 is Midi Type 0 id
	FileHdrType0 = Codes{0x00, 0x00}
	// FileHdrType1 is Midi Type 1 id
	FileHdrType1 = Codes{0x00, 0x01}
)

// File is a midi file
type File struct {
	ticks  int
	tracks []*Track
}

// NewFile returns a new midi file
// ticks - Number of ticks per beat, defaults to 128
func NewFile(ticks int, tracks ...*Track) (*File, error) {
	if ticks == 0 {
		ticks = DefaultTicks
	}
	if ticks < 0 || ticks >= (1<<15) || ticks%1 != 0 {
		return nil, errors.New("ticks per beat must be an integer between 1 and 32767")
	}
	return &File{
		ticks:  ticks,
		tracks: tracks,
	}, nil
}

// AddTrack adds a track to the file
func (f *File) AddTrack(t *Track) {
	if t != nil {
		f.tracks = append(f.tracks, t)
	} else {
		f.tracks = append(f.tracks, NewTrack())
	}
}

// Bytes returns the serialized file
func (f *File) Bytes() Codes {
	trackCount := byte(len(f.tracks))

	// prepare the file header
	bytes := append(Codes(FileHdrChunkID), FileHdrChunkSize...)

	// set Midi type based on number of tracks
	if trackCount > 1 {
		bytes = append(bytes, FileHdrType1...)
	} else {
		bytes = append(bytes, FileHdrType0...)
	}

	// add the number of tracks (2 bytes)
	bytes = append(bytes, GetCodes(string(trackCount), 2)...)
	// add the number of ticks per beat
	bytes = append(bytes, byte(f.ticks/256), byte(f.ticks%256))

	// iterate over the tracks, converting to bytes too
	for _, track := range f.tracks {
		bytes = append(bytes, track.Bytes()...)
	}

	return bytes
}

// Save will save the midi file under filename
func (f *File) Save(filename string) error {
	if !strings.HasSuffix(filename, ".mid") {
		filename += ".mid"
	}
	ff, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = ff.Write([]byte(f.Bytes()))
	if err != nil {
		return err
	}
	return ff.Close()
}
