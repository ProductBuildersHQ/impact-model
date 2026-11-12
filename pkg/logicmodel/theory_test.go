package logicmodel

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// wideView has n activities, one with a long title, under one outcome.
func wideView(n int) *report.View {
	v := &report.View{
		Document: report.DocumentView{Title: "Network"},
		Filter:   "pw-net",
		Pathways: []report.PathwayView{{ID: "pw-net", Title: "Network and Evidence"}},
		Outcomes: []report.OutcomeView{{ID: "oc-served", Title: "Residents served", Level: model.LevelImpact}},
	}
	v.Interventions = append(v.Interventions, report.InterventionView{
		ID: "int-long", Title: "Deliver capacity building trainings, workshops, and cohort programs for nonprofit staff", Targets: []string{"oc-served"},
	})
	for i := 1; i < n; i++ {
		id := string(rune('a' + i))
		v.Interventions = append(v.Interventions, report.InterventionView{ID: "int-" + id, Title: "Activity " + id, Targets: []string{"oc-served"}})
	}
	return v
}

// Every box is one width, and labels are wrapped to suit it; a wider stage
// gets narrower boxes.
func TestBoxesShareAWidthChosenByTheWidestStage(t *testing.T) {
	for _, c := range []struct {
		activities int
		want       boxSize
	}{
		{2, boxSize{330, 26}},
		{3, boxSize{330, 26}},
		{4, boxSize{290, 22}},
		{9, boxSize{270, 20}},
	} {
		g := Build(wideView(c.activities)).TheoryOfChange()
		if got := g.sizeFor(); got != c.want {
			t.Errorf("%d activities: size %+v, want %+v", c.activities, got, c.want)
		}
		spec := g.Spec()
		for _, st := range g.Stages {
			cls := spec.Classes[string(st.Stage)]
			if cls.Width == nil || *cls.Width != c.want.width || cls.FontSize == nil || *cls.FontSize != boxFont {
				t.Errorf("%d activities: class %s = %+v", c.activities, st.Stage, cls)
			}
		}
		for _, ct := range spec.Containers {
			for _, n := range ct.Nodes {
				for _, line := range strings.Split(n.Label, "\n") {
					if len(line) > c.want.wrap {
						t.Errorf("%d activities: line %q is longer than %d", c.activities, line, c.want.wrap)
					}
				}
			}
		}
	}
}

// A pathway's theory of change is titled with the pathway and drawn with ELK,
// bottom-up, with wrapped labels.
func TestTheoryOfChangeDiagramNamesItsPathway(t *testing.T) {
	text, err := TheoryOfChangeDiagram(wideView(2), FormatD2)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"## Theory of change: Network and Evidence\n",
		"layout-engine: elk",
		"direction: up",
		`"Deliver capacity building\ntrainings, workshops, and\ncohort programs for\nnonprofit staff"`,
	} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q:\n%s", want, text)
		}
	}
	mm, err := TheoryOfChangeDiagram(wideView(2), FormatMermaid)
	if err != nil || !strings.Contains(mm, "```mermaid") || !strings.Contains(mm, "<br/>") {
		t.Errorf("mermaid (%v):\n%s", err, mm)
	}
}

// The whole outcome model comes first, then each pathway.
func TestTheoryOfChangeDiagrams(t *testing.T) {
	text, err := TheoryOfChangeDiagrams(view(), []*report.View{wideView(2)}, FormatD2)
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
