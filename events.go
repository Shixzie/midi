package midi

import "errors"

// Event represents a midi event, let it be
// MetaEvent or NormalEvent
type Event interface {
	SetTime(ticks int)
	Bytes() Codes
}

// EventType is self explanatory
type EventType byte

const (

	// MidiEvents

	// EventNoteOff midi event
	EventNoteOff EventType = 0x80
	// EventNoteOn midi event
	EventNoteOn EventType = 0x90
	// EventAfterTouch midi event
	EventAfterTouch EventType = 0xA0
	// EventController midi event
	EventController EventType = 0xB0
	// EventProgramChange midi event
	EventProgramChange EventType = 0xC0
	// EventChannelAfterTouch midi event
	EventChannelAfterTouch EventType = 0xD0
	// EventPitchBend midi event
	EventPitchBend EventType = 0xE0

	// MetaEvents

	// EventSequence midi event
	EventSequence EventType = iota
	// EventText midi event
	EventText
	// EventCopyright midi event
	EventCopyright
	// EventTrackName midi event
	EventTrackName
	// EventInstrument midi event
	EventInstrument
	// EventLyric midi event
	EventLyric
	// EventMarker midi event
	EventMarker
	// EventCuePoint midi event
	EventCuePoint
	// EventChannelPrefix midi event
	EventChannelPrefix EventType = 0x20
	// EventEndOfTrack midi event
	EventEndOfTrack EventType = 0x2f
	// EventTempo midi event
	EventTempo EventType = 0x51
	// EventSmpte midi event
	EventSmpte EventType = 0x54
	// EventTimeSig midi event
	EventTimeSig EventType = 0x58
	// EventKeySig midi event
	EventKeySig EventType = 0x59
	// EventSeqEvent midi event
	EventSeqEvent EventType = 0x7f
)

// NormalEvent is a single midi event
type NormalEvent struct {
	time           []byte
	_type          EventType
	channel        int
	param1, param2 byte
}

// NewEvent returns a new midi event
func NewEvent(time []byte, _type EventType, channel int, param1, param2 byte) (*NormalEvent, error) {
	if channel < 0 || channel > 15 {
		return nil, errors.New("channel out of bounds")
	}
	if _type < EventNoteOff || _type > EventPitchBend {
		return nil, errors.New("unknown event type")
	}
	if len(time) == 0 {
		time = []byte{0}
	}
	return &NormalEvent{
		time:    time,
		_type:   _type,
		channel: channel,
		param1:  param1,
		param2:  param2,
	}, nil
}

// SetTime sets the time for the event in ticks since the
// previous event
func (e *NormalEvent) SetTime(ticks int) {
	e.time = TranslateTickTime(ticks)
}

// Bytes returns the serielized event
func (e *NormalEvent) Bytes() Codes {
	typeChannel := e._type
	if typeChannel == 0 {
		typeChannel = EventType(e.channel & 0xF)
	}

	bytes := []byte{}

	bytes = append(bytes, e.time...)
	bytes = append(bytes, byte(typeChannel))
	bytes = append(bytes, e.param1)

	if e.param2 != 0 {
		bytes = append(bytes, e.param2)
	}

	return Codes(bytes)
}

// MetaEvent is a single meta event on a midi file
type MetaEvent struct {
	time  []byte
	_type EventType
	// data must be a string or []byte
	data interface{}
}

// NewMetaEvent returns a new meta event, data must be string, []byte or Timing
func NewMetaEvent(time []byte, _type EventType, data interface{}) (*MetaEvent, error) {
	switch data.(type) {
	case string, []byte, Timing:
	default:
		return nil, errors.New("invalid data type")
	}
	switch _type {
	case EventSequence, EventText, EventCopyright, EventTrackName, EventInstrument,
		EventLyric, EventMarker, EventCuePoint, EventChannelPrefix, EventEndOfTrack,
		EventTempo, EventSmpte, EventTimeSig, EventKeySig, EventSeqEvent:
	default:
		return nil, errors.New("invalid meta type")
	}
	if len(time) == 0 {
		time = []byte{0}
	}
	return &MetaEvent{
		time:  time,
		_type: _type,
		data:  data,
	}, nil
}

// SetTime sets the time for the event in ticks since the
// previous event
func (e *MetaEvent) SetTime(ticks int) {
	e.time = TranslateTickTime(ticks)
}

// Bytes returns the serielized event
func (e *MetaEvent) Bytes() Codes {
	bytes := []byte{}

	bytes = append(bytes, e.time...)
	bytes = append(bytes, byte(0xFF), byte(e._type))
	if v, ok := e.data.([]byte); ok {
		bytes = append(bytes, byte(len(v)))
		bytes = append(bytes, v...)
	} else if v, ok := e.data.(string); ok {
		bytes = append(bytes, byte(len(v)))
		bytes = append(bytes, []byte(v)...)
	} else if v, ok := e.data.(Timing); ok {
		bytes = append(bytes, 0x1, byte(v))
	} else {
		bytes = append(bytes, 0)
	}

	return Codes(bytes)
}
