package core

// A Sha1 represents a 40-character alphanumberic SHA-1 hash checksum.
type Sha1 string

// IsValid returns true if the underlying string that represents the SHA-1
// checksum is exactly 40 characters long and is lowercase alphanumeric.
// Otherwise, it returns false.
func (sha Sha1) IsValid() bool {
	return len(sha) == 40 && allRunesMatch(string(sha), isAlphanumeric)
}

func isAlphanumeric(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
}

func allRunesMatch(s string, f func(rune) bool) bool {
	for _, r := range s {
		if !f(r) {
			return false
		}
	}
	return true
}
