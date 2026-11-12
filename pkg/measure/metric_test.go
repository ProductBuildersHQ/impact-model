package measure

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func rollupDoc(agg model.Aggregation) *model.Model {
	d := model.NewModel("t", "p")
	d.Metrics = []model.Metric{
		{ID: "met", Name: "M", Definition: "D", Formula: "F", Aggregation: agg, WeightMetricID: weightIf(agg)},
		{ID: "met-w", Name: "W", Definition: "D", Formula: "F", Aggregation: model.AggregationSum},
	}
	d.Pathways = []model.Pathway{{ID: "pw-1", Title: "1"}, {ID: "pw-2", Title: "2"}, {ID: "pw-3", Title: "3"}}
	d.Indicators = []model.Indicator{
		{ID: "total", Measures: "oc", Name: "t", MetricID: "met", RollupFrom: []string{"a", "b", "c"}},
		{ID: "a", Measures: "oc", Name: "a", MetricID: "met", PathwayIDs: []string{"pw-1"}},
		{ID: "b", Measures: "oc", Name: "b", MetricID: "met", PathwayIDs: []string{"pw-2"}},
		{ID: "c", Measures: "oc", Name: "c", MetricID: "met", PathwayIDs: []string{"pw-3"}},
		{ID: "wa", Measures: "oc", Name: "wa", MetricID: "met-w", PathwayIDs: []string{"pw-1"}},
		{ID: "wb", Measures: "oc", Name: "wb", MetricID: "met-w", PathwayIDs: []string{"pw-2"}},
		{ID: "wc", Measures: "oc", Name: "wc", MetricID: "met-w", PathwayIDs: []string{"pw-3"}},
	}
	return d
}

func weightIf(agg model.Aggregation) string {
	if agg == model.AggregationWeightedMean {
		return "met-w"
	}
	return ""
}

func values(pairs ...any) []model.Observation {
	var out []model.Observation
	for i := 0; i+2 < len(pairs); i += 3 {
		v := pairs[i+2].(float64)
		out = append(out, model.Observation{IndicatorID: pairs[i].(string), Date: pairs[i+1].(string), Value: &v})
	}
	return out
}

func TestRollupAggregations(t *testing.T) {
	asOf := day("2026-12-31")
	obs := values(
		"a", "2026-12-01", 10.0, "b", "2026-12-01", 20.0, "c", "2026-12-01", 20.0,
		"wa", "2026-12-01", 1.0, "wb", "2026-12-01", 1.0, "wc", "2026-12-01", 2.0,
	)
	tests := []struct {
		agg  model.Aggregation
		want float64
	}{
		{model.AggregationSum, 50},
		{model.AggregationMean, 50.0 / 3},
		{model.AggregationWeightedMean, (10*1 + 20*1 + 20*2) / 4.0},
		{model.AggregationMedian, 20},
		{model.AggregationMode, 20},
		{model.AggregationMin, 10},
		{model.AggregationMax, 20},
		{model.AggregationLatest, 20},
	}
	for _, tt := range tests {
		t.Run(string(tt.agg), func(t *testing.T) {
			d := rollupDoc(tt.agg)
			res, ok := Rollup(d, d.Indicators[0], obs, asOf)
			if !ok || res.Value == nil || math.Abs(*res.Value-tt.want) > 1e-9 {
				t.Fatalf("got %+v ok=%v, want %v", res, ok, tt.want)
			}
			if res.State != model.StateBaselineOnly || len(res.Missing) != 0 {
				t.Errorf("state %s missing %v", res.State, res.Missing)
			}
		})
	}
}

func TestRollupMedianOfAnEvenCount(t *testing.T) {
	d := rollupDoc(model.AggregationMedian)
	d.Indicators[0].RollupFrom = []string{"a", "b"}
	res, _ := Rollup(d, d.Indicators[0], values("a", "2026-12-01", 10.0, "b", "2026-12-01", 20.0), day("2026-12-31"))
	if res.Value == nil || *res.Value != 15 {
		t.Errorf("median = %v", res.Value)
	}
}

func TestRollupLatestAndTheDateCutoff(t *testing.T) {
	d := rollupDoc(model.AggregationLatest)
	obs := values("a", "2026-12-01", 10.0, "b", "2026-12-20", 77.0, "c", "2026-12-10", 5.0, "b", "2027-02-01", 999.0)
	res, _ := Rollup(d, d.Indicators[0], obs, day("2026-12-31"))
	if res.Value == nil || *res.Value != 77 {
		t.Errorf("latest = %v (a later observation must be ignored)", res.Value)
	}
}

func TestRollupReportsMissingSources(t *testing.T) {
	d := rollupDoc(model.AggregationSum)
	res, ok := Rollup(d, d.Indicators[0], values("a", "2026-12-01", 10.0), day("2026-12-31"))
	if !ok || res.Value != nil || res.State != model.StateNoData || !slices.Equal(res.Missing, []string{"b", "c"}) {
		t.Errorf("result = %+v", res)
	}
	// A weighted mean missing a weight is missing that source.
	w := rollupDoc(model.AggregationWeightedMean)
	res, _ = Rollup(w, w.Indicators[0], values("a", "2026-12-01", 1.0, "b", "2026-12-01", 1.0, "c", "2026-12-01", 1.0, "wa", "2026-12-01", 1.0, "wb", "2026-12-01", 1.0), day("2026-12-31"))
	if res.Value != nil || !slices.Equal(res.Missing, []string{"c"}) {
		t.Errorf("missing weight: %+v", res)
	}
}

func TestRollupDirectObservationsAreReportedAlongside(t *testing.T) {
	d := rollupDoc(model.AggregationSum)
	obs := values("a", "2026-12-01", 1.0, "b", "2026-12-01", 2.0, "c", "2026-12-01", 3.0, "total", "2026-12-15", 100.0)
	res, _ := Rollup(d, d.Indicators[0], obs, day("2026-12-31"))
	if !res.Direct || res.DirectValue == nil || *res.DirectValue != 100 || res.Value == nil || *res.Value != 6 {
		t.Errorf("direct value must take precedence and the computed one stay visible: %+v", res)
	}
}

func TestRollupDoesNotApply(t *testing.T) {
	d := rollupDoc(model.AggregationNone)
	if _, ok := Rollup(d, d.Indicators[0], nil, time.Now()); ok {
		t.Error("a metric that cannot be combined rolled up")
	}
	s := rollupDoc(model.AggregationSum)
	if _, ok := Rollup(s, s.Indicators[1], nil, time.Now()); ok {
		t.Error("an indicator with no rollupFrom rolled up")
	}
}
