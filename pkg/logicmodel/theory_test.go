package logicmodel

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// wideView has five activities, one with a long title, under one outcome.
func wideView() *report.View {
	v := &report.View{
		Document: report.DocumentView{Title: "Network"},
		Filter:   "pw-net",
		Pathways: []report.PathwayView{{ID: "pw-net", Title: "Network and Evidence"}},
		Outcomes: []report.OutcomeView{{ID: "oc-served", Title: "Residents served", Level: model.LevelImpact}},
	}
	for _, id := range []string{"a", "b", "c", "d"} {
		v.Interventions = append(v.Interventions, report.InterventionView{ID: "int-" + id, Title: "Activity " + id, Targets: []string{"oc-served"}})
	}
	v.Interventions = append(v.Interventions, report.InterventionView{
		ID: "int-long", Title: "Deliver capacity building trainings, workshops, and cohort programs for nonprofit staff", Targets: []string{"oc-served"},
	})
	return v
}

// The theory of change is drawn bottom-up: a row per stage, impact first, at
// most three boxes side by side, with labels wrapped and enlarged.
func TestTheorySpecIsBottomUp(t *testing.T) {
	g := Build(wideView()).TheoryOfChange()
	spec := g.TheorySpec()
	if spec.GridRows != len(g.Stages) || spec.GridColumns != 0 || len(spec.Containers) != len(g.Stages) {
		t.Fatalf("rows %d, columns %d, containers %d for %d stages", spec.GridRows, spec.GridColumns, len(spec.Containers), len(g.Stages))
	}
	if spec.Containers[0].ID != string(StageImpact) || spec.Containers[len(spec.Containers)-1].ID != string(StageActivities) {
		t.Errorf("rows run %s to %s, want impact to activities", spec.Containers[0].ID, spec.Containers[len(spec.Containers)-1].ID)
	}
	acts := spec.Containers[len(spec.Containers)-1]
	if len(acts.Nodes) != 5 || acts.GridColumns != theoryPerRow {
		t.Errorf("activities: %d nodes in %d columns, want 5 in %d", len(acts.Nodes), acts.GridColumns, theoryPerRow)
	}
	if spec.Containers[0].GridColumns != 1 {
		t.Errorf("one impact outcome takes %d columns", spec.Containers[0].GridColumns)
	}
	var long string
	for _, n := range acts.Nodes {
		if n.ID == "int-long" {
			long = n.Label
		}
	}
	for _, line := range strings.Split(long, "\n") {
		if len(line) > theoryWrap {
			t.Errorf("line %q is longer than %d", line, theoryWrap)
		}
	}
	if strings.Count(long, "\n") < 2 {
		t.Errorf("long label not wrapped: %q", long)
	}
	if f := spec.Classes[string(StageActivities)].FontSize; f == nil || *f != theoryFont {
		t.Errorf("activity font = %v, want %d", f, theoryFont)
	}
	if len(spec.Edges) != len(g.Edges) {
		t.Errorf("%d edges for %d graph edges", len(spec.Edges), len(g.Edges))
	}
}

// A pathway's theory of change is titled with the pathway, and the logic
// model keeps its columns.
func TestTheoryOfChangeDiagramNamesItsPathway(t *testing.T) {
	text, err := TheoryOfChangeDiagram(wideView(), FormatD2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(text, "## Theory of change: Network and Evidence\n") || !strings.Contains(text, "grid-rows: 2") {
		t.Errorf("unexpected diagram:\n%s", text)
	}
	if !strings.Contains(text, `"Deliver capacity building\ntrainings, workshops, and\ncohort programs for\nnonprofit staff"`) {
		t.Errorf("label not wrapped in the D2:\n%s", text)
	}
	if lm := Build(wideView()).D2(); strings.Contains(lm, "grid-rows") || !strings.Contains(lm, "grid-columns") {
		t.Errorf("the logic model lost its columns:\n%s", lm)
	}
	mm, err := TheoryOfChangeDiagram(wideView(), FormatMermaid)
	if err != nil || !strings.Contains(mm, "```mermaid") || !strings.Contains(mm, "<br/>") {
		t.Errorf("mermaid (%v):\n%s", err, mm)
	}
}

// The whole outcome model comes first, then each pathway.
func TestTheoryOfChangeDiagrams(t *testing.T) {
	text, err := TheoryOfChangeDiagrams(view(), []*report.View{wideView()}, FormatD2)
	if err != nil {
		t.Fatal(err)
	}
	first := strings.Index(text, "## Theory of change on one page")
	second := strings.Index(text, "## Theory of change: Network and Evidence")
	if first != 0 || second < 0 || strings.Count(text, "```d2") != 2 {
		t.Errorf("unexpected diagrams:\n%s", text)
	}
	if _, err := TheoryOfChangeDiagrams(view(), nil, Format("svg")); err == nil {
		t.Error("an unknown format is accepted")
	}
}
