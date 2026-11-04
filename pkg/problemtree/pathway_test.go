package problemtree

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
)

// pathwayModel has two programs that share one cause.
func pathwayModel() *impact.Model {
	d := model(
		[]impact.Cause{
			{ID: "c-shared", Statement: "Thin staffing", Level: impact.CauseImmediate, PathwayIDs: []string{"pw-food", "pw-exercise"}},
			// Its parent belongs only to pw-food, so in pw-exercise it leads straight to the problem.
			{ID: "c-food", Statement: "No cooking skills", Level: impact.CauseUnderlying, Parents: []string{"c-shared"}, PathwayIDs: []string{"pw-food"}},
			{ID: "c-gym", Statement: "No safe parks", Level: impact.CauseUnderlying, PathwayIDs: []string{"pw-exercise"}},
			{ID: "c-model", Statement: "Low incomes", Level: impact.CauseStructural},
		},
		[]impact.Effect{
			{ID: "e-food", Statement: "Diabetes", PathwayIDs: []string{"pw-food"}},
		})
	d.Pathways = []impact.Pathway{
		{ID: "pw-food", Title: "Food", ProblemStatement: "Families eat poorly", OutcomeIDs: []string{"oc-food"}},
		{ID: "pw-exercise", Title: "Exercise"},
	}
	d.Outcomes = []impact.Outcome{
		{ID: "oc-food", Title: "Families cook at home", AddressesCauses: []string{"c-shared"}},
		{ID: "oc-other", Title: "Staff are hired", AddressesCauses: []string{"c-shared"}},
	}
	return d
}

func nodeIDs(tr *Tree) map[string]bool {
	out := map[string]bool{}
	for _, r := range tr.Rows {
		for _, n := range r.Nodes {
			out[n.ID] = true
		}
	}
	return out
}

func TestBuildForKeepsThePathwaysOwnCausesAndEffects(t *testing.T) {
	d := pathwayModel()
	food, err := BuildFor(d, "pw-food")
	if err != nil {
		t.Fatal(err)
	}
	got := nodeIDs(food)
	for _, id := range []string{ProblemID, "c-shared", "c-food", "e-food"} {
		if !got[id] {
			t.Errorf("pw-food tree is missing %s", id)
		}
	}
	for _, id := range []string{"c-gym", "c-model"} {
		if got[id] {
			t.Errorf("pw-food tree has %s, which does not name it", id)
		}
	}
	if r := rowOf(t, food, "c-food"); r != -2 {
		t.Errorf("c-food is in row %d, want -2, below its parent c-shared", r)
	}
	for _, r := range food.Rows {
		if r.Index == 0 && r.Nodes[0].Label != "Families eat poorly" {
			t.Errorf("core problem is %q, want the pathway's own statement", r.Nodes[0].Label)
		}
	}
}

func TestBuildForCutsParentsOutsideThePathway(t *testing.T) {
	d := pathwayModel()
	// Make c-gym name c-food, which is not in pw-exercise.
	d.Causes[2].Parents = []string{"c-food"}
	ex, err := BuildFor(d, "pw-exercise")
	if err != nil {
		t.Fatal(err)
	}
	if r := rowOf(t, ex, "c-gym"); r != -1 {
		t.Errorf("c-gym is in row %d; with its parent outside the pathway it should lead to the problem (row -1)", r)
	}
	for _, r := range ex.Rows {
		if r.Index == 0 && r.Nodes[0].Label != "Poor eating health" {
			t.Errorf("a pathway with no problem statement should use the model's, got %q", r.Nodes[0].Label)
		}
	}
}

func TestBuildForUnknownAndWholeModel(t *testing.T) {
	d := pathwayModel()
	if _, err := BuildFor(d, "pw-ghost"); err == nil {
		t.Error("an unknown pathway should be an error")
	}
	whole, err := BuildFor(d, "")
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeIDs(whole); len(got) != 6 {
		t.Errorf("the whole-model tree should hold every cause and effect and the problem, got %v", got)
	}
}

func TestReportForListsOnlyThePathwaysOutcomes(t *testing.T) {
	out, err := ReportFor(pathwayModel(), "pw-food", FormatMermaid)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Problem tree: Food", "Core problem: Families eat poorly", "`oc-food`", "`c-food`"} {
		if !strings.Contains(out, want) {
			t.Errorf("report is missing %q:\n%s", want, out)
		}
	}
	// oc-other addresses the shared cause from outside any pathway.
	if !strings.Contains(out, "Staff are hired (`oc-other`, outside any pathway)") {
		t.Errorf("an outcome from outside the pathway should be labelled:\n%s", out)
	}
	if strings.Index(out, "`oc-food`") > strings.Index(out, "`oc-other`") {
		t.Errorf("the pathway's own outcome should come first:\n%s", out)
	}
	for _, unwanted := range []string{"`c-gym`"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("report should not mention %s, which is outside the pathway's tree:\n%s", unwanted, out)
		}
	}
}
