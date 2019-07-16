package builder

import (
	"bytes"
	"io"

	"github.com/dbarney/professor/types"
)

type filter struct {
	w io.Writer
}

type updater struct {
	sha string
	b   *Builder
}

func (f filter) Write(p []byte) (int, error) {
	lines := bytes.Split(p, []byte("\n"))
	skipped := 0
	for i := 0; i < len(lines); {
		l := lines[i]
		match := targetMatch.FindSubmatch(l)
		if len(match) == 0 {
			i++
			continue
		}
		skipped += len(l) + 1 // +1 because of the new line character
		lines = append(lines[:i], lines[i+1:]...)
	}
	p1 := bytes.Join(lines, []byte("\n"))
	n, err := f.w.Write(p1)
	return n + skipped, err
}

func (u updater) Write(p []byte) (int, error) {
	lines := bytes.Split(p, []byte("\n"))
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		match := targetMatch.FindSubmatch(l)
		if len(match) != 0 {
			u.b.update(types.Target, u.sha, match[2])
			continue
		}
		c := make([]byte, len(l))
		// need a copy because everything is backed by the same byte buffer.
		// and without sending a copy, we will have really messed up
		// log lines.
		copy(c, l)
		u.b.update(types.Log, u.sha, c)
	}
	return len(p), nil
}
