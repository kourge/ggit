package core

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// A GitMode is the 32-bit unsigned integer that represents the file mode of a
// stored tree entry, such as a blob or another tree. It is usually represented
// in a 6-digit octal format, where the first three digits indicate the nature
// of the entry (e.g. regular file, directory, etc), and the last three digits
// are standard Unix permission bits.
//
// Despite the last three digits being Unix permission bits, the only four
// permissions are allowed: 000, 644, 664, and 755.
type GitMode uint32

var _ Encoder = GitMode(0)

const (
	GitModeNull    GitMode = 0000 << 9
	GitModeDir     GitMode = 0040 << 9
	GitModeRegular GitMode = 0100 << 9
	GitModeSymlink GitMode = 0120 << 9
	GitModeGitlink GitMode = 0160 << 9

	GitModeNullPerm      GitMode = 0000
	GitModeExecutable    GitMode = 0755
	GitModeReadWritable  GitMode = 0644
	GitModeGroupWritable GitMode = 0664
)

// String returns this GitMode in form of a six-digit, left zero-padded octal
// number.
func (mode GitMode) String() string {
	return fmt.Sprintf("%06o", mode)
}

// Reader returns an io.Reader that formats the mode into an octal string. Note
// that this string is not left zero-padded to six characters.
func (mode GitMode) Reader() io.Reader {
	return strings.NewReader(fmt.Sprintf("%o", mode))
}

// GitModeFromString attempts to convert a string to a GitMode. If the string
// is not a properly formatted octal number, it returns an error. This string
// should be 6 characters long, all of them digits, to avoid any surprises
// that may result from octal number parsing.
func GitModeFromString(s string) (GitMode, error) {
	if mode, err := strconv.ParseInt(s, 8, 32); err != nil {
		return 0, err
	} else {
		return GitMode(mode), nil
	}
}
