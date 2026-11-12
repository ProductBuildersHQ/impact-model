package prismx

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/vsem"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
)

func TestVSEMExportValidates(t *testing.T) {
	doc, st := example(t)
	for name, status := range map[string]*impact.Status{"with status": st, "without": nil} {
		out, err := ToVSEM(doc, status)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if errs := vsem.Errors(out.Validate(vsem.DefaultValidationOptions())); len(errs) != 0 {
			t.Errorf("%s: the exported VSEM does not validate: %+v", name, errs)
		}
	}
}

func TestVSEMMapping(t *testing.T) {
	doc, st := example(t)
	out, err := ToVSEM(doc, st)
	if err != nil {
		t.Fatal(err)
	}

	if out.Vision != doc.Context.Vision || out.Mission != doc.Context.Mission {
		t.Errorf("vision, mission = %q, %q; want context.vision and context.mission", out.Vision, out.Mission)
	}
	if h := doc.Context.TimeHorizon; out.Metadata.Period != h.Start+" to "+h.End {
		t.Errorf("period = %q, want the time horizon", out.Metadata.Period)
	}

	// Each pathway that is not rejected is a strategy, in file order.
	var pathways []string
	for _, p := range doc.Pathways {
		if !p.Rejected() {
			pathways = append(pathways, p.ID)
		}
	}
	var strategies []string
	for _, s := range out.Strategies {
		strategies = append(strategies, s.ID)
	}
	if !slices.Equal(strategies, pathways) {
		t.Fatalf("strategies = %v, want the pathways %v", strategies, pathways)
	}

	learning := out.Strategies[0]
	p, _ := doc.Pathway(learning.ID)
	if len(learning.Resources) != len(p.Resources) || learning.Resources[0].Allocation != p.Resources[0].Allocation {
		t.Errorf("resources = %+v, want the pathway's %+v", learning.Resources, p.Resources)
	}
	if owner, _ := doc.Actor(p.OwnerID); learning.Owner != owner.Name {
		t.Errorf("owner = %q, want the actor's name %q", learning.Owner, owner.Name)
	}

	// Execution is the strategy's interventions, with kind, status, and dates.
	work := out.ExecutionFor(learning.ID)
	var ids []string
	for _, e := range work {
		ids = append(ids, e.ID)
	}
	if !slices.Equal(ids, p.InterventionIDs) {
		t.Fatalf("execution for %s = %v, want its interventions %v", learning.ID, ids, p.InterventionIDs)
	}
	for _, e := range work {
		in, _ := doc.Intervention(e.ID)
		if e.StartDate != in.StartDate || e.EndDate != in.EndDate {
			t.Errorf("%s dates = %s..%s, want %s..%s", e.ID, e.StartDate, e.EndDate, in.StartDate, in.EndDate)
		}
		wantKind := vsem.KindInitiative
		if in.Type == impact.InterventionProgram {
			wantKind = vsem.KindProgram
		}
		if e.Kind != wantKind {
			t.Errorf("%s kind = %q, want %q for type %q", e.ID, e.Kind, wantKind, in.Type)
		}
		if e.Status == "" && in.Status != "" && in.Status != impact.InterventionPaused {
			t.Errorf("%s has no status, from %q", e.ID, in.Status)
		}
	}

	// The North Star leads the strategy's metrics, and guardrails follow.
	metrics := out.MetricsFor(learning.ID)
	if len(metrics) == 0 || metrics[0].ID != p.Focus.NorthStarID {
		t.Fatalf("first metric of %s is not its North Star %q: %+v", learning.ID, p.Focus.NorthStarID, metrics)
	}
	if !strings.HasPrefix(metrics[0].Description, "North Star") {
		t.Errorf("North Star description = %q", metrics[0].Description)
	}
	for _, g := range p.Focus.GuardrailIDs {
		if !slices.ContainsFunc(metrics, func(m vsem.Metric) bool { return m.ID == g }) {
			t.Errorf("guardrail %s is not a metric of %s", g, learning.ID)
		}
	}
	ns, _ := doc.Indicator(p.Focus.NorthStarID)
	if owner, _ := doc.Actor(ns.OwnerID); metrics[0].Owner != owner.Name {
		t.Errorf("North Star owner = %q, want %q", metrics[0].Owner, owner.Name)
	}
	if metrics[0].Current == "" || metrics[0].Status == "" {
		t.Errorf("with status, the North Star should carry its current value and status: %+v", metrics[0])
	}

	// The vision's metrics are indicators of impact-level outcomes.
	vision := out.VisionMetrics()
	if len(vision) == 0 {
		t.Fatal("no vision metrics; the example has impact-level indicators")
	}
	for _, m := range vision {
		ind, _ := doc.Indicator(m.ID)
		o, _ := doc.Outcome(ind.Measures)
		if o.Level != impact.LevelImpact {
			t.Errorf("vision metric %s measures %s at level %s, want impact", m.ID, o.ID, o.Level)
		}
	}
}

func TestVSEMWithoutStatusHasNoCurrentValues(t *testing.T) {
	doc, _ := example(t)
	out, err := ToVSEM(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range out.Metrics {
		if m.Current != "" || m.Progress != 0 || m.Status != "" {
			t.Errorf("%s carries status without observations: %+v", m.ID, m)
		}
	}
}

// An intervention or indicator that serves two pathways is listed once, under
// the first, so the VSEM's IDs stay unique.
func TestVSEMListsSharedElementsOnce(t *testing.T) {
	doc, st := example(t)
	first, second := &doc.Pathways[0], &doc.Pathways[1]
	shared := first.InterventionIDs[0]
	second.InterventionIDs = append(second.InterventionIDs, shared)
	second.Focus.GuardrailIDs = append(second.Focus.GuardrailIDs, first.Focus.NorthStarID)

	out, err := ToVSEM(doc, st)
	if err != nil {
		t.Fatal(err)
	}
	if errs := vsem.Errors(out.Validate(nil)); len(errs) != 0 {
		t.Fatalf("shared elements broke the VSEM: %+v", errs)
	}
	for _, e := range out.Execution {
		if e.ID == shared && e.StrategyID != first.ID {
			t.Errorf("%s is listed under %s, want the first pathway %s", shared, e.StrategyID, first.ID)
		}
	}
}

func TestVSEMRefusesWhatItCannotDraw(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(d *impact.Model)
		want   string
	}{
		{"no vision or mission", func(d *impact.Model) { d.Context.Vision, d.Context.Mission = "", "" }, "context.vision or context.mission"},
		{"only rejected pathways", func(d *impact.Model) {
			var rejected []impact.Pathway
			for _, p := range d.Pathways {
				if p.Rejected() {
					rejected = append(rejected, p)
				}
			}
			d.Pathways = rejected
		}, "a pathway"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, _ := example(t)
			tt.mutate(doc)
			_, err := ToVSEM(doc, nil)
			if !errors.Is(err, ErrNotVSEMReady) || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want ErrNotVSEMReady naming %q", err, tt.want)
			}
		})
	}
}

func TestVSEMAcceptsAMissionWithoutAVision(t *testing.T) {
	doc, _ := example(t)
	doc.Context.Vision = ""
	doc.Context.Mission = "Strengthen local nonprofits."
	out, err := ToVSEM(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out.Statement() != doc.Context.Mission {
		t.Errorf("statement = %q, want the mission", out.Statement())
	}
}

func TestFrequencyLabel(t *testing.T) {
	for in, want := range map[string]string{
		"P1D": "Daily", "P1W": "Weekly", "P7D": "Weekly", "P1M": "Monthly",
		"P3M": "Quarterly", "P6M": "Every six months", "P1Y": "Yearly", "P12M": "Yearly", "P2W": "P2W", "": "",
	} {
		if got := frequencyLabel(in); got != want {
			t.Errorf("frequencyLabel(%q) = %q, want %q", in, got, want)
		}
	}
}

// The intermediary example sets every field the strict options ask for, so its
// export shows that a complete model gives a complete VSEM.
func TestVSEMExportOfACompleteModelMeetsStrictOptions(t *testing.T) {
	doc, st := example(t)
	out, err := ToVSEM(doc, st)
	if err != nil {
		t.Fatal(err)
	}
	if errs := out.Validate(vsem.StrictValidationOptions()); len(errs) != 0 {
		t.Errorf("strict validation: %+v", errs)
	}
}
