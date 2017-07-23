package midi

import (
	"errors"
	"strconv"
)

var (
	// TrackStartBytes are the bytes at the start of each track
	TrackStartBytes = Codes{0x4d, 0x54, 0x72, 0x6b}
	// TrackEndBytes are the bytes at the end of each track
	TrackEndBytes = Codes{0x00, 0xFF, 0x2F, 0x00}
)

const (
	// DefaultVolume for midi tracks
	DefaultVolume = 90
	// DefaultDuration for midi tracks
	DefaultDuration = 128
	// DefaultChannel for midi tracks
	DefaultChannel = 0
)

// Track is a midi track
type Track struct {
	events []Event
}

// NewTrack returns a new midi track
func NewTrack() *Track {
	return &Track{}
}

// AddEvent adds an event to the track
func (t *Track) AddEvent(e Event) error {
	if e == nil {
		return errors.New("can't add nil event to track")
	}
	t.events = append(t.events, e)
	return nil
}

// AddNoteOn adds a note-on event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) AddNoteOn(channel int, p Pitchier, time []byte, velocity int) *Track {
	p2 := velocity
	if p2 == 0 {
		p2 = DefaultVolume
	}
	p1, err := EnsurePitch(p)
	if err != nil {
		return nil
	}
	e, err := NewEvent(time, EventNoteOn, channel, byte(p1), byte(p2))
	if err != nil {
		return nil
	}
	t.events = append(t.events, e)
	return t
}

// NoteOn adds a note-on event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) NoteOn(channel int, p Pitchier, time []byte, velocity int) *Track {
	return t.AddNoteOn(channel, p, time, velocity)
}

// AddNoteOff adds a note-off event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) AddNoteOff(channel int, p Pitchier, time []byte, velocity int) *Track {
	p2 := velocity
	if p2 == 0 {
		p2 = DefaultVolume
	}
	p1, err := EnsurePitch(p)
	if err != nil {
		return nil
	}
	e, err := NewEvent(time, EventNoteOff, channel, byte(p1), byte(p2))
	if err != nil {
		return nil
	}
	t.events = append(t.events, e)
	return t
}

// NoteOff adds a note-off event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) NoteOff(channel int, p Pitchier, time []byte, velocity int) *Track {
	return t.AddNoteOff(channel, p, time, velocity)
}

// AddNote adds a note-on and -off event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// dur      - The duration of the note, is ticks
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) AddNote(channel int, p Pitchier, dur int, time []byte, velocity int) *Track {
	t.AddNoteOn(channel, p, time, velocity)
	if dur != 0 {
		t.AddNoteOff(channel, p, TranslateTickTime(dur), velocity)
	}
	return t
}

// Note adds a note-on and -off event to the track
// channel  - The channel to add the event to
// p        - The pitch of the note {Note|Pitch}
// dur      - The duration of the note, is ticks
// time     - The number of ticks since the previous event, default is 0
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) Note(channel int, p Pitchier, dur int, time []byte, velocity int) *Track {
	return t.AddNote(channel, p, dur, time, velocity)
}

// AddChord adds a note-on and -off event to the track for each
// pitch is chord
// channel  - The channel to add the event to
// chord    - The pitches {Note|Pitch}
// dur      - The duration of the note, is ticks
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) AddChord(channel int, chord []Pitchier, dur, velocity int) *Track {
	for _, note := range chord {
		t.NoteOn(channel, note, nil, velocity)
	}
	for i, note := range chord {
		if i == 0 {
			t.NoteOff(channel, note, TranslateTickTime(dur), velocity)
		} else {
			t.NoteOff(channel, note, nil, 0)
		}
	}
	return t
}

// Chord adds a note-on and -off event to the track for each
// pitch is chord
// channel  - The channel to add the event to
// chord    - The pitches {Note|Pitch}
// dur      - The duration of the note, is ticks
// velocity - The velocity the note was released, default is DefaultVolume
func (t *Track) Chord(channel int, chord []Pitchier, dur, velocity int) *Track {
	return t.AddChord(channel, chord, dur, velocity)
}

// SetInstrument sets the instrument for the track
// channel    - The channel to add the event to
// instrument - The instrument to set it to
// time       - The number of ticks since the previous event, default is 0
func (t *Track) SetInstrument(channel int, instrument byte, time []byte) *Track {
	e, err := NewEvent(time, EventProgramChange, channel, instrument, 0)
	if err != nil {
		return nil
	}
	t.events = append(t.events, e)
	return t
}

// Instrument sets the instrument for the track
// channel    - The channel to add the event to
// instrument - The instrument to set it to
// time       - The number of ticks since the previous event
func (t *Track) Instrument(channel int, instrument byte, time []byte) *Track {
	return t.SetInstrument(channel, instrument, time)
}

// SetTempo sets the tempo for the track
// bpm  - The new beats per minute
// time - The number of ticks since the previous event, default is 0
func (t *Track) SetTempo(bpm Timing, time []byte) *Track {
	e, _ := NewMetaEvent(time, EventTempo, MpqnFromBpm(bpm))
	t.events = append(t.events, e)
	return t
}

// Tempo sets the tempo for the track
// bpm  - The new beats per minute
// time - The number of ticks since the previous event, default is 0
func (t *Track) Tempo(bpm Timing, time []byte) *Track {
	return t.SetTempo(bpm, time)
}

// Bytes returns the serialized track
func (t *Track) Bytes() Codes {
	trackLength := 0
	bytes := Codes{}

	addEventBytes := func(e Event) {
		b := e.Bytes()
		trackLength += len(b)
		bytes = append(bytes, b...)
	}

	for _, event := range t.events {
		addEventBytes(event)
	}

	// Add the end-of-track bytes to the sum of bytes for the track, since
	// they are counted (unlike the start-of-track ones)
	trackLength += len(TrackEndBytes)

	s := strconv.Itoa(trackLength)
	if len(s) > 16 {
		s = s[:16]
	}

	lengthBytes := GetCodes(s, 4)

	tmp := append(TrackStartBytes, lengthBytes...)

	bytes = append(bytes, TrackEndBytes...)

	return append(tmp, bytes...)
}
