package zipfix

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
	"time"
)

func archive(t *testing.T, at time.Time) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, name := range []string{"[Content_Types].xml", "word/document.xml", "word/media/a.svg"} {
		f, err := w.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, Modified: at})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte("content of " + name)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestNormalizeMakesTimesIrrelevant(t *testing.T) {
	a, err := Normalize(archive(t, time.Date(2026, 11, 12, 6, 0, 0, 0, time.UTC)))
	if err != nil {
		t.Fatal(err)
	}
	b, err := Normalize(archive(t, time.Date(2027, 3, 1, 9, 30, 0, 0, time.UTC)))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("archives built at different times still differ")
	}
	r, err := zip.NewReader(bytes.NewReader(a), int64(len(a)))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.File) != 3 || r.File[1].Name != "word/document.xml" {
		t.Fatalf("entries changed: %d", len(r.File))
	}
	for _, f := range r.File {
		if !f.Modified.Equal(Epoch) {
			t.Errorf("%s: time %v", f.Name, f.Modified)
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		got, err := io.ReadAll(rc)
		if err != nil {
			t.Fatal(err)
		}
		if err := rc.Close(); err != nil {
			t.Fatal(err)
		}
		if string(got) != "content of "+f.Name {
			t.Errorf("%s: content changed to %q", f.Name, got)
		}
	}
}

func TestNormalizeRejectsNonZip(t *testing.T) {
	if _, err := Normalize([]byte("not a zip")); err == nil {
		t.Error("non-zip data was accepted")
	}
}
