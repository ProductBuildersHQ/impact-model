// Package zipfix makes a zip archive, such as a .docx or .pptx file, the same
// bytes every time it is built from the same content. Pandoc stamps some
// entries with the time it runs; zipfix rewrites every entry with one fixed
// time and keeps the entries, their order, and their contents.
package zipfix

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// Epoch is the time every entry is given: the earliest a zip file can record.
var Epoch = time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)

// Normalize returns the archive with every entry's time set to Epoch and no
// extra fields, which can carry further timestamps.
func Normalize(data []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("zipfix: %w", err)
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, f := range r.File {
		hdr := f.FileHeader
		hdr.Modified = Epoch
		hdr.Extra = nil
		dst, err := w.CreateHeader(&hdr)
		if err != nil {
			return nil, fmt.Errorf("zipfix: %s: %w", f.Name, err)
		}
		src, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("zipfix: %s: %w", f.Name, err)
		}
		_, err = io.Copy(dst, src) //nolint:gosec // entries come from pandoc's own output, not untrusted input
		if cerr := src.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return nil, fmt.Errorf("zipfix: %s: %w", f.Name, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zipfix: %w", err)
	}
	return buf.Bytes(), nil
}

// NormalizeFile normalizes the archive at path in place.
func NormalizeFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("zipfix: %w", err)
	}
	out, err := Normalize(data)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("zipfix: %w", err)
	}
	if err := os.WriteFile(path, out, info.Mode().Perm()); err != nil {
		return fmt.Errorf("zipfix: %w", err)
	}
	return nil
}
