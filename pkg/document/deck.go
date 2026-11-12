package document

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed deck.html
var deckTemplate string

// slideBudget is roughly how many lines of text fit on one slide. Longer
// content is carried onto further slides.
const slideBudget = 14

// lineWidth is roughly how many characters fit on one line of a slide, and
// cellWidth how many of a table row's characters fit on one line, since a
// row's text wraps within its columns.
const (
	lineWidth = 110
	cellWidth = 70
)

// lineBreakMark stands in for a <br> while Markdown is converted.
const lineBreakMark = "\uE000"

var lineBreaks = strings.NewReplacer("<br>", lineBreakMark, "<br/>", lineBreakMark, "<br />", lineBreakMark)

// DeckOptions control the HTML presentation.
type DeckOptions struct {
	// SVG returns the rendered SVG of a diagram. A diagram without one is
	// shown as its D2 source.
	SVG func(*Diagram) ([]byte, bool)
}

// slide is one slide as the deck's script reads it.
type slide struct {
	Kind     string   `json:"kind"`
	Kicker   string   `json:"kicker,omitempty"`
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	Number   int      `json:"number,omitempty"`
	Items    []string `json:"items,omitempty"`
	HTML     string   `json:"html,omitempty"`
	SVG      string   `json:"svg,omitempty"`
}

// deckData is everything the deck's script needs.
type deckData struct {
	Title  string  `json:"title"`
	Slides []slide `json:"slides"`
}

// Deck writes the document as a single-file HTML presentation built with
// React: a title slide, a contents slide, then for each section a divider and
// its content, one heading per slide, with each diagram on a slide of its
// own. Long text and tables run on to further slides. Arrow keys, Page
// Up/Down, Space, Home, and End move between slides, and the address keeps
// the slide number. Opened with ?print, it lays every slide out as a fixed
// 1280×720 page and sets window.__deckReady, the deck-pdf/v1 contract, so a
// headless browser can print it to PDF. React is loaded from cdnjs; all
// content is in the file.
func (d *Document) Deck(opts DeckOptions) ([]byte, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	// Raw HTML stays escaped, except the line breaks the table renderers
	// write in cells, which are carried through as a private-use character.
	toHTML := func(src string) (string, error) {
		var buf bytes.Buffer
		if err := md.Convert([]byte(lineBreaks.Replace(src)), &buf); err != nil {
			return "", fmt.Errorf("document: render Markdown: %w", err)
		}
		return strings.ReplaceAll(buf.String(), lineBreakMark, "<br>"), nil
	}

	data := deckData{Title: d.Title}
	data.Slides = append(data.Slides, slide{Kind: "title", Title: d.Title, Subtitle: d.Subtitle, Kicker: d.Stamp})
	if len(d.Sections) > 1 {
		contents := slide{Kind: "agenda", Title: "Contents"}
		for _, s := range d.Sections {
			contents.Items = append(contents.Items, s.Title)
		}
		data.Slides = append(data.Slides, contents)
	}
	for i, s := range d.Sections {
		data.Slides = append(data.Slides, slide{Kind: "section", Number: i + 1, Title: s.Title})
		// lead titles text that comes before any heading of its own: the
		// diagram it follows, or else the section.
		lead := s.Title
		for _, p := range s.Parts {
			if p.Diagram != nil {
				lead = plain(p.Diagram.Title)
				sl := slide{Kind: "diagram", Kicker: s.Title, Title: lead}
				if svg, ok := svgFor(opts, p.Diagram); ok {
					sl.SVG = svg
				} else {
					h, err := toHTML("```d2\n" + p.Diagram.D2 + "```\n")
					if err != nil {
						return nil, err
					}
					sl.Kind, sl.HTML = "content", h
				}
				data.Slides = append(data.Slides, sl)
				continue
			}
			for _, c := range chunks(p.Markdown, s.Title, lead) {
				for _, page := range paginate(c.body) {
					h, err := toHTML(page)
					if err != nil {
						return nil, err
					}
					data.Slides = append(data.Slides, slide{Kind: "content", Kicker: c.kicker, Title: c.title, HTML: h})
				}
			}
		}
	}
	js, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("document: encode deck: %w", err)
	}
	out := strings.NewReplacer(
		"{{TITLE}}", html.EscapeString(d.Title),
		"{{DATA}}", string(js),
	).Replace(deckTemplate)
	return []byte(out), nil
}

func svgFor(opts DeckOptions, d *Diagram) (string, bool) {
	if opts.SVG == nil {
		return "", false
	}
	svg, ok := opts.SVG(d)
	if !ok || len(svg) == 0 {
		return "", false
	}
	return stripXMLProlog(string(svg)), true
}

// stripXMLProlog drops an XML declaration, which is not allowed inside HTML.
func stripXMLProlog(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "<?xml") {
		if i := strings.Index(s, "?>"); i >= 0 {
			s = strings.TrimSpace(s[i+2:])
		}
	}
	return s
}

// chunk is the text under one heading.
type chunk struct {
	kicker, title, body string
}

var headingRE = regexp.MustCompile(`^(#{1,6}) +(.*?)\s*#*\s*$`)

// chunks splits Markdown at its headings. The text before the first heading
// is titled lead, under the section's title. A heading with nothing under it, such as
// "Strategy" over one heading per strategy, becomes the kicker of the deeper
// headings that follow it; otherwise the kicker is the section's title.
func chunks(md, sectionTitle, lead string) []chunk {
	type open struct {
		level int
		title string
	}
	var out []chunk
	var parents []open
	cur, level := chunk{title: lead}, 0
	if lead != sectionTitle {
		cur.kicker = sectionTitle
	}
	var body []string
	inFence := false
	kicker := func() string {
		if len(parents) == 0 {
			return sectionTitle
		}
		return parents[len(parents)-1].title
	}
	flush := func() {
		text := strings.TrimSpace(strings.Join(body, "\n"))
		if text != "" {
			out = append(out, chunk{kicker: cur.kicker, title: cur.title, body: text + "\n"})
		} else if level > 0 {
			parents = append(parents, open{level: level, title: cur.title})
		}
		body = nil
	}
	for _, line := range strings.Split(md, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
		}
		if m := headingRE.FindStringSubmatch(line); m != nil && !inFence {
			flush()
			level = len(m[1])
			for len(parents) > 0 && parents[len(parents)-1].level >= level {
				parents = parents[:len(parents)-1]
			}
			cur = chunk{kicker: kicker(), title: plain(m[2])}
			continue
		}
		body = append(body, line)
	}
	flush()
	return out
}

// paginate splits text into pages of about slideBudget lines, at blank lines
// or, inside a table, between rows, repeating the table's header.
func paginate(text string) []string {
	var pages, cur []string
	used := 0
	emit := func() {
		if len(cur) > 0 {
			pages = append(pages, strings.Join(cur, "\n\n")+"\n")
			cur, used = nil, 0
		}
	}
	for _, blk := range splitBlocks(text) {
		lines := strings.Split(blk, "\n")
		if !isTable(lines) || len(lines) <= 2 {
			h := height(blk)
			if used > 0 && used+h > slideBudget {
				emit()
			}
			cur = append(cur, blk)
			used += h
			continue
		}
		header, rows := lines[:2], lines[2:]
		for len(rows) > 0 {
			// Start a new page rather than put a header and one row at the
			// foot of a page.
			if used > 0 && used+rowHeight(header[0])+rowHeight(rows[0]) > slideBudget {
				emit()
			}
			part := append([]string{}, header...)
			used += rowHeight(header[0])
			for len(rows) > 0 && (len(part) == 2 || used+rowHeight(rows[0]) <= slideBudget) {
				part = append(part, rows[0])
				used += rowHeight(rows[0])
				rows = rows[1:]
			}
			cur = append(cur, strings.Join(part, "\n"))
			if len(rows) > 0 {
				emit()
			}
		}
	}
	emit()
	return pages
}

// rowHeight estimates how many slide lines a table row takes once its
// cells wrap.
func rowHeight(row string) int { return 1 + len(row)/cellWidth }

// splitBlocks splits Markdown at blank lines, keeping fenced blocks whole.
func splitBlocks(text string) []string {
	var blocks, cur []string
	inFence := false
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
		}
		if t == "" && !inFence {
			if len(cur) > 0 {
				blocks = append(blocks, strings.Join(cur, "\n"))
				cur = nil
			}
			continue
		}
		cur = append(cur, line)
	}
	if len(cur) > 0 {
		blocks = append(blocks, strings.Join(cur, "\n"))
	}
	return blocks
}

func isTable(lines []string) bool {
	return len(lines) >= 2 && strings.HasPrefix(strings.TrimSpace(lines[0]), "|") &&
		strings.HasPrefix(strings.TrimSpace(lines[1]), "|") && strings.Contains(lines[1], "--")
}

// height estimates how many slide lines a block takes.
func height(blk string) int {
	h := 0
	for _, line := range strings.Split(blk, "\n") {
		h += 1 + len(line)/lineWidth
	}
	return h
}

var inlineMarks = strings.NewReplacer("**", "", "__", "", "`", "")

// plain removes inline Markdown marks from a heading.
func plain(s string) string { return oneLine(inlineMarks.Replace(s)) }
