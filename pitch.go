package midi

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Pitch represents a single pitch on a midi file
type Pitch int

// Note represents a single note of a midi file
type Note string

// Pitchier is just an utility interface to
// accept both Notes and Pitches
type Pitchier interface {
	pitchie() byte
}

func (Pitch) pitchie() byte { return 0 }
func (Note) pitchie() byte  { return 1 }

// LetterPitches maps letters to it's pitch
var LetterPitches = map[Note]Pitch{
	"c":  12,
	"c#": 13,
	"d":  14,
	"d#": 15,
	"e":  16,
	"f":  17,
	"f#": 18,
	"g":  19,
	"g#": 20,
	"a":  21,
	"a#": 22,
	"b":  23,
}

// PitchLetters maps pitchs to it's note
var PitchLetters map[Pitch]Note

// FlattenedNotes maps shar notes to flattened notes
var FlattenedNotes = map[Note]Note{
	"a#": "bb",
	"c#": "db",
	"d#": "eb",
	"f#": "gb",
	"g#": "ab",
}

// ([a-g])(#+|b+)?([0-9]+)$
var noteParser *regexp.Regexp

// EnsurePitch ensures that the given argument is converted to a MIDI pitch.
// Note that it may already be one (including a purely numeric string)
func EnsurePitch(p Pitchier) (Pitch, error) {
	if p == nil {
		return 0, errors.New("nil pitchier")
	}
	if p.pitchie() == 0 {
		return p.(Pitch), nil
	}
	return PitchFromNote(p.(Note))
}

// PitchFromNote converts a symbolic note name (e.g. "c4")
// to a numeric MIDI pitch (e.g. 60, middle C)
func PitchFromNote(n Note) (Pitch, error) {
	match := noteParser.FindString(string(n))
	if len(match) == 0 {
		return -1, fmt.Errorf("invalid note %q", n)
	}
	note := match[0]
	octave, accLen, acc := int64(0), 0, -1
	if len(match) == 2 {
		octave, _ = strconv.ParseInt(string(match[1]), 10, 64)
	} else {
		if match[1] == '#' {
			acc = 1
			accLen = 1
		}
		octave, _ = strconv.ParseInt(string(match[2]), 10, 64)
	}
	return (12 * Pitch(octave)) + LetterPitches[Note(note)] + Pitch(acc*accLen), nil
}

// NoteFromPitch convert a numeric MIDI pitch value (e.g. 60)
// to a symbolic note name (e.g. "c4")
func NoteFromPitch(p Pitch, returnFlattened bool) Note {
	octave := float64(0)
	if p > 23 {
		// p is on octave 1 or more
		octave = math.Floor(float64(p/12)) - 1
		p = p - Pitch(octave*12)
	}
	note := PitchLetters[p]
	if returnFlattened && strings.Contains(string(note), "#") {
		note = FlattenedNotes[note]
	}
	return Note(fmt.Sprintf("%v%v", note, octave))
}
