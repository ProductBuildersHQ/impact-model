package document_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/document"
)

func TestDeck(t *testing.T) {
	d := build(t, document.Options{Stamp: "Example stamp"})
	page, err := d.Deck(document.DeckOptions{SVG: func(g *document.Diagram) ([]byte, bool) {
		return []byte(`<?xml version="1.0"?><svg id="` + g.ID + `"></svg>`), true
	}})
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	for _, want := range []string{"<title>", "/ajax/libs/react/18.3.1/", "/ajax/libs/react-dom/18.3.1/", "data-deck-page", "window.__deckReady = true", "deck-page-content"} {
		if !strings.Contains(html, want) {
			t.Errorf("deck has no %q", want)
		}
	}
	m := regexp.MustCompile(`(?s)<script id="deck-data" type="application/json">(.*?)</script>`).FindStringSubmatch(html)
	if m == nil {
		t.Fatal("deck has no data block")
	}
	if strings.Contains(m[1], "<") {
		t.Error("the data block holds a raw '<', which could end the script early")
	}
	var data struct {
		Slides []struct {
			Kind, Title, HTML, SVG string
		}
	}
	if err := json.Unmarshal([]byte(m[1]), &data); err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	for _, s := range data.Slides {
		kinds[s.Kind]++
		if strings.Contains(s.SVG, "<?xml") {
			t.Error("an SVG keeps its XML declaration")
		}
		if strings.Contains(s.HTML, "<script") {
			t.Error("slide HTML holds a script")
		}
	}
	if kinds["title"] != 1 || kinds["section"] != len(d.Sections) || kinds["diagram"] != len(d.Diagrams()) || kinds["content"] == 0 {
		t.Errorf("slide kinds = %v", kinds)
	}
	again, err := d.Deck(document.DeckOptions{SVG: func(g *document.Diagram) ([]byte, bool) {
		return []byte(`<?xml version="1.0"?><svg id="` + g.ID + `"></svg>`), true
	}})
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != html {
		t.Error("the deck is not deterministic")
	}
}
