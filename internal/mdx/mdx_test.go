package mdx

import (
	"reflect"
	"testing"
)

const sample = "# Title\n\nIntro.\n\n## Part\n\n```text\n# not a heading\n```\n\n### Detail\n"

func TestStripTitle(t *testing.T) {
	got := StripTitle(sample)
	want := "Intro.\n\n## Part\n\n```text\n# not a heading\n```\n\n### Detail\n"
	if got != want {
		t.Errorf("StripTitle:\n got %q\nwant %q", got, want)
	}
	if body := "Just text.\n"; StripTitle(body) != body {
		t.Errorf("StripTitle changed text that has no title: %q", StripTitle(body))
	}
}

func TestShiftKeepsFencesAndCaps(t *testing.T) {
	got := Shift(sample, 2)
	want := "### Title\n\nIntro.\n\n#### Part\n\n```text\n# not a heading\n```\n\n##### Detail\n"
	if got != want {
		t.Errorf("Shift +2:\n got %q\nwant %q", got, want)
	}
	if got := Shift("##### Deep\n", 3); got != "###### Deep\n" {
		t.Errorf("Shift past 6 = %q", got)
	}
	if got := Shift("## Up\n", -5); got != "# Up\n" {
		t.Errorf("Shift below 1 = %q", got)
	}
	if Shift(sample, 0) != sample {
		t.Error("Shift 0 changed the text")
	}
	if got := Shift("#hashtag and #### not\n", 1); got != "#hashtag and #### not\n" {
		t.Errorf("Shift touched a non-heading: %q", got)
	}
}

func TestRebase(t *testing.T) {
	if got := Rebase("## A\n\n#### B\n", 3); got != "### A\n\n##### B\n" {
		t.Errorf("Rebase = %q", got)
	}
	if got := Rebase("No headings.\n", 3); got != "No headings.\n" {
		t.Errorf("Rebase without headings = %q", got)
	}
}

func TestSplitFences(t *testing.T) {
	md := "## Logic model: A\n\n```d2\na -> b\n```\n\nNote.\n\n```mermaid\nx\n```\n"
	got := SplitFences(md, "d2")
	want := []Part{
		{Text: "## Logic model: A\n"},
		{Lang: "d2", Source: "a -> b\n", Title: "Logic model: A"},
		{Text: "Note.\n\n```mermaid\nx\n```\n"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SplitFences:\n got %#v\nwant %#v", got, want)
	}
}
