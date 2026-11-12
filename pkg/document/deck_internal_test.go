package document

import (
	"strings"
	"testing"
)

func TestChunksTitleTextByHeading(t *testing.T) {
	md := "Lead text.\n\n### Strategy\n\n#### 1. One\n\nFirst.\n\n#### 2. Two\n\nSecond.\n\n### Metrics\n\nAll.\n\n```text\n### not a heading\n```\n"
	got := chunks(md, "VSEM", "VSEM")
	want := []chunk{
		{kicker: "", title: "VSEM", body: "Lead text.\n"},
		{kicker: "Strategy", title: "1. One", body: "First.\n"},
		{kicker: "Strategy", title: "2. Two", body: "Second.\n"},
		{kicker: "VSEM", title: "Metrics", body: "All.\n\n```text\n### not a heading\n```\n"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d chunks: %+v", len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("chunk %d = %+v, want %+v", i, got[i], want[i])
		}
	}
	after := chunks("A table.\n", "Logic Models", "Capacity Building")
	if len(after) != 1 || after[0].title != "Capacity Building" || after[0].kicker != "Logic Models" {
		t.Errorf("text after a diagram = %+v", after)
	}
}

func TestPaginateSplitsLongTablesAndRepeatsTheHeader(t *testing.T) {
	var b strings.Builder
	b.WriteString("| A | B |\n|---|---|\n")
	for i := 0; i < 40; i++ {
		b.WriteString("| row | value |\n")
	}
	pages := paginate("Intro.\n\n" + b.String())
	if len(pages) < 3 {
		t.Fatalf("40 rows fit on %d pages", len(pages))
	}
	rows := 0
	for i, p := range pages {
		if !strings.Contains(p, "| A | B |\n|---|---|") {
			t.Errorf("page %d has no table header:\n%s", i+1, p)
		}
		rows += strings.Count(p, "| row | value |")
	}
	if rows != 40 {
		t.Errorf("pages hold %d rows, want 40", rows)
	}
	if !strings.HasPrefix(pages[0], "Intro.") {
		t.Errorf("first page lost the intro:\n%s", pages[0])
	}
}

func TestPaginateKeepsShortTextOnOnePage(t *testing.T) {
	if pages := paginate("One.\n\nTwo.\n\n- a\n- b\n"); len(pages) != 1 {
		t.Errorf("short text split into %d pages", len(pages))
	}
}

func TestLineBreaksSurviveWithoutRawHTML(t *testing.T) {
	d := &Document{Title: "T", Sections: []Section{{Kind: SectionVSEM, Title: "VSEM", Parts: []Part{{
		Markdown: "### Work\n\n| Work |\n|---|\n| Name<br>Detail <script>x</script> |\n",
	}}}}}
	page, err := d.Deck(DeckOptions{})
	if err != nil {
		t.Fatal(err)
	}
	s := string(page)
	if !strings.Contains(s, `Name\u003cbr\u003eDetail`) {
		t.Error("the line break in a cell is lost")
	}
	if strings.Contains(s, `\u003cscript\u003ex`) {
		t.Error("raw HTML other than a line break got through")
	}
}
