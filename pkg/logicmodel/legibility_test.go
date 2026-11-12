package logicmodel_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grokify/d2vision"
	"github.com/grokify/d2vision/generate"
	"github.com/grokify/d2vision/legibility"

	"github.com/ProductBuildersHQ/outcome-model/pkg/logicmodel"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

// Each pathway's theory of change reads at 9 points or more on a portrait
// Word page once d2 draws it, as the document export does. The test needs the
// d2 tool, and is skipped without it.
func TestPathwayTheoryOfChangeReadsOnAPortraitPage(t *testing.T) {
	tool, err := exec.LookPath("d2")
	if err != nil {
		t.Skip("d2 is not installed")
	}
	doc, err := modelio.LoadDocument("../../examples/intermediary/outcome-model.yaml")
	if err != nil {
		t.Fatal(err)
	}
	g, err := service.BuildGraph(doc)
	if err != nil {
		t.Fatal(err)
	}
	views, err := report.PathwayViews(doc, g, nil)
	if err != nil {
		t.Fatal(err)
	}
	page, _ := legibility.TargetByName("docx-portrait")
	dir := t.TempDir()
	for _, v := range views {
		code := generate.NewGenerator().Generate(logicmodel.Build(v).TheoryOfChange().TheorySpec())
		in, out := filepath.Join(dir, v.Filter+".d2"), filepath.Join(dir, v.Filter+".svg")
		if err := os.WriteFile(in, []byte(code), 0o600); err != nil {
			t.Fatal(err)
		}
		if msg, err := exec.Command(tool, "--pad", "20", in, out).CombinedOutput(); err != nil {
			t.Fatalf("%s: d2 failed: %v\n%s", v.Filter, err, msg)
		}
		d, err := d2vision.ParseFile(out)
		if err != nil {
			t.Fatal(err)
		}
		r, err := legibility.Analyze(d, page)
		if err != nil {
			t.Fatal(err)
		}
		if !r.Readable {
			t.Errorf("%s: smallest label %.1f pt on a portrait page: %v", v.Filter, r.SmallestPt, r.Advice)
		}
	}
}
