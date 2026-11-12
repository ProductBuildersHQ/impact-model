package edit

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/validate"
)

func focusDoc(t *testing.T, ms ...model.OneMetric) *model.Model {
	t.Helper()
	doc := loadFull(t)
	doc.Pathways[0].Focus = &model.Focus{OneMetrics: ms}
	return doc
}

func applyFocus(t *testing.T, doc *model.Model, m model.OneMetric) *model.Model {
	t.Helper()
	p, err := FocusOn(doc, doc.Pathways[0].ID, m)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Apply(doc, p, Options{})
	if err != nil || !res.Applied {
		t.Fatalf("focus rejected: %v %s %v", err, res.Reason, res.IssuesIntroduced)
	}
	if validate.HasErrors(validate.Document(res.Document)) {
		t.Fatalf("result has errors: %v", validate.Document(res.Document))
	}
	return res.Document
}

// A new focus closes the open one the day before it starts, and both stay on
// record.
func TestFocusOnClosesTheCurrentFocus(t *testing.T) {
	ind := func(doc *model.Model, i int) string { return doc.Indicators[i].ID }
	doc := focusDoc(t)
	first := model.OneMetric{IndicatorID: ind(doc, 0), From: "2027-01-01", Stage: "activation"}
	doc.Pathways[0].Focus.OneMetrics = []model.OneMetric{first}
	next := model.OneMetric{IndicatorID: ind(doc, 1), From: "2027-04-01", Reason: "Retention is now the bottleneck"}

	got := applyFocus(t, doc, next).Pathways[0].Focus.OneMetrics
	first.To = "2027-03-31"
	if want := []model.OneMetric{first, next}; !reflect.DeepEqual(got, want) {
		t.Errorf("one metrics = %+v, want %+v", got, want)
	}
}

// A pathway without a focus, or with one but no one metric yet, gets its first.
func TestFocusOnStartsTheFirstFocus(t *testing.T) {
	doc := focusDoc(t)
	doc.Pathways[0].Focus = nil
	m := model.OneMetric{IndicatorID: doc.Indicators[0].ID, From: "2027-01-01"}
	if got := applyFocus(t, doc, m).Pathways[0].Focus.OneMetrics; !reflect.DeepEqual(got, []model.OneMetric{m}) {
		t.Errorf("without a focus: %+v", got)
	}
	doc = focusDoc(t)
	if got := applyFocus(t, doc, m).Pathways[0].Focus.OneMetrics; !reflect.DeepEqual(got, []model.OneMetric{m}) {
		t.Errorf("with an empty focus: %+v", got)
	}
}

func TestFocusOnRefusesOverlapsAndBadInput(t *testing.T) {
	doc := focusDoc(t,
		model.OneMetric{IndicatorID: "a", From: "2026-01-01", To: "2026-06-30"},
		model.OneMetric{IndicatorID: "b", From: "2027-01-01"})
	id := doc.Indicators[0].ID
	for name, m := range map[string]model.OneMetric{
		"before the open one starts": {IndicatorID: id, From: "2026-12-01"},
		"on the day it starts":       {IndicatorID: id, From: "2027-01-01"},
		"overlapping a closed one":   {IndicatorID: id, From: "2026-03-01", To: "2026-08-01"},
		"not a date":                 {IndicatorID: id, From: "April"},
		"no indicator":               {From: "2027-06-01"},
	} {
		if _, err := FocusOn(doc, doc.Pathways[0].ID, m); !errors.Is(err, ErrPatch) {
			t.Errorf("%s: err = %v, want ErrPatch", name, err)
		}
	}
	if _, err := FocusOn(doc, "pw-ghost", model.OneMetric{IndicatorID: id, From: "2027-06-01"}); !errors.Is(err, ErrPatch) {
		t.Errorf("unknown pathway: err = %v", err)
	}
}
