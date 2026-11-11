package prismx

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/mel"

	"github.com/ProductBuildersHQ/impact-model/pkg/impact"
)

func TestMELExportValidates(t *testing.T) {
	doc, st := example(t)
	for name, status := range map[string]*impact.Status{"with status": st, "without": nil} {
		out, err := ToMEL(doc, status)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if errs := mel.Errors(out.Validate(mel.DefaultValidationOptions())); len(errs) != 0 {
			t.Errorf("%s: the exported MEL plan does not validate: %+v", name, errs)
		}
	}
}

func TestMELMapping(t *testing.T) {
	doc, st := example(t)
	out, err := ToMEL(doc, st)
	if err != nil {
		t.Fatal(err)
	}

	// Results are the outcomes that are not rejected, in file order and at
	// the same level.
	var outcomes []string
	for _, o := range doc.Outcomes {
		if !o.Rejected() {
			outcomes = append(outcomes, o.ID)
		}
	}
	var results []string
	for _, r := range out.Results {
		results = append(results, r.ID)
		o, _ := doc.Outcome(r.ID)
		if r.Level != string(o.Level) || r.Statement != o.Title {
			t.Errorf("result %s = %q at %s, want %q at %s", r.ID, r.Statement, r.Level, o.Title, o.Level)
		}
	}
	if !slices.Equal(results, outcomes) {
		t.Fatalf("results = %v, want the outcomes %v", results, outcomes)
	}

	// A link leads to its target's result; a feedback link does not.
	for _, l := range doc.Links {
		from := out.Result(l.From)
		if from == nil || l.Rejected() {
			continue
		}
		if _, ok := doc.Outcome(l.To); !ok {
			continue
		}
		has := slices.Contains(from.ContributesTo, l.To)
		if l.Feedback && has {
			t.Errorf("feedback link %s should not make %s contribute to %s", l.ID, l.From, l.To)
		}
		if !l.Feedback && !has {
			t.Errorf("link %s: %s should contribute to %s", l.ID, l.From, l.To)
		}
	}

	// Activities are the interventions, producing their targets, with dates.
	for _, a := range out.Activities {
		in, _ := doc.Intervention(a.ID)
		if a.StartDate != in.StartDate || a.EndDate != in.EndDate {
			t.Errorf("activity %s dates = %s..%s, want %s..%s", a.ID, a.StartDate, a.EndDate, in.StartDate, in.EndDate)
		}
		for _, id := range a.ResultIDs {
			if !slices.Contains(in.Targets, id) {
				t.Errorf("activity %s produces %s, which it does not target", a.ID, id)
			}
		}
	}

	// An indicator measures what it measured in the model; its quality checks
	// and caveats become data quality, and its frequency is in words.
	ind := out.Indicator("ind-practice-maturity")
	if ind == nil {
		t.Fatal("ind-practice-maturity is missing")
	}
	src, _ := doc.Indicator(ind.ID)
	if ind.Measures != src.Measures || ind.Frequency != "Every six months" || ind.Current == "" {
		t.Errorf("indicator = %+v", ind)
	}
	if dq := ind.DataQuality; dq == nil || len(dq.Checks) != 1 || len(dq.Limitations) != len(src.Caveats) {
		t.Errorf("data quality = %+v", ind.DataQuality)
	}

	// An assumption's test is how it is monitored; the results that depend on
	// it include a link's target.
	asm := findAssumption(out, "asm-skills-applied")
	if asm == nil || asm.Monitoring == "" || len(asm.ResultIDs) == 0 {
		t.Errorf("assumption = %+v", asm)
	}

	// Evaluations, learning questions, and review meetings carry over; a link
	// a question names becomes both of its ends.
	if len(out.Evaluations) != len(doc.Evaluations) || len(out.Reviews) != len(doc.ReviewMeetings) {
		t.Errorf("%d evaluations and %d reviews, want %d and %d", len(out.Evaluations), len(out.Reviews), len(doc.Evaluations), len(doc.ReviewMeetings))
	}
	q := findQuestion(out, "lq-skills-applied")
	if q == nil {
		t.Fatal("lq-skills-applied is missing")
	}
	l, _ := doc.Link("ln-skills-practice")
	for _, want := range []string{"asm-skills-applied", l.From, l.To} {
		if !slices.Contains(q.RelatedIDs, want) {
			t.Errorf("question related IDs %v lack %s", q.RelatedIDs, want)
		}
	}
	if slices.Contains(q.RelatedIDs, "ln-skills-practice") {
		t.Error("a MEL plan has no links, so the link ID should not be written")
	}

	// Reports are the audiences: needs, channels as what they receive, how
	// often, and the indicators they see, headline first.
	board, _ := doc.Audience("aud-board")
	var rep *mel.Report
	for i := range out.Reports {
		if out.Reports[i].ID == "aud-board" {
			rep = &out.Reports[i]
		}
	}
	if rep == nil || rep.Needs != board.Needs || rep.Frequency != "Quarterly" || rep.Product == "" {
		t.Fatalf("board report = %+v", rep)
	}
	if len(rep.IndicatorIDs) == 0 || rep.IndicatorIDs[0] != board.Headline[0].IndicatorID {
		t.Errorf("board indicators %v should start with the headline %s", rep.IndicatorIDs, board.Headline[0].IndicatorID)
	}
	if !strings.Contains(out.Budget, "staff minutes a week") {
		t.Errorf("budget = %q", out.Budget)
	}
}

// A rejected element and every reference to it are left out.
func TestMELLeavesOutRejectedElements(t *testing.T) {
	doc, _ := example(t)
	doc.LearningQuestions[0].Provenance = &impact.Provenance{Basis: impact.BasisRejected}
	rejected := doc.LearningQuestions[0].ID
	out, err := ToMEL(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if findQuestion(out, rejected) != nil {
		t.Errorf("rejected question %s was exported", rejected)
	}
	for _, r := range out.Reviews {
		if slices.Contains(r.QuestionIDs, rejected) {
			t.Errorf("review %s still takes up the rejected question", r.ID)
		}
	}
	if errs := mel.Errors(out.Validate(nil)); len(errs) != 0 {
		t.Errorf("dangling references: %+v", errs)
	}
}

func TestMELWithoutStatusHasNoCurrentValues(t *testing.T) {
	doc, _ := example(t)
	out, err := ToMEL(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, ind := range out.Indicators {
		if ind.Current != "" {
			t.Errorf("%s carries a current value without observations", ind.ID)
		}
	}
}

func TestMELRefusesAModelWithoutOutcomes(t *testing.T) {
	doc, _ := example(t)
	for i := range doc.Outcomes {
		doc.Outcomes[i].Provenance = &impact.Provenance{Basis: impact.BasisRejected}
	}
	if _, err := ToMEL(doc, nil); !errors.Is(err, ErrNotMELReady) || !strings.Contains(err.Error(), "an outcome") {
		t.Errorf("err = %v, want ErrNotMELReady naming an outcome", err)
	}
}

func TestIndicatorType(t *testing.T) {
	for in, want := range map[impact.MeasurementType]string{
		impact.MeasurementQuantitative: mel.TypeQuantitative, impact.MeasurementThemeCount: mel.TypeQuantitative,
		impact.MeasurementQualitative: mel.TypeQualitative, impact.MeasurementHybrid: "", "": "",
	} {
		if got := indicatorType(in); got != want {
			t.Errorf("indicatorType(%q) = %q, want %q", in, got, want)
		}
	}
}

func findAssumption(p *mel.Plan, id string) *mel.Assumption {
	for i := range p.Assumptions {
		if p.Assumptions[i].ID == id {
			return &p.Assumptions[i]
		}
	}
	return nil
}

func findQuestion(p *mel.Plan, id string) *mel.LearningQuestion {
	for i := range p.Learning {
		if p.Learning[i].ID == id {
			return &p.Learning[i]
		}
	}
	return nil
}
