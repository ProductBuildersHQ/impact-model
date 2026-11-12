// Package document assembles one human-readable document from an outcome
// model: a chosen list of sections, in a chosen order, each a projection of
// the same model, such as the theory of change, the logic models, a VSEM, and
// a MEL plan. The document is format-neutral: Markdown text with D2 diagrams
// kept apart, so it can be written as Markdown, converted to Word, or shown as
// an HTML presentation. The same model and options always give the same
// document.
package document

import (
	"fmt"
	"sort"
	"strings"
)

// SectionKind names one section a document can include.
type SectionKind string

// The sections a document can include.
const (
	SectionBackground     SectionKind = "background"
	SectionTheoryOfChange SectionKind = "theory-of-change"
	SectionProblemTrees   SectionKind = "problem-trees"
	SectionLogicModels    SectionKind = "logic-models"
	SectionOutcomeMap     SectionKind = "outcome-map"
	SectionLogframe       SectionKind = "logframe"
	SectionMetrics        SectionKind = "metrics"
	SectionVSEM           SectionKind = "vsem"
	SectionOGSM           SectionKind = "ogsm"
	SectionMEL            SectionKind = "mel"
	SectionDMAIC          SectionKind = "dmaic"
	SectionGlossary       SectionKind = "glossary"
)

// SectionKinds returns every section kind, in the order the full profile
// uses.
func SectionKinds() []SectionKind {
	return []SectionKind{
		SectionBackground, SectionTheoryOfChange, SectionProblemTrees, SectionLogicModels, SectionOutcomeMap,
		SectionLogframe, SectionMetrics, SectionVSEM, SectionOGSM, SectionMEL, SectionDMAIC,
		SectionGlossary,
	}
}

// Valid reports whether k is a known section kind.
func (k SectionKind) Valid() bool {
	for _, s := range SectionKinds() {
		if s == k {
			return true
		}
	}
	return false
}

// Title is the section's heading in the document.
func (k SectionKind) Title() string {
	switch k {
	case SectionBackground:
		return "Background"
	case SectionTheoryOfChange:
		return "Theory of Change"
	case SectionProblemTrees:
		return "Problem Trees"
	case SectionLogicModels:
		return "Logic Models"
	case SectionOutcomeMap:
		return "Outcome Map"
	case SectionLogframe:
		return "Logframe"
	case SectionMetrics:
		return "Indicators"
	case SectionVSEM:
		return "VSEM: Vision, Strategy, Execution, Metrics"
	case SectionOGSM:
		return "OGSM: Objective, Goals, Strategies, Measures"
	case SectionMEL:
		return "Monitoring, Evaluation, and Learning Plan"
	case SectionDMAIC:
		return "DMAIC Improvement Plan"
	case SectionGlossary:
		return "Glossary"
	}
	return string(k)
}

// Purpose says what a section shows a reader, in terms of the organization's
// mission and vision. The background section lists it for the sections that
// follow.
func (k SectionKind) Purpose() string {
	switch k {
	case SectionBackground:
		return "Who the organization is, its vision and mission, and its programs"
	case SectionTheoryOfChange:
		return "How the programs' activities lead, step by step, to the impact in the vision, for the whole organization and for each program, with the reasoning, assumptions, and risks"
	case SectionProblemTrees:
		return "The problem the mission responds to, its causes, and its effects"
	case SectionLogicModels:
		return "For each program, what goes in, what it does, what it produces, and the change it aims for"
	case SectionOutcomeMap:
		return "The outcomes and the preconditions each one depends on"
	case SectionLogframe:
		return "Each result with its indicators, how they are checked, and the assumptions it rests on"
	case SectionMetrics:
		return "How each indicator is defined, collected, and reported"
	case SectionVSEM:
		return "The vision, and the strategies, execution, and metrics that carry it out, on one page"
	case SectionOGSM:
		return "The objective, the goals that say it is reached, and the strategies and measures for each"
	case SectionMEL:
		return "How progress is monitored, evaluated, and learned from, and who hears about it"
	case SectionDMAIC:
		return "Which processes to measure and improve, and how"
	case SectionGlossary:
		return "The terms this document uses"
	}
	return ""
}

// Profiles are named lists of sections for common readers. A profile is only
// a starting point: any list of sections, in any order, can be given instead.
var Profiles = map[string][]SectionKind{
	// nonprofit is the document a nonprofit shares with its board, funders,
	// and staff: why change happens, what each program does, where the
	// organization is going, and how it will learn.
	"nonprofit": {SectionBackground, SectionTheoryOfChange, SectionLogicModels, SectionVSEM, SectionMEL, SectionGlossary},
	// strategy is a leadership plan: the theory, then the two strategy
	// formats.
	"strategy": {SectionBackground, SectionTheoryOfChange, SectionVSEM, SectionOGSM},
	// improvement is for process improvement work: the problem and its
	// causes, then the DMAIC plan.
	"improvement": {SectionBackground, SectionProblemTrees, SectionDMAIC},
	// evaluation is for evaluators and funders: the logic, the logframe, the
	// indicators, and the MEL plan.
	"evaluation": {SectionBackground, SectionTheoryOfChange, SectionLogicModels, SectionLogframe, SectionMetrics, SectionMEL},
	// full is every section.
	"full": SectionKinds(),
}

// DefaultProfile is the profile used when no sections are given.
const DefaultProfile = "nonprofit"

// ProfileNames returns the profile names, sorted.
func ProfileNames() []string {
	names := make([]string, 0, len(Profiles))
	for n := range Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// ParseSections reads a comma-separated list of section kinds, keeping its
// order. A section may appear only once.
func ParseSections(list string) ([]SectionKind, error) {
	var out []SectionKind
	seen := map[SectionKind]bool{}
	for _, f := range strings.Split(list, ",") {
		k := SectionKind(strings.TrimSpace(f))
		if k == "" {
			continue
		}
		if !k.Valid() {
			return nil, fmt.Errorf("unknown section %q (want one of %s)", k, sectionNames())
		}
		if seen[k] {
			return nil, fmt.Errorf("section %q is listed twice", k)
		}
		seen[k] = true
		out = append(out, k)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no sections listed (want some of %s)", sectionNames())
	}
	return out, nil
}

// ProfileSections returns the sections of a named profile.
func ProfileSections(name string) ([]SectionKind, error) {
	s, ok := Profiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown profile %q (want one of %s)", name, strings.Join(ProfileNames(), ", "))
	}
	return append([]SectionKind(nil), s...), nil
}

func sectionNames() string {
	var names []string
	for _, k := range SectionKinds() {
		names = append(names, string(k))
	}
	return strings.Join(names, ", ")
}

// Document is an outcome model assembled for people to read.
type Document struct {
	// Title is the outcome model's name.
	Title string
	// Subtitle says what the document is, such as the profile it follows.
	Subtitle string
	// Stamp, when not empty, records the date and source the document
	// describes.
	Stamp string
	// Sections are in the order they were asked for.
	Sections []Section
	// Skipped lists the sections that were asked for but could not be drawn
	// from this outcome model, and why.
	Skipped []Skipped
}

// Section is one part of the document.
type Section struct {
	Kind  SectionKind
	Title string
	// Parts are the section's text and diagrams, in order. Headings in the
	// text start at level 3, under the section's own level-2 heading.
	Parts []Part
}

// Part is Markdown text or a diagram; exactly one is set.
type Part struct {
	Markdown string
	Diagram  *Diagram
}

// Diagram is a diagram written in D2.
type Diagram struct {
	// ID is unique in the document, such as logic-models-2, and stable for
	// the same model and sections.
	ID string
	// Title is the heading the diagram appears under.
	Title string
	// D2 is the diagram source.
	D2 string
}

// Skipped is a section left out of the document.
type Skipped struct {
	Kind   SectionKind
	Reason string
}

// Diagrams returns every diagram in the document, in order.
func (d *Document) Diagrams() []*Diagram {
	var out []*Diagram
	for _, s := range d.Sections {
		for _, p := range s.Parts {
			if p.Diagram != nil {
				out = append(out, p.Diagram)
			}
		}
	}
	return out
}
