package model

import "slices"

// ScorecardObjectives returns the outcomes that are objectives of the
// model's Balanced Scorecard: those with a scorecard perspective that review
// has not rejected, in file order.
func (d *Model) ScorecardObjectives() []Outcome {
	var out []Outcome
	for _, o := range d.Outcomes {
		if o.ScorecardPerspective != "" && !o.Rejected() {
			out = append(out, o)
		}
	}
	return out
}

// ScorecardContributesTo returns the scorecard objectives the outcome leads
// to most directly: the strategy map's arrows from it. It follows causal links
// that are not rejected or feedback, and precondition groups, through
// outcomes that are not objectives, and stops at the first objective on each
// path. The result is in file order and never includes the outcome itself.
func (d *Model) ScorecardContributesTo(outcomeID string) []string {
	reached := d.scorecardReach(d.successors(outcomeID))
	return slices.DeleteFunc(reached, func(id string) bool { return id == outcomeID })
}

// ScorecardObjectivesServed returns the scorecard objectives an intervention
// serves: those among its targets, and otherwise the nearest objectives its
// targets and its own links lead to, in file order.
func (d *Model) ScorecardObjectivesServed(in Intervention) []string {
	return d.scorecardReach(append(slices.Clone(in.Targets), d.successors(in.ID)...))
}

// scorecardReach walks forward from the given outcomes and returns the first
// scorecard objective on each path, the starting outcomes included, in file
// order.
func (d *Model) scorecardReach(start []string) []string {
	objective := map[string]bool{}
	for _, o := range d.ScorecardObjectives() {
		objective[o.ID] = true
	}
	seen := map[string]bool{}
	found := map[string]bool{}
	queue := slices.Clone(start)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		seen[id] = true
		if o, ok := d.Outcome(id); !ok || o.Rejected() {
			continue
		}
		if objective[id] {
			found[id] = true
			continue
		}
		queue = append(queue, d.successors(id)...)
	}
	var out []string
	for _, o := range d.ScorecardObjectives() {
		if found[o.ID] {
			out = append(out, o.ID)
		}
	}
	return out
}

// successors returns what an outcome or intervention leads to directly:
// the targets of its causal links that are not rejected or feedback, and of
// the precondition groups it is a member of.
func (d *Model) successors(id string) []string {
	var out []string
	for _, l := range d.Links {
		if l.From == id && !l.Rejected() && !l.Feedback {
			out = append(out, l.To)
		}
	}
	for _, g := range d.Groups {
		if slices.Contains(g.Members, id) && !g.Rejected() {
			out = append(out, g.Target)
		}
	}
	return out
}
