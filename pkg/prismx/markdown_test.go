package prismx

import (
	"errors"
	"strings"
	"testing"
)

func TestProjectionMarkdownEmbeds(t *testing.T) {
	doc, st := example(t)
	render := map[string]func(MarkdownOptions) (string, error){
		"vsem": func(o MarkdownOptions) (string, error) { return VSEMMarkdown(doc, st, o) },
		"ogsm": func(o MarkdownOptions) (string, error) { return OGSMMarkdown(doc, st, o) },
		"mel":  func(o MarkdownOptions) (string, error) { return MELMarkdown(doc, st, o) },
	}
	for name, f := range render {
		standalone, err := f(MarkdownOptions{})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !strings.HasPrefix(standalone, "# ") {
			t.Errorf("%s: a stand-alone document does not start with its title:\n%.200s", name, standalone)
		}
		embedded, err := f(MarkdownOptions{Embed: true})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for _, line := range strings.Split(embedded, "\n") {
			if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") {
				t.Errorf("%s: an embedded document keeps a top heading: %q", name, line)
				break
			}
		}
		if !strings.Contains(embedded, "\n### ") && !strings.HasPrefix(embedded, "### ") {
			t.Errorf("%s: an embedded document has no level-3 heading", name)
		}
	}
}

func TestDMAICMarkdownReportsReadiness(t *testing.T) {
	doc, st := example(t)
	if _, _, err := DMAICMarkdown(doc, st, nil, MarkdownOptions{Embed: true}); !errors.Is(err, ErrNotDMAICReady) {
		t.Errorf("a model without DMAIC pathways: %v", err)
	}
}
