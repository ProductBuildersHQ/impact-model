package document

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/outcome-model/internal/mdx"
	"github.com/ProductBuildersHQ/outcome-model/pkg/logicmodel"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/prismx"
	"github.com/ProductBuildersHQ/outcome-model/pkg/problemtree"
	"github.com/ProductBuildersHQ/outcome-model/pkg/report"
)

// sectionLevel is the heading level of a section; its content starts one
// level below.
const sectionLevel = 2

// Input is what a document is built from.
type Input struct {
	// Model is the outcome model.
	Model *model.Model
	// Whole is the view of the whole outcome model, and Pathways one view per
	// pathway (report.BuildView and report.PathwayViews).
	Whole    *report.View
	Pathways []*report.View
	// Status and Observations are optional. With them, the strategy and MEL
	// sections carry current values.
	Status       *model.Status
	Observations []*model.ObservationSet
}

// Options choose what the document holds.
type Options struct {
	// Sections are the sections to include, in order. Empty means the
	// default profile.
	Sections []SectionKind
	// Subtitle replaces the default subtitle.
	Subtitle string
	// Stamp is shown under the title when not empty.
	Stamp string
	// SkipUnavailable leaves out, and lists in Document.Skipped, a strategy
	// section the outcome model cannot be drawn as, such as an OGSM for a
	// model with no objective. Without it, such a section is an error.
	SkipUnavailable bool
}

// ErrSectionUnavailable marks a section the outcome model cannot be drawn as.
// The wrapped error says what the model needs.
var ErrSectionUnavailable = errors.New("section unavailable")

// Build assembles the document.
func Build(in Input, opts Options) (*Document, error) {
	if in.Model == nil || in.Whole == nil {
		return nil, errors.New("document: no outcome model")
	}
	sections := opts.Sections
	if len(sections) == 0 {
		sections = Profiles[DefaultProfile]
	}
	doc := &Document{
		Title:    oneLine(in.Model.Metadata.Name),
		Subtitle: opts.Subtitle,
		Stamp:    opts.Stamp,
	}
	if doc.Title == "" {
		doc.Title = oneLine(in.Whole.Document.Title)
	}
	if doc.Subtitle == "" {
		doc.Subtitle = "Outcome model: " + titles(sections)
	}
	seen := map[SectionKind]bool{}
	backgroundAt := -1
	for _, k := range sections {
		if !k.Valid() {
			return nil, fmt.Errorf("document: unknown section %q", k)
		}
		if seen[k] {
			return nil, fmt.Errorf("document: section %q is listed twice", k)
		}
		seen[k] = true
		if k == SectionBackground {
			// The background names the sections that follow it, so it is
			// written once they are known.
			backgroundAt = len(doc.Sections)
			doc.Sections = append(doc.Sections, Section{})
			continue
		}
		body, err := sectionBody(in, k)
		if errors.Is(err, ErrSectionUnavailable) && opts.SkipUnavailable {
			doc.Skipped = append(doc.Skipped, Skipped{Kind: k, Reason: reason(err)})
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("document: section %s: %w", k, err)
		}
		doc.Sections = append(doc.Sections, section(k, body))
	}
	if backgroundAt >= 0 {
		var following []SectionKind
		for _, s := range doc.Sections[backgroundAt+1:] {
			following = append(following, s.Kind)
		}
		doc.Sections[backgroundAt] = section(SectionBackground, background(in, following))
	}
	return doc, nil
}

// section turns a section's Markdown into parts, its headings starting one
// level below the section's own, its D2 blocks made diagrams.
func section(k SectionKind, body string) Section {
	s := Section{Kind: k, Title: k.Title()}
	body = mdx.Rebase(body, sectionLevel+1)
	n := 0
	for _, p := range mdx.SplitFences(body, "d2") {
		if p.Lang == "" {
			s.Parts = append(s.Parts, Part{Markdown: p.Text})
			continue
		}
		n++
		title := p.Title
		if title == "" {
			title = s.Title
		}
		s.Parts = append(s.Parts, Part{Diagram: &Diagram{
			ID:    fmt.Sprintf("%s-%d", k, n),
			Title: title,
			D2:    p.Source,
		}})
	}
	return s
}

// sectionBody renders the Markdown of one section, with D2 diagrams as fenced
// blocks.
func sectionBody(in Input, k SectionKind) (string, error) {
	switch k {
	case SectionTheoryOfChange:
		diagram, err := logicmodel.TheoryOfChangeDiagrams(in.Whole, in.Pathways, logicmodel.FormatD2)
		if err != nil {
			return "", err
		}
		return diagram + "\n" + mdx.StripTitle(report.Narrative(in.Whole, in.Model, in.Status)), nil
	case SectionProblemTrees:
		return problemTrees(in.Model)
	case SectionLogicModels:
		return logicModels(in)
	case SectionOutcomeMap:
		return report.OutcomeMaps(in.Whole, in.Pathways), nil
	case SectionLogframe:
		return report.Logframes(in.Whole, in.Pathways), nil
	case SectionMetrics:
		return mdx.StripTitle(report.MetricDefinitions(in.Whole)), nil
	case SectionGlossary:
		return mdx.StripTitle(report.Glossary(in.Whole)), nil
	case SectionVSEM:
		text, err := prismx.VSEMMarkdown(in.Model, in.Status, embed)
		return text, unavailable(err, prismx.ErrNotVSEMReady)
	case SectionOGSM:
		text, err := prismx.OGSMMarkdown(in.Model, in.Status, embed)
		return text, unavailable(err, prismx.ErrNotOGSMReady)
	case SectionMEL:
		text, err := prismx.MELMarkdown(in.Model, in.Status, embed)
		return text, unavailable(err, prismx.ErrNotMELReady)
	case SectionDMAIC:
		text, omissions, err := prismx.DMAICMarkdown(in.Model, in.Status, in.Observations, embed)
		if err != nil {
			return "", unavailable(err, prismx.ErrNotDMAICReady)
		}
		if len(omissions) > 0 {
			var b strings.Builder
			b.WriteString(text)
			b.WriteString("\nLeft out of the DMAIC plan:\n\n")
			for _, o := range omissions {
				fmt.Fprintf(&b, "- `%s`: %s\n", o.IndicatorID, oneLine(o.Reason))
			}
			text = b.String()
		}
		return text, nil
	}
	return "", fmt.Errorf("unknown section %q", k)
}

// problemTrees renders the whole model's problem tree and each pathway's own.
func problemTrees(doc *model.Model) (string, error) {
	ids, err := report.ProblemTreeIDs(doc, "")
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "No causes or effects have been recorded for the core problem.\n", nil
	}
	var b strings.Builder
	for _, id := range ids {
		text, err := problemtree.ReportFor(doc, id, problemtree.FormatD2)
		if err != nil {
			return "", err
		}
		b.WriteString(text)
	}
	return b.String(), nil
}

// logicModels renders, for each pathway (or the whole model when there are
// none), the layered diagram and then the logic model table, under one heading.
func logicModels(in Input) (string, error) {
	views := in.Pathways
	if len(views) == 0 {
		views = []*report.View{in.Whole}
	}
	var b strings.Builder
	for _, v := range views {
		diagram, err := logicmodel.Diagrams(v, nil, logicmodel.FormatD2)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "## %s\n\n", oneLine(report.PathwayTitle(v)))
		b.WriteString(mdx.StripTitle(diagram))
		b.WriteString("\n")
		b.WriteString(mdx.StripTitle(report.LogicModels(v, nil)))
		b.WriteString("\n")
	}
	return b.String(), nil
}

// embed places a strategy projection under the section's heading.
var embed = prismx.MarkdownOptions{Embed: true}

// unavailable marks err as ErrSectionUnavailable when it is the readiness
// error of the projection. A nil err stays nil.
func unavailable(err, ready error) error {
	if err != nil && errors.Is(err, ready) {
		return fmt.Errorf("%w: %w", ErrSectionUnavailable, err)
	}
	return err
}

// reason is the message of an unavailable section's error, without the
// marker.
func reason(err error) string {
	return strings.TrimPrefix(err.Error(), ErrSectionUnavailable.Error()+": ")
}

func titles(sections []SectionKind) string {
	short := make([]string, len(sections))
	for i, k := range sections {
		t := k.Title()
		if i := strings.Index(t, ":"); i > 0 {
			t = t[:i]
		}
		short[i] = t
	}
	return strings.Join(short, ", ")
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }
