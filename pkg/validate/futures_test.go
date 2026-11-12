package validate

import (
	"os"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
)

func future(id, actor, futureRef string) *model.Model {
	d := model.NewModel(id, "A problem.")
	d.ActorRegistry = "shared-actors"
	d.Perspective = &model.Perspective{ActorID: actor, Role: model.PerspectiveProponent, FutureRef: futureRef}
	d.Outcomes = []model.Outcome{{ID: "oc-a", Title: "A", Level: model.LevelImpact}}
	return d
}

func codes(issues []Issue) []string {
	var out []string
	for _, i := range issues {
		out = append(out, i.Code)
	}
	return out
}

func TestFuturesAcceptsThreeFuturesSharingOneRegistry(t *testing.T) {
	issues := Futures(future("a", "act-1", "ufm:one"), future("b", "act-2", "ufm:two"), future("c", "act-3", "ufm:three"))
	if len(issues) != 0 {
		t.Errorf("issues: %v", issues)
	}
}

func TestFuturesFindsProblems(t *testing.T) {
	tests := []struct {
		name   string
		docs   func() []*model.Model
		code   string
		sev    Severity
		source string
	}{
		{"no perspective", func() []*model.Model {
			d := future("a", "x", "ufm:one")
			d.Perspective = nil
			return []*model.Model{d}
		}, "VAL-FUT-001", SeverityWarning, "future[a]"},
		{"no future", func() []*model.Model { return []*model.Model{future("a", "x", "")} }, "VAL-FUT-001", SeverityWarning, "future[a]"},
		{"badly formed future", func() []*model.Model { return []*model.Model{future("a", "x", "open-wins")} }, "VAL-FUT-002", SeverityError, "future[a]"},
		{"different registries", func() []*model.Model {
			b := future("b", "y", "ufm:two")
			b.ActorRegistry = "other-actors"
			return []*model.Model{future("a", "x", "ufm:one"), b}
		}, "VAL-FUT-003", SeverityError, "futures"},
		{"one names no registry", func() []*model.Model {
			b := future("b", "y", "ufm:two")
			b.ActorRegistry = ""
			return []*model.Model{future("a", "x", "ufm:one"), b}
		}, "VAL-FUT-003", SeverityWarning, "futures"},
		{"same actor, same future", func() []*model.Model {
			return []*model.Model{future("a", "x", "ufm:one"), future("b", "x", "ufm:one")}
		}, "VAL-FUT-004", SeverityError, "future[b]"},
		{"reference to a missing element", func() []*model.Model {
			a := future("a", "x", "ufm:one")
			a.Alignments = []model.Alignment{{ElementID: "oc-a", Target: model.ExternalRef{Framework: model.FrameworkCustom, Ref: "model:b#oc-ghost"}}}
			return []*model.Model{a, future("b", "y", "ufm:two")}
		}, "VAL-FUT-006", SeverityError, "future[a]"},
		{"malformed reference", func() []*model.Model {
			a := future("a", "x", "ufm:one")
			a.Alignments = []model.Alignment{{ElementID: "oc-a", Target: model.ExternalRef{Framework: model.FrameworkCustom, Ref: "model:Not A Slug#x"}}}
			return []*model.Model{a}
		}, "VAL-FUT-005", SeverityError, "future[a]"},
		{"different markets", func() []*model.Model {
			a, b := future("a", "x", "ufm:one"), future("b", "y", "ufm:two")
			a.Context.Market = &model.ExternalRef{Framework: model.FrameworkMarketSpec, Ref: "market:one"}
			b.Context.Market = &model.ExternalRef{Framework: model.FrameworkMarketSpec, Ref: "market:two"}
			return []*model.Model{a, b}
		}, "VAL-FUT-007", SeverityWarning, "future[b]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, i := range Futures(tt.docs()...) {
				if i.Code == tt.code && i.Severity == tt.sev && i.Source == tt.source {
					return
				}
			}
			t.Fatalf("want %s %s from %s; got %v", tt.sev, tt.code, tt.source, codes(Futures(tt.docs()...)))
		})
	}
}

func TestReferencesToTheoriesOutsideTheSetAreNotChecked(t *testing.T) {
	a := future("a", "x", "ufm:one")
	a.Alignments = []model.Alignment{{ElementID: "oc-a", Target: model.ExternalRef{Framework: model.FrameworkCustom, Ref: "model:elsewhere#oc-anything"}}}
	if issues := Futures(a); len(issues) != 0 {
		t.Errorf("issues: %v", issues)
	}
}

// The competing-futures example is three outcome models and one registry that
// validate with no change to the schema.
func TestCompetingFuturesExample(t *testing.T) {
	data, err := os.ReadFile("../../examples/competing-futures/actors.yaml")
	if err != nil {
		t.Fatal(err)
	}
	reg, err := modelio.ParseActorRegistry(data, modelio.FormatYAML)
	if err != nil {
		t.Fatal(err)
	}
	var docs []*model.Model
	for _, name := range []string{"future-open", "future-bundle", "future-extension"} {
		d, err := modelio.LoadDocument("../../examples/competing-futures/" + name + ".yaml")
		if err != nil {
			t.Fatal(err)
		}
		if issues := Document(d, WithActorRegistry(reg)); len(issues) != 0 {
			t.Errorf("%s: %v", name, issues)
		}
		// Without the registry the actors cannot be resolved, which is a warning, not an error.
		without := Document(d)
		if HasErrors(without) || len(without) == 0 {
			t.Errorf("%s without its registry: %v", name, codes(without))
		}
		docs = append(docs, d)
	}
	if issues := Futures(docs...); len(issues) != 0 {
		t.Errorf("the futures disagree: %v", issues)
	}
	// The three perspectives are three different actors of the one registry.
	actors := map[string]bool{}
	for _, d := range docs {
		actors[d.Perspective.ActorID] = true
		if _, ok := regActor(reg, d.Perspective.ActorID); !ok {
			t.Errorf("%s: perspective actor %q is not in the registry", d.ID, d.Perspective.ActorID)
		}
	}
	if len(actors) != 3 || !strings.HasPrefix(docs[0].Perspective.FutureRef, "ufm:") {
		t.Errorf("perspectives = %v", actors)
	}
}

func regActor(r *model.ActorRegistry, id string) (model.Actor, bool) {
	for _, a := range r.Actors {
		if a.ID == id {
			return a, true
		}
	}
	return model.Actor{}, false
}
