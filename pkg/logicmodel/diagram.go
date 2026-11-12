package logicmodel

import (
	"strings"

	"github.com/grokify/d2vision/generate"
	"github.com/grokify/d2vision/legibility"
)

// The look of the diagram. A stage heading has no box of its own, and each
// stage's items share a color.
const (
	stageClass = "stage"
	// columnGap leaves room between columns for the arrows.
	columnGap = 90
	// itemGap spaces the items in a column.
	itemGap = 24
)

var stageStyle = map[StageID]generate.StyleSpec{
	StageInputs:     {Fill: "#fff1dc", Stroke: "#e08a1e"},
	StageActivities: {Fill: "#e4f3e5", Stroke: "#2e7d32"},
	StageOutputs:    {Fill: "#ece4fb", Stroke: "#7e57c2"},
	StageShortTerm:  {Fill: "#e1f3fd", Stroke: "#2d9cdb"},
	StageMediumTerm: {Fill: "#dfe8f6", Stroke: "#3f5fa8"},
	StageLongTerm:   {Fill: "#d3dcee", Stroke: "#2c4a8a"},
	StageImpact:     {Fill: "#c5cfe8", Stroke: "#1f3a7a"},
}

// Spec describes the diagram for d2vision: one container for each stage, laid
// out as a row of columns, a node for each item, and the edges between them.
// D2 and Mermaid are both written from it.
func (g *Graph) Spec() *generate.DiagramSpec {
	round := generate.IntPtr(12)
	spec := &generate.DiagramSpec{
		Direction:     "right",
		GridColumns:   len(g.Stages),
		HorizontalGap: columnGap,
		Classes: map[string]generate.StyleSpec{
			stageClass: {Fill: "transparent", StrokeWidth: generate.IntPtr(0), FontSize: generate.IntPtr(22), Bold: generate.BoolPtr(true)},
		},
	}
	for _, st := range g.Stages {
		style := stageStyle[st.Stage]
		style.BorderRadius = round
		spec.Classes[string(st.Stage)] = style
	}
	for _, st := range g.Stages {
		c := generate.ContainerSpec{
			ID:          string(st.Stage),
			Label:       st.Stage.Name(),
			Class:       stageClass,
			GridColumns: 1,
			VerticalGap: itemGap,
		}
		for _, n := range st.Nodes {
			c.Nodes = append(c.Nodes, generate.NodeSpec{ID: diagramID(n.ID), Label: n.Label, Class: string(st.Stage)})
		}
		spec.Containers = append(spec.Containers, c)
	}
	spec.Edges = g.edgeSpecs()
	return spec
}

// The theory of change is drawn bottom-up, a row per stage with the impact on
// top, so it suits a portrait page. Labels are wrapped and a long stage takes
// more than one row of boxes, which keeps the diagram narrow enough for its
// labels to read at 9 points or more on a portrait Word page, as checked with
// d2vision's legibility analysis.
const (
	// theoryWrap is the most characters on a line of a label.
	theoryWrap = 26
	// theoryFont is the label size in pixels; a stage heading is larger.
	theoryFont = 24
	// theoryPerRow is the most boxes side by side in a stage.
	theoryPerRow = 3
	// theoryGap spaces the rows, and the boxes within a stage.
	theoryGap = 24
)

// TheorySpec describes the theory of change for d2vision: a row for each
// stage, from activities at the bottom to impact at the top, with wrapped
// labels. D2 and Mermaid are both written from it.
func (g *Graph) TheorySpec() *generate.DiagramSpec {
	round := generate.IntPtr(12)
	spec := &generate.DiagramSpec{
		GridRows:    len(g.Stages),
		VerticalGap: theoryGap,
		Classes: map[string]generate.StyleSpec{
			stageClass: {Fill: "transparent", StrokeWidth: generate.IntPtr(0), FontSize: generate.IntPtr(theoryFont + 4), Bold: generate.BoolPtr(true)},
		},
	}
	for _, st := range g.Stages {
		style := stageStyle[st.Stage]
		style.BorderRadius = round
		style.FontSize = generate.IntPtr(theoryFont)
		spec.Classes[string(st.Stage)] = style
	}
	for i := len(g.Stages) - 1; i >= 0; i-- {
		st := g.Stages[i]
		c := generate.ContainerSpec{
			ID:            string(st.Stage),
			Label:         st.Stage.Name(),
			Class:         stageClass,
			GridColumns:   min(len(st.Nodes), theoryPerRow),
			HorizontalGap: theoryGap,
			VerticalGap:   theoryGap,
		}
		for _, n := range st.Nodes {
			c.Nodes = append(c.Nodes, generate.NodeSpec{
				ID:    diagramID(n.ID),
				Label: legibility.Wrap(n.Label, theoryWrap),
				Class: string(st.Stage),
			})
		}
		spec.Containers = append(spec.Containers, c)
	}
	spec.Edges = g.edgeSpecs()
	return spec
}

// edgeSpecs joins each edge's nodes through their stage containers.
func (g *Graph) edgeSpecs() []generate.EdgeSpec {
	stageOf := map[string]StageID{}
	for _, st := range g.Stages {
		for _, n := range st.Nodes {
			stageOf[n.ID] = st.Stage
		}
	}
	var out []generate.EdgeSpec
	for _, e := range g.Edges {
		out = append(out, generate.EdgeSpec{
			From: string(stageOf[e.From]) + "." + diagramID(e.From),
			To:   string(stageOf[e.To]) + "." + diagramID(e.To),
		})
	}
	return out
}

// D2 writes the diagram as D2. Laid out with a grid, the stages stay in their
// columns, in order.
func (g *Graph) D2() string {
	return generate.NewGenerator().Generate(g.Spec())
}

// Mermaid writes the diagram as a Mermaid flowchart. Mermaid has no grid, so it
// places the columns itself.
func (g *Graph) Mermaid() (string, error) {
	return generate.NewGenerator().Mermaid(g.Spec())
}

// diagramID makes an element ID usable as a node name: a dot would read as a
// step into a container.
func diagramID(id string) string { return strings.ReplaceAll(id, ".", "_") }
