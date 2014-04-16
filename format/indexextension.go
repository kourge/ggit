package format

type indexExtensionHeader struct {
	Signature [4]byte
	Size      uint32
}

type indexExtension struct {
	indexExtensionHeader
	Data []byte
}

func (h indexExtensionHeader) Optional() bool {
	c := h.Signature[0]
	return c >= 'A' && c <= 'Z'
}
