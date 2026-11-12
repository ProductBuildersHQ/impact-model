package lint

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
)

// Each case seeds one defect in the evaluation plan or learning agenda, checks
// the rule reports it, repairs it, and checks the rule falls quiet.
func TestLearningRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	critical := func(d *impact.Model) { d.Assumptions[0].Criticality = impact.CriticalityCritical }
	question := func(d *impact.Model, related ...string) {
		d.LearningQuestions = append(d.LearningQuestions, impact.LearningQuestion{ID: "lq-1", Question: "Why?", RelatedIDs: related})
	}
	horizon := func(d *impact.Model) {
		d.Context.TimeHorizon = &impact.TimeHorizon{Start: "2026-01-01", End: "2027-12-31"}
	}
	tests := []struct {
		rule, element string
		defect, fix   func(d *impact.Model)
	}{
		{"IMP-LRN-001", "asm-1", func(d *impact.Model) { critical(d); question(d, "oc-a") },
			func(d *impact.Model) { d.LearningQuestions[0].RelatedIDs = []string{"asm-1"} }},
		{"IMP-LRN-001", "asm-1", func(d *impact.Model) {
			critical(d)
			d.Evaluations = []impact.Evaluation{{ID: "eval-1", Title: "Final"}}
		}, func(d *impact.Model) { d.Evaluations[0].RelatedIDs = []string{"asm-1"} }},
		{"IMP-LRN-001", "asm-1", func(d *impact.Model) { critical(d); question(d) },
			func(d *impact.Model) { d.Assumptions[0].Status = impact.AssumptionHolding }},
		{"IMP-LRN-002", "lq-1", func(d *impact.Model) { question(d) },
			func(d *impact.Model) {
				d.ReviewMeetings = []impact.ReviewMeeting{{ID: "mtg-1", Name: "Quarterly", QuestionIDs: []string{"lq-1"}}}
			}},
		{"IMP-LRN-002", "lq-1", func(d *impact.Model) { question(d) },
			func(d *impact.Model) {
				d.LearningQuestions[0].Status, d.LearningQuestions[0].Answer = impact.QuestionAnswered, "Yes"
			}},
		{"IMP-LRN-003", "lq-1", func(d *impact.Model) {
			horizon(d)
			question(d)
			d.LearningQuestions[0].DueDate = "2028-06-30"
		}, func(d *impact.Model) { d.LearningQuestions[0].DueDate = "2027-06-30" }},
		{"IMP-LRN-003", "eval-1", func(d *impact.Model) {
			horizon(d)
			d.Evaluations = []impact.Evaluation{{ID: "eval-1", Title: "Final", StartDate: "2027-10-01", EndDate: "2028-01-31"}}
		}, func(d *impact.Model) { d.Evaluations[0].EndDate = "2027-12-31" }},
		{"IMP-ASM-004", "asm-1", func(d *impact.Model) {
			d.Assumptions[0].NeedsTesting, d.Assumptions[0].Test = true, "Ask providers"
		}, func(d *impact.Model) { question(d, "asm-1") }},
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

// A model without an evaluation plan or learning agenda is not asked to
// tie its critical assumptions to one.
func TestLearningRulesAreQuietWithoutAPlan(t *testing.T) {
	d := chainDoc()
	d.Assumptions[0].Criticality = impact.CriticalityCritical
	for _, rule := range []string{"IMP-LRN-001", "IMP-LRN-002", "IMP-LRN-003"} {
		if got := fires(t, d, rule); len(got) != 0 {
			t.Errorf("%s reported %v on a model with no evaluations or learning questions", rule, got)
		}
	}
}
