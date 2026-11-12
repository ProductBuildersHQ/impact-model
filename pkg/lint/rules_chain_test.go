package lint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// chainDoc is the clean base with one program pathway that owns the short-term
// outcomes, an audience, and indicators of the kinds the rules tell apart.
func chainDoc() *model.Model {
	d := clean()
	d.Pathways = []model.Pathway{{
		ID: "pw-1", Title: "Program", Kind: model.PathwayProgram,
		OutcomeIDs: []string{"oc-s1", "oc-a"}, InterventionIDs: []string{"int-1"},
	}}
	d.Audiences = []model.Audience{{ID: "aud-1", Name: "Board"}}
	return d
}

func ptr[T any](v T) *T { return &v }

// Each case seeds one defect, checks the rule reports it, repairs it, and checks
// the rule falls quiet. These rules check the chain from what is done to what
// changes, the numbers each audience is told, and the North Star.
func TestChainAndHeadlineRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	headline := func(d *model.Model, ids ...string) {
		d.Audiences[0].Headline = nil
		for _, id := range ids {
			d.Audiences[0].Headline = append(d.Audiences[0].Headline, model.HeadlineMetric{IndicatorID: id})
		}
	}
	forAudience := func(i *model.Indicator) { i.AudienceIDs = []string{"aud-1"} }
	focus := func(d *model.Model, f model.Focus) { d.Pathways[0].Focus = &f }

	tests := []struct {
		rule, element string
		defect, fix   func(d *model.Model)
	}{
		// ---- from activity to impact
		{"OM-INT-003", "int-1", func(d *model.Model) {},
			func(d *model.Model) { d.Interventions[0].Outputs = []string{"Mixers held"} }},
		{"OM-INT-003", "int-1", func(d *model.Model) {},
			func(d *model.Model) {
				d.Outcomes = append(d.Outcomes, model.Outcome{ID: "oc-out", Title: "Mixers held", Level: model.LevelOutput})
				d.Interventions[0].Targets = append(d.Interventions[0].Targets, "oc-out")
			}},
		{"OM-INT-004", "int-3", func(d *model.Model) {
			d.Outcomes = append(d.Outcomes, model.Outcome{ID: "oc-dead", Title: "A dead end", Level: model.LevelShortTerm})
			d.Interventions = append(d.Interventions, model.Intervention{ID: "int-3", Title: "Do it", Targets: []string{"oc-dead"}})
		}, func(d *model.Model) {
			d.Links = append(d.Links, model.CausalLink{ID: "ln-dead-a", From: "oc-dead", To: "oc-a", Rationale: "Because it helps.", AssumptionIDs: []string{"asm-1"}})
		}},
		{"OM-INT-005", "ln-s1-a", func(d *model.Model) {
			for i := range d.Links {
				if d.Links[i].ID == "ln-s1-a" {
					d.Links[i].InterventionIDs = []string{"int-2"} // int-2 targets oc-s2, neither end of this link
				}
			}
		}, func(d *model.Model) {
			for i := range d.Links {
				if d.Links[i].ID == "ln-s1-a" {
					d.Links[i].InterventionIDs = []string{"int-1"}
				}
			}
		}},
		{"OM-INT-006", "int-1", func(d *model.Model) {},
			func(d *model.Model) { addIndicator(d, "ind-did", "int-1") }},
		{"OM-INT-007", "int-1", func(d *model.Model) {},
			func(d *model.Model) { d.Interventions[0].Resources = "Two staff and a room" }},
		{"OM-INT-008", "int-1", func(d *model.Model) {
			d.Context.TimeHorizon = &model.TimeHorizon{Start: "2026-01-01", End: "2027-12-31"}
			d.Interventions[0].StartDate, d.Interventions[0].EndDate = "2026-03-01", "2028-06-30"
		}, func(d *model.Model) { d.Interventions[0].EndDate = "2027-06-30" }},

		// ---- the numbers each audience is told
		{"OM-AUD-002", "aud-1", func(d *model.Model) { outcome(d, "oc-impact").AudienceIDs = []string{"aud-1"} },
			func(d *model.Model) { addIndicator(d, "ind-n", "oc-s1", forAudience) }},
		{"OM-HLN-001", "aud-1", func(d *model.Model) {
			var ids []string
			for i := 0; i < maxHeadline+1; i++ {
				id := fmt.Sprintf("ind-h%d", i)
				addIndicator(d, id, "oc-s1", forAudience)
				ids = append(ids, id)
			}
			headline(d, ids...)
		}, func(d *model.Model) { d.Audiences[0].Headline = d.Audiences[0].Headline[:maxHeadline] }},
		{"OM-HLN-002", "aud-1", func(d *model.Model) {
			addIndicator(d, "ind-h", "oc-s1")
			headline(d, "ind-h")
		}, func(d *model.Model) { forAudience(indicator(d, "ind-h")) }},
		{"OM-HLN-003", "aud-1", func(d *model.Model) {
			addIndicator(d, "ind-h", "oc-s1", forAudience)
			headline(d, "ind-h")
		}, func(d *model.Model) { indicator(d, "ind-h").DataAvailability = model.DataInternalExisting }},
		{"OM-HLN-004", "aud-1", func(d *model.Model) {
			for i := 0; i < curateAbove+1; i++ {
				addIndicator(d, fmt.Sprintf("ind-r%d", i), "oc-s1", forAudience)
			}
		}, func(d *model.Model) { headline(d, "ind-r0", "ind-r1") }},

		// ---- the North Star
		{"OM-NSM-001", "pw-1", func(d *model.Model) {
			d.Outcomes = append(d.Outcomes, model.Outcome{ID: "oc-out", Title: "Mixers held", Level: model.LevelOutput})
			d.Pathways[0].OutcomeIDs = append(d.Pathways[0].OutcomeIDs, "oc-out")
			addIndicator(d, "ind-ns", "oc-out")
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) {
			addIndicator(d, "ind-ns2", "oc-s1")
			d.Pathways[0].Focus.NorthStarID = "ind-ns2"
		}},
		{"OM-NSM-001", "pw-1", func(d *model.Model) {
			addIndicator(d, "ind-ns", "int-1")
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) {
			addIndicator(d, "ind-ns2", "oc-s1")
			d.Pathways[0].Focus.NorthStarID = "ind-ns2"
		}},
		{"OM-NSM-002", "pw-1", func(d *model.Model) {
			d.Ceiling = &model.Ceiling{Level: model.LevelMediumTerm}
			addIndicator(d, "ind-ns", "oc-impact")
			d.Pathways[0].OutcomeIDs = append(d.Pathways[0].OutcomeIDs, "oc-impact")
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) {
			addIndicator(d, "ind-ns2", "oc-a")
			d.Pathways[0].Focus.NorthStarID = "ind-ns2"
		}},
		{"OM-NSM-003", "pw-1", func(d *model.Model) {
			addIndicator(d, "ind-ns", "oc-s1")
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) {
			ns := indicator(d, "ind-ns")
			ns.Baseline = &model.Baseline{Value: ptr(10.0), Date: "2026-01-01"}
			ns.Targets = []model.Target{{Value: ptr(20.0), Date: "2027-01-01"}}
		}},
		{"OM-NSM-004", "pw-1", func(d *model.Model) {
			addIndicator(d, "ind-ns", "oc-s1")
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) {
			addIndicator(d, "ind-g", "oc-a")
			d.Pathways[0].Focus.GuardrailIDs = []string{"ind-g"}
		}},
		{"OM-NSM-005", "pw-1", func(d *model.Model) {
			addIndicator(d, "ind-ns", "oc-s2") // oc-s2 is not listed by pw-1
			focus(d, model.Focus{NorthStarID: "ind-ns"})
		}, func(d *model.Model) { d.Pathways[0].OutcomeIDs = append(d.Pathways[0].OutcomeIDs, "oc-s2") }},
		{"OM-NSM-006", "pw-1", func(d *model.Model) {
			addIndicator(d, "ind-ns", "oc-s1")
			focus(d, model.Focus{NorthStarID: "ind-ns", GuardrailIDs: []string{"ind-ns"}})
		}, func(d *model.Model) {
			addIndicator(d, "ind-g", "oc-a")
			d.Pathways[0].Focus.GuardrailIDs = []string{"ind-g"}
		}},
		{"OM-NSM-007", "pw-1", func(d *model.Model) {},
			func(d *model.Model) {
				addIndicator(d, "ind-ns", "oc-s1")
				focus(d, model.Focus{NorthStarID: "ind-ns"})
			}},

		// ---- framing
		{"OM-CTX-001", "", func(d *model.Model) {},
			func(d *model.Model) { d.Context.Mission = "Why we exist." }},
		{"OM-CTX-002", "", func(d *model.Model) {
			d.Outcomes = append(d.Outcomes, model.Outcome{ID: "oc-impact-2", Title: "A second impact", Level: model.LevelImpact})
		}, func(d *model.Model) { d.Context.Objective = &model.Objective{Statement: "One objective"} }},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s/%d", tt.rule, tt.element, i), func(t *testing.T) {
			d := chainDoc()
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

// A document that uses none of the new fields is not asked for them: the
// North Star, guardrail, and headline rules are about things it does not have.
func TestNewRulesAreQuietWhenTheFieldsAreAbsent(t *testing.T) {
	d := clean()
	for _, rule := range []string{"OM-HLN-001", "OM-HLN-002", "OM-HLN-003", "OM-HLN-004",
		"OM-NSM-001", "OM-NSM-002", "OM-NSM-003", "OM-NSM-004", "OM-NSM-005", "OM-NSM-006", "OM-NSM-007", "OM-AUD-002"} {
		if got := fires(t, d, rule); len(got) != 0 {
			t.Errorf("%s reported %v on a document with no audiences, pathways, or focus", rule, got)
		}
	}
}
