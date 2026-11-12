package prismx

import (
	dmaicrender "github.com/grokify/prism-roadmap/goals/dmaic/render"
	dmaicmd "github.com/grokify/prism-roadmap/goals/dmaic/render/markdown"
	melrender "github.com/grokify/prism-roadmap/goals/mel/render"
	melmd "github.com/grokify/prism-roadmap/goals/mel/render/markdown"
	ogsmrender "github.com/grokify/prism-roadmap/goals/ogsm/render"
	ogsmmd "github.com/grokify/prism-roadmap/goals/ogsm/render/markdown"
	vsemrender "github.com/grokify/prism-roadmap/goals/vsem/render"
	vsemmd "github.com/grokify/prism-roadmap/goals/vsem/render/markdown"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// MarkdownOptions control how a projection is written as Markdown by
// prism-roadmap's renderers.
type MarkdownOptions struct {
	// Embed leaves out the projection's own title and metadata line and
	// moves its headings one level down, for a caller that places it under a
	// heading of its own.
	Embed bool
}

// headingOffset is how far an embedded projection's headings move down.
func (o MarkdownOptions) headingOffset() int {
	if o.Embed {
		return 1
	}
	return 0
}

// VSEMMarkdown draws the outcome model as a VSEM (see ToVSEM) and writes it as
// Markdown.
func VSEMMarkdown(doc *model.Model, st *model.Status, o MarkdownOptions) (string, error) {
	v, err := ToVSEM(doc, st)
	if err != nil {
		return "", err
	}
	opts := vsemrender.DefaultOptions()
	opts.OmitHeader, opts.HeadingOffset = o.Embed, o.headingOffset()
	return text(vsemmd.New().Render(v, opts))
}

// OGSMMarkdown draws the outcome model as an OGSM (see ToOGSM) and writes it
// as Markdown.
func OGSMMarkdown(doc *model.Model, st *model.Status, o MarkdownOptions) (string, error) {
	g, err := ToOGSM(doc, st)
	if err != nil {
		return "", err
	}
	opts := ogsmrender.DefaultOptions()
	opts.OmitHeader, opts.HeadingOffset = o.Embed, o.headingOffset()
	return text(ogsmmd.New().Render(g, opts))
}

// MELMarkdown draws the outcome model as a MEL plan (see ToMEL) and writes it
// as Markdown, with the indicator reference sheets.
func MELMarkdown(doc *model.Model, st *model.Status, o MarkdownOptions) (string, error) {
	p, err := ToMEL(doc, st)
	if err != nil {
		return "", err
	}
	opts := melrender.DefaultOptions()
	opts.OmitHeader, opts.HeadingOffset = o.Embed, o.headingOffset()
	return text(melmd.New().Render(p, opts))
}

// DMAICMarkdown draws the outcome model's DMAIC pathways as a DMAIC document
// (see ToDMAIC) and writes it as Markdown. The omissions are the indicators
// left out, as ToDMAIC reports them.
func DMAICMarkdown(doc *model.Model, st *model.Status, sets []*model.ObservationSet, o MarkdownOptions) (string, []Omission, error) {
	d, omissions, err := ToDMAIC(doc, st, sets)
	if err != nil {
		return "", omissions, err
	}
	opts := dmaicrender.DefaultOptions()
	opts.OmitHeader, opts.HeadingOffset = o.Embed, o.headingOffset()
	s, err := text(dmaicmd.New().Render(d, opts))
	return s, omissions, err
}

func text(b []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return string(b), nil
}
