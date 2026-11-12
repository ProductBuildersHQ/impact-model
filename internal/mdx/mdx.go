// Package mdx adjusts generated Markdown so one report can sit inside a larger
// document: it finds ATX headings and fenced code blocks, drops a report's own
// title, moves headings to a given level, and splits out fenced diagrams.
// Text inside a fence is never treated as a heading.
package mdx

import (
	"slices"
	"strings"
)

// maxLevel is the deepest heading Markdown has.
const maxLevel = 6

// line is one line of a document and what kind of line it is.
type line struct {
	text  string
	level int  // heading level, 0 when the line is not a heading
	fence bool // the line opens, closes, or sits inside a fenced block
}

// scan splits md into lines and marks headings and fenced lines.
func scan(md string) []line {
	raw := strings.Split(md, "\n")
	out := make([]line, len(raw))
	inFence, marker := false, ""
	for i, t := range raw {
		trim := strings.TrimLeft(t, " ")
		switch {
		case inFence:
			out[i] = line{text: t, fence: true}
			if strings.HasPrefix(trim, marker) && strings.TrimSpace(trim[len(marker):]) == "" {
				inFence = false
			}
		case strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~"):
			inFence, marker = true, trim[:3]
			out[i] = line{text: t, fence: true}
		default:
			out[i] = line{text: t, level: headingLevel(t)}
		}
	}
	return out
}

// headingLevel returns the level of an ATX heading, or 0.
func headingLevel(t string) int {
	n := 0
	for n < len(t) && t[n] == '#' {
		n++
	}
	if n == 0 || n > maxLevel || n == len(t) || t[n] != ' ' {
		return 0
	}
	return n
}

func join(lines []line) string {
	parts := make([]string, len(lines))
	for i, l := range lines {
		parts[i] = l.text
	}
	return strings.Join(parts, "\n")
}

// StripTitle removes the first line when it is a heading, with the blank
// lines after it. A report's title repeats the heading its section already has.
func StripTitle(md string) string {
	lines := scan(strings.TrimLeft(md, "\n"))
	if len(lines) == 0 || lines[0].level == 0 {
		return md
	}
	i := 1
	for i < len(lines) && strings.TrimSpace(lines[i].text) == "" {
		i++
	}
	return join(lines[i:])
}

// MinLevel returns the level of the highest heading in md, or 0 when it has
// none.
func MinLevel(md string) int {
	minimum := 0
	for _, l := range scan(md) {
		if l.level > 0 && (minimum == 0 || l.level < minimum) {
			minimum = l.level
		}
	}
	return minimum
}

// Shift moves every heading n levels deeper (shallower when n is negative),
// keeping each between 1 and 6.
func Shift(md string, n int) string {
	if n == 0 {
		return md
	}
	lines := scan(md)
	for i, l := range lines {
		if l.level == 0 {
			continue
		}
		level := min(max(l.level+n, 1), maxLevel)
		lines[i].text = strings.Repeat("#", level) + l.text[l.level:]
	}
	return join(lines)
}

// Rebase shifts md so that its highest heading is at level base. Markdown
// without headings is returned as it is.
func Rebase(md string, base int) string {
	if m := MinLevel(md); m > 0 {
		return Shift(md, base-m)
	}
	return md
}

// Part is a piece of a document: either Markdown text or one fenced block in
// a diagram language.
type Part struct {
	// Text is the Markdown, when Lang is empty.
	Text string
	// Lang is the fence's language, such as d2, when the part is a diagram.
	Lang string
	// Source is the diagram source, without the fence lines.
	Source string
	// Title is the text of the last heading before the diagram.
	Title string
}

// SplitFences splits md around the fenced blocks whose language is one of
// langs. Text parts keep their blank lines trimmed; empty ones are dropped.
func SplitFences(md string, langs ...string) []Part {
	var parts []Part
	var text, src []string
	lang, title := "", ""
	flush := func() {
		if t := strings.TrimSpace(strings.Join(text, "\n")); t != "" {
			parts = append(parts, Part{Text: t + "\n"})
		}
		text = nil
	}
	for _, l := range scan(md) {
		trim := strings.TrimSpace(l.text)
		switch {
		case lang != "":
			if trim == "```" {
				parts = append(parts, Part{Lang: lang, Source: strings.Join(src, "\n") + "\n", Title: title})
				lang, src = "", nil
				continue
			}
			src = append(src, l.text)
		case l.fence && strings.HasPrefix(trim, "```") && slices.Contains(langs, strings.TrimSpace(trim[3:])):
			flush()
			lang = strings.TrimSpace(trim[3:])
		default:
			if l.level > 0 {
				title = strings.TrimSpace(l.text[l.level:])
			}
			text = append(text, l.text)
		}
	}
	flush()
	return parts
}
