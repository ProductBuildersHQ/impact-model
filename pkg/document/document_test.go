package document_test

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/document"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

var update = flag.Bool("update", false, "rewrite golden files")

const example = "../../examples/intermediary/outcome-model.yaml"

func load(t *testing.T) *model.Model {
	t.Helper()
	doc, err := modelio.LoadDocument(example)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func build(t *testing.T, opts document.Options) *document.Document {
	t.Helper()
	d, err := service.Document(load(t), nil, nil, opts)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func golden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with -update to create it)", path, err)
	}
	if string(want) != got {
		t.Errorf("%s differs from the golden file (run with -update to accept)", name)
	}
}

func TestNonprofitDocumentGolden(t *testing.T) {
	d := build(t, document.Options{Stamp: "Example stamp"})
	golden(t, "intermediary-nonprofit.md", d.Markdown(document.MarkdownOptions{Contents: true}))
}

func TestSectionsFollowTheGivenOrder(t *testing.T) {
	order := []document.SectionKind{document.SectionMEL, document.SectionTheoryOfChange, document.SectionOGSM}
	d := build(t, document.Options{Sections: order})
	var got []document.SectionKind
	for _, s := range d.Sections {
		got = append(got, s.Kind)
	}
	if !slices.Equal(got, order) {
		t.Errorf("sections = %v, want %v", got, order)
	}
}

func TestEverySectionNestsUnderItsHeading(t *testing.T) {
	d := build(t, document.Options{Sections: document.Profiles["full"], SkipUnavailable: true})
	heading := regexp.MustCompile(`(?m)^(#{1,6}) `)
	for _, s := range d.Sections {
		for _, p := range s.Parts {
			for _, m := range heading.FindAllStringSubmatch(p.Markdown, -1) {
				if len(m[1]) < 3 {
					t.Errorf("section %s has a level-%d heading; content starts at level 3", s.Kind, len(m[1]))
				}
			}
		}
	}
}

func TestDiagramsAreSeparateAndUniquelyNamed(t *testing.T) {
	d := build(t, document.Options{Sections: document.Profiles["full"], SkipUnavailable: true})
	seen := map[string]bool{}
	for _, g := range d.Diagrams() {
		if seen[g.ID] {
			t.Errorf("diagram ID %s is used twice", g.ID)
		}
		seen[g.ID] = true
		if g.D2 == "" || g.Title == "" {
			t.Errorf("diagram %s has no source or no title", g.ID)
		}
	}
	for _, want := range []string{"theory-of-change-1", "logic-models-1", "problem-trees-1"} {
		if !seen[want] {
			t.Errorf("no diagram %s among %v", want, seen)
		}
	}
	md := d.Markdown(document.MarkdownOptions{})
	if strings.Count(md, "```d2") != len(d.Diagrams()) {
		t.Error("Markdown without image paths does not keep every diagram as D2")
	}
	md = d.Markdown(document.MarkdownOptions{
		DiagramPath: func(g *document.Diagram) string { return "img/" + g.ID + ".svg" },
		ImageAttrs:  "width=100%",
	})
	if strings.Contains(md, "```d2") || !strings.Contains(md, "(img/theory-of-change-1.svg){width=100%}") {
		t.Error("Markdown with image paths does not reference the rendered diagrams")
	}
}

func TestUnavailableSectionIsAnErrorUnlessSkipped(t *testing.T) {
	sections := []document.SectionKind{document.SectionTheoryOfChange, document.SectionDMAIC}
	_, err := service.Document(load(t), nil, nil, document.Options{Sections: sections})
	if !errors.Is(err, document.ErrSectionUnavailable) {
		t.Fatalf("want ErrSectionUnavailable, got %v", err)
	}
	d := build(t, document.Options{Sections: sections, SkipUnavailable: true})
	if len(d.Sections) != 1 || len(d.Skipped) != 1 || d.Skipped[0].Kind != document.SectionDMAIC || d.Skipped[0].Reason == "" {
		t.Errorf("sections %d, skipped %+v", len(d.Sections), d.Skipped)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	opts := document.Options{Sections: document.Profiles["full"], SkipUnavailable: true}
	a, b := build(t, opts), build(t, opts)
	if a.Markdown(document.MarkdownOptions{}) != b.Markdown(document.MarkdownOptions{}) {
		t.Error("two builds differ")
	}
}

func TestParseSectionsAndProfiles(t *testing.T) {
	got, err := document.ParseSections(" vsem, mel ,theory-of-change")
	if err != nil || !slices.Equal(got, []document.SectionKind{document.SectionVSEM, document.SectionMEL, document.SectionTheoryOfChange}) {
		t.Errorf("ParseSections = %v, %v", got, err)
	}
	for _, bad := range []string{"vsem,vsem", "vsem,nope", " , "} {
		if _, err := document.ParseSections(bad); err == nil {
			t.Errorf("ParseSections(%q) accepted", bad)
		}
	}
	for _, name := range document.ProfileNames() {
		s, err := document.ProfileSections(name)
		if err != nil || len(s) == 0 {
			t.Errorf("profile %s: %v", name, err)
		}
		for _, k := range s {
			if !k.Valid() {
				t.Errorf("profile %s lists unknown section %s", name, k)
			}
		}
	}
	if _, err := document.ProfileSections("nope"); err == nil {
		t.Error("an unknown profile is accepted")
	}
	if !slices.Equal(document.Profiles["full"], document.SectionKinds()) {
		t.Error("the full profile is not every section")
	}
}
