package prismx

import (
	"errors"
	"fmt"
	"math"

	"github.com/grokify/prism-roadmap/goals/vsem"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// ErrNotVSEMReady is returned, wrapped, when an impact model lacks what a VSEM
// requires: a vision or mission, and at least one pathway.
var ErrNotVSEMReady = errors.New("the impact model cannot be drawn as a VSEM")

// ToVSEM exports the impact model as a VSEM (Vision, Strategy, Execution,
// Metrics).
//
//   - Vision is context.vision, and Mission context.mission.
//   - Each pathway that is not rejected is a strategy, with the resources it
//     commits.
//   - Execution is each strategy's interventions, with their owners, status,
//     and dates. An intervention of type program is a program; any other is an
//     initiative.
//   - A strategy's metrics are its focus: the North Star, the one metric that
//     matters, then the guardrails. The vision's metrics are the indicators of
//     impact-level outcomes not already listed under a strategy.
//
// VSEM IDs are unique, so an intervention or indicator that serves more than
// one pathway is listed under the first. st may be nil; the metrics then have
// no current value, progress, or status. ToVSEM refuses a model without a
// vision or mission, or without a pathway.
func ToVSEM(doc *model.Model, st *model.Status) (*vsem.VSEM, error) {
	if doc == nil {
		return nil, fmt.Errorf("prismx: nil document")
	}
	var missing []string
	if doc.Context.Vision == "" && doc.Context.Mission == "" {
		missing = append(missing, "context.vision or context.mission")
	}
	var pathways []model.Pathway
	for _, p := range doc.Pathways {
		if !p.Rejected() {
			pathways = append(pathways, p)
		}
	}
	if len(pathways) == 0 {
		missing = append(missing, "a pathway")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: add %s", ErrNotVSEMReady, joinAnd(missing))
	}

	statuses := map[string]model.IndicatorStatus{}
	if st != nil {
		for _, s := range st.Indicators {
			statuses[s.IndicatorID] = s
		}
	}

	out := &vsem.VSEM{
		Metadata: &vsem.Metadata{
			ID:        doc.ID,
			Name:      doc.Metadata.Name,
			Author:    doc.Metadata.Owner,
			Team:      doc.Metadata.Team,
			Version:   doc.Metadata.Version,
			Status:    vsemStatus(doc.Metadata.Status),
			ParentID:  doc.Parent,
			CreatedAt: doc.Metadata.CreatedAt,
			UpdatedAt: doc.Metadata.UpdatedAt,
		},
		Vision:  doc.Context.Vision,
		Mission: doc.Context.Mission,
	}
	if h := doc.Context.TimeHorizon; h != nil {
		out.Metadata.Period = h.Start + " to " + h.End
	}

	listed := map[string]bool{} // interventions and indicators already placed
	for _, p := range pathways {
		s := vsem.Strategy{ID: p.ID, Name: p.Title, Description: p.Description, Owner: actorName(doc, p.OwnerID)}
		for _, r := range p.Resources {
			s.Resources = append(s.Resources, vsem.Resource{Name: r.Name, Allocation: r.Allocation, Description: r.Description})
		}
		out.Strategies = append(out.Strategies, s)

		for _, id := range p.InterventionIDs {
			in, ok := doc.Intervention(id)
			if !ok || in.Rejected() || listed[in.ID] {
				continue
			}
			listed[in.ID] = true
			out.Execution = append(out.Execution, vsemExecution(doc, *in, p.ID))
		}
		if f := p.Focus; f != nil {
			for _, r := range focusMeasures(f) {
				ind, ok := doc.Indicator(r.id)
				if !ok || ind.Rejected() || listed[ind.ID] {
					continue
				}
				listed[ind.ID] = true
				out.Metrics = append(out.Metrics, vsemMetric(doc, *ind, p.ID, r.role, statuses[ind.ID], st != nil))
			}
		}
	}

	impactOutcomes := map[string]bool{}
	for _, o := range doc.ImpactOutcomes() {
		if !o.Rejected() {
			impactOutcomes[o.ID] = true
		}
	}
	for _, ind := range doc.Indicators {
		if !impactOutcomes[ind.Measures] || ind.Rejected() || listed[ind.ID] {
			continue
		}
		listed[ind.ID] = true
		out.Metrics = append(out.Metrics, vsemMetric(doc, ind, "", "", statuses[ind.ID], st != nil))
	}
	return out, nil
}

func vsemExecution(doc *model.Model, in model.Intervention, strategyID string) vsem.Execution {
	kind := vsem.KindInitiative
	if in.Type == model.InterventionProgram {
		kind = vsem.KindProgram
	}
	return vsem.Execution{
		ID:            in.ID,
		Name:          in.Title,
		Description:   in.Description,
		StrategyID:    strategyID,
		Kind:          kind,
		Owner:         actorName(doc, in.OwnerID),
		Status:        executionStatus(in.Status),
		StartDate:     in.StartDate,
		EndDate:       in.EndDate,
		ExternalLinks: initiativeLinks(in.InitiativeRefs),
	}
}

// vsemMetric draws a VSEM metric from an indicator. role names its place in
// the strategy's focus, such as "North Star"; it is empty for a vision metric.
func vsemMetric(doc *model.Model, ind model.Indicator, strategyID, role string, s model.IndicatorStatus, withStatus bool) vsem.Metric {
	m := vsem.Metric{
		ID:          ind.ID,
		Name:        ind.Name,
		Description: describe(role, ind.Definition),
		StrategyID:  strategyID,
		Owner:       actorName(doc, ind.OwnerID),
		Unit:        ind.Unit,
		Direction:   string(ind.Direction),
		Frequency:   frequencyLabel(ind.Frequency),
		Source:      ind.DataSource,
	}
	if ind.Baseline != nil {
		m.Baseline = measureValue(ind.Baseline.Value, ind.Baseline.Level)
	}
	if len(ind.Targets) > 0 {
		final := ind.Targets[0]
		for _, t := range ind.Targets[1:] {
			if t.Date > final.Date {
				final = t
			}
		}
		m.Target = measureValue(final.Value, final.Level)
		m.DueDate = final.Date
	}
	if !withStatus {
		return m
	}
	if s.Latest != nil {
		m.Current = measureValue(s.Latest.Value, s.Latest.Level)
	}
	if s.Progress != nil {
		m.Progress = math.Round(math.Max(0, math.Min(*s.Progress, 1))*100) / 100
	}
	m.Status = measureStatus(s.State)
	return m
}

// describe joins a focus role and an indicator's definition.
func describe(role, definition string) string {
	switch {
	case role == "":
		return definition
	case definition == "":
		return role
	default:
		return role + ": " + definition
	}
}

// frequencyLabel writes a frequency of one day, week, month, quarter, six
// months, or year in words, and any other ISO 8601 duration as given.
func frequencyLabel(iso string) string {
	switch iso {
	case "P1D":
		return "Daily"
	case "P1W", "P7D":
		return "Weekly"
	case "P1M":
		return "Monthly"
	case "P3M":
		return "Quarterly"
	case "P6M":
		return "Every six months"
	case "P1Y", "P12M":
		return "Yearly"
	}
	return iso
}

// executionStatus maps an intervention's status; a paused intervention has no
// VSEM status and is left blank.
func executionStatus(s model.InterventionStatus) string {
	switch s {
	case model.InterventionPlanned:
		return vsem.ExecutionPlanned
	case model.InterventionInProgress:
		return vsem.ExecutionInProgress
	case model.InterventionCompleted:
		return vsem.ExecutionCompleted
	case model.InterventionCancelled:
		return vsem.ExecutionCancelled
	}
	return ""
}

func vsemStatus(s string) string {
	switch s {
	case "active":
		return vsem.StatusActive
	case "approved":
		return vsem.StatusApproved
	case "in_review":
		return vsem.StatusInReview
	case "completed", "deprecated":
		return vsem.StatusCompleted
	case "archived":
		return vsem.StatusArchived
	default:
		return vsem.StatusDraft
	}
}
