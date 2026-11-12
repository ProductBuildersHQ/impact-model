package document

import (
	"fmt"
	"strings"
)

// MarkdownOptions control how a document is written as Markdown.
type MarkdownOptions struct {
	// DiagramPath, when set, gives the path of the rendered image for a
	// diagram, and the diagram is written as an image reference. When nil,
	// diagrams are written as fenced D2 blocks.
	DiagramPath func(*Diagram) string
	// ImageAttrs, such as "width=100%", is added to each diagram image as a
	// Pandoc attribute block. Leave it empty for plain Markdown.
	ImageAttrs string
	// Contents adds a list of the sections after the title.
	Contents bool
}

// Markdown writes the document as one Markdown file: the title, then each
// section under a level-2 heading.
func (d *Document) Markdown(opts MarkdownOptions) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", d.Title)
	if d.Subtitle != "" {
		fmt.Fprintf(&b, "%s\n\n", d.Subtitle)
	}
	if d.Stamp != "" {
		fmt.Fprintf(&b, "*%s*\n\n", d.Stamp)
	}
	if opts.Contents && len(d.Sections) > 1 {
		b.WriteString("## Contents\n\n")
		for i, s := range d.Sections {
			fmt.Fprintf(&b, "%d. %s\n", i+1, s.Title)
		}
		b.WriteString("\n")
	}
	for _, s := range d.Sections {
		fmt.Fprintf(&b, "## %s\n\n", s.Title)
		for _, p := range s.Parts {
			switch {
			case p.Diagram != nil && opts.DiagramPath != nil:
				attrs := ""
				if opts.ImageAttrs != "" {
					attrs = "{" + opts.ImageAttrs + "}"
				}
				fmt.Fprintf(&b, "![%s](%s)%s\n\n", altText(p.Diagram.Title), opts.DiagramPath(p.Diagram), attrs)
			case p.Diagram != nil:
				fmt.Fprintf(&b, "```d2\n%s```\n\n", p.Diagram.D2)
			default:
				b.WriteString(strings.TrimRight(p.Markdown, "\n"))
				b.WriteString("\n\n")
			}
		}
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// altText keeps an image's alt text on one line and free of brackets.
func altText(s string) string {
	s = strings.NewReplacer("[", "(", "]", ")").Replace(s)
	return oneLine(s)
}
