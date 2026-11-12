package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/dmaic"
	"github.com/grokify/prism-roadmap/goals/mel"
	"github.com/grokify/prism-roadmap/goals/ogsm"
	"github.com/grokify/prism-roadmap/goals/okr"
	"github.com/grokify/prism-roadmap/goals/v2mom"
	"github.com/grokify/prism-roadmap/goals/vsem"
)

func TestExportOKRAndJourneys(t *testing.T) {
	out := filepath.Join(t.TempDir(), "okr.json")
	r := run(t, "", "export", "okr", exTheory, "-O", exObs, "-o", out)
	wantCode(t, r, ExitOK)
	doc, err := okr.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := okr.Errors(doc.Validate(okr.DefaultValidationOptions())); len(errs) != 0 {
		t.Errorf("exported OKR does not validate: %+v", errs)
	}
	if doc.Metadata.Period != "2026-09-30" {
		t.Errorf("period = %q", doc.Metadata.Period)
	}

	// Without observations the key results have no current value.
	bare := run(t, "", "export", "okr", exTheory)
	wantCode(t, bare, ExitOK)
	if strings.Contains(bare.stdout, `"current"`) {
		t.Error("an export with no observations has current values")
	}

	j := run(t, "", "export", "journey", exTheory, "-O", exObs)
	wantCode(t, j, ExitOK)
	var journeys []map[string]any
	if err := json.Unmarshal([]byte(j.stdout), &journeys); err != nil || len(journeys) == 0 {
		t.Errorf("journeys: %v %.200s", err, j.stdout)
	}
	if again := run(t, "", "export", "journey", exTheory, "-O", exObs); again.stdout != j.stdout {
		t.Error("two exports differ")
	}
	wantCode(t, run(t, "", "export", "okr", exTheory+".none"), ExitUsageErr)
}

func TestExportV2MOM(t *testing.T) {
	out := filepath.Join(t.TempDir(), "v2mom.json")
	r := run(t, "", "export", "v2mom", exTheory, "-O", exObs, "-o", out)
	wantCode(t, r, ExitOK)
	doc, err := v2mom.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := v2mom.Errors(doc.Validate(v2mom.DefaultValidationOptions())); len(errs) != 0 {
		t.Errorf("exported V2MOM does not validate: %+v", errs)
	}
	if doc.Vision == "" || len(doc.Values) == 0 || len(doc.Methods) == 0 {
		t.Errorf("vision, values, or methods missing: %+v", doc)
	}
	first := run(t, "", "export", "v2mom", exTheory, "-O", exObs)
	if again := run(t, "", "export", "v2mom", exTheory, "-O", exObs); again.stdout != first.stdout {
		t.Error("two exports differ")
	}

	// A model without a vision, values, or pathways is refused, naming them.
	refused := run(t, "", "export", "v2mom", "../../examples/market-adoption/impact-model.yaml")
	wantCode(t, refused, ExitUsageErr)
	for _, want := range []string{"cannot be drawn as a V2MOM", "context.vision", "context.values"} {
		if !strings.Contains(refused.err, want) {
			t.Errorf("refusal does not say %q: %s", want, refused.err)
		}
	}
}

func TestExportVSEM(t *testing.T) {
	out := filepath.Join(t.TempDir(), "vsem.json")
	r := run(t, "", "export", "vsem", exTheory, "-O", exObs, "-o", out)
	wantCode(t, r, ExitOK)
	doc, err := vsem.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := vsem.Errors(doc.Validate(vsem.DefaultValidationOptions())); len(errs) != 0 {
		t.Errorf("exported VSEM does not validate: %+v", errs)
	}
	if doc.Vision == "" || len(doc.Strategies) == 0 || len(doc.Execution) == 0 || len(doc.Metrics) == 0 {
		t.Errorf("vision, strategies, execution, or metrics missing: %+v", doc)
	}
	first := run(t, "", "export", "vsem", exTheory, "-O", exObs)
	if again := run(t, "", "export", "vsem", exTheory, "-O", exObs); again.stdout != first.stdout {
		t.Error("two exports differ")
	}

	// A model without a vision, a mission, or pathways is refused, naming them.
	refused := run(t, "", "export", "vsem", "../../examples/market-adoption/impact-model.yaml")
	wantCode(t, refused, ExitUsageErr)
	for _, want := range []string{"cannot be drawn as a VSEM", "context.vision or context.mission", "a pathway"} {
		if !strings.Contains(refused.err, want) {
			t.Errorf("refusal does not say %q: %s", want, refused.err)
		}
	}
}

func TestExportMEL(t *testing.T) {
	out := filepath.Join(t.TempDir(), "mel.json")
	r := run(t, "", "export", "mel", exTheory, "-O", exObs, "-o", out)
	wantCode(t, r, ExitOK)
	plan, err := mel.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := mel.Errors(plan.Validate(mel.DefaultValidationOptions())); len(errs) != 0 {
		t.Errorf("exported MEL plan does not validate: %+v", errs)
	}
	if len(plan.Results) == 0 || len(plan.Indicators) == 0 || len(plan.Evaluations) == 0 || len(plan.Learning) == 0 || len(plan.Reviews) == 0 {
		t.Errorf("results, indicators, evaluations, learning questions, or reviews missing")
	}
	first := run(t, "", "export", "mel", exTheory, "-O", exObs)
	if again := run(t, "", "export", "mel", exTheory, "-O", exObs); again.stdout != first.stdout {
		t.Error("two exports differ")
	}
}

func TestExportOGSM(t *testing.T) {
	out := filepath.Join(t.TempDir(), "ogsm.json")
	r := run(t, "", "export", "ogsm", exTheory, "-O", exObs, "-o", out)
	wantCode(t, r, ExitOK)
	plan, err := ogsm.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := ogsm.Errors(plan.Validate(ogsm.StrictValidationOptions())); len(errs) != 0 {
		t.Errorf("exported OGSM does not meet the strict options: %+v", errs)
	}
	if plan.Objective.Statement == "" || len(plan.Goals) == 0 || len(plan.Strategies) == 0 || len(plan.Measures) == 0 || len(plan.Actions) == 0 {
		t.Errorf("objective, goals, strategies, measures, or actions missing")
	}
	first := run(t, "", "export", "ogsm", exTheory, "-O", exObs)
	if again := run(t, "", "export", "ogsm", exTheory, "-O", exObs); again.stdout != first.stdout {
		t.Error("two exports differ")
	}
}

func TestExportDMAIC(t *testing.T) {
	const catalog = "../../examples/course-catalog/impact-model.yaml"
	const catalogObs = "../../examples/course-catalog/observations-2026-h1.yaml"
	out := filepath.Join(t.TempDir(), "dmaic.json")
	r := run(t, "", "export", "dmaic", catalog, "-O", catalogObs, "-o", out)
	wantCode(t, r, ExitOK)
	doc, err := dmaic.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if errs := dmaic.Errors(doc.Validate(dmaic.DefaultValidationOptions())); len(errs) != 0 {
		t.Errorf("exported DMAIC does not validate: %+v", errs)
	}
	if len(doc.Categories) != 1 || len(doc.Categories[0].Metrics) == 0 {
		t.Fatalf("categories: %+v", doc.Categories)
	}
	// An indicator that cannot be a metric is named on standard error.
	if !strings.Contains(r.err, "note: ind-service-quality is left out") {
		t.Errorf("no omission note: %q", r.err)
	}
	first := run(t, "", "export", "dmaic", catalog, "-O", catalogObs)
	if again := run(t, "", "export", "dmaic", catalog, "-O", catalogObs); again.stdout != first.stdout {
		t.Error("two exports differ")
	}

	// A model with no pathway marked for DMAIC is refused.
	refused := run(t, "", "export", "dmaic", exTheory)
	wantCode(t, refused, ExitUsageErr)
	if !strings.Contains(refused.err, "projections: [dmaic]") {
		t.Errorf("refusal does not say how to mark a pathway: %s", refused.err)
	}
}
