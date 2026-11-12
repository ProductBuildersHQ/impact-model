package lint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// treeDoc is the clean base with one effect of the core problem.
func treeDoc() *model.Model {
	d := clean()
	d.Effects = []model.Effect{{ID: "eff-1", Statement: "People miss appointments"}}
	return d
}

func TestProblemTreeRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	relieve := func(d *model.Model, outcomeID string) {
		for i := range d.Outcomes {
			if d.Outcomes[i].ID == outcomeID {
				d.Outcomes[i].AddressesEffects = []string{"eff-1"}
			}
		}
	}
	tests := []struct {
		rule, element string
		defect, fix   func(d *model.Model)
	}{
		{"IMP-EFF-001", "eff-1", func(d *model.Model) {}, func(d *model.Model) { relieve(d, "oc-impact") }},
		{"IMP-EFF-001", "eff-1", func(d *model.Model) {},
			func(d *model.Model) { d.Effects[0].OutOfScope = &model.OutOfScope{Reason: "Others handle it."} }},
		{"IMP-EFF-002", "", func(d *model.Model) { relieve(d, "oc-s1") }, func(d *model.Model) { relieve(d, "oc-impact") }},
		// An immediate cause listing a structural one as its parent is backwards.
		{"IMP-CAU-003", "cause-near", func(d *model.Model) {
			d.Causes = []model.Cause{
				{ID: "cause-deep", Statement: "No public transport", Level: model.CauseStructural},
				{ID: "cause-near", Statement: "People miss buses", Level: model.CauseImmediate, Parents: []string{"cause-deep"}},
			}
		}, func(d *model.Model) {
			d.Causes[1].Parents = nil
			d.Causes[0].Parents = []string{"cause-near"}
		}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s/%d", tt.rule, i), func(t *testing.T) {
			d := treeDoc()
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

// Rejected causes and effects are not asked to be addressed, and a model with
// no effects is not asked about them.
func TestProblemTreeRulesSkipRejectedAndAbsent(t *testing.T) {
	d := treeDoc()
	d.Effects[0].Provenance = &model.Provenance{Basis: model.BasisRejected}
	d.Causes = []model.Cause{{ID: "cause-1", Statement: "No buses", Level: model.CauseStructural,
		Provenance: &model.Provenance{Basis: model.BasisRejected}}}
	for _, rule := range []string{"IMP-EFF-001", "IMP-EFF-002", "IMP-CAU-001"} {
		if got := fires(t, d, rule); len(got) != 0 {
			t.Errorf("%s reported %v for rejected elements", rule, got)
		}
	}
	if got := fires(t, clean(), "IMP-EFF-002"); len(got) != 0 {
		t.Errorf("IMP-EFF-002 reported %v on a model with no effects", got)
	}
}
