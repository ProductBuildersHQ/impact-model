package prismx

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/ogsm"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
)

func TestOGSMExportValidates(t *testing.T) {
	doc, st := example(t)
	for name, status := range map[string]*impact.Status{"with status": st, "without": nil} {
		out, err := ToOGSM(doc, status)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if errs := ogsm.Errors(out.Validate(ogsm.DefaultValidationOptions())); len(errs) != 0 {
			t.Errorf("%s: the exported OGSM does not validate: %+v", name, errs)
		}
	}
	// The example names its objective, sets targets and owners, and gives
	// every pathway a focus and interventions, so it meets the strict options.
	out, err := ToOGSM(doc, st)
	if err != nil {
		t.Fatal(err)
	}
	if errs := ogsm.Errors(out.Validate(ogsm.StrictValidationOptions())); len(errs) != 0 {
		t.Errorf("the example does not meet the strict options: %+v", errs)
	}
}

func TestOGSMMapping(t *testing.T) {
	doc, st := example(t)
	out, err := ToOGSM(doc, st)
	if err != nil {
		t.Fatal(err)
	}

	// The objective is context.objective, due at the end of the horizon.
	if out.Objective.Statement != doc.Context.Objective.Statement {
		t.Errorf("objective = %q, want %q", out.Objective.Statement, doc.Context.Objective.Statement)
	}
	if h := doc.Context.TimeHorizon; out.Objective.DueDate != h.End || out.Metadata.Period != h.Start+" to "+h.End {
		t.Errorf("due %s, period %q, want the horizon %s to %s", out.Objective.DueDate, out.Metadata.Period, h.Start, h.End)
	}

	// Goals are the indicators of the objective's outcomes.
	objective := doc.Context.Objective.OutcomeIDs
	if len(out.Goals) == 0 {
		t.Fatal("no goals")
	}
	for _, g := range out.Goals {
		ind, ok := doc.Indicator(g.ID)
		if !ok || !slices.Contains(objective, ind.Measures) {
			t.Errorf("goal %s does not measure an objective outcome", g.ID)
		}
		if g.Target == "" || g.DueDate == "" {
			t.Errorf("goal %s has no final target: %+v", g.ID, g)
		}
	}

	// Strategies are the pathways that are not rejected, in file order, and
	// each serves a goal, since every pathway leads to the objective.
	var pathways, strategies []string
	for _, p := range doc.Pathways {
		if !p.Rejected() {
			pathways = append(pathways, p.ID)
		}
	}
	for _, s := range out.Strategies {
		strategies = append(strategies, s.ID)
		if len(s.GoalIDs) == 0 {
			t.Errorf("strategy %s serves no goal", s.ID)
		}
	}
	if !slices.Equal(strategies, pathways) {
		t.Fatalf("strategies = %v, want the pathways %v", strategies, pathways)
	}

	// A measure is a focus indicator of its strategy's pathway, never a goal,
	// and its frequency is in words.
	for _, m := range out.Measures {
		p, _ := doc.Pathway(m.StrategyID)
		var focus []string
		for _, r := range focusMeasures(p.Focus) {
			focus = append(focus, r.id)
		}
		if !slices.Contains(focus, m.ID) {
			t.Errorf("measure %s is not in the focus of %s", m.ID, m.StrategyID)
		}
		if out.Goal(m.ID) != nil {
			t.Errorf("measure %s is also a goal", m.ID)
		}
		if strings.HasPrefix(m.Frequency, "P") {
			t.Errorf("measure %s frequency %q is not in words", m.ID, m.Frequency)
		}
	}

	// An action is an intervention of its strategy's pathway, with its dates.
	for _, a := range out.Actions {
		p, _ := doc.Pathway(a.StrategyID)
		in, _ := doc.Intervention(a.ID)
		if !slices.Contains(p.InterventionIDs, a.ID) {
			t.Errorf("action %s is not an intervention of %s", a.ID, a.StrategyID)
		}
		if a.StartDate != in.StartDate || a.EndDate != in.EndDate {
			t.Errorf("action %s dates = %s..%s, want %s..%s", a.ID, a.StartDate, a.EndDate, in.StartDate, in.EndDate)
		}
	}
}

func TestOGSMLeavesOutRejectedPathways(t *testing.T) {
	doc, _ := example(t)
	rejected := doc.Pathways[0].ID
	doc.Pathways[0].Provenance = &impact.Provenance{Basis: impact.BasisRejected}
	out, err := ToOGSM(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out.Strategy(rejected) != nil {
		t.Errorf("rejected pathway %s became a strategy", rejected)
	}
	if errs := ogsm.Errors(out.Validate(nil)); len(errs) != 0 {
		t.Errorf("dangling references: %+v", errs)
	}
}

func TestOGSMWithoutStatusHasNoCurrentValues(t *testing.T) {
	doc, _ := example(t)
	out, err := ToOGSM(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range out.Goals {
		if g.Current != "" || g.Status != "" {
			t.Errorf("goal %s carries a current value or status without observations", g.ID)
		}
	}
	for _, m := range out.Measures {
		if m.Current != "" || m.Status != "" {
			t.Errorf("measure %s carries a current value or status without observations", m.ID)
		}
	}
}

func TestOGSMRefusesAModelWithoutAnObjectiveOrPathway(t *testing.T) {
	doc, _ := example(t)
	doc.Context.Objective = nil // the example has two impact outcomes
	for i := range doc.Pathways {
		doc.Pathways[i].Provenance = &impact.Provenance{Basis: impact.BasisRejected}
	}
	_, err := ToOGSM(doc, nil)
	if !errors.Is(err, ErrNotOGSMReady) {
		t.Fatalf("err = %v, want ErrNotOGSMReady", err)
	}
	for _, want := range []string{"context.objective", "a pathway"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, want it to name %s", err, want)
		}
	}
	if _, err := ToOGSM(nil, nil); err == nil {
		t.Error("a nil document was accepted")
	}
}

func TestOGSMStatuses(t *testing.T) {
	for in, want := range map[impact.InterventionStatus]string{
		impact.InterventionPlanned: ogsm.ActionPlanned, impact.InterventionInProgress: ogsm.ActionInProgress,
		impact.InterventionCompleted: ogsm.ActionCompleted, impact.InterventionCancelled: ogsm.ActionCancelled,
		impact.InterventionPaused: "",
	} {
		if got := actionStatus(in); got != want {
			t.Errorf("actionStatus(%q) = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{
		"draft": ogsm.StatusDraft, "": ogsm.StatusDraft, "active": ogsm.StatusActive, "deprecated": ogsm.StatusCompleted,
	} {
		if got := ogsmStatus(in); got != want {
			t.Errorf("ogsmStatus(%q) = %q, want %q", in, got, want)
		}
	}
}
