package format

import (
	"errors"
	"io"
)

var (
	ErrWriteLimitReached = errors.New("LimitedWriter write limit reached")
)

// LimitWriter returns an io.Writer that writes to w but stops with
// ErrWriteLimitReached after n bytes. The underlying implementation is
// a *LimitedWriter.
func LimitWriter(w io.Writer, n int64) *LimitedWriter {
	return &LimitedWriter{w, n, ErrWriteLimitReached}
}

// SilentLimitWriter returns an io.Writer that writes to w but stops with after
// n bytes. When the limit of n bytes is reached, the error returned is nil.
// This specifically breaks the contract of the Write method on the io.Writer
// interface, which states that it must return a non-nil error if it returns
// n < len(p). The underlying implementation is a *LimitedWriter.
func SilentLimitWriter(w io.Writer, n int64) *LimitedWriter {
	return &LimitedWriter{w, n, nil}
}

// A LimitedWriter writes to W but limits the total amount of data written to
// just N bytes. Each call to Write updates N to reflect the new amount
// remaining. 
type LimitedWriter struct {
	W io.Writer
	N int64
	limitError error
}

func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.N >= int64(len(p)) {
		l.N -= int64(len(p))
		return l.W.Write(p)
	} else {
		l.W.Write(p[:l.N])
		return int(l.N), l.limitError
	}
}
