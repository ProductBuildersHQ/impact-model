package measure

import (
	"math"
	"testing"
	"time"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func day(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t
}

func f(v float64) *float64 { return &v }

// growth is an indicator that should rise from 0 (2026-01-01) to 100 by
// 2026-12-31, passing 50 at 2026-07-01 (2026-07-01 is day 181 of 364, about
// half the year).
func growth() model.Indicator {
	return model.Indicator{
		ID: "ind", Measures: "oc", Name: "n", Direction: model.DirectionIncrease, Frequency: "P3M",
		Baseline: &model.Baseline{Value: f(0), Date: "2026-01-01"},
		Targets:  []model.Target{{Value: f(50), Date: "2026-07-01"}, {Value: f(100), Date: "2026-12-31"}},
	}
}

func obs(date string, v float64) model.Observation {
	return model.Observation{IndicatorID: "ind", Date: date, Value: f(v)}
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestStatusBranches(t *testing.T) {
	declining := growth()
	declining.Direction = model.DirectionDecrease
	declining.Baseline.Value = f(100)
	declining.Targets = []model.Target{{Value: f(50), Date: "2026-07-01"}, {Value: f(0), Date: "2026-12-31"}}

	noTargets := growth()
	noTargets.Targets = nil
	noBaseline := growth()
	noBaseline.Baseline = nil
	noBoth := noTargets
	noBoth.Baseline = nil
	noFreq := growth()
	noFreq.Frequency = ""

	tests := []struct {
		name string
		ind  model.Indicator
		obs  []model.Observation
		asOf string
		want model.State
		hlth string
	}{
		{"baseline only: no targets", noTargets, []model.Observation{obs("2026-03-01", 10)}, "2026-04-01", model.StateBaselineOnly, ""},
		{"baseline only: no targets or observations", noTargets, nil, "2026-04-01", model.StateBaselineOnly, ""},
		{"no targets, baseline, nothing else is no data", noBoth, nil, "2026-04-01", model.StateNoData, ""},
		{"baseline only: observations but no baseline", noBaseline, []model.Observation{obs("2026-03-01", 10)}, "2026-04-01", model.StateBaselineOnly, ""},
		{"not started: nothing yet, before the first target", growth(), nil, "2026-03-01", model.StateNotStarted, ""},
		{"no data: nothing, and the first target date has passed", growth(), nil, "2026-08-01", model.StateNoData, ""},
		{"stale: last observation older than twice the frequency", growth(), []model.Observation{obs("2026-01-15", 5)}, "2026-08-01", model.StateStale, "yellow"},
		{"fresh enough: just inside twice the frequency", noFreq, []model.Observation{obs("2026-01-15", 5)}, "2026-08-01", model.StateOffTrack, "red"},
		{"achieved", growth(), []model.Observation{obs("2026-12-01", 100)}, "2026-12-15", model.StateAchieved, "green"},
		{"achieved beyond target", growth(), []model.Observation{obs("2026-12-01", 400)}, "2026-12-15", model.StateAchieved, "green"},
		{"on track: at expected progress", growth(), []model.Observation{obs("2026-06-20", 50)}, "2026-07-01", model.StateOnTrack, "green"},
		{"on track: within tolerance behind expected", growth(), []model.Observation{obs("2026-06-20", 42)}, "2026-07-01", model.StateOnTrack, "green"},
		{"at risk: behind by more than tolerance", growth(), []model.Observation{obs("2026-06-20", 35)}, "2026-07-01", model.StateAtRisk, "yellow"},
		{"off track: far behind", growth(), []model.Observation{obs("2026-06-20", 10)}, "2026-07-01", model.StateOffTrack, "red"},
		{"decrease direction: on track", declining, []model.Observation{obs("2026-06-20", 50)}, "2026-07-01", model.StateOnTrack, "green"},
		{"decrease direction: moving the wrong way is off track", declining, []model.Observation{obs("2026-06-20", 130)}, "2026-07-01", model.StateOffTrack, "red"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateIndicator(tt.ind, tt.obs, Options{AsOf: day(tt.asOf)})
			if got.State != tt.want || got.Health != tt.hlth {
				t.Errorf("state %s health %q; want %s %q (progress %v expected %v)", got.State, got.Health, tt.want, tt.hlth, got.Progress, got.Expected)
			}
		})
	}
}

func TestProgressAndExpected(t *testing.T) {
	got := EvaluateIndicator(growth(), []model.Observation{obs("2026-05-01", 25)}, Options{AsOf: day("2026-05-02")})
	if got.Progress == nil || !near(*got.Progress, 0.25) {
		t.Fatalf("progress = %v, want 0.25", got.Progress)
	}
	// 2026-05-02 is 121 days after 2026-01-01; the first target (0.5) is 181 days out.
	if want := 0.5 * 121.0 / 181.0; got.Expected == nil || !near(*got.Expected, want) {
		t.Errorf("expected = %v, want %v", got.Expected, want)
	}
	if got.ObservationCount != 1 || got.Latest == nil || got.Latest.Date != "2026-05-01" {
		t.Errorf("latest = %+v, count %d", got.Latest, got.ObservationCount)
	}

	after := EvaluateIndicator(growth(), []model.Observation{obs("2027-01-15", 100)}, Options{AsOf: day("2027-02-01")})
	if after.Expected == nil || *after.Expected != 1 {
		t.Errorf("expected after the last target = %v, want 1", after.Expected)
	}
	before := EvaluateIndicator(growth(), []model.Observation{obs("2025-12-01", 0)}, Options{AsOf: day("2025-12-15")})
	if before.Expected == nil || *before.Expected != 0 {
		t.Errorf("expected before the baseline date = %v, want 0", before.Expected)
	}
}

func TestProgressIsCappedAndNeverNegative(t *testing.T) {
	over := EvaluateIndicator(growth(), []model.Observation{obs("2026-12-01", 1000)}, Options{AsOf: day("2026-12-15")})
	if over.Progress == nil || *over.Progress != 1.5 {
		t.Errorf("capped progress = %v, want 1.5", over.Progress)
	}
	under := EvaluateIndicator(growth(), []model.Observation{obs("2026-06-01", -40)}, Options{AsOf: day("2026-06-15")})
	if under.Progress == nil || *under.Progress != 0 {
		t.Errorf("progress below the baseline = %v, want 0", under.Progress)
	}
}

func TestObservationSelection(t *testing.T) {
	other := model.Observation{IndicatorID: "someone-else", Date: "2026-06-01", Value: f(99)}
	future := obs("2026-09-01", 90)
	bad := model.Observation{IndicatorID: "ind", Date: "not-a-date", Value: f(70)}
	noValue := model.Observation{IndicatorID: "ind", Date: "2026-06-10"}
	observations := []model.Observation{obs("2026-05-01", 20), other, future, bad, noValue, obs("2026-06-01", 30), obs("2026-06-01", 35)}
	got := EvaluateIndicator(growth(), observations, Options{AsOf: day("2026-06-15")})
	if got.ObservationCount != 3 {
		t.Errorf("count = %d, want 3 (another indicator, a future date, a bad date, and no value are all ignored)", got.ObservationCount)
	}
	if got.Latest == nil || *got.Latest.Value != 35 {
		t.Errorf("latest = %+v, want the later of the two on 2026-06-01", got.Latest)
	}
}

func TestMaintainIndicator(t *testing.T) {
	ind := model.Indicator{
		ID: "ind", Measures: "oc", Name: "n", Direction: model.DirectionMaintain, Tolerance: 2,
		Baseline: &model.Baseline{Value: f(90), Date: "2026-01-01"},
		Targets:  []model.Target{{Value: f(90), Date: "2026-12-31"}},
	}
	in := EvaluateIndicator(ind, []model.Observation{obs("2026-12-01", 91)}, Options{AsOf: day("2026-12-15")})
	if in.State != model.StateAchieved {
		t.Errorf("within the band: %s", in.State)
	}
	out := EvaluateIndicator(ind, []model.Observation{obs("2026-12-01", 80)}, Options{AsOf: day("2026-12-15")})
	if out.State != model.StateOffTrack {
		t.Errorf("outside the band: %s", out.State)
	}
}

func TestTargetEqualToBaseline(t *testing.T) {
	ind := growth()
	ind.Targets = []model.Target{{Value: f(0), Date: "2026-12-31"}}
	held := EvaluateIndicator(ind, []model.Observation{obs("2026-12-01", 5)}, Options{AsOf: day("2026-12-15")})
	if held.State != model.StateAchieved {
		t.Errorf("target equal to baseline, value not below: %s", held.State)
	}
	ind.Direction = model.DirectionDecrease
	rose := EvaluateIndicator(ind, []model.Observation{obs("2026-12-01", 5)}, Options{AsOf: day("2026-12-15")})
	if rose.State == model.StateAchieved {
		t.Errorf("a decrease indicator that rose cannot have achieved a flat target: %s", rose.State)
	}
}

func TestDirectionFollowsTheTargetWhenUnstated(t *testing.T) {
	ind := growth()
	ind.Direction = ""
	got := EvaluateIndicator(ind, []model.Observation{obs("2026-06-20", 50)}, Options{AsOf: day("2026-07-01")})
	if got.State != model.StateOnTrack || got.Progress == nil || !near(*got.Progress, 0.5) {
		t.Errorf("state %s progress %v", got.State, got.Progress)
	}
}

func TestQualitativeIndicator(t *testing.T) {
	ind := model.Indicator{
		ID: "ind", Measures: "oc", Name: "n", Levels: []string{"none", "informal", "guidance", "rule"},
		Baseline: &model.Baseline{Level: "none", Date: "2026-01-01"},
		Targets:  []model.Target{{Level: "guidance", Date: "2026-12-31"}},
	}
	lvl := func(date, level string) model.Observation {
		return model.Observation{IndicatorID: "ind", Date: date, Level: level}
	}
	half := EvaluateIndicator(ind, []model.Observation{lvl("2026-12-01", "informal")}, Options{AsOf: day("2026-12-15")})
	if half.Progress == nil || !near(*half.Progress, 0.5) || half.State != model.StateOffTrack {
		t.Errorf("halfway: progress %v state %s", half.Progress, half.State)
	}
	done := EvaluateIndicator(ind, []model.Observation{lvl("2026-12-01", "rule")}, Options{AsOf: day("2026-12-15")})
	if done.State != model.StateAchieved {
		t.Errorf("beyond the target level: %s", done.State)
	}
	unknown := EvaluateIndicator(ind, []model.Observation{lvl("2026-12-01", "bogus")}, Options{AsOf: day("2026-12-15")})
	if unknown.ObservationCount != 0 {
		t.Errorf("a level the indicator does not have counted as an observation")
	}
}

func TestToleranceOption(t *testing.T) {
	o := []model.Observation{obs("2026-06-20", 35)}
	if st := EvaluateIndicator(growth(), o, Options{AsOf: day("2026-07-01")}).State; st != model.StateAtRisk {
		t.Fatalf("default tolerance: %s", st)
	}
	if st := EvaluateIndicator(growth(), o, Options{AsOf: day("2026-07-01"), Tolerance: 0.2}).State; st != model.StateOnTrack {
		t.Errorf("tolerance 0.2: %s", st)
	}
}

func TestStaleUsesCalendarMonths(t *testing.T) {
	ind := growth() // quarterly, so stale after six months
	o := []model.Observation{obs("2026-01-31", 5)}
	if st := EvaluateIndicator(ind, o, Options{AsOf: day("2026-07-30")}).State; st == model.StateStale {
		t.Errorf("stale a day early")
	}
	if st := EvaluateIndicator(ind, o, Options{AsOf: day("2026-08-05")}).State; st != model.StateStale {
		t.Errorf("not stale after six months: %s", st)
	}
}

func TestAddDuration(t *testing.T) {
	base := day("2026-01-15")
	for iso, want := range map[string]string{"P1M": "2026-03-15", "P1W": "2026-01-29", "P1Y": "2028-01-15", "PT36H": "2026-01-18", "P1Y6M": "2029-01-15"} {
		got, ok := addDuration(base, iso, 2)
		if !ok || got.Format(time.DateOnly) != want {
			// PT36H twice is 72 hours, three days.
			t.Errorf("addDuration(%q, 2) = %v, %v; want %s", iso, got.Format(time.DateOnly), ok, want)
		}
	}
	for _, bad := range []string{"", "P", "monthly", "P0D"} {
		if _, ok := addDuration(base, bad, 1); ok {
			t.Errorf("%q accepted", bad)
		}
	}
}

func TestHealth(t *testing.T) {
	want := map[model.State]string{
		model.StateAchieved: "green", model.StateOnTrack: "green", model.StateAtRisk: "yellow", model.StateStale: "yellow",
		model.StateOffTrack: "red", model.StateNotStarted: "", model.StateNoData: "", model.StateBaselineOnly: "", model.StateUnmeasured: "",
	}
	for st, h := range want {
		if got := Health(st); got != h {
			t.Errorf("Health(%s) = %q, want %q", st, got, h)
		}
	}
	for _, st := range model.AllStates() {
		if _, ok := want[st]; !ok {
			t.Errorf("state %s has no health mapping in this test", st)
		}
	}
}
