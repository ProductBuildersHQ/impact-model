package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDocumentMarkdownKeepsD2WithoutTools(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	out := filepath.Join(t.TempDir(), "doc.md")
	r := run(t, "", "export", "document", exTheory, "--diagrams", "d2", "--no-stamp", "--out", out)
	if r.code != ExitOK {
		t.Fatalf("exit %d: %s", r.code, r.err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	md := string(data)
	for _, want := range []string{"## Theory of Change", "## Logic Models", "## VSEM: Vision, Strategy, Execution, Metrics", "## Monitoring, Evaluation, and Learning Plan", "## Glossary", "```d2"} {
		if !strings.Contains(md, want) {
			t.Errorf("document has no %q", want)
		}
	}
	if strings.Index(md, "\n## VSEM") > strings.Index(md, "\n## Monitoring") {
		t.Error("sections are out of the profile's order")
	}
}

func TestExportDocumentSectionsAndUnavailable(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "doc.md")
	r := run(t, "", "export", "document", exTheory, "--diagrams", "d2", "--sections", "ogsm,dmaic", "--out", out)
	if r.code != ExitUsageErr || !strings.Contains(r.err, "--skip-unavailable") {
		t.Fatalf("DMAIC without a DMAIC pathway: exit %d: %s", r.code, r.err)
	}
	r = run(t, "", "export", "document", exTheory, "--diagrams", "d2", "--sections", "ogsm,dmaic", "--skip-unavailable", "--out", out)
	if r.code != ExitOK || !strings.Contains(r.err, "note: section dmaic is left out") {
		t.Fatalf("exit %d: %s", r.code, r.err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "## OGSM") || strings.Contains(string(data), "## DMAIC") {
		t.Error("the document does not hold just the OGSM")
	}
}

func TestExportDocumentUsageErrors(t *testing.T) {
	dir := t.TempDir()
	for name, args := range map[string][]string{
		"bad extension":       {"--out", filepath.Join(dir, "x.pdf")},
		"both lists":          {"--out", filepath.Join(dir, "x.md"), "--profile", "full", "--sections", "vsem"},
		"unknown section":     {"--out", filepath.Join(dir, "x.md"), "--sections", "vsem,nope"},
		"unknown profile":     {"--out", filepath.Join(dir, "x.md"), "--profile", "nope"},
		"d2 source in html":   {"--out", filepath.Join(dir, "x.html"), "--diagrams", "d2"},
		"reference doc in md": {"--out", filepath.Join(dir, "x.md"), "--reference-doc", exTheory},
	} {
		r := run(t, "", append([]string{"export", "document", exTheory}, args...)...)
		if r.code != ExitUsageErr {
			t.Errorf("%s: exit %d, want %d: %s", name, r.code, ExitUsageErr, r.err)
		}
	}
}

func TestExportDocumentNeedsD2ForDiagrams(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	out := filepath.Join(t.TempDir(), "deck.html")
	r := run(t, "", "export", "document", exTheory, "--out", out)
	if r.code != ExitToolMissing || !strings.Contains(r.err, "d2") {
		t.Fatalf("exit %d: %s", r.code, r.err)
	}
	if _, err := os.Stat(out); err == nil {
		t.Error("a file was written although d2 is missing")
	}
}

func TestExportDocumentRendersDiagramsWithD2(t *testing.T) {
	if _, err := exec.LookPath("d2"); err != nil {
		t.Skip("d2 is not installed")
	}
	dir := t.TempDir()
	md := filepath.Join(dir, "doc.md")
	if r := run(t, "", "export", "document", exTheory, "--no-stamp", "--out", md); r.code != ExitOK {
		t.Fatalf("md: exit %d: %s", r.code, r.err)
	}
	text, err := os.ReadFile(md)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(text), "](doc-diagrams/theory-of-change-1.svg)") {
		t.Error("the Markdown does not reference the rendered diagram")
	}
	svg, err := os.ReadFile(filepath.Join(dir, "doc-diagrams", "theory-of-change-1.svg"))
	if err != nil || !strings.Contains(string(svg), "<svg") {
		t.Fatalf("no rendered SVG: %v", err)
	}
	deck := filepath.Join(dir, "deck.html")
	if r := run(t, "", "export", "document", exTheory, "--profile", "strategy", "--out", deck); r.code != ExitOK {
		t.Fatalf("html: exit %d: %s", r.code, r.err)
	}
	page, err := os.ReadFile(deck)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(page), "data-deck-page") || !strings.Contains(string(page), `\u003csvg`) {
		t.Error("the deck has no slides or no inline SVG")
	}
}
