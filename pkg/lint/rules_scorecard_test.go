package lint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
)

// scorecardDoc places the clean base's outcomes in perspectives so that its
// strategy map runs upward, owns them, and measures each with a lead or lag
// indicator.
func scorecardDoc() *impact.Model {
	d := clean()
	d.Actors = []impact.Actor{{ID: "act-1", Name: "Network", Type: impact.ActorOrganization}}
	perspective := map[string]impact.ScorecardPerspective{
		"oc-impact": impact.ScorecardMission, "oc-a": impact.ScorecardCustomer, "oc-s1": impact.ScorecardInternalProcess,
	}
	for i, o := range d.Outcomes {
		if p, ok := perspective[o.ID]; ok {
			d.Outcomes[i].ScorecardPerspective = p
			d.Outcomes[i].OwnerID = "act-1"
			lead := impact.LagIndicator
			if p == impact.ScorecardInternalProcess {
				lead = impact.LeadIndicator
			}
			d.Indicators = append(d.Indicators, impact.Indicator{ID: "ind-" + o.ID, Measures: o.ID, Name: o.Title, LeadLag: lead})
		}
	}
	return d
}

func TestScorecardRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	outcome := func(d *impact.Model, id string) *impact.Outcome { o, _ := d.Outcome(id); return o }
	tests := []struct {
		rule, element string
		defect, fix   func(d *impact.Model)
	}{
		{"IMP-BSC-001", "oc-s1",
			func(d *impact.Model) { outcome(d, "oc-s1").ScorecardPerspective = impact.ScorecardMission },
			func(d *impact.Model) { outcome(d, "oc-s1").ScorecardPerspective = impact.ScorecardLearningGrowth }},
		{"IMP-BSC-002", "oc-a",
			func(d *impact.Model) { outcome(d, "oc-a").OwnerID = "" },
			func(d *impact.Model) { outcome(d, "oc-a").OwnerID = "act-1" }},
		{"IMP-BSC-003", "oc-b",
			func(d *impact.Model) { outcome(d, "oc-b").ScorecardPerspective = impact.ScorecardCustomer },
			func(d *impact.Model) {
				d.Indicators = append(d.Indicators, impact.Indicator{ID: "ind-b", Measures: "oc-b", Name: "Schedules shared"})
			}},
		{"IMP-BSC-004", "ind-oc-a",
			func(d *impact.Model) { d.Indicators[1].LeadLag = "" },
			func(d *impact.Model) { d.Indicators[1].LeadLag = impact.LagIndicator }},
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
	for _, rule := range []string{"IMP-BSC-001", "IMP-BSC-002", "IMP-BSC-003", "IMP-BSC-004"} {
		if got := fires(t, clean(), rule); len(got) != 0 {
			t.Errorf("%s reported %v on a model with no scorecard perspectives", rule, got)
		}
	}
}
