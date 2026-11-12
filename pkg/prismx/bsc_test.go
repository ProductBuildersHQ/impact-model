package prismx

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/bsc"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
)

func TestBSCExportValidates(t *testing.T) {
	doc, st := example(t)
	for name, status := range map[string]*impact.Status{"with status": st, "without": nil} {
		out, err := ToBSC(doc, status)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if errs := bsc.Errors(out.Validate(bsc.DefaultValidationOptions())); len(errs) != 0 {
			t.Errorf("%s: the exported scorecard does not validate: %+v", name, errs)
		}
	}
	// The example places owned and measured outcomes so that its strategy map
	// runs upward and mixes lead and lag measures: it meets the strict
	// options.
	out, err := ToBSC(doc, st)
	if err != nil {
		t.Fatal(err)
	}
	if errs := bsc.Errors(out.Validate(bsc.StrictValidationOptions())); len(errs) != 0 {
		t.Errorf("the example does not meet the strict options: %+v", errs)
	}
}

func TestBSCMapping(t *testing.T) {
	doc, st := example(t)
	out, err := ToBSC(doc, st)
	if err != nil {
		t.Fatal(err)
	}
	if out.Mission != doc.Context.Mission || out.Vision != doc.Context.Vision {
		t.Errorf("mission %q and vision %q differ from the context", out.Mission, out.Vision)
	}

	// The perspectives are those in use, top first, mission above financial.
	var kinds []string
	for _, p := range out.Perspectives {
		kinds = append(kinds, p.Kind)
	}
	if want := []string{bsc.KindMission, bsc.KindFinancial, bsc.KindCustomer, bsc.KindInternalProcess}; !slices.Equal(kinds, want) {
		t.Errorf("perspective kinds = %v, want %v", kinds, want)
	}

	// Objectives are the placed outcomes, in file order, owned, with the
	// strategy map from the causal graph.
	var placed, objectives []string
	for _, o := range doc.ScorecardObjectives() {
		placed = append(placed, o.ID)
	}
	for _, o := range out.Objectives {
		objectives = append(objectives, o.ID)
		if o.Owner == "" {
			t.Errorf("objective %s has no owner", o.ID)
		}
		if !slices.Equal(o.ContributesTo, doc.ScorecardContributesTo(o.ID)) {
			t.Errorf("objective %s contributes to %v, want %v", o.ID, o.ContributesTo, doc.ScorecardContributesTo(o.ID))
		}
	}
	if !slices.Equal(objectives, placed) {
		t.Fatalf("objectives = %v, want the placed outcomes %v", objectives, placed)
	}
	// Learning cohorts lead, through skills that are not objectives, to the
	// leaders' practice.
	if c := out.Objective("oc-out-cohorts"); c == nil || !slices.Contains(c.ContributesTo, "oc-im-leaders-practice") {
		t.Errorf("cohorts objective = %+v, want it to contribute to the leaders' practice", c)
	}

	// Themes are the pathways that hold an objective; a rejected pathway is
	// not one.
	for _, th := range out.Themes {
		p, _ := doc.Pathway(th.ID)
		if p == nil || p.Rejected() || len(out.ObjectivesInTheme(th.ID)) == 0 {
			t.Errorf("theme %s is not a pathway holding an objective", th.ID)
		}
	}
	if out.Theme("pw-certification") != nil {
		t.Error("the rejected certification pathway became a theme")
	}

	// Measures are the objectives' indicators, typed lead or lag.
	for _, m := range out.Measures {
		ind, _ := doc.Indicator(m.ID)
		if ind.Measures != m.ObjectiveID || m.Type != string(ind.LeadLag) {
			t.Errorf("measure %s = %+v", m.ID, m)
		}
	}

	// Initiatives serve the objectives their targets are or lead to.
	for _, in := range out.Initiatives {
		src, _ := doc.Intervention(in.ID)
		if !slices.Equal(in.ObjectiveIDs, doc.ScorecardObjectivesServed(*src)) || in.Budget != src.Resources {
			t.Errorf("initiative %s = %+v", in.ID, in)
		}
	}
	if in := findInitiative(out, "int-liaison-channel"); in == nil || !slices.Equal(in.ObjectiveIDs, []string{"oc-lt-collective-voice"}) {
		t.Errorf("the liaison channel should serve the collective voice through officials' engagement: %+v", in)
	}
}

func TestBSCWithoutStatusHasNoCurrentValues(t *testing.T) {
	doc, _ := example(t)
	out, err := ToBSC(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range out.Measures {
		if m.Current != "" || m.Status != "" {
			t.Errorf("measure %s carries a current value or status without observations", m.ID)
		}
	}
}

func TestBSCRefusesAModelWithoutPerspectives(t *testing.T) {
	doc, _ := example(t)
	for i := range doc.Outcomes {
		doc.Outcomes[i].ScorecardPerspective = ""
	}
	if _, err := ToBSC(doc, nil); !errors.Is(err, ErrNotBSCReady) || !strings.Contains(err.Error(), "scorecardPerspective") {
		t.Errorf("err = %v, want ErrNotBSCReady naming scorecardPerspective", err)
	}
	if _, err := ToBSC(nil, nil); err == nil {
		t.Error("a nil document was accepted")
	}
}

func TestBSCStatuses(t *testing.T) {
	for in, want := range map[impact.InterventionStatus]string{
		impact.InterventionPlanned: bsc.InitiativePlanned, impact.InterventionInProgress: bsc.InitiativeInProgress,
		impact.InterventionCompleted: bsc.InitiativeCompleted, impact.InterventionCancelled: bsc.InitiativeCancelled,
		impact.InterventionPaused: "",
	} {
		if got := bscInitiativeStatus(in); got != want {
			t.Errorf("bscInitiativeStatus(%q) = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{"draft": bsc.StatusDraft, "active": bsc.StatusActive, "archived": bsc.StatusArchived} {
		if got := bscStatus(in); got != want {
			t.Errorf("bscStatus(%q) = %q, want %q", in, got, want)
		}
	}
	if p := perspective(impact.ScorecardLearningGrowth); p.Kind != bsc.KindLearningGrowth || p.Question == "" {
		t.Errorf("learning and growth perspective = %+v", p)
	}
}

func findInitiative(b *bsc.BSC, id string) *bsc.Initiative {
	for i := range b.Initiatives {
		if b.Initiatives[i].ID == id {
			return &b.Initiatives[i]
		}
	}
	return nil
}
