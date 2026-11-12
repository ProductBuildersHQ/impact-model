package prismx

import (
	"errors"
	"fmt"
	"slices"

	"github.com/grokify/prism-roadmap/goals/bsc"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
)

// ErrNotBSCReady is returned, wrapped, when an impact model lacks what a
// Balanced Scorecard requires: an outcome placed in a scorecard perspective.
var ErrNotBSCReady = errors.New("the impact model cannot be drawn as a Balanced Scorecard")

// ToBSC exports the impact model as a Balanced Scorecard.
//
//   - The mission and vision are context.mission and context.vision.
//   - The perspectives are those the outcomes use, top first: mission,
//     financial, customer, internal process, learning and growth.
//   - Each outcome with a scorecard perspective is a strategic objective,
//     owned by its ownerActorId. It contributes to the nearest objectives its
//     outcome leads to, through causal links that are not feedback and
//     precondition groups, which makes the strategy map.
//   - Each pathway with an objective is a strategic theme, holding the
//     objectives among its outcomes.
//   - An objective's measures are its indicators, lead or lag.
//   - Each intervention is an initiative serving the objectives among its
//     targets, or else the nearest objectives they lead to.
//
// st may be nil; measures then have no current value, progress, or status.
// ToBSC refuses a model with no outcome in a scorecard perspective.
func ToBSC(doc *impact.Model, st *impact.Status) (*bsc.BSC, error) {
	if doc == nil {
		return nil, fmt.Errorf("prismx: nil document")
	}
	objectives := doc.ScorecardObjectives()
	if len(objectives) == 0 {
		return nil, fmt.Errorf("%w: add a scorecardPerspective to an outcome", ErrNotBSCReady)
	}

	statuses := map[string]impact.IndicatorStatus{}
	if st != nil {
		for _, s := range st.Indicators {
			statuses[s.IndicatorID] = s
		}
	}

	out := &bsc.BSC{
		Metadata: &bsc.Metadata{
			ID:        doc.ID,
			Name:      doc.Metadata.Name,
			Author:    doc.Metadata.Owner,
			Team:      doc.Metadata.Team,
			Version:   doc.Metadata.Version,
			Status:    bscStatus(doc.Metadata.Status),
			CreatedAt: doc.Metadata.CreatedAt,
			UpdatedAt: doc.Metadata.UpdatedAt,
		},
		Mission: doc.Context.Mission,
		Vision:  doc.Context.Vision,
	}
	if h := doc.Context.TimeHorizon; h != nil {
		out.Metadata.Period = h.Start + " to " + h.End
	}

	used := map[impact.ScorecardPerspective]bool{}
	for _, o := range objectives {
		used[o.ScorecardPerspective] = true
	}
	for _, p := range impact.AllScorecardPerspectives() {
		if used[p] {
			out.Perspectives = append(out.Perspectives, perspective(p))
		}
	}

	themeOf := map[string][]string{} // outcome ID to the pathways that hold it
	for _, p := range doc.Pathways {
		if p.Rejected() {
			continue
		}
		held := false
		for _, o := range objectives {
			if slices.Contains(p.OutcomeIDs, o.ID) {
				themeOf[o.ID] = append(themeOf[o.ID], p.ID)
				held = true
			}
		}
		if held {
			out.Themes = append(out.Themes, bsc.Theme{
				ID: p.ID, Name: p.Title, Description: p.Description, Owner: actorName(doc, p.OwnerID),
			})
		}
	}

	for _, o := range objectives {
		out.Objectives = append(out.Objectives, bsc.Objective{
			ID: o.ID, Name: o.Title, Description: o.Description, PerspectiveID: perspective(o.ScorecardPerspective).ID,
			ThemeIDs: themeOf[o.ID], Owner: actorName(doc, o.OwnerID), ContributesTo: doc.ScorecardContributesTo(o.ID),
		})
		for _, ind := range doc.Indicators {
			if ind.Measures != o.ID || ind.Rejected() {
				continue
			}
			n := ogsmNumberFrom(doc, ind, statuses[ind.ID], st != nil)
			out.Measures = append(out.Measures, bsc.Measure{
				ID: ind.ID, Name: ind.Name, Description: ind.Definition, ObjectiveID: o.ID, Type: string(ind.LeadLag),
				Owner: n.owner, Baseline: n.baseline, Target: n.target, Current: n.current, Unit: ind.Unit,
				Direction: string(ind.Direction), DueDate: n.dueDate, Frequency: frequencyLabel(ind.Frequency),
				Source: ind.DataSource, Progress: n.progress, Status: n.status,
			})
		}
	}

	for _, in := range doc.Interventions {
		if in.Rejected() {
			continue
		}
		served := doc.ScorecardObjectivesServed(in)
		if len(served) == 0 {
			continue
		}
		out.Initiatives = append(out.Initiatives, bsc.Initiative{
			ID: in.ID, Name: in.Title, Description: in.Description, ObjectiveIDs: served,
			Owner: actorName(doc, in.OwnerID), StartDate: in.StartDate, EndDate: in.EndDate, Budget: in.Resources,
			Status: bscInitiativeStatus(in.Status), ExternalLinks: initiativeLinks(in.InitiativeRefs),
		})
	}
	return out, nil
}

// perspective returns the scorecard perspective for a model perspective: one
// of prism-roadmap's four classic perspectives, or the mission perspective.
func perspective(p impact.ScorecardPerspective) bsc.Perspective {
	kind := map[impact.ScorecardPerspective]string{
		impact.ScorecardFinancial:       bsc.KindFinancial,
		impact.ScorecardCustomer:        bsc.KindCustomer,
		impact.ScorecardInternalProcess: bsc.KindInternalProcess,
		impact.ScorecardLearningGrowth:  bsc.KindLearningGrowth,
	}[p]
	for _, d := range bsc.DefaultPerspectives() {
		if d.Kind == kind {
			return d
		}
	}
	return bsc.Perspective{ID: "mission", Name: "Mission", Kind: bsc.KindMission,
		Question: "To fulfil our mission, what must change for the people we serve?"}
}

// bscInitiativeStatus maps an intervention's status; a paused intervention has
// no initiative status and is left blank.
func bscInitiativeStatus(s impact.InterventionStatus) string {
	switch s {
	case impact.InterventionPlanned:
		return bsc.InitiativePlanned
	case impact.InterventionInProgress:
		return bsc.InitiativeInProgress
	case impact.InterventionCompleted:
		return bsc.InitiativeCompleted
	case impact.InterventionCancelled:
		return bsc.InitiativeCancelled
	}
	return ""
}

func bscStatus(s string) string {
	switch s {
	case "active":
		return bsc.StatusActive
	case "approved":
		return bsc.StatusApproved
	case "in_review":
		return bsc.StatusInReview
	case "completed", "deprecated":
		return bsc.StatusCompleted
	case "archived":
		return bsc.StatusArchived
	default:
		return bsc.StatusDraft
	}
}
