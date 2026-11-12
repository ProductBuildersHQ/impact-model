package measure

import (
	"fmt"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// Diff lists what changed between two status snapshots of the same outcome model:
// the state of indicators, links, and assumptions, and both states and the
// confidence of outcomes (their pathway state is the pathway's status).
// Changes come in a fixed order: indicators, outcomes, links, assumptions,
// each in the order of the later snapshot and then elements that only the
// earlier one had. An element that appears or disappears is a change in its
// presence.
func Diff(prev, cur *model.Status) *model.StatusDiff {
	d := &model.StatusDiff{ModelID: cur.ModelID, FromAsOf: prev.AsOf, ToAsOf: cur.AsOf, Changes: []model.StatusChange{}}
	add := func(id string, kind model.ElementKind, field, from, to string) {
		if from != to {
			d.Changes = append(d.Changes, model.StatusChange{ElementID: id, ElementKind: kind, Field: field, From: from, To: to})
		}
	}

	diffKind(d, model.ElementIndicator, ids(cur.Indicators, func(s model.IndicatorStatus) string { return s.IndicatorID }),
		ids(prev.Indicators, func(s model.IndicatorStatus) string { return s.IndicatorID }),
		func(id string) {
			a, b := findBy(prev.Indicators, id, func(s model.IndicatorStatus) string { return s.IndicatorID }), findBy(cur.Indicators, id, func(s model.IndicatorStatus) string { return s.IndicatorID })
			add(id, model.ElementIndicator, "state", string(a.State), string(b.State))
		})
	diffKind(d, model.ElementOutcome, ids(cur.Outcomes, func(s model.OutcomeStatus) string { return s.OutcomeID }),
		ids(prev.Outcomes, func(s model.OutcomeStatus) string { return s.OutcomeID }),
		func(id string) {
			a, b := findBy(prev.Outcomes, id, func(s model.OutcomeStatus) string { return s.OutcomeID }), findBy(cur.Outcomes, id, func(s model.OutcomeStatus) string { return s.OutcomeID })
			add(id, model.ElementOutcome, "state", string(a.State), string(b.State))
			add(id, model.ElementOutcome, "pathwayState", string(a.PathwayState), string(b.PathwayState))
			add(id, model.ElementOutcome, "confidence", confidence(a.Confidence), confidence(b.Confidence))
		})
	diffKind(d, model.ElementLink, ids(cur.Links, func(s model.LinkStatus) string { return s.LinkID }),
		ids(prev.Links, func(s model.LinkStatus) string { return s.LinkID }),
		func(id string) {
			a, b := findBy(prev.Links, id, func(s model.LinkStatus) string { return s.LinkID }), findBy(cur.Links, id, func(s model.LinkStatus) string { return s.LinkID })
			add(id, model.ElementLink, "state", string(a.State), string(b.State))
		})
	diffKind(d, model.ElementAssumption, ids(cur.Assumptions, func(s model.AssumptionState) string { return s.AssumptionID }),
		ids(prev.Assumptions, func(s model.AssumptionState) string { return s.AssumptionID }),
		func(id string) {
			a, b := findBy(prev.Assumptions, id, func(s model.AssumptionState) string { return s.AssumptionID }), findBy(cur.Assumptions, id, func(s model.AssumptionState) string { return s.AssumptionID })
			add(id, model.ElementAssumption, "status", string(a.Status), string(b.Status))
		})
	return d
}

// diffKind walks the later snapshot's elements, then those only the earlier
// one had, recording presence changes and calling compare for elements in
// both.
func diffKind(d *model.StatusDiff, kind model.ElementKind, cur, prev []string, compare func(id string)) {
	inPrev := map[string]bool{}
	for _, id := range prev {
		inPrev[id] = true
	}
	inCur := map[string]bool{}
	for _, id := range cur {
		inCur[id] = true
		if inPrev[id] {
			compare(id)
		} else {
			d.Changes = append(d.Changes, model.StatusChange{ElementID: id, ElementKind: kind, Field: "presence", To: "added"})
		}
	}
	for _, id := range prev {
		if !inCur[id] {
			d.Changes = append(d.Changes, model.StatusChange{ElementID: id, ElementKind: kind, Field: "presence", From: "present", To: "removed"})
		}
	}
}

func ids[T any](items []T, id func(T) string) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = id(it)
	}
	return out
}

func findBy[T any](items []T, id string, idOf func(T) string) T {
	for _, it := range items {
		if idOf(it) == id {
			return it
		}
	}
	var zero T
	return zero
}

// confidence renders a confidence to two places, the precision at which a
// change is worth reporting.
func confidence(c *float64) string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *c)
}
