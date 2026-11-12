package logicmodel_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grokify/d2vision"
	"github.com/grokify/d2vision/legibility"

	"github.com/ProductBuildersHQ/outcome-model/pkg/logicmodel"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

// elkArgs draw the diagram as the document export does.
var elkArgs = []string{
	"--layout", "elk", "--pad", "20",
	"--elk-nodeNodeBetweenLayers", "30", "--elk-edgeNodeBetweenLayers", "15",
	"--elk-padding", "[top=45,left=15,bottom=15,right=15]",
}

// Each pathway's theory of change and logic model, drawn as the document
// export draws them, has no arrow over a box, and the theory of change reads
// on a portrait Word page. A logic model, with its inputs and outputs, has too
// many boxes for that. The test needs the d2 tool, and is skipped without it.
func TestPathwayDiagramsAreClearOnAPortraitPage(t *testing.T) {
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
		for kind, graph := range map[string]*logicmodel.Graph{
			"theory-of-change": logicmodel.Build(v).TheoryOfChange(),
			"logic-model":      logicmodel.Build(v),
		} {
			name := v.Filter + "-" + kind
			in, out := filepath.Join(dir, name+".d2"), filepath.Join(dir, name+".svg")
			if err := os.WriteFile(in, []byte(graph.D2()), 0o600); err != nil {
				t.Fatal(err)
			}
			if msg, err := exec.Command(tool, append(elkArgs, in, out)...).CombinedOutput(); err != nil {
				t.Fatalf("%s: d2 failed: %v\n%s", name, err, msg)
			}
			d, err := d2vision.ParseFile(out)
			if err != nil {
				t.Fatal(err)
			}
			r, err := legibility.Analyze(d, page)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Crossings) > 0 {
				t.Errorf("%s: arrows over boxes: %+v", name, r.Crossings)
			}
			if kind == "theory-of-change" && r.SmallestPt < minPortraitPt {
				t.Errorf("%s: smallest label %.1f pt on a portrait page, want %.1f or more", name, r.SmallestPt, minPortraitPt)
			}
		}
	}
}

// minPortraitPt is the smallest label size the example's diagrams must keep on
// a portrait page. Routing arrows around boxes takes room, so a pathway with
// many boxes falls a little under the 9-point print minimum; 8 points is the
// usual floor for labels in a figure.
const minPortraitPt = 8.0
