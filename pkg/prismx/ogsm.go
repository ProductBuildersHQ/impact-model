package prismx

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/grokify/prism-roadmap/goals/ogsm"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
)

// ErrNotOGSMReady is returned, wrapped, when an impact model lacks what an
// OGSM requires: one clear objective and at least one pathway.
var ErrNotOGSMReady = errors.New("the impact model cannot be drawn as an OGSM")

// ToOGSM exports the impact model as an OGSM (Objective, Goals, Strategies,
// Measures).
//
//   - The objective is context.objective, or the model's one impact outcome,
//     due at the end of the time horizon.
//   - Goals are the indicators of the outcomes the objective stands for.
//   - Each pathway that is not rejected is a strategy. It serves the goals
//     whose outcomes it contains or leads to, through causal links and
//     precondition groups.
//   - A strategy's measures are its focus: the North Star, the one metric that
//     matters, then the guardrails, leaving out any that is already a goal.
//   - A strategy's actions are its interventions.
//
// OGSM IDs are unique, so an intervention or indicator that serves more than
// one pathway is listed under the first. st may be nil; goals and measures
// then have no current value, progress, or status. ToOGSM refuses a model
// whose objective is ambiguous, or that has no pathway.
func ToOGSM(doc *impact.Model, st *impact.Status) (*ogsm.OGSM, error) {
	if doc == nil {
		return nil, fmt.Errorf("prismx: nil document")
	}
	var missing []string
	statement, ok := doc.ObjectiveStatement()
	if !ok {
		missing = append(missing, "context.objective")
	}
	var pathways []impact.Pathway
	for _, p := range doc.Pathways {
		if !p.Rejected() {
			pathways = append(pathways, p)
		}
	}
	if len(pathways) == 0 {
		missing = append(missing, "a pathway")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: add %s", ErrNotOGSMReady, joinAnd(missing))
	}

	statuses := map[string]impact.IndicatorStatus{}
	if st != nil {
		for _, s := range st.Indicators {
			statuses[s.IndicatorID] = s
		}
	}
	number := func(ind impact.Indicator) ogsmNumber {
		return ogsmNumberFrom(doc, ind, statuses[ind.ID], st != nil)
	}

	out := &ogsm.OGSM{
		Metadata: &ogsm.Metadata{
			ID:        doc.ID,
			Name:      doc.Metadata.Name,
			Author:    doc.Metadata.Owner,
			Team:      doc.Metadata.Team,
			Version:   doc.Metadata.Version,
			Status:    ogsmStatus(doc.Metadata.Status),
			CreatedAt: doc.Metadata.CreatedAt,
			UpdatedAt: doc.Metadata.UpdatedAt,
		},
		Objective: ogsm.Objective{Statement: statement},
	}
	if h := doc.Context.TimeHorizon; h != nil {
		out.Metadata.Period = h.Start + " to " + h.End
		out.Objective.DueDate = h.End
	}

	listed := map[string]bool{}     // indicators and interventions already placed
	goalOf := map[string][]string{} // outcome ID to the goals that measure it
	for _, o := range doc.ObjectiveOutcomes() {
		for _, ind := range doc.Indicators {
			if ind.Measures != o.ID || ind.Rejected() || listed[ind.ID] {
				continue
			}
			listed[ind.ID] = true
			n := number(ind)
			out.Goals = append(out.Goals, ogsm.Goal{
				ID: ind.ID, Name: ind.Name, Description: ind.Definition, Owner: n.owner,
				Baseline: n.baseline, Target: n.target, Current: n.current, Unit: ind.Unit,
				Direction: string(ind.Direction), DueDate: n.dueDate, Progress: n.progress, Status: n.status,
			})
			goalOf[o.ID] = append(goalOf[o.ID], ind.ID)
		}
	}

	for _, p := range pathways {
		reach := reachable(doc, p.OutcomeIDs)
		var serves []string
		for _, o := range doc.ObjectiveOutcomes() {
			if reach[o.ID] {
				serves = append(serves, goalOf[o.ID]...)
			}
		}
		out.Strategies = append(out.Strategies, ogsm.Strategy{
			ID: p.ID, Name: p.Title, Description: p.Description, Owner: actorName(doc, p.OwnerID), GoalIDs: serves,
		})
		if f := p.Focus; f != nil {
			for _, r := range focusMeasures(f) {
				ind, ok := doc.Indicator(r.id)
				if !ok || ind.Rejected() || listed[ind.ID] {
					continue
				}
				listed[ind.ID] = true
				n := number(*ind)
				out.Measures = append(out.Measures, ogsm.Measure{
					ID: ind.ID, Name: ind.Name, Description: describe(r.role, ind.Definition), StrategyID: p.ID,
					Owner: n.owner, Baseline: n.baseline, Target: n.target, Current: n.current, Unit: ind.Unit,
					Direction: string(ind.Direction), DueDate: n.dueDate, Frequency: frequencyLabel(ind.Frequency),
					Source: ind.DataSource, Progress: n.progress, Status: n.status,
				})
			}
		}
		for _, id := range p.InterventionIDs {
			in, ok := doc.Intervention(id)
			if !ok || in.Rejected() || listed[in.ID] {
				continue
			}
			listed[in.ID] = true
			out.Actions = append(out.Actions, ogsm.Action{
				ID: in.ID, Name: in.Title, Description: in.Description, StrategyID: p.ID,
				Owner: actorName(doc, in.OwnerID), StartDate: in.StartDate, EndDate: in.EndDate,
				Status: actionStatus(in.Status), ExternalLinks: initiativeLinks(in.InitiativeRefs),
			})
		}
	}
	return out, nil
}

// ogsmNumber holds the fields goals and measures share.
type ogsmNumber struct {
	owner, baseline, target, current, dueDate, status string
	progress                                          float64
}

func ogsmNumberFrom(doc *impact.Model, ind impact.Indicator, s impact.IndicatorStatus, withStatus bool) ogsmNumber {
	n := ogsmNumber{owner: actorName(doc, ind.OwnerID)}
	if ind.Baseline != nil {
		n.baseline = measureValue(ind.Baseline.Value, ind.Baseline.Level)
	}
	if len(ind.Targets) > 0 {
		final := ind.Targets[0]
		for _, t := range ind.Targets[1:] {
			if t.Date > final.Date {
				final = t
			}
		}
		n.target, n.dueDate = measureValue(final.Value, final.Level), final.Date
	}
	if !withStatus {
		return n
	}
	if s.Latest != nil {
		n.current = measureValue(s.Latest.Value, s.Latest.Level)
	}
	if s.Progress != nil {
		n.progress = math.Round(math.Max(0, math.Min(*s.Progress, 1))*100) / 100
	}
	n.status = measureStatus(s.State)
	return n
}

// reachable returns the outcomes the given outcomes are or lead to, through
// causal links that are not rejected or feedback, and precondition groups.
func reachable(doc *impact.Model, from []string) map[string]bool {
	seen := map[string]bool{}
	queue := slices.Clone(from)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		seen[id] = true
		for _, l := range doc.Links {
			if l.From == id && !l.Rejected() && !l.Feedback {
				queue = append(queue, l.To)
			}
		}
		for _, g := range doc.Groups {
			if slices.Contains(g.Members, id) && !g.Rejected() {
				queue = append(queue, g.Target)
			}
		}
	}
	return seen
}

// actionStatus maps an intervention's status; a paused intervention has no
// OGSM action status and is left blank.
func actionStatus(s impact.InterventionStatus) string {
	switch s {
	case impact.InterventionPlanned:
		return ogsm.ActionPlanned
	case impact.InterventionInProgress:
		return ogsm.ActionInProgress
	case impact.InterventionCompleted:
		return ogsm.ActionCompleted
	case impact.InterventionCancelled:
		return ogsm.ActionCancelled
	}
	return ""
}

func ogsmStatus(s string) string {
	switch s {
	case "active":
		return ogsm.StatusActive
	case "approved":
		return ogsm.StatusApproved
	case "in_review":
		return ogsm.StatusInReview
	case "completed", "deprecated":
		return ogsm.StatusCompleted
	case "archived":
		return ogsm.StatusArchived
	default:
		return ogsm.StatusDraft
	}
}
