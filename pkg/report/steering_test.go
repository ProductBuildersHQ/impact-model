package report

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func steeringInput() StatusInput {
	f := func(v float64) *float64 { return &v }
	doc := &model.Model{
		Indicators: []model.Indicator{
			{ID: "ind-retention", Name: "Retention"},
			{ID: "ind-activation", Name: "Activation"},
			{ID: "ind-support", Name: "Support load"},
			{ID: "ind-quality", Name: "Quality"},
			{ID: "ind-reach", Name: "Reach"},
		},
		Pathways: []model.Pathway{
			{ID: "pw-app", Title: "App", Focus: &model.Focus{
				NorthStarID: "ind-retention", OneMetrics: []model.OneMetric{
					{IndicatorID: "ind-quality", From: "2026-10-01", To: "2026-12-31"},
					{IndicatorID: "ind-activation", From: "2027-01-01"},
				},
				GuardrailIDs: []string{"ind-support", "ind-quality"},
			}},
			{ID: "pw-early", Title: "Early product", Focus: &model.Focus{
				NorthStarID: "ind-reach", OneMetrics: []model.OneMetric{{IndicatorID: "ind-reach", From: "2027-01-01"}},
			}},
			{ID: "pw-dropped", Title: "Dropped", Focus: &model.Focus{NorthStarID: "ind-retention"},
				Provenance: &model.Provenance{Basis: model.BasisRejected}},
			{ID: "pw-none", Title: "No focus"},
		},
	}
	st := &model.Status{AsOf: "2027-03-31", Indicators: []model.IndicatorStatus{
		{IndicatorID: "ind-retention", State: model.StateOnTrack, Progress: f(0.5), Expected: f(0.4)},
		{IndicatorID: "ind-activation", State: model.StateAtRisk, Progress: f(0.2)},
		{IndicatorID: "ind-support", State: model.StateOffTrack},
		{IndicatorID: "ind-quality", State: model.StateOnTrack},
	}}
	return StatusInput{Doc: doc, Status: st}
}

func TestSteeringListsNorthStarThenOneMetricThenGuardrailsNeedingAttention(t *testing.T) {
	got := Steering(steeringInput())
	want := []SteeringMetric{
		{PathwayID: "pw-app", Pathway: "App", Role: RoleNorthStar, IndicatorID: "ind-retention", Indicator: "Retention", State: "on_track", Progress: "50% (40% expected)"},
		{PathwayID: "pw-app", Pathway: "App", Role: RoleOneMetric, IndicatorID: "ind-activation", Indicator: "Activation", State: "at_risk", Progress: "20%"},
		{PathwayID: "pw-app", Pathway: "App", Role: RoleGuardrail, IndicatorID: "ind-support", Indicator: "Support load", State: "off_track", Progress: "—"},
		{PathwayID: "pw-early", Pathway: "Early product", Role: RoleNorthStarAndOne, IndicatorID: "ind-reach", Indicator: "Reach", Progress: "—"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Steering =\n%+v\nwant\n%+v", got, want)
	}
}

func TestSteeringSectionIsLeftOutWithoutAFocus(t *testing.T) {
	in := steeringInput()
	if s := steering(in); !strings.HasPrefix(s, "## What each program steers by") {
		t.Errorf("section = %q", s)
	}
	in.Doc.Pathways = in.Doc.Pathways[3:]
	if s := steering(in); s != "" {
		t.Errorf("a model without a focus got a section: %q", s)
	}
}

// The one metric shown is the one whose period covers the snapshot's date.
func TestSteeringShowsTheOneMetricOfTheSnapshotsDate(t *testing.T) {
	in := steeringInput()
	in.Status.AsOf = "2026-11-15"
	var got []string
	for _, r := range Steering(in) {
		if r.PathwayID == "pw-app" {
			got = append(got, r.Role+":"+r.IndicatorID)
		}
	}
	want := []string{RoleNorthStar + ":ind-retention", RoleOneMetric + ":ind-quality", RoleGuardrail + ":ind-support"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("rows = %v, want %v", got, want)
	}
}
