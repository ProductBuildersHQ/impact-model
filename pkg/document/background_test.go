package document_test

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/document"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

func backgroundOf(t *testing.T, doc *model.Model, sections ...document.SectionKind) string {
	t.Helper()
	d, err := service.Document(doc, nil, nil, document.Options{Sections: sections, SkipUnavailable: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range d.Sections {
		if s.Kind == document.SectionBackground {
			var b strings.Builder
			for _, p := range s.Parts {
				b.WriteString(p.Markdown)
			}
			return b.String()
		}
	}
	t.Fatal("no background section")
	return ""
}

// The background introduces the organization, its vision, mission, and
// values, and its programs with their activities, all from the model.
func TestBackgroundIntroducesTheOrganization(t *testing.T) {
	text := backgroundOf(t, load(t), document.SectionBackground, document.SectionTheoryOfChange)
	for _, want := range []string{
		"### About Riverbend Nonprofit Network\n",
		"Riverbend Nonprofit Network began in 2015",
		"\n\nThe network runs leadership cohorts", // the second paragraph stays one
		"**Whom it serves.** About 150 member nonprofits",
		"**Period of this model.** January 2026 to December 2028.",
		"### Vision, mission, and values\n",
		"**Vision.** A county where every resident",
		"**Mission.** Strengthen the nonprofits",
		"- **Members first.** Every program answers",
		"### Programs\n",
		"Riverbend Nonprofit Network works through three programs.",
		"#### Learning Exchange\n",
		"**Activities:**",
		"### How this document serves the mission and vision\n",
		"reached through two impact outcomes:",
		"For this period the model sets one objective: Member nonprofits are strong",
		"answers for its results up to and including its long-term outcomes.",
		"| Theory of Change | How the programs' activities lead",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q:\n%s", want, text)
		}
	}
}

// The table lists only the sections that follow the background and are in
// the document, not ones skipped or placed before it.
func TestBackgroundListsTheSectionsThatFollow(t *testing.T) {
	text := backgroundOf(t, load(t),
		document.SectionGlossary, document.SectionBackground, document.SectionVSEM, document.SectionDMAIC)
	if !strings.Contains(text, "| VSEM: Vision, Strategy, Execution, Metrics |") {
		t.Errorf("VSEM not listed:\n%s", text)
	}
	for _, not := range []string{"| Glossary |", "| DMAIC Improvement Plan |"} {
		if strings.Contains(text, not) {
			t.Errorf("listed %q:\n%s", not, text)
		}
	}
}

// Without a background, the perspective actor's description introduces the
// organization; without values described, they are one sentence.
func TestBackgroundFallsBack(t *testing.T) {
	doc := load(t)
	doc.Context.Background = ""
	for i, a := range doc.Actors {
		if a.ID == doc.Perspective.ActorID {
			doc.Actors[i].Description = "A membership network of nonprofits"
		}
	}
	for i := range doc.Context.Values {
		doc.Context.Values[i].Description = ""
	}
	doc.Pathways[0].Kind = model.PathwayStrategy
	text := backgroundOf(t, doc, document.SectionBackground)
	for _, want := range []string{
		"Riverbend Nonprofit Network is a membership network of nonprofits.",
		"**Values.** Members first, Stronger together, and Candor about results.",
		"works through one strategy and two programs.",
		"(strategy)\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "| Section |") {
		t.Errorf("a background with nothing after it lists sections:\n%s", text)
	}

	doc.Perspective = nil
	if text := backgroundOf(t, doc, document.SectionBackground); !strings.Contains(text, "### About the organization\n") {
		t.Errorf("no perspective:\n%s", text)
	}
}

// Every profile opens with the background, and every section kind says what
// it shows.
func TestProfilesOpenWithTheBackground(t *testing.T) {
	for _, name := range document.ProfileNames() {
		s, err := document.ProfileSections(name)
		if err != nil || s[0] != document.SectionBackground {
			t.Errorf("profile %s starts with %v (%v)", name, s, err)
		}
	}
	for _, k := range document.SectionKinds() {
		if k.Purpose() == "" {
			t.Errorf("%s has no purpose", k)
		}
	}
}
