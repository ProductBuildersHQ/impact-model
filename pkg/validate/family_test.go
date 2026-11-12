package validate

import (
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

func TestPathwayIssues(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(d *model.Model)
		code   string
		path   string
	}{
		{"parent is an element ref", func(d *model.Model) { d.Parent = "model:umbrella#oc-x" }, "VAL-PTH-001", "/parent"},
		{"parent is malformed", func(d *model.Model) { d.Parent = "umbrella" }, "VAL-PTH-001", "/parent"},
		{"parent is this document", func(d *model.Model) { d.Parent = "model:" + d.ID }, "VAL-PTH-001", "/parent"},
		{"childModel malformed", func(d *model.Model) { d.Pathways[0].ChildModel = "model:Bad" }, "VAL-PTH-001", "/pathways/0/childModel"},
		{"childModel is an element", func(d *model.Model) { d.Pathways[0].ChildModel = "model:child#oc-1" }, "VAL-PTH-001", "/pathways/0/childModel"},
		{"contributesTo malformed", func(d *model.Model) { d.Outcomes[0].ContributesTo = []string{"nope"} }, "VAL-PTH-002", "/outcomes/0/contributesTo/0"},
		{"contributesTo a whole document", func(d *model.Model) { d.Outcomes[0].ContributesTo = []string{"model:agent-ecosystem-trust"} }, "VAL-PTH-002", "/outcomes/0/contributesTo/0"},
		{"contributesTo this document", func(d *model.Model) { d.Outcomes[0].ContributesTo = []string{"model:" + d.ID + "#oc-idp-support"} }, "VAL-PTH-002", "/outcomes/0/contributesTo/0"},
		{"contributesTo outside the parent", func(d *model.Model) { d.Outcomes[0].ContributesTo = []string{"model:somewhere-else#oc-x"} }, "VAL-PTH-003", "/outcomes/0/contributesTo/0"},
		{"pathway without title", func(d *model.Model) { d.Pathways[0].Title = "" }, "VAL-PTH-005", "/pathways/0/title"},
		{"pathway outcome unknown", func(d *model.Model) { d.Pathways[0].OutcomeIDs[0] = "oc-ghost" }, "VAL-REF-001", "/pathways/0/outcomeIds/0"},
		{"pathway outcome is an actor", func(d *model.Model) { d.Pathways[0].OutcomeIDs[0] = "act-idp" }, "VAL-REF-002", "/pathways/0/outcomeIds/0"},
		{"pathway intervention is an outcome", func(d *model.Model) { d.Pathways[0].InterventionIDs[0] = "oc-idp-support" }, "VAL-REF-002", "/pathways/0/interventionIds/0"},
		{"pathway owner unknown", func(d *model.Model) { d.Pathways[0].OwnerID = "act-ghost" }, "VAL-REF-001", "/pathways/0/ownerActorId"},
		{"pathway resource without a name", func(d *model.Model) {
			d.Pathways[0].Resources = []model.Resource{{Name: "Engineering"}, {Allocation: "2 people"}}
		}, "VAL-PTH-008", "/pathways/0/resources/1/name"},
		{"pathway kind invalid", func(d *model.Model) { d.Pathways[0].Kind = "region" }, "VAL-ENUM-001", "/pathways/0/kind"},
		{"duplicate pathway id", func(d *model.Model) { d.Pathways[0].ID = "oc-idp-support" }, "VAL-ID-002", "/pathways/0/id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := fullDoc(t)
			tt.mutate(doc)
			for _, i := range Document(doc) {
				if i.Code == tt.code && i.Path == tt.path {
					return
				}
			}
			t.Fatalf("expected %s at %s; got:%s", tt.code, tt.path, describe(Document(doc)))
		})
	}
}

// umbrella is a parent outcome model the full fixture elaborates.
func umbrella() *model.Model {
	d := model.NewModel("agent-ecosystem-trust", "Agents cannot yet be trusted across organizations.")
	d.Outcomes = []model.Outcome{
		{ID: "oc-trusted-agents", Title: "Agents are trusted across organizations", Level: model.LevelImpact},
	}
	d.Pathways = []model.Pathway{
		{ID: "pw-authz", Title: "Authorization", ChildModel: "model:open-agent-auth"},
	}
	return d
}

func TestFamilyResolvesAgainstParent(t *testing.T) {
	if issues := Family(umbrella(), fullDoc(t)); len(issues) != 0 {
		t.Fatalf("matching family has issues:%s", describe(issues))
	}
}

func TestFamilyIssues(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(parent, child *model.Model)
		code   string
		sev    Severity
		source string
	}{
		{"child names no parent", func(_, c *model.Model) { c.Parent = "" }, "VAL-FAM-001", SeverityError, "child[open-agent-auth]"},
		{"child names another parent", func(_, c *model.Model) { c.Parent = "model:someone-else" }, "VAL-FAM-001", SeverityError, "child[open-agent-auth]"},
		{"child parent malformed", func(_, c *model.Model) { c.Parent = "oops" }, "VAL-FAM-001", SeverityError, "child[open-agent-auth]"},
		{"outcome missing in parent", func(p, _ *model.Model) { p.Outcomes[0].ID = "oc-renamed" }, "VAL-FAM-002", SeverityError, "child[open-agent-auth]"},
		{"reference is not an outcome", func(_, c *model.Model) {
			c.Outcomes[0].ContributesTo = []string{"model:agent-ecosystem-trust#pw-authz"}
		}, "VAL-FAM-003", SeverityError, "child[open-agent-auth]"},
		{"no parent pathway names the child", func(p, _ *model.Model) { p.Pathways[0].ChildModel = "" }, "VAL-FAM-004", SeverityWarning, "child[open-agent-auth]"},
		{"parent pathway names an unsupplied child", func(p, _ *model.Model) {
			p.Pathways = append(p.Pathways, model.Pathway{ID: "pw-other", Title: "Other", ChildModel: "model:not-supplied"})
		}, "VAL-FAM-005", SeverityWarning, "parent"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, child := umbrella(), fullDoc(t)
			tt.mutate(parent, child)
			issues := Family(parent, child)
			for _, i := range issues {
				if i.Code == tt.code {
					if i.Severity != tt.sev || i.Source != tt.source {
						t.Fatalf("%s = severity %s source %q, want %s %q", tt.code, i.Severity, i.Source, tt.sev, tt.source)
					}
					return
				}
			}
			t.Fatalf("expected %s; got:%s", tt.code, describe(issues))
		})
	}

	child := fullDoc(t)
	if issues := Family(umbrella(), child, child); len(issues) == 0 || issues[len(issues)-1].Code != "VAL-FAM-006" {
		t.Fatalf("duplicate child not reported:%s", describe(issues))
	}
}
