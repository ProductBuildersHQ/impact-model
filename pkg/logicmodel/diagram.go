package logicmodel

import (
	"strings"

	"github.com/grokify/d2vision/generate"
	"github.com/grokify/d2vision/legibility"
)

// The diagram is drawn bottom-up, a band per stage with the impact on top, so
// it suits a portrait page. ELK lays it out and routes every arrow around the
// boxes. Boxes are one width, so they line up in columns, and labels are
// wrapped to suit it. The sizes keep a pathway's labels readable on a portrait
// Word page, as checked with d2vision's legibility analysis.
const (
	stageClass = "stage"
	// LayoutEngine is the D2 layout engine the diagram is drawn with.
	LayoutEngine = "elk"
	// boxFont is the label size in pixels; a stage heading is larger.
	boxFont = 24
)

// boxSize is the width of every box, in pixels, and the line length its
// labels are wrapped to, in characters.
type boxSize struct{ width, wrap int }

// boxSizes are chosen by the most boxes in one stage, which sit side by side:
// a wide stage gets narrower boxes so the diagram stays narrow enough to read.
var boxSizes = []struct {
	upTo int
	size boxSize
}{
	{3, boxSize{width: 330, wrap: 26}},
	{4, boxSize{width: 290, wrap: 22}},
	{0, boxSize{width: 270, wrap: 20}}, // more than four
}

var stageStyle = map[StageID]generate.StyleSpec{
	StageInputs:     {Fill: "#fff1dc", Stroke: "#e08a1e"},
	StageActivities: {Fill: "#e4f3e5", Stroke: "#2e7d32"},
	StageOutputs:    {Fill: "#ece4fb", Stroke: "#7e57c2"},
	StageShortTerm:  {Fill: "#e1f3fd", Stroke: "#2d9cdb"},
	StageMediumTerm: {Fill: "#dfe8f6", Stroke: "#3f5fa8"},
	StageLongTerm:   {Fill: "#d3dcee", Stroke: "#2c4a8a"},
	StageImpact:     {Fill: "#c5cfe8", Stroke: "#1f3a7a"},
}

// stageBand is the light band behind each stage, so it is clear which stage a
// box belongs to when arrows within a stage stack its boxes.
var stageBand = generate.StyleSpec{
	Fill:        "#f4f6fa",
	Stroke:      "#dde3ec",
	StrokeWidth: generate.IntPtr(1),
	FontSize:    generate.IntPtr(boxFont + 4),
	Bold:        generate.BoolPtr(true),
}

// sizeFor picks the box size for the widest stage of the graph.
func (g *Graph) sizeFor() boxSize {
	widest := 0
	for _, st := range g.Stages {
		widest = max(widest, len(st.Nodes))
	}
	for _, s := range boxSizes {
		if s.upTo == 0 || widest <= s.upTo {
			return s.size
		}
	}
	return boxSizes[len(boxSizes)-1].size
}

// Spec describes the diagram for d2vision: a band for each stage, from the
// first stage at the bottom to impact on top, a box for each item, and the
// arrows between them, drawn with ELK. D2 and Mermaid are both written from
// it.
func (g *Graph) Spec() *generate.DiagramSpec {
	size := g.sizeFor()
	spec := &generate.DiagramSpec{
		LayoutEngine: LayoutEngine,
		Direction:    "up",
		Classes:      map[string]generate.StyleSpec{stageClass: stageBand},
	}
	for _, st := range g.Stages {
		style := stageStyle[st.Stage]
		style.BorderRadius = generate.IntPtr(12)
		style.FontSize = generate.IntPtr(boxFont)
		style.Width = generate.IntPtr(size.width)
		spec.Classes[string(st.Stage)] = style
	}
	for _, st := range g.Stages {
		c := generate.ContainerSpec{
			ID:        string(st.Stage),
			Label:     st.Stage.Name(),
			Class:     stageClass,
			Direction: "right",
		}
		for _, n := range st.Nodes {
			c.Nodes = append(c.Nodes, generate.NodeSpec{
				ID:    diagramID(n.ID),
				Label: legibility.Wrap(n.Label, size.wrap),
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

// D2 writes the diagram as D2, naming ELK as its layout engine.
func (g *Graph) D2() string {
	return generate.NewGenerator().Generate(g.Spec())
}

// Mermaid writes the diagram as a Mermaid flowchart, left to right with each
// stage's boxes stacked, which suits the landscape slides it is shown on.
// Mermaid places the boxes itself.
func (g *Graph) Mermaid() (string, error) {
	spec := g.Spec()
	spec.Direction = "right"
	for i := range spec.Containers {
		spec.Containers[i].Direction = "down"
	}
	return generate.NewGenerator().Mermaid(spec)
}

// diagramID makes an element ID usable as a node name: a dot would read as a
// step into a container.
func diagramID(id string) string { return strings.ReplaceAll(id, ".", "_") }
