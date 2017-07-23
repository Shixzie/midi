package midi

import (
	"math"
)

// Timing is used to represent both
// MPQN - microseconds per quarter note
// BPM  - beats per minute
type Timing int64

// MicrosecondsPerMinute constant
const MicrosecondsPerMinute = 60000000

// MpqnFromBpm converts beats per minute (BPM) to
// microseconds per quarter note (MPQN)
func MpqnFromBpm(bpm Timing) Timing {
	return Timing(math.Floor(MicrosecondsPerMinute / float64(bpm)))
}

// BpmFromMpqn converts microseconds per quarter note (MPQN) to
// beats per minute (BPM)
func BpmFromMpqn(mpqn Timing) Timing {
	return Timing(math.Floor(MicrosecondsPerMinute / float64(mpqn)))
}

// TranslateTickTime translates number of ticks to MIDI timestamp format
// returning a []byte with the time values
func TranslateTickTime(ticks int) []byte {
	buffer := ticks & 0x7F

	for ticks != 0 {
		ticks = ticks >> 7
		buffer <<= 8
		buffer |= ((ticks & 0x7F) | 0x80)
	}

	bList := []byte{}

	for {
		bList = append(bList, byte(buffer&0xFF))
		if buffer&0x80 != 0 {
			buffer >>= 8
		} else {
			break
		}
	}

	return bList
}
