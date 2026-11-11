package problemtree

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
)

// Format is how the diagram in a report is written.
type Format string

// The diagram formats.
const (
	FormatMermaid Format = "mermaid"
	FormatD2      Format = "d2"
)

// Report renders the problem tree as Markdown: the diagram in format f, then a
// table of the causes and one of the effects, each with what it leads to and
// the outcomes that address it, or why it is set aside.
func Report(doc *impact.Model, f Format) (string, error) {
	return ReportFor(doc, "", f)
}

// ReportFor renders the problem tree of one pathway (see BuildFor), or of the
// whole model when pathwayID is empty. A pathway's tables list only its own
// causes and effects. Its own outcomes are listed first as addressing them,
// then any other outcome, labelled with the pathway it belongs to.
func ReportFor(doc *impact.Model, pathwayID string, f Format) (string, error) {
	t, err := BuildFor(doc, pathwayID)
	if err != nil {
		return "", err
	}
	heading, problem := "Problem tree", doc.Context.ProblemStatement
	var ownOutcomes []string
	if pathwayID != "" {
		p, _ := doc.Pathway(pathwayID) // BuildFor checked it exists
		heading, problem, ownOutcomes = "Problem tree: "+oneLine(p.Title), ProblemStatement(doc, p), p.OutcomeIDs
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", heading)
	if t.Empty() {
		fmt.Fprintf(&b, "No causes or effects have been recorded for the core problem: %s\n\n", oneLine(problem))
		return b.String(), nil
	}
	if pathwayID != "" {
		fmt.Fprintf(&b, "Core problem: %s\n\n", oneLine(problem))
	}
	inTreeIDs := map[string]bool{}
	for _, row := range t.Rows {
		for _, n := range row.Nodes {
			inTreeIDs[n.ID] = true
		}
	}
	var diagram string
	switch f {
	case FormatMermaid:
		var err error
		if diagram, err = t.Mermaid(); err != nil {
			return "", err
		}
	case FormatD2:
		diagram = t.D2()
	default:
		return "", fmt.Errorf("problemtree: unknown diagram format %q", f)
	}
	fmt.Fprintf(&b, "Causes (roots) below lead up to the core problem, which leads up to its effects (branches).\n\n```%s\n%s```\n\n", f, diagram)

	// In a pathway's tree its own outcomes come first; an outcome of another
	// pathway, or of none, follows with where it belongs, so a shared cause
	// does not look unaddressed.
	addressedBy := map[string][]string{}
	for _, own := range []bool{true, false} {
		for _, o := range doc.Outcomes {
			if o.Rejected() || (pathwayID != "" && slices.Contains(ownOutcomes, o.ID) != own) || (pathwayID == "" && !own) {
				continue
			}
			name := fmt.Sprintf("%s (`%s`)", cell(o.Title), o.ID)
			if pathwayID != "" && !own {
				name = fmt.Sprintf("%s (`%s`, %s)", cell(o.Title), o.ID, outcomeHome(doc, o.ID))
			}
			for _, id := range o.AddressesCauses {
				addressedBy[id] = append(addressedBy[id], name)
			}
			for _, id := range o.AddressesEffects {
				addressedBy[id] = append(addressedBy[id], name)
			}
		}
	}
	statement := map[string]string{ProblemID: "the core problem"}
	for _, c := range doc.Causes {
		statement[c.ID] = c.Statement
	}
	for _, e := range doc.Effects {
		statement[e.ID] = e.Statement
	}
	targets := func(ids []string) string {
		var out []string
		for _, id := range ids {
			if inTreeIDs[id] {
				out = append(out, fmt.Sprintf("%s (`%s`)", cell(statement[id]), id))
			}
		}
		if len(out) == 0 {
			return "the core problem"
		}
		return strings.Join(out, "; ")
	}
	answer := func(id string, aside *impact.OutOfScope) string {
		switch {
		case aside != nil:
			return "set aside: " + cell(aside.Reason)
		case len(addressedBy[id]) == 0:
			return "—"
		}
		return strings.Join(addressedBy[id], "; ")
	}

	var causes, effects strings.Builder
	for _, c := range doc.Causes {
		if inTreeIDs[c.ID] {
			fmt.Fprintf(&causes, "| %s (`%s`) | %s | %s | %s |\n", cell(c.Statement), c.ID, c.Level, targets(c.Parents), answer(c.ID, c.OutOfScope))
		}
	}
	for _, e := range doc.Effects {
		if inTreeIDs[e.ID] {
			fmt.Fprintf(&effects, "| %s (`%s`) | %s | %s |\n", cell(e.Statement), e.ID, targets(e.Parents), answer(e.ID, e.OutOfScope))
		}
	}
	if causes.Len() > 0 {
		b.WriteString("### Causes\n\n| Cause | Level | Leads to | Addressed by |\n|-------|-------|----------|--------------|\n")
		b.WriteString(causes.String() + "\n")
	}
	if effects.Len() > 0 {
		b.WriteString("### Effects\n\n| Effect | Follows from | Relieved by |\n|--------|--------------|-------------|\n")
		b.WriteString(effects.String() + "\n")
	}
	return b.String(), nil
}

// outcomeHome names the pathways an outcome belongs to, for a label such as
// "in Network and Evidence", or says it belongs to none.
func outcomeHome(doc *impact.Model, outcomeID string) string {
	var titles []string
	for _, p := range doc.Pathways {
		if !p.Rejected() && slices.Contains(p.OutcomeIDs, outcomeID) {
			titles = append(titles, oneLine(p.Title))
		}
	}
	if len(titles) == 0 {
		return "outside any pathway"
	}
	return "in " + cell(strings.Join(titles, ", "))
}

// cell makes text safe for a Markdown table cell.
func cell(s string) string { return strings.ReplaceAll(oneLine(s), "|", `\|`) }
