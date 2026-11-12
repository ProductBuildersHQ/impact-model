package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/prismx"
)

func TestExportPRISM(t *testing.T) {
	doc, err := modelio.LoadDocument("../../examples/course-catalog/outcome-model.yaml")
	if err != nil {
		t.Fatal(err)
	}
	// The scorecard needs an outcome in a perspective, which this example
	// does not place; the other exports ignore it.
	doc.Outcomes[0].ScorecardPerspective = model.ScorecardMission
	for _, kind := range PRISMExportKinds {
		out, err := ExportPRISM(kind, doc, nil, nil)
		if err != nil || out.Kind != kind || out.Document == nil {
			t.Errorf("%s: %+v %v", kind, out, err)
		}
	}
	if _, err := ExportPRISM("mystery", doc, nil, nil); err == nil || !strings.Contains(err.Error(), "okr, journey, v2mom, vsem, dmaic, mel, ogsm, bsc") {
		t.Errorf("unknown kind: %v", err)
	}

	// A refused DMAIC export carries no document but still lists what it left out.
	noNumbers := *doc
	noNumbers.Indicators = append([]model.Indicator(nil), doc.Indicators...)
	for i := range noNumbers.Indicators {
		noNumbers.Indicators[i].Baseline = nil
	}
	out, err := ExportPRISM("dmaic", &noNumbers, nil, nil)
	if !errors.Is(err, prismx.ErrNotDMAICReady) || out.Document != nil || len(out.Omissions) == 0 {
		t.Errorf("refused DMAIC export: %+v %v", out, err)
	}
}
