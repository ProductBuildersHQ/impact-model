package report

import (
	"fmt"
	"slices"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// Steering roles, as a reader sees them.
const (
	RoleNorthStar       = "North Star"
	RoleOneMetric       = "One metric that matters"
	RoleNorthStarAndOne = "North Star and one metric that matters"
	RoleGuardrail       = "Guardrail"
)

// SteeringMetric is one indicator a pathway is steered by, with its status.
type SteeringMetric struct {
	PathwayID   string
	Pathway     string
	Role        string
	IndicatorID string
	Indicator   string
	// State is the indicator's state in the snapshot, or empty when the
	// snapshot has no status for it.
	State string
	// Progress is how far the indicator has come toward its target and how
	// far it was expected to be, as text, or "—".
	Progress string
}

// Steering returns, for each pathway review has not rejected, in file order,
// what it is steered by: its North Star, then its one metric that matters on
// the snapshot's date when that is a different indicator, then each guardrail that needs
// attention (off track, at risk, stale, or without data). A guardrail on
// track stays out, so the list stays short.
func Steering(in StatusInput) []SteeringMetric {
	byInd := in.indicatorStatus()
	var out []SteeringMetric
	add := func(p model.Pathway, role, id string) {
		ind, ok := in.Doc.Indicator(id)
		if !ok {
			return // validation reports the unknown reference
		}
		m := SteeringMetric{PathwayID: p.ID, Pathway: p.Title, Role: role, IndicatorID: id, Indicator: ind.Name, Progress: "—"}
		if s, ok := byInd[id]; ok {
			m.State = string(s.State)
			m.Progress = progressText(s)
		}
		out = append(out, m)
	}
	for _, p := range in.Doc.Pathways {
		f := p.Focus
		if p.Rejected() || f == nil {
			continue
		}
		om, hasOM := f.OneMetricAt(in.Status.AsOf)
		switch {
		case hasOM && f.NorthStarID == om.IndicatorID:
			add(p, RoleNorthStarAndOne, f.NorthStarID)
		default:
			if f.NorthStarID != "" {
				add(p, RoleNorthStar, f.NorthStarID)
			}
			if hasOM {
				add(p, RoleOneMetric, om.IndicatorID)
			}
		}
		for _, id := range f.GuardrailIDs {
			if id == f.NorthStarID || (hasOM && id == om.IndicatorID) {
				continue // a lint rule reports it; it is listed once above
			}
			if s, ok := byInd[id]; ok && needsAttention(s.State) {
				add(p, RoleGuardrail, id)
			}
		}
	}
	return out
}

// needsAttention reports whether a guardrail in this state should be shown.
func needsAttention(s model.State) bool {
	return slices.Contains([]model.State{model.StateOffTrack, model.StateAtRisk, model.StateStale, model.StateNoData}, s)
}

func progressText(s model.IndicatorStatus) string {
	switch {
	case s.Progress != nil && s.Expected != nil:
		return fmt.Sprintf("%s (%s expected)", pct(s.Progress), pct(s.Expected))
	case s.Progress != nil:
		return pct(s.Progress)
	}
	return "—"
}

// steering renders the executive summary's section on what each pathway is
// steered by, or nothing when no pathway names a North Star or a one metric
// that matters.
func steering(in StatusInput) string {
	rows := Steering(in)
	if len(rows) == 0 {
		return ""
	}
	s := "## What each program steers by\n\n" +
		"| Program | Role | Indicator | Status | Progress |\n|---------|------|-----------|--------|----------|\n"
	for _, r := range rows {
		s += fmt.Sprintf("| %s | %s | %s | %s | %s |\n", MDCell(r.Pathway), r.Role, MDCell(r.Indicator), statusLabel(r.State), r.Progress)
	}
	return s + "\nA guardrail appears only when it needs attention.\n\n"
}
