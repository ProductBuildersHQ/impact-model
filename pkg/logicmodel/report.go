package logicmodel

import (
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// Format is how a diagram is written.
type Format string

// The formats a diagram can be written in.
const (
	FormatMermaid Format = "mermaid"
	FormatD2      Format = "d2"
)

// Diagrams renders the layered logic model diagram of each pathway, or of the
// whole outcome model when it has no pathways, each under a heading and in a fenced
// block that names its format.
func Diagrams(whole *report.View, pathways []*report.View, f Format) (string, error) {
	views := pathways
	if len(views) == 0 {
		views = []*report.View{whole}
	}
	var b strings.Builder
	for _, v := range views {
		g := Build(v)
		var text string
		switch f {
		case FormatMermaid:
			var err error
			if text, err = g.Mermaid(); err != nil {
				return "", err
			}
		case FormatD2:
			text = g.D2()
		default:
			return "", fmt.Errorf("logicmodel: unknown diagram format %q", f)
		}
		fmt.Fprintf(&b, "## Layered logic model: %s\n\n```%s\n%s```\n\n", strings.Join(strings.Fields(g.Title), " "), f, text)
		b.WriteString(feedbackNote(g))
	}
	return b.String(), nil
}

// TheoryOfChangeDiagram renders the theory of change of a view: for the whole
// outcome model, every pathway's activities and outcomes together on one page;
// for a pathway's view, that pathway's. It is drawn bottom-up (see
// Graph.Spec) in a fenced block that names its format, under a heading,
// then the feedback loops it leaves out.
func TheoryOfChangeDiagram(v *report.View, f Format) (string, error) {
	g := Build(v).TheoryOfChange()
	var text string
	switch f {
	case FormatMermaid:
		var err error
		if text, err = g.Mermaid(); err != nil {
			return "", err
		}
	case FormatD2:
		text = g.D2()
	default:
		return "", fmt.Errorf("logicmodel: unknown diagram format %q", f)
	}
	heading := "Theory of change on one page"
	if v.Filter != "" {
		heading = "Theory of change: " + strings.Join(strings.Fields(g.Title), " ")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n```%s\n%s```\n\n", heading, f, text)
	b.WriteString(feedbackNote(g))
	return b.String(), nil
}

// TheoryOfChangeDiagrams renders the whole outcome model's theory of change on
// one page, then each pathway's, which are small enough to read on a page.
func TheoryOfChangeDiagrams(whole *report.View, pathways []*report.View, f Format) (string, error) {
	var b strings.Builder
	for _, v := range append([]*report.View{whole}, pathways...) {
		text, err := TheoryOfChangeDiagram(v, f)
		if err != nil {
			return "", err
		}
		b.WriteString(text)
	}
	return b.String(), nil
}

// feedbackNote lists the feedback loops the diagram leaves out.
func feedbackNote(g *Graph) string {
	if len(g.Feedback) == 0 {
		return ""
	}
	label := map[string]string{}
	for _, st := range g.Stages {
		for _, n := range st.Nodes {
			label[n.ID] = n.Label
		}
	}
	var b strings.Builder
	b.WriteString("Feedback, not drawn:\n\n")
	for _, e := range g.Feedback {
		fmt.Fprintf(&b, "- %s leads back to %s\n", label[e.From], label[e.To])
	}
	b.WriteString("\n")
	return b.String()
}
