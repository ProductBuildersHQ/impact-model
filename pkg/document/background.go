package document

import (
	"fmt"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// background renders the Background section for a reader new to the
// organization: who it is, its vision, mission, and values, its programs, and
// how the sections that follow serve its mission and vision. Everything is
// drawn from the outcome model, so the same model always gives the same text.
// following are the sections after the background, in order.
func background(in Input, following []SectionKind) string {
	doc, ctx := in.Model, in.Model.Context
	org := organization(in)
	var b strings.Builder

	fmt.Fprintf(&b, "## About %s\n\n", org.name)
	about := paragraphs(ctx.Background)
	if len(about) == 0 && org.description != "" {
		about = []string{introduce(org.name, org.description)}
	}
	for _, p := range about {
		b.WriteString(p + "\n\n")
	}
	if ctx.TargetPopulation != "" {
		fmt.Fprintf(&b, "**Whom it serves.** %s\n\n", sentence(ctx.TargetPopulation))
	}
	if ctx.Geography != "" {
		fmt.Fprintf(&b, "**Where.** %s\n\n", sentence(ctx.Geography))
	}
	if period := period(ctx.TimeHorizon); period != "" {
		fmt.Fprintf(&b, "**Period of this model.** %s.\n\n", period)
	}

	if ctx.Vision != "" || ctx.Mission != "" || len(ctx.Values) > 0 {
		b.WriteString("## Vision, mission, and values\n\n")
		if ctx.Vision != "" {
			fmt.Fprintf(&b, "**Vision.** %s\n\n", sentence(ctx.Vision))
		}
		if ctx.Mission != "" {
			fmt.Fprintf(&b, "**Mission.** %s\n\n", sentence(ctx.Mission))
		}
		b.WriteString(values(ctx.Values))
	}

	b.WriteString(programs(doc, org.name))
	b.WriteString(howItServes(in, org.name, following))
	return b.String()
}

// org is the organization whose outcome model this is.
type org struct{ name, description string }

// organization names the perspective's actor, or else says "the
// organization".
func organization(in Input) org {
	if p := in.Model.Perspective; p != nil {
		for _, a := range in.Model.Actors {
			if a.ID == p.ActorID && a.Name != "" {
				return org{name: oneLine(a.Name), description: a.Description}
			}
		}
	}
	return org{name: "the organization"}
}

// values lists the values: as a sentence when none is described, otherwise
// one per line with its description.
func values(vs []model.Value) string {
	if len(vs) == 0 {
		return ""
	}
	described := false
	names := make([]string, len(vs))
	for i, v := range vs {
		names[i] = oneLine(v.Name)
		described = described || strings.TrimSpace(v.Description) != ""
	}
	if !described {
		return fmt.Sprintf("**Values.** %s.\n\n", series(names))
	}
	var b strings.Builder
	b.WriteString("**Values.**\n\n")
	for i, v := range vs {
		if d := oneLine(v.Description); d != "" {
			fmt.Fprintf(&b, "- **%s.** %s\n", names[i], sentence(d))
		} else {
			fmt.Fprintf(&b, "- **%s**\n", names[i])
		}
	}
	b.WriteString("\n")
	return b.String()
}

// programs describes each pathway that is not rejected, with the activities
// it carries out.
func programs(doc *model.Model, name string) string {
	var ps []model.Pathway
	for _, p := range doc.Pathways {
		if !p.Rejected() {
			ps = append(ps, p)
		}
	}
	if len(ps) == 0 {
		return ""
	}
	titles := map[string]string{}
	for _, iv := range doc.Interventions {
		if !iv.Rejected() {
			titles[iv.ID] = oneLine(iv.Title)
		}
	}
	var b strings.Builder
	b.WriteString("## Programs\n\n")
	fmt.Fprintf(&b, "%s works through %s. Each is a pathway in the theory of change that follows.\n\n",
		capitalize(name), countKinds(ps))
	for _, p := range ps {
		heading := oneLine(p.Title)
		if p.Kind != "" && p.Kind != model.PathwayProgram {
			heading += " (" + string(p.Kind) + ")"
		}
		fmt.Fprintf(&b, "### %s\n\n", heading)
		if d := oneLine(p.Description); d != "" {
			b.WriteString(sentence(d) + "\n\n")
		}
		var acts []string
		for _, id := range p.InterventionIDs {
			if t, ok := titles[id]; ok {
				acts = append(acts, t)
			}
		}
		if len(acts) > 0 {
			b.WriteString("**Activities:**\n\n")
			for _, a := range acts {
				b.WriteString("- " + a + "\n")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

// countKinds says how many pathways there are of each kind, such as "three
// programs and one strategy", in the order the kinds first appear.
func countKinds(ps []model.Pathway) string {
	var order []model.PathwayKind
	n := map[model.PathwayKind]int{}
	for _, p := range ps {
		k := p.Kind
		if k == "" {
			k = model.PathwayProgram
		}
		if n[k] == 0 {
			order = append(order, k)
		}
		n[k]++
	}
	parts := make([]string, len(order))
	for i, k := range order {
		noun := string(k)
		if n[k] != 1 {
			noun += "s"
		}
		parts[i] = number(n[k]) + " " + noun
	}
	return series(parts)
}

// howItServes ties the outcome model to the mission and vision, and says what
// each section that follows shows.
func howItServes(in Input, name string, following []SectionKind) string {
	doc, ctx := in.Model, in.Model.Context
	var b strings.Builder
	b.WriteString("## How this document serves the mission and vision\n\n")
	b.WriteString("This document is drawn from one outcome model: a single description of the change " +
		name + " works toward, the steps that lead to it, and how progress is measured. ")
	b.WriteString("Every section is generated from that model, so the views agree with each other and change together when the model changes.\n\n")

	var impacts []string
	for _, o := range in.Whole.Outcomes {
		if o.Level == model.LevelImpact {
			impacts = append(impacts, oneLine(o.Title))
		}
	}
	if ctx.Vision != "" && len(impacts) > 0 {
		fmt.Fprintf(&b, "The vision is where the work leads. In the model it is reached through %s:\n\n",
			plural(len(impacts), "one impact outcome", number(len(impacts))+" impact outcomes"))
		for _, t := range impacts {
			b.WriteString("- " + t + "\n")
		}
		b.WriteString("\n")
	}
	if ctx.Mission != "" && len(doc.Pathways) > 0 {
		b.WriteString("The mission is carried out through the programs above. Each is a pathway in the model, from its activities through short-, medium-, and long-term outcomes to the impact.\n\n")
	}
	if o := doc.Context.Objective; o != nil {
		if s := objectiveText(o, in.Whole); s != "" {
			fmt.Fprintf(&b, "For this period the model sets one objective: %s\n\n", sentence(s))
		}
	}
	if c := doc.Ceiling; c != nil && c.Level != "" {
		fmt.Fprintf(&b, "%s answers for its results up to and including its %ss. Above that, it contributes to change alongside others and tracks it without claiming it as its own.\n\n",
			capitalize(name), strings.ToLower(report.LevelLabel(c.Level)))
	}

	if len(following) > 0 {
		b.WriteString("The sections that follow look at the same model from different angles:\n\n")
		b.WriteString("| Section | What it shows |\n|---------|---------------|\n")
		for _, k := range following {
			fmt.Fprintf(&b, "| %s | %s |\n", k.Title(), k.Purpose())
		}
		b.WriteString("\n")
	}
	return b.String()
}

// objectiveText is the objective's statement, or the title of its one outcome.
func objectiveText(o *model.Objective, v *report.View) string {
	if s := oneLine(o.Statement); s != "" {
		return s
	}
	if len(o.OutcomeIDs) == 1 {
		for _, oc := range v.Outcomes {
			if oc.ID == o.OutcomeIDs[0] {
				return oneLine(oc.Title)
			}
		}
	}
	return ""
}

// period writes a time horizon as months, such as "July 2025 to June 2027".
func period(t *model.TimeHorizon) string {
	if t == nil {
		return ""
	}
	start, err1 := time.Parse(time.DateOnly, t.Start)
	end, err2 := time.Parse(time.DateOnly, t.End)
	if err1 != nil || err2 != nil {
		return ""
	}
	return start.Format("January 2006") + " to " + end.Format("January 2006")
}

// paragraphs splits text at blank lines, each paragraph on one line.
func paragraphs(text string) []string {
	var out []string
	for _, p := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n\n") {
		if p = oneLine(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// sentence ends text with a full stop unless it already ends a sentence.
func sentence(s string) string {
	s = oneLine(s)
	if s == "" || strings.HasSuffix(s, ".") || strings.HasSuffix(s, "!") || strings.HasSuffix(s, "?") {
		return s
	}
	return s + "."
}

// series joins items as "a", "a and b", or "a, b, and c".
func series(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	}
	return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
}

// introduce turns an actor's description into a sentence about it: "The
// convener ..." becomes "Thrive Alliance is the convener ...".
func introduce(name, description string) string {
	d := sentence(description)
	for _, article := range []string{"The ", "A ", "An "} {
		if strings.HasPrefix(d, article) {
			return name + " is " + strings.ToLower(article[:1]) + d[1:]
		}
	}
	return d
}

// capitalize raises the first letter.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// number writes small counts as words.
func number(n int) string {
	words := []string{"no", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	if n >= 0 && n < len(words) {
		return words[n]
	}
	return fmt.Sprint(n)
}

// plural picks the singular or plural form for n.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
