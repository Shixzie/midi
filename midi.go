package midi

import (
	"encoding/hex"
	"regexp"
)

func init() {
	PitchLetters = make(map[Pitch]Note, len(LetterPitches))
	for note, pitch := range LetterPitches {
		PitchLetters[pitch] = note
	}

	noteParser = regexp.MustCompile("([a-g])(#+|b+)?([0-9]+)$")
}

// Codes for the file
type Codes []byte

// Converts an array of bytes to a string of hexadecimal characters.
// Prepares it to be converted into a base64 string
func (c Codes) String() string {
	return hex.EncodeToString(c)
}

// GetCodes converts a string of hexadecimal values to an array of bytes,
// It can also add remaining "0" nibbles in order to have enough bytes in the
// array as the finalBytes param
func GetCodes(str string, finalBytes int) Codes {
	if finalBytes != 0 {
		for len(str)/2 < finalBytes {
			str = "0" + str
		}
	}
	codes, _ := hex.DecodeString(str)
	return Codes(codes)
}
