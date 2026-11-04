package report

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
	"github.com/ProductBuildersHQ/impact-model/pkg/problemtree"
)

// DMAICPhaseUse records which DMAIC phases an indicator serves, by the
// convention DMAICPlan documents. Measure is always true: every metric needs
// a baseline.
type DMAICPhaseUse struct {
	Define, Measure, Analyze, Improve, Control bool
}

// DMAICUse derives the DMAIC phases an indicator serves within pathway p:
//
//   - Define: it measures the pathway itself, an outcome at the long-term or
//     impact level, or an outcome that relieves an effect of the problem, so
//     it sizes the gap the effort is meant to close.
//   - Measure: always.
//   - Analyze: it is broken down by group (disaggregation or target groups),
//     so it can show where the gap is largest.
//   - Improve: it measures one of the pathway's interventions, or an output,
//     short-term, or medium-term outcome that one of them is aimed at, so it
//     moves soon after a change is tried.
//   - Control: it is measured on a schedule at least once a year, so a gain
//     can be checked for whether it lasts.
func DMAICUse(doc *impact.Model, p *impact.Pathway, ind impact.Indicator) DMAICPhaseUse {
	u := DMAICPhaseUse{Measure: true}
	u.Analyze = len(ind.Disaggregation) > 0 || len(ind.TargetGroupIDs) > 0
	if n, ok := impact.TimesPerYear(ind.Frequency); ok && n >= 0.99 {
		u.Control = true
	}
	if ind.Measures == p.ID {
		u.Define = true
	}
	if slices.Contains(p.InterventionIDs, ind.Measures) {
		u.Improve = true
	}
	if o, ok := doc.Outcome(ind.Measures); ok {
		switch o.Level {
		case impact.LevelLongTerm, impact.LevelImpact:
			u.Define = true
		default:
			u.Improve = aimedAt(doc, p, o.ID)
		}
		if len(o.AddressesEffects) > 0 {
			u.Define = true
		}
	}
	return u
}

// aimedAt reports whether a live intervention of the pathway targets outcome.
func aimedAt(doc *impact.Model, p *impact.Pathway, outcome string) bool {
	for _, id := range p.InterventionIDs {
		in, ok := doc.Intervention(id)
		if ok && !in.Rejected() && in.Status != impact.InterventionCancelled && slices.Contains(in.Targets, outcome) {
			return true
		}
	}
	return false
}

// DMAICNeeds lists what an indicator still needs before the DMAIC export can
// chart it with control limits: a numeric baseline and target (without them
// the export leaves it out), a frequency and a data source (to measure it on
// a schedule), and at least impact.DMAICMinPoints measurements a year.
func DMAICNeeds(ind impact.Indicator) []string {
	var out []string
	if ind.Baseline == nil || ind.Baseline.Value == nil {
		out = append(out, "baseline")
	}
	hasTarget := false
	for _, t := range ind.Targets {
		hasTarget = hasTarget || t.Value != nil
	}
	if !hasTarget {
		out = append(out, "numeric target")
	}
	if ind.DataSource == "" && (ind.Collection == nil || ind.Collection.System == "") {
		out = append(out, "data source")
	}
	n, ok := impact.TimesPerYear(ind.Frequency)
	switch {
	case !ok:
		out = append(out, "frequency")
	case n < impact.DMAICMinPoints:
		out = append(out, fmt.Sprintf("%d+ measurements a year for control limits", impact.DMAICMinPoints))
	}
	return out
}

// DMAICPlan renders the DMAIC plan of an impact model: for each pathway
// marked `projections: [dmaic]`, or for the one pathway named, the problem
// the effort improves (Define), the interventions it can change (Improve),
// and a table of its indicators with the DMAIC phases each serves and what
// each still needs before `export dmaic` can chart it. It needs no
// observations, so a team can plan DMAIC before any data is collected.
func DMAICPlan(doc *impact.Model, pathway string) (string, error) {
	var pathways []*impact.Pathway
	if pathway != "" {
		p, ok := doc.Pathway(pathway)
		if !ok {
			return "", fmt.Errorf("report: unknown pathway %q", pathway)
		}
		pathways = append(pathways, p)
	} else {
		for i := range doc.Pathways {
			p := &doc.Pathways[i]
			if !p.Rejected() && slices.Contains(p.Projections, impact.ProjectionDMAIC) {
				pathways = append(pathways, p)
			}
		}
	}

	var b strings.Builder
	b.WriteString("## DMAIC plan\n\n")
	if len(pathways) == 0 {
		b.WriteString("No pathway is marked for DMAIC. Add `projections: [dmaic]` to a pathway that is measured often enough to improve with DMAIC, or name one with `--pathway`.\n\n")
		return b.String(), nil
	}
	b.WriteString("DMAIC (Define, Measure, Analyze, Improve, Control) is a cycle for improving a program that is measured often. ")
	b.WriteString("One metric can serve several phases. A metric serves **Define** when it sizes the problem (a long-term or impact outcome, or one that relieves an effect), ")
	b.WriteString("**Measure** always, **Analyze** when it is broken down by group, **Improve** when it moves soon after an intervention changes, ")
	b.WriteString("and **Control** when it is measured at least once a year. ")
	fmt.Fprintf(&b, "`export dmaic` includes a metric once it has a numeric baseline and target, and draws control limits once it has %d measurements.\n\n", impact.DMAICMinPoints)

	for _, p := range pathways {
		if err := dmaicPathway(&b, doc, p); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func dmaicPathway(b *strings.Builder, doc *impact.Model, p *impact.Pathway) error {
	fmt.Fprintf(b, "### %s\n\n", mdText(p.Title))

	fmt.Fprintf(b, "**Define.** Core problem: %s\n\n", mdText(problemtree.ProblemStatement(doc, p)))
	tree, err := problemtree.BuildFor(doc, p.ID)
	if err != nil {
		return err
	}
	var causes []string
	for i := len(tree.Rows) - 1; i >= 0; i-- { // deepest causes first
		for _, n := range tree.Rows[i].Nodes {
			if n.Kind == problemtree.KindCause && !n.SetAside {
				if c, ok := doc.Cause(n.ID); ok {
					causes = append(causes, fmt.Sprintf("%s (`%s`, %s)", mdText(c.Statement), c.ID, c.Level))
				}
			}
		}
	}
	if len(causes) > 0 {
		b.WriteString("Root causes to analyze, deepest first:\n\n")
		for _, c := range causes {
			fmt.Fprintf(b, "- %s\n", c)
		}
		b.WriteString("\n")
	}

	var interventions []string
	for _, id := range p.InterventionIDs {
		in, ok := doc.Intervention(id)
		if !ok || in.Rejected() || in.Status == impact.InterventionCancelled {
			continue
		}
		s := fmt.Sprintf("%s (`%s`)", mdText(in.Title), in.ID)
		if in.Status != "" {
			s += ", " + strings.ReplaceAll(string(in.Status), "_", " ")
		}
		interventions = append(interventions, s)
	}
	if len(interventions) > 0 {
		b.WriteString("**Improve.** The interventions an improvement would change:\n\n")
		for _, s := range interventions {
			fmt.Fprintf(b, "- %s\n", s)
		}
		b.WriteString("\n")
	}

	var rows strings.Builder
	total, withBaseline, withTarget, chartable := 0, 0, 0, 0
	for _, ind := range doc.Indicators {
		if ind.Rejected() || !doc.IndicatorServes(ind, p.ID) {
			continue
		}
		total++
		u := DMAICUse(doc, p, ind)
		needs := DMAICNeeds(ind)
		if !slices.Contains(needs, "baseline") {
			withBaseline++
		}
		if !slices.Contains(needs, "numeric target") {
			withTarget++
		}
		if len(needs) == 0 {
			chartable++
		}
		mark := func(on bool) string {
			if on {
				return "✓"
			}
			return ""
		}
		fmt.Fprintf(&rows, "| %s (`%s`) | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			MDCell(ind.Name), ind.ID, mark(u.Define), mark(u.Measure), mark(u.Analyze), mark(u.Improve), mark(u.Control),
			orDash(MDCell(frequencyWords(ind.Frequency))), orDash(string(ind.DataAvailability)), orDash(strings.Join(needs, "; ")))
	}
	if total == 0 {
		b.WriteString("No indicator serves this pathway yet.\n\n")
		return nil
	}
	fmt.Fprintf(b, "**Metrics.** %d metrics: %d have a baseline, %d a numeric target, and %d have everything control limits need.\n\n", total, withBaseline, withTarget, chartable)
	b.WriteString("| Metric | D | M | A | I | C | Cadence | Data | Still needs |\n")
	b.WriteString("|--------|---|---|---|---|---|---------|------|-------------|\n")
	b.WriteString(rows.String())
	b.WriteString("\n")
	return nil
}

// frequencyWords names the common ISO 8601 cadences in words and passes any
// other through as written.
func frequencyWords(iso string) string {
	names := map[string]string{
		"P1D": "daily", "P1W": "weekly", "P2W": "every two weeks", "P1M": "monthly",
		"P3M": "quarterly", "P6M": "every six months", "P1Y": "yearly", "P2Y": "every two years",
	}
	if n, ok := names[iso]; ok {
		return n
	}
	return iso
}
