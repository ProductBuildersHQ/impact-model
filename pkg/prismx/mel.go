package prismx

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/grokify/prism-roadmap/goals/mel"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// ErrNotMELReady is returned, wrapped, when an outcome model has no outcome
// that review has not rejected, so a MEL plan would have no results.
var ErrNotMELReady = errors.New("the outcome model cannot be drawn as a MEL plan")

// ToMEL exports the outcome model as a monitoring, evaluation, and learning
// plan.
//
//   - Results are the outcomes, at the same five levels. A result contributes
//     to the outcomes its causal links and precondition groups lead to;
//     feedback links, which are intentional loops, are left out.
//   - Activities are the interventions, producing the outcomes they target.
//   - Indicators are the indicators, each measuring its outcome (a result),
//     intervention (an activity), or assumption.
//   - Assumptions are the assumptions; the results that depend on one are the
//     outcomes it applies to, the targets of links it applies to, and the
//     targets of interventions it applies to. How it is monitored is its test.
//   - Evaluations, the learning agenda, and reviews are the model's
//     evaluations, learning questions, and review meetings. A causal link they
//     name is written as both of its ends, since a MEL plan has no links.
//   - Reports are the audiences: what each needs to know, what it receives (its
//     channels), how often, and the indicators it sees, headline first.
//
// Elements review rejected are left out, with every reference to them. st may
// be nil; the indicators then have no current value.
func ToMEL(doc *model.Model, st *model.Status) (*mel.Plan, error) {
	if doc == nil {
		return nil, fmt.Errorf("prismx: nil document")
	}
	x := melExport{doc: doc, live: map[string]bool{}}
	for _, o := range doc.Outcomes {
		if !o.Rejected() {
			x.live[o.ID] = true
		}
	}
	if len(x.live) == 0 {
		return nil, fmt.Errorf("%w: add an outcome", ErrNotMELReady)
	}
	for _, in := range doc.Interventions {
		if !in.Rejected() {
			x.live[in.ID] = true
		}
	}
	for _, a := range doc.Assumptions {
		if !a.Rejected() {
			x.live[a.ID] = true
		}
	}
	for _, ind := range doc.Indicators {
		if !ind.Rejected() && x.live[ind.Measures] {
			x.live[ind.ID] = true
		}
	}
	for _, q := range doc.LearningQuestions {
		if !q.Rejected() {
			x.live[q.ID] = true
		}
	}

	current := map[string]string{}
	if st != nil {
		for _, s := range st.Indicators {
			if s.Latest != nil {
				current[s.IndicatorID] = measureValue(s.Latest.Value, s.Latest.Level)
			}
		}
	}

	out := &mel.Plan{
		Metadata: &mel.Metadata{
			ID:        doc.ID,
			Name:      doc.Metadata.Name,
			Author:    doc.Metadata.Owner,
			Team:      doc.Metadata.Team,
			Version:   doc.Metadata.Version,
			Status:    melStatus(doc.Metadata.Status),
			ParentID:  doc.Parent,
			CreatedAt: doc.Metadata.CreatedAt,
			UpdatedAt: doc.Metadata.UpdatedAt,
		},
	}
	if h := doc.Context.TimeHorizon; h != nil {
		out.Metadata.Period = h.Start + " to " + h.End
	}
	if b := doc.Budget; b != nil && b.StaffMinutesPerWeek > 0 {
		out.Budget = strconv.FormatFloat(b.StaffMinutesPerWeek, 'f', -1, 64) + " staff minutes a week for measurement"
	}

	for _, o := range doc.Outcomes {
		if !x.live[o.ID] {
			continue
		}
		out.Results = append(out.Results, mel.Result{
			ID: o.ID, Statement: o.Title, Level: string(o.Level), Description: o.Description,
			ContributesTo: x.leadsTo(o.ID),
		})
	}
	for _, in := range doc.Interventions {
		if !x.live[in.ID] {
			continue
		}
		out.Activities = append(out.Activities, mel.Activity{
			ID: in.ID, Name: in.Title, Description: in.Description,
			ResultIDs: x.keep(in.Targets), Owner: actorName(doc, in.OwnerID),
			StartDate: in.StartDate, EndDate: in.EndDate, Status: activityStatus(in.Status),
		})
	}
	for _, ind := range doc.Indicators {
		if x.live[ind.ID] {
			out.Indicators = append(out.Indicators, x.indicator(ind, current[ind.ID]))
		}
	}
	for _, a := range doc.Assumptions {
		if !x.live[a.ID] {
			continue
		}
		watch := slices.Clone(x.keep(a.IndicatorIDs))
		for _, ind := range doc.Indicators {
			if ind.Measures == a.ID && x.live[ind.ID] && !slices.Contains(watch, ind.ID) {
				watch = append(watch, ind.ID)
			}
		}
		out.Assumptions = append(out.Assumptions, mel.Assumption{
			ID: a.ID, Statement: a.Statement, ResultIDs: x.dependents(a.AppliesTo),
			Monitoring: a.Test, IndicatorIDs: watch,
		})
	}
	for _, e := range doc.Evaluations {
		if e.Rejected() {
			continue
		}
		out.Evaluations = append(out.Evaluations, mel.Evaluation{
			ID: e.ID, Title: e.Title, Type: string(e.Type), Questions: slices.Clone(e.Questions),
			Method: e.Method, Evaluator: string(e.Evaluator), Owner: actorName(doc, e.OwnerID),
			StartDate: e.StartDate, EndDate: e.EndDate, Status: evaluationStatus(e.Status),
			Budget: e.Budget, RelatedIDs: x.related(e.RelatedIDs),
		})
	}
	for _, q := range doc.LearningQuestions {
		if !x.live[q.ID] {
			continue
		}
		out.Learning = append(out.Learning, mel.LearningQuestion{
			ID: q.ID, Question: q.Question, Rationale: q.Rationale, Methods: slices.Clone(q.Methods),
			Owner: actorName(doc, q.OwnerID), DueDate: q.DueDate, Status: questionStatus(q.Status),
			Answer: q.Answer, RelatedIDs: x.related(q.RelatedIDs),
		})
	}
	for _, m := range doc.ReviewMeetings {
		out.Reviews = append(out.Reviews, mel.Review{
			ID: m.ID, Name: m.Name, Purpose: m.Purpose, Frequency: frequencyLabel(m.Frequency),
			Participants: slices.Clone(m.Participants), Owner: actorName(doc, m.OwnerID),
			NextDate: m.NextDate, QuestionIDs: x.keep(m.QuestionIDs),
		})
	}
	for _, a := range doc.Audiences {
		out.Reports = append(out.Reports, mel.Report{
			ID: a.ID, Audience: a.Name, Needs: a.Needs, Product: strings.Join(a.Channels, ", "),
			Frequency: frequencyLabel(a.Frequency), Owner: actorName(doc, a.OwnerID),
			IndicatorIDs: x.reported(a),
		})
	}
	return out, nil
}

// melExport carries what the MEL export has decided to keep.
type melExport struct {
	doc  *model.Model
	live map[string]bool // IDs of exported elements
}

// keep returns the IDs that are exported, in order, or nil.
func (x melExport) keep(ids []string) []string {
	var out []string
	for _, id := range ids {
		if x.live[id] && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out
}

// leadsTo returns the outcomes an outcome contributes to: the targets of its
// live causal links, then the targets of the precondition groups it belongs to.
// Feedback links are left out: they are intentional loops, and a results
// framework is a hierarchy.
func (x melExport) leadsTo(id string) []string {
	var to []string
	for _, l := range x.doc.Links {
		if l.From == id && !l.Rejected() && !l.Feedback {
			if o, ok := x.doc.Outcome(l.To); ok && x.live[o.ID] {
				to = append(to, o.ID)
			}
		}
	}
	for _, g := range x.doc.Groups {
		if slices.Contains(g.Members, id) && !g.Rejected() {
			to = append(to, g.Target)
		}
	}
	return x.keep(to)
}

// dependents returns the results that depend on what an assumption applies
// to: an outcome itself, a link's target, or an intervention's targets.
func (x melExport) dependents(appliesTo []string) []string {
	var ids []string
	for _, id := range appliesTo {
		switch ref, _ := x.doc.Element(id); ref.Kind {
		case model.ElementOutcome:
			ids = append(ids, id)
		case model.ElementLink:
			if l, ok := x.doc.Link(id); ok && !l.Rejected() {
				ids = append(ids, l.To)
			}
		case model.ElementIntervention:
			if in, ok := x.doc.Intervention(id); ok {
				ids = append(ids, in.Targets...)
			}
		}
	}
	var out []string
	for _, id := range x.keep(ids) {
		if _, ok := x.doc.Outcome(id); ok {
			out = append(out, id)
		}
	}
	return out
}

// related returns the MEL IDs for what an evaluation or learning question
// names. A causal link becomes both of its ends.
func (x melExport) related(ids []string) []string {
	var out []string
	for _, id := range ids {
		if l, ok := x.doc.Link(id); ok {
			if !l.Rejected() {
				out = append(out, l.From, l.To)
			}
			continue
		}
		out = append(out, id)
	}
	return x.keep(out)
}

// reported returns the indicators an audience sees: its headline list in
// order, then the other indicators reported to it, in file order.
func (x melExport) reported(a model.Audience) []string {
	var ids []string
	for _, h := range a.Headline {
		ids = append(ids, h.IndicatorID)
	}
	for _, ind := range x.doc.Indicators {
		if slices.Contains(ind.AudienceIDs, a.ID) {
			ids = append(ids, ind.ID)
		}
	}
	return x.keep(ids)
}

func (x melExport) indicator(ind model.Indicator, current string) mel.Indicator {
	m := mel.Indicator{
		ID: ind.ID, Name: ind.Name, Measures: ind.Measures, Definition: ind.Definition,
		Type: indicatorType(ind.MeasurementType), Unit: ind.Unit, Direction: string(ind.Direction),
		Disaggregation: slices.Clone(ind.Disaggregation), DataSource: ind.DataSource,
		CollectionMethod: ind.Method, Frequency: frequencyLabel(ind.Frequency),
		Owner: actorName(x.doc, ind.OwnerID), Current: current,
	}
	if c := ind.Collection; c != nil {
		if m.DataSource == "" {
			m.DataSource = c.System
		}
		if m.CollectionMethod == "" {
			m.CollectionMethod = strings.ReplaceAll(string(c.Mode), "_", " ")
		}
	}
	if b := ind.Baseline; b != nil {
		m.Baseline = &mel.Baseline{Value: measureValue(b.Value, b.Level), Date: b.Date, Source: b.Source}
	}
	for _, t := range ind.Targets {
		m.Targets = append(m.Targets, mel.Target{Value: measureValue(t.Value, t.Level), Date: t.Date, Label: t.Label})
	}
	if qa := ind.QualityAssurance; qa != nil || len(ind.Caveats) > 0 {
		m.DataQuality = &mel.DataQuality{Limitations: slices.Clone(ind.Caveats)}
		if qa != nil {
			m.DataQuality.Checks = slices.Clone(qa.Checks)
			m.DataQuality.LastAssessed = qa.LastAssessed
		}
	}
	return m
}

// indicatorType maps a measurement type; a theme count is a number, and a
// hybrid indicator has no single MEL type, so it is left blank.
func indicatorType(t model.MeasurementType) string {
	switch t {
	case model.MeasurementQuantitative, model.MeasurementThemeCount:
		return mel.TypeQuantitative
	case model.MeasurementQualitative:
		return mel.TypeQualitative
	}
	return ""
}

func activityStatus(s model.InterventionStatus) string {
	switch s {
	case model.InterventionPlanned:
		return mel.ActivityPlanned
	case model.InterventionInProgress:
		return mel.ActivityInProgress
	case model.InterventionPaused:
		return mel.ActivityPaused
	case model.InterventionCompleted:
		return mel.ActivityCompleted
	case model.InterventionCancelled:
		return mel.ActivityCancelled
	}
	return ""
}

func evaluationStatus(s model.EvaluationStatus) string {
	switch s {
	case model.EvaluationPlanned:
		return mel.EvaluationPlanned
	case model.EvaluationInProgress:
		return mel.EvaluationInProgress
	case model.EvaluationCompleted:
		return mel.EvaluationCompleted
	case model.EvaluationCancelled:
		return mel.EvaluationCancelled
	}
	return ""
}

func questionStatus(s model.QuestionStatus) string {
	switch s {
	case model.QuestionOpen:
		return mel.QuestionOpen
	case model.QuestionAnswered:
		return mel.QuestionAnswered
	case model.QuestionRetired:
		return mel.QuestionRetired
	}
	return ""
}

func melStatus(s string) string {
	switch s {
	case "active":
		return mel.StatusActive
	case "approved":
		return mel.StatusApproved
	case "in_review":
		return mel.StatusInReview
	case "completed", "deprecated":
		return mel.StatusCompleted
	case "archived":
		return mel.StatusArchived
	default:
		return mel.StatusDraft
	}
}
