package logicmodel

import (
	"slices"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// view is a small outcome model: two activities share a resource, one lists an
// output, and the outcomes run from output level to impact.
func view() *report.View {
	return &report.View{
		Document: report.DocumentView{Title: "Whole outcome model"},
		Interventions: []report.InterventionView{
			{ID: "int-a", Title: "Run workshops", Resources: "Two staff", Outputs: []string{"Workshops held"}, Targets: []string{"oc-skills"}},
			{ID: "int-b", Title: "Coach leaders", Resources: "Two staff", Targets: []string{"oc-reach"}},
		},
		Outcomes: []report.OutcomeView{
			{ID: "oc-reach", Title: "Members reached", Level: model.LevelOutput},
			{ID: "oc-skills", Title: "Skills gained", Level: model.LevelShortTerm},
			{ID: "oc-practice", Title: "Practice applied", Level: model.LevelMediumTerm, Basis: report.BasisView{Basis: model.BasisInferred, Marked: true}},
			{ID: "oc-strong", Title: "Members are strong", Level: model.LevelLongTerm},
			{ID: "oc-served", Title: "Residents served", Level: model.LevelImpact},
		},
		Links: []report.LinkView{
			{ID: "l1", From: "oc-reach", To: "oc-skills"},
			{ID: "l2", From: "oc-skills", To: "oc-practice"},
			{ID: "l3", From: "oc-practice", To: "oc-strong"},
			{ID: "l4", From: "oc-strong", To: "oc-served"},
			{ID: "l5", From: "oc-strong", To: "oc-elsewhere"}, // an outcome of another pathway
			{ID: "l6", From: "oc-skills", To: "oc-skills"},    // a loop back to itself
		},
	}
}

func hasEdge(g *Graph, from, to string) bool { return slices.Contains(g.Edges, Edge{from, to}) }

func ids(ns []Node) []string {
	var out []string
	for _, n := range ns {
		out = append(out, n.ID)
	}
	return out
}

func TestStagesRunFromInputsToImpact(t *testing.T) {
	g := Build(view())
	var got []StageID
	for _, s := range g.Stages {
		got = append(got, s.Stage)
	}
	want := []StageID{StageInputs, StageActivities, StageOutputs, StageShortTerm, StageMediumTerm, StageLongTerm, StageImpact}
	if !slices.Equal(got, want) {
		t.Errorf("stages = %v, want %v", got, want)
	}
	if g.Title != "Whole outcome model" {
		t.Errorf("title = %q", g.Title)
	}
}

func TestEveryKindOfEdgeComesFromTheTheory(t *testing.T) {
	g := Build(view())
	for _, e := range []Edge{
		{"_input1", "int-a"}, {"_input1", "int-b"}, // one resource feeds two activities
		{"int-a", "_output1"},     // an activity leads to the output it lists
		{"_output1", "oc-skills"}, // which leads to the outcome the activity targets
		{"int-b", "oc-reach"},     // an activity leads to an output-level outcome it targets
		{"oc-reach", "oc-skills"}, // links join outcomes
		{"oc-skills", "oc-practice"}, {"oc-practice", "oc-strong"}, {"oc-strong", "oc-served"},
	} {
		if !hasEdge(g, e.From, e.To) {
			t.Errorf("missing edge %s -> %s; edges are %v", e.From, e.To, g.Edges)
		}
	}
	if hasEdge(g, "int-a", "oc-skills") {
		t.Error("an activity that lists outputs should lead to its targets through them, not around them")
	}
}

// A link back to an earlier stage is feedback. It is kept out of the edges so
// it is not drawn across the boxes in between.
func TestFeedbackIsKeptApart(t *testing.T) {
	v := view()
	v.Links = append(v.Links, report.LinkView{ID: "loop", From: "oc-strong", To: "oc-skills", Feedback: true})
	g := Build(v)
	if hasEdge(g, "oc-strong", "oc-skills") {
		t.Error("a feedback link was drawn as an edge")
	}
	if !slices.Contains(g.Feedback, Edge{"oc-strong", "oc-skills"}) || len(g.Feedback) != 1 {
		t.Errorf("feedback = %v", g.Feedback)
	}
	note := feedbackNote(g)
	if !strings.Contains(note, "Members are strong leads back to Skills gained") {
		t.Errorf("note:\n%s", note)
	}
	if feedbackNote(Build(view())) != "" {
		t.Error("An outcome model with no feedback gets no note")
	}
}

func TestEdgesToNowhereAndToSelfAreDropped(t *testing.T) {
	g := Build(view())
	for _, e := range g.Edges {
		if e.From == e.To || strings.Contains(e.To, "elsewhere") || strings.Contains(e.From, "elsewhere") {
			t.Errorf("edge %v should not be drawn", e)
		}
	}
	seen := map[Edge]bool{}
	for _, e := range g.Edges {
		if seen[e] {
			t.Errorf("edge %v is repeated", e)
		}
		seen[e] = true
	}
}

func TestSharedTextIsOneNode(t *testing.T) {
	g := Build(view())
	var inputs []Node
	for _, s := range g.Stages {
		if s.Stage == StageInputs {
			inputs = s.Nodes
		}
	}
	if len(inputs) != 1 || inputs[0].Label != "Two staff" {
		t.Errorf("two activities that state the same resource should share one input: %v", inputs)
	}
}

func TestMarkedElementsSayWhatTheyAre(t *testing.T) {
	g := Build(view())
	for _, s := range g.Stages {
		for _, n := range s.Nodes {
			if n.ID == "oc-practice" && n.Label != "Practice applied (inferred)" {
				t.Errorf("label = %q", n.Label)
			}
		}
	}
}

func TestStagesWithNothingInThemAreLeftOut(t *testing.T) {
	v := &report.View{Outcomes: []report.OutcomeView{{ID: "a", Title: "A", Level: model.LevelShortTerm}, {ID: "b", Title: "B", Level: model.LevelImpact}}}
	g := Build(v)
	if len(g.Stages) != 2 || g.Stages[0].Stage != StageShortTerm || g.Stages[1].Stage != StageImpact {
		t.Errorf("stages = %+v", g.Stages)
	}
}

// Nodes in a later stage are ordered by where the nodes that lead to them sit,
// so arrows do not cross when they need not.
func TestNodesAreOrderedToAvoidCrossings(t *testing.T) {
	v := &report.View{Outcomes: []report.OutcomeView{
		{ID: "s1", Title: "S1", Level: model.LevelShortTerm},
		{ID: "s2", Title: "S2", Level: model.LevelShortTerm},
		// Written in the opposite order to the nodes that lead to them.
		{ID: "m2", Title: "M2", Level: model.LevelMediumTerm},
		{ID: "m1", Title: "M1", Level: model.LevelMediumTerm},
	}, Links: []report.LinkView{{From: "s1", To: "m1"}, {From: "s2", To: "m2"}}}
	g := Build(v)
	var mediums []string
	for _, s := range g.Stages {
		if s.Stage == StageMediumTerm {
			mediums = ids(s.Nodes)
		}
	}
	if !slices.Equal(mediums, []string{"m1", "m2"}) {
		t.Errorf("order = %v, want [m1 m2]", mediums)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	a, b := Build(view()), Build(view())
	if len(a.Edges) != len(b.Edges) || !slices.Equal(a.Edges, b.Edges) {
		t.Error("the same view gave different edges")
	}
	for i := range a.Stages {
		if !slices.Equal(ids(a.Stages[i].Nodes), ids(b.Stages[i].Nodes)) {
			t.Error("the same view gave different node order")
		}
	}
}

func TestTheoryOfChangeJoinsActivitiesToOutcomes(t *testing.T) {
	full := Build(view())
	toc := full.TheoryOfChange()
	for _, st := range toc.Stages {
		if st.Stage == StageInputs || st.Stage == StageOutputs {
			t.Errorf("theory of change keeps the %s stage", st.Stage)
		}
	}
	kept := map[string]bool{}
	for _, st := range toc.Stages {
		for _, n := range st.Nodes {
			kept[n.ID] = true
		}
	}
	for _, e := range append(slices.Clone(toc.Edges), toc.Feedback...) {
		if !kept[e.From] || !kept[e.To] {
			t.Errorf("edge %s -> %s touches a removed node", e.From, e.To)
		}
	}
	// int-a produced "Workshops held", which led on to oc-skills: the
	// theory of change joins the activity to that outcome directly.
	if !slices.Contains(toc.Edges, Edge{From: "int-a", To: "oc-skills"}) {
		t.Errorf("int-a is not joined to oc-skills through its output: %v", toc.Edges)
	}
	seen := map[Edge]bool{}
	for _, e := range toc.Edges {
		if seen[e] {
			t.Errorf("edge %v is drawn twice", e)
		}
		seen[e] = true
	}
	if len(full.Stages) == len(toc.Stages) {
		t.Error("TheoryOfChange changed the graph it was called on, or removed nothing")
	}
}

func TestTheoryOfChangeDiagramIsOneBlock(t *testing.T) {
	text, err := TheoryOfChangeDiagram(view(), FormatD2)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(text, "```d2"); n != 1 {
		t.Errorf("want one D2 block, got %d:\n%s", n, text)
	}
	if !strings.HasPrefix(text, "## Theory of change on one page\n") || strings.Contains(text, "inputs:") {
		t.Errorf("unexpected diagram:\n%s", text)
	}
	if _, err := TheoryOfChangeDiagram(view(), Format("svg")); err == nil {
		t.Error("an unknown format is accepted")
	}
}
