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
				NorthStarID: "ind-retention", OneMetricID: "ind-activation",
				GuardrailIDs: []string{"ind-support", "ind-quality"},
			}},
			{ID: "pw-early", Title: "Early product", Focus: &model.Focus{
				NorthStarID: "ind-reach", OneMetricID: "ind-reach",
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
