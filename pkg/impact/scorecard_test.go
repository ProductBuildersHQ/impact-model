package impact

import (
	"slices"
	"testing"
)

// scorecardModel: out (process) -> x (no perspective) -> keep (customer) ->
// revenue (financial); out also feeds a group with y whose target is keep;
// keep loops back to x by feedback, and a rejected link and a rejected
// objective are ignored.
func scorecardModel() *Model {
	rejected := &Provenance{Basis: BasisRejected}
	return &Model{
		Outcomes: []Outcome{
			{ID: "revenue", ScorecardPerspective: ScorecardFinancial},
			{ID: "keep", ScorecardPerspective: ScorecardCustomer},
			{ID: "x"},
			{ID: "y"},
			{ID: "out", ScorecardPerspective: ScorecardInternalProcess},
			{ID: "dropped", ScorecardPerspective: ScorecardCustomer, Provenance: rejected},
		},
		Links: []CausalLink{
			{ID: "l1", From: "out", To: "x"},
			{ID: "l2", From: "x", To: "keep"},
			{ID: "l3", From: "keep", To: "revenue"},
			{ID: "l4", From: "keep", To: "x", Feedback: true},
			{ID: "l5", From: "out", To: "revenue", Provenance: rejected},
			{ID: "l6", From: "out", To: "dropped"},
			{ID: "l7", From: "int", To: "y"},
		},
		Groups: []PreconditionGroup{{ID: "g", Members: []string{"y", "out"}, Target: "keep"}},
		Interventions: []Intervention{
			{ID: "int", Targets: []string{"x"}},
			{ID: "direct", Targets: []string{"out", "revenue"}},
		},
	}
}

func TestScorecardObjectives(t *testing.T) {
	d := scorecardModel()
	var ids []string
	for _, o := range d.ScorecardObjectives() {
		ids = append(ids, o.ID)
	}
	if want := []string{"revenue", "keep", "out"}; !slices.Equal(ids, want) {
		t.Errorf("objectives = %v, want %v", ids, want)
	}
}

func TestScorecardContributesTo(t *testing.T) {
	d := scorecardModel()
	for id, want := range map[string][]string{
		"out":     {"keep"},    // through x and the group; rejected link and objective skipped
		"keep":    {"revenue"}, // the feedback link back to x is not followed
		"revenue": nil,
		"x":       {"keep"},
	} {
		if got := d.ScorecardContributesTo(id); !slices.Equal(got, want) {
			t.Errorf("ScorecardContributesTo(%s) = %v, want %v", id, got, want)
		}
	}
}

func TestScorecardObjectivesServed(t *testing.T) {
	d := scorecardModel()
	if got := d.ScorecardObjectivesServed(d.Interventions[0]); !slices.Equal(got, []string{"keep"}) {
		t.Errorf("served by int = %v, want keep through x and its own link to y", got)
	}
	if got := d.ScorecardObjectivesServed(d.Interventions[1]); !slices.Equal(got, []string{"revenue", "out"}) {
		t.Errorf("served by direct = %v, want its own targets in file order", got)
	}
}
