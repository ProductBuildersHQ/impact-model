package report

import (
	"slices"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func TestDMAICPlanGolden(t *testing.T) {
	out, err := DMAICPlan(fullDoc(t), "")
	if err != nil {
		t.Fatal(err)
	}
	golden(t, "dmaic-plan.md", out)
}

func TestDMAICPlanWithoutMarkedPathways(t *testing.T) {
	doc := fullDoc(t)
	for i := range doc.Pathways {
		doc.Pathways[i].Projections = nil
	}
	out, err := DMAICPlan(doc, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No pathway is marked for DMAIC") {
		t.Errorf("want the note that nothing is marked:\n%s", out)
	}
	// A named pathway is planned whether or not it is marked.
	named, err := DMAICPlan(doc, doc.Pathways[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(named, "### "+doc.Pathways[0].Title) {
		t.Errorf("a named pathway should be planned:\n%s", named)
	}
	if _, err := DMAICPlan(doc, "pw-ghost"); err == nil {
		t.Error("an unknown pathway should be an error")
	}
}

func TestDMAICUseAndNeeds(t *testing.T) {
	v := 3.0
	doc := model.NewModel("m", "Problem")
	doc.Outcomes = []model.Outcome{
		{ID: "oc-long", Title: "Long", Level: model.LevelLongTerm},
		{ID: "oc-short", Title: "Short", Level: model.LevelShortTerm},
		{ID: "oc-relief", Title: "Relief", Level: model.LevelMediumTerm, AddressesEffects: []string{"e-1"}},
	}
	doc.Interventions = []model.Intervention{{ID: "int-a", Title: "A", Targets: []string{"oc-short"}}}
	p := &model.Pathway{ID: "pw", Title: "P", OutcomeIDs: []string{"oc-long", "oc-short", "oc-relief"}, InterventionIDs: []string{"int-a"}}

	tests := []struct {
		ind  model.Indicator
		want DMAICPhaseUse
	}{
		{model.Indicator{Measures: "oc-long"}, DMAICPhaseUse{Define: true, Measure: true}},
		{model.Indicator{Measures: "oc-short", Frequency: "P1M", Disaggregation: []string{"size"}}, DMAICPhaseUse{Measure: true, Analyze: true, Improve: true, Control: true}},
		{model.Indicator{Measures: "oc-relief", Frequency: "P2Y"}, DMAICPhaseUse{Define: true, Measure: true}},
		{model.Indicator{Measures: "int-a", Frequency: "P1Y"}, DMAICPhaseUse{Measure: true, Improve: true, Control: true}},
	}
	for _, tt := range tests {
		if got := DMAICUse(doc, p, tt.ind); got != tt.want {
			t.Errorf("DMAICUse(%s, %s) = %+v, want %+v", tt.ind.Measures, tt.ind.Frequency, got, tt.want)
		}
	}

	ready := model.Indicator{Baseline: &model.Baseline{Value: &v}, Targets: []model.Target{{Value: &v, Date: "2027-01-01"}}, Frequency: "P2W", DataSource: "CRM"}
	if needs := DMAICNeeds(ready); len(needs) != 0 {
		t.Errorf("a fortnightly indicator with baseline, target, and source needs nothing, got %v", needs)
	}
	ready.Frequency = "P3M"
	if needs := DMAICNeeds(ready); len(needs) != 1 || !strings.Contains(needs[0], "control limits") {
		t.Errorf("a quarterly indicator needs more measurements, got %v", needs)
	}
	if needs := DMAICNeeds(model.Indicator{}); !slices.Equal(needs, []string{"baseline", "numeric target", "data source", "frequency"}) {
		t.Errorf("an empty indicator needs everything, got %v", needs)
	}
}
