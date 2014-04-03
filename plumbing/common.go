package plumbing

import (
	"os"

	"github.com/kourge/ggit/core"
)

const (
	DefaultZlibCompressionLevel int         = 6
	DefaultObjectFileMode       os.FileMode = 0444
)

// Errorf is a wrapper around errors.New(fmt.Sprintf(format, rest...)).
var Errorf = core.Errorf
