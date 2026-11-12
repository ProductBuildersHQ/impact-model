package lint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// scorecardDoc places the clean base's outcomes in perspectives so that its
// strategy map runs upward, owns them, and measures each with a lead or lag
// indicator.
func scorecardDoc() *model.Model {
	d := clean()
	d.Actors = []model.Actor{{ID: "act-1", Name: "Network", Type: model.ActorOrganization}}
	perspective := map[string]model.ScorecardPerspective{
		"oc-impact": model.ScorecardMission, "oc-a": model.ScorecardCustomer, "oc-s1": model.ScorecardInternalProcess,
	}
	for i, o := range d.Outcomes {
		if p, ok := perspective[o.ID]; ok {
			d.Outcomes[i].ScorecardPerspective = p
			d.Outcomes[i].OwnerID = "act-1"
			lead := model.LagIndicator
			if p == model.ScorecardInternalProcess {
				lead = model.LeadIndicator
			}
			d.Indicators = append(d.Indicators, model.Indicator{ID: "ind-" + o.ID, Measures: o.ID, Name: o.Title, LeadLag: lead})
		}
	}
	return d
}

func TestScorecardRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	outcome := func(d *model.Model, id string) *model.Outcome { o, _ := d.Outcome(id); return o }
	tests := []struct {
		rule, element string
		defect, fix   func(d *model.Model)
	}{
		{"OM-BSC-001", "oc-s1",
			func(d *model.Model) { outcome(d, "oc-s1").ScorecardPerspective = model.ScorecardMission },
			func(d *model.Model) { outcome(d, "oc-s1").ScorecardPerspective = model.ScorecardLearningGrowth }},
		{"OM-BSC-002", "oc-a",
			func(d *model.Model) { outcome(d, "oc-a").OwnerID = "" },
			func(d *model.Model) { outcome(d, "oc-a").OwnerID = "act-1" }},
		{"OM-BSC-003", "oc-b",
			func(d *model.Model) { outcome(d, "oc-b").ScorecardPerspective = model.ScorecardCustomer },
			func(d *model.Model) {
				d.Indicators = append(d.Indicators, model.Indicator{ID: "ind-b", Measures: "oc-b", Name: "Schedules shared"})
			}},
		{"OM-BSC-004", "ind-oc-a",
			func(d *model.Model) { d.Indicators[1].LeadLag = "" },
			func(d *model.Model) { d.Indicators[1].LeadLag = model.LagIndicator }},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s/%d", tt.rule, i), func(t *testing.T) {
			d := scorecardDoc()
			if got := fires(t, d, tt.rule); len(got) != 0 {
				t.Fatalf("%s reports %v on the clean scorecard", tt.rule, got)
			}
			tt.defect(d)
			if !slices.Contains(fires(t, d, tt.rule), tt.element) {
				t.Fatalf("%s did not report %q after the defect; reported %v", tt.rule, tt.element, fires(t, d, tt.rule))
			}
			tt.fix(d)
			if slices.Contains(fires(t, d, tt.rule), tt.element) {
				t.Errorf("%s still reports %q after the repair", tt.rule, tt.element)
			}
		})
	}
}

// A model with no scorecard perspectives is not asked for scorecard fields.
func TestScorecardRulesAreQuietWithoutPerspectives(t *testing.T) {
	for _, rule := range []string{"OM-BSC-001", "OM-BSC-002", "OM-BSC-003", "OM-BSC-004"} {
		if got := fires(t, clean(), rule); len(got) != 0 {
			t.Errorf("%s reported %v on a model with no scorecard perspectives", rule, got)
		}
	}
}
