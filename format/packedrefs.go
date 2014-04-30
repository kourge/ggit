package format

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/kourge/ggit/core"
)

// The PackedRefs type represents the on-disk format used when multiple loose
// refs are packed into one file.
type PackedRefs struct {
	Refs []Ref
}

var _ core.EncodeDecoder = &PackedRefs{}

// Reader returns an io.Reader that yields refs in packed form.
func (p *PackedRefs) Reader() io.Reader {
	readers := make([]io.Reader, len(p.Refs)*4+1)
	readers[0] = strings.NewReader("# pack-refs with: peeled fully-peeled \n")

	for i, ref := range p.Refs {
		readers[i*4+1] = ref.Sha1.Reader()
		readers[i*4+2] = bytes.NewReader([]byte{' '})
		readers[i*4+3] = strings.NewReader(ref.Name)
		readers[i*4+4] = bytes.NewReader([]byte{'\n'})
	}

	return io.MultiReader(readers...)
}

// Decode reads from an io.Reader, presumably a packed refs file, and parses
// all the refs within it, ignoring comments and whitespace.
func (p *PackedRefs) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	atEof := false
	for !atEof {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			atEof = true
		} else if err != nil {
			return err
		}

		if pound := strings.IndexRune(line, '#'); pound != -1 {
			line = line[:pound]
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		sha1, name := parts[0], parts[1]

		hash, err := core.Sha1FromString(sha1)
		if err != nil {
			return err
		}

		p.Refs = append(p.Refs, Ref{Sha1: hash, Name: name})
	}

	return nil
}

// Sha1ForName returns the Sha1 of the packed ref with the given name, or an
// empty Sha1 if no packed ref of that name exists.
func (p *PackedRefs) Sha1ForName(name string) core.Sha1 {
	for _, ref := range p.Refs {
		if ref.Name == name {
			return ref.Sha1
		}
	}

	return core.Sha1{}
}
