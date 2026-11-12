package measure

import (
	"math"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/graph"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func set(assessments ...model.AssumptionAssessment) []*model.ObservationSet {
	return []*model.ObservationSet{{Assessments: assessments}}
}

func TestEffectiveAssumptions(t *testing.T) {
	doc := model.NewModel("t", "p")
	doc.Assumptions = []model.Assumption{
		{ID: "a1"}, {ID: "a2", Status: model.AssumptionHolding}, {ID: "a3", Status: model.AssumptionHolding},
		{ID: "a4", Provenance: &model.Provenance{Basis: model.BasisRejected}},
	}
	got := EffectiveAssumptions(doc, set(
		model.AssumptionAssessment{AssumptionID: "a2", Date: "2026-03-01", Status: model.AssumptionAtRisk},
		model.AssumptionAssessment{AssumptionID: "a2", Date: "2026-06-01", Status: model.AssumptionInvalidated},
		model.AssumptionAssessment{AssumptionID: "a3", Date: "2027-01-01", Status: model.AssumptionInvalidated},
	), day("2026-07-01"))
	want := []model.AssumptionState{
		{AssumptionID: "a1", Status: model.AssumptionUntested, Source: model.SourceDocument},
		{AssumptionID: "a2", Status: model.AssumptionInvalidated, Source: model.SourceAssessment, AssessedAt: "2026-06-01"},
		{AssumptionID: "a3", Status: model.AssumptionHolding, Source: model.SourceDocument},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func near6(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func twoOutcomes() *model.Model {
	d := model.NewModel("t", "p")
	d.Outcomes = []model.Outcome{{ID: "low", Level: model.LevelShortTerm}, {ID: "top", Level: model.LevelImpact}}
	d.Evidence = []model.Evidence{
		{ID: "ev-strong", Strength: model.StrengthStrong, Direction: model.EvidenceSupports},
		{ID: "ev-weak", Strength: model.StrengthWeak},
		{ID: "ev-against", Strength: model.StrengthStrong, Direction: model.EvidenceContradicts},
	}
	return d
}

func TestLinkConfidenceIsRuleBased(t *testing.T) {
	w := model.DefaultStatusOptions().Weights // 0.4 / 0.4 / 0.2
	tests := []struct {
		name    string
		link    model.CausalLink
		assumed map[string]model.AssumptionStatus
		want    float64
		state   model.LinkState
	}{
		{"nothing at all: no assumptions score 0.5", model.CausalLink{}, nil, 0.4 * 0.5, model.LinkStateOK},
		{"rationale only", model.CausalLink{Rationale: "r"}, nil, 0.4*0 + 0.4*0.5 + 0.2, model.LinkStateOK},
		{"strong evidence, validated, rationale", model.CausalLink{Rationale: "r", EvidenceIDs: []string{"ev-strong"}, AssumptionIDs: []string{"a"}}, map[string]model.AssumptionStatus{"a": model.AssumptionValidated}, 1, model.LinkStateOK},
		{"best supporting evidence counts", model.CausalLink{EvidenceIDs: []string{"ev-weak", "ev-strong"}}, nil, 0.4*1 + 0.4*0.5, model.LinkStateOK},
		{"contradicting evidence takes half off", model.CausalLink{EvidenceIDs: []string{"ev-strong", "ev-against"}}, nil, 0.4*0.5 + 0.4*0.5, model.LinkStateOK},
		{"contradiction floors at zero", model.CausalLink{EvidenceIDs: []string{"ev-against"}}, nil, 0.4 * 0.5, model.LinkStateOK},
		{"lowest assumption wins", model.CausalLink{AssumptionIDs: []string{"a", "b"}}, map[string]model.AssumptionStatus{"a": model.AssumptionValidated, "b": model.AssumptionHolding}, 0.4 * 0.8, model.LinkStateOK},
		{"at risk assumption", model.CausalLink{AssumptionIDs: []string{"a"}}, map[string]model.AssumptionStatus{"a": model.AssumptionAtRisk}, 0.4 * 0.25, model.LinkStateAtRisk},
		{"invalidated assumption breaks the link", model.CausalLink{AssumptionIDs: []string{"a", "b"}}, map[string]model.AssumptionStatus{"a": model.AssumptionInvalidated, "b": model.AssumptionValidated}, 0, model.LinkStateBroken},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := twoOutcomes()
			tt.link.ID, tt.link.From, tt.link.To = "ln", "low", "top"
			d.Links = []model.CausalLink{tt.link}
			var states []model.AssumptionState
			for id, s := range tt.assumed {
				states = append(states, model.AssumptionState{AssumptionID: id, Status: s})
			}
			got := LinkStatuses(d, states, w)
			if len(got) != 1 || got[0].State != tt.state || !near6(got[0].Confidence, tt.want) {
				t.Errorf("got %+v, want %s %.4f", got, tt.state, tt.want)
			}
		})
	}
}

func TestAssumptionsThatNameALinkInAppliesToCount(t *testing.T) {
	d := twoOutcomes()
	d.Links = []model.CausalLink{{ID: "ln", From: "low", To: "top"}}
	d.Assumptions = []model.Assumption{{ID: "a", AppliesTo: []string{"ln"}}}
	got := LinkStatuses(d, []model.AssumptionState{{AssumptionID: "a", Status: model.AssumptionInvalidated}}, model.DefaultStatusOptions().Weights)
	if got[0].State != model.LinkStateBroken {
		t.Errorf("state = %s", got[0].State)
	}
}

func TestWeightsAreOptions(t *testing.T) {
	d := twoOutcomes()
	d.Links = []model.CausalLink{{ID: "ln", From: "low", To: "top", Rationale: "r"}}
	got := LinkStatuses(d, nil, model.ConfidenceWeights{Evidence: 0, Assumptions: 0, Rationale: 1})
	if !near6(got[0].Confidence, 1) {
		t.Errorf("confidence = %v", got[0].Confidence)
	}
}

// chain3: a -> b -> top, with an assumption on the first link.
func chain3() (*model.Model, *graph.Graph) {
	d := model.NewModel("t", "p")
	d.Outcomes = []model.Outcome{{ID: "a", Title: "a", Level: model.LevelShortTerm}, {ID: "b", Title: "b", Level: model.LevelMediumTerm}, {ID: "top", Title: "top", Level: model.LevelImpact}, {ID: "side", Title: "side", Level: model.LevelShortTerm}}
	d.Links = []model.CausalLink{
		{ID: "ln-ab", From: "a", To: "b", Rationale: "r", AssumptionIDs: []string{"asm"}},
		{ID: "ln-bt", From: "b", To: "top", Rationale: "r"},
		{ID: "ln-st", From: "side", To: "top", Rationale: "r"},
	}
	d.Assumptions = []model.Assumption{{ID: "asm", Statement: "s", Criticality: model.CriticalityHigh}}
	g, err := graph.Build(d)
	if err != nil {
		panic(err)
	}
	return d, g
}

func TestInvalidatedAssumptionFlagsEveryDownstreamOutcome(t *testing.T) {
	d, g := chain3()
	states := []model.AssumptionState{{AssumptionID: "asm", Status: model.AssumptionInvalidated}}
	p := Propagate(d, g, states)
	for _, id := range []string{"b", "top"} {
		if !slices.Equal(p.Broken[id], []string{"asm"}) || p.Cap[id] != model.StateAtRisk {
			t.Errorf("%s: broken %v cap %s", id, p.Broken[id], p.Cap[id])
		}
	}
	for _, id := range []string{"a", "side"} {
		if _, flagged := p.Broken[id]; flagged {
			t.Errorf("%s is not downstream of the broken link", id)
		}
	}
	pathway := map[string]model.State{"a": model.StateAchieved, "b": model.StateAchieved, "top": model.StateAchieved, "side": model.StateAchieved}
	p.Apply(pathway)
	if pathway["b"] != model.StateAtRisk || pathway["top"] != model.StateAtRisk || pathway["a"] != model.StateAchieved || pathway["side"] != model.StateAchieved {
		t.Errorf("after capping: %v", pathway)
	}
}

func TestAtRiskAssumptionsCapAtRisk(t *testing.T) {
	d, g := chain3()
	p := Propagate(d, g, []model.AssumptionState{{AssumptionID: "asm", Status: model.AssumptionAtRisk}})
	if p.Cap["top"] != model.StateAtRisk {
		t.Errorf("cap = %s", p.Cap["top"])
	}
	for _, s := range []model.AssumptionStatus{model.AssumptionHolding, model.AssumptionValidated, model.AssumptionUntested} {
		if p := Propagate(d, g, []model.AssumptionState{{AssumptionID: "asm", Status: s}}); len(p.Broken) != 0 {
			t.Errorf("%s flagged outcomes: %v", s, p.Broken)
		}
	}
}

func TestCriticalInvalidatedAssumptionCapsOffTrack(t *testing.T) {
	d, g := chain3()
	d.Assumptions[0].Criticality = model.CriticalityCritical
	p := Propagate(d, g, []model.AssumptionState{{AssumptionID: "asm", Status: model.AssumptionInvalidated}})
	if p.Cap["b"] != model.StateOffTrack || p.Cap["top"] != model.StateOffTrack {
		t.Errorf("caps = %v", p.Cap)
	}
	// Critical but only at risk is still just at risk.
	p = Propagate(d, g, []model.AssumptionState{{AssumptionID: "asm", Status: model.AssumptionAtRisk}})
	if p.Cap["top"] != model.StateAtRisk {
		t.Errorf("critical at-risk cap = %s", p.Cap["top"])
	}
}

func TestAssumptionsOnOutcomesAndInterventionsPropagate(t *testing.T) {
	d, _ := chain3()
	d.Links[0].AssumptionIDs = nil
	d.Assumptions = []model.Assumption{{ID: "on-outcome", Statement: "s", AppliesTo: []string{"b"}}, {ID: "on-int", Statement: "s", AppliesTo: []string{"int"}}}
	d.Interventions = []model.Intervention{{ID: "int", Title: "i", Targets: []string{"side"}}}
	g, err := graph.Build(d)
	if err != nil {
		t.Fatal(err)
	}
	p := Propagate(d, g, []model.AssumptionState{
		{AssumptionID: "on-outcome", Status: model.AssumptionInvalidated}, {AssumptionID: "on-int", Status: model.AssumptionAtRisk},
	})
	if !slices.Equal(p.Broken["b"], []string{"on-outcome"}) || !slices.Equal(p.Broken["top"], []string{"on-int", "on-outcome"}) || !slices.Equal(p.Broken["side"], []string{"on-int"}) {
		t.Errorf("broken = %v", p.Broken)
	}
}

func TestPathAndOutcomeConfidence(t *testing.T) {
	d, g := chain3()
	links := []model.LinkStatus{{LinkID: "ln-ab", Confidence: 0.5}, {LinkID: "ln-bt", Confidence: 0.8}, {LinkID: "ln-st", Confidence: 0.3}}
	outcomes, impacts := Confidence(d, g, links)
	want := map[string]float64{"top": 1, "b": 0.8, "a": 0.4, "side": 0.3}
	for id, w := range want {
		if !near6(outcomes[id], w) {
			t.Errorf("confidence(%s) = %v, want %v", id, outcomes[id], w)
		}
	}
	// The impact's headline is its best complete chain: a->b->top is 0.4, side->top is 0.3.
	if !near6(impacts["top"], 0.4) {
		t.Errorf("impact confidence = %v, want 0.4", impacts["top"])
	}

	// Two links between the same outcomes: the stronger one counts.
	d.Links = append(d.Links, model.CausalLink{ID: "ln-bt2", From: "b", To: "top"})
	g, err := graph.Build(d)
	if err != nil {
		t.Fatal(err)
	}
	links = append(links, model.LinkStatus{LinkID: "ln-bt2", Confidence: 0.9})
	outcomes, _ = Confidence(d, g, links)
	if !near6(outcomes["b"], 0.9) {
		t.Errorf("parallel links: %v", outcomes["b"])
	}
}

func TestOutcomeWithoutAPathHasNoConfidence(t *testing.T) {
	d := model.NewModel("t", "p")
	d.Outcomes = []model.Outcome{{ID: "top", Title: "t", Level: model.LevelImpact}, {ID: "lost", Title: "l", Level: model.LevelShortTerm}}
	g, err := graph.Build(d)
	if err != nil {
		t.Fatal(err)
	}
	outcomes, impacts := Confidence(d, g, nil)
	if _, ok := outcomes["lost"]; ok {
		t.Error("a dead end has a confidence")
	}
	if outcomes["top"] != 1 {
		t.Errorf("impact outcome = %v", outcomes["top"])
	}
	if _, ok := impacts["top"]; ok {
		t.Error("an impact nothing leads to has a headline confidence")
	}
}
