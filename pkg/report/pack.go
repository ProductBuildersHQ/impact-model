package report

import (
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/problemtree"
)

// DefinitionKinds lists the definition report kinds that need no
// observations, in the order the definitions pack presents them.
var DefinitionKinds = []string{"pathway", "logic-model", "logic-model-diagram", "logic-model-layered", "logic-model-d2", "problem-tree", "problem-tree-d2", "problem-tree-json", "dmaic-plan", "outcome-map", "logframe", "metrics", "metric-selection", "collection-plan", "glossary", "assessment", "definitions"}

// PathwayDiagrams renders the causal-pathway diagram of the whole outcome model,
// then one per pathway, each as a titled Mermaid block.
func PathwayDiagrams(whole *View, pathways []*View) string {
	var b strings.Builder
	b.WriteString("## Causal pathways\n\n")
	diagram := func(title string, v *View) {
		fmt.Fprintf(&b, "### %s\n\n```mermaid\n%s```\n\n", mdText(title), Mermaid(v, MermaidOptions{}))
	}
	diagram("The theory of change on one page", whole)
	for _, v := range pathways {
		diagram(PathwayTitle(v), v)
	}
	return b.String()
}

// Definitions renders the definition reports as one document, the pack a
// workshop reviews before any data exists: the problem trees (see
// problemTrees), then the causal pathways, logic models, outcome maps,
// logframes, metric definitions, and glossary. whole and pathways are views of
// doc.
func Definitions(whole *View, pathways []*View, doc *model.Model) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: definitions\n\n", mdText(whole.Document.Title))
	trees, err := problemTrees(doc, whole.Filter)
	if err != nil {
		return "", err
	}
	for _, id := range trees {
		tree, err := problemtree.ReportFor(doc, id, problemtree.FormatMermaid)
		if err != nil {
			return "", err
		}
		b.WriteString(tree)
	}
	b.WriteString(PathwayDiagrams(whole, pathways))
	b.WriteString(LogicModels(whole, pathways))
	b.WriteString(OutcomeMaps(whole, pathways))
	b.WriteString(Logframes(whole, pathways))
	b.WriteString(MetricDefinitions(whole))
	b.WriteString(Glossary(whole))
	return b.String(), nil
}

// problemTrees lists the problem trees the definitions pack draws, by pathway
// ID, "" being the whole model's. For the whole model: its tree, then the tree
// of each pathway that has causes or effects of its own. For one pathway: its
// own tree, or the whole model's when it has none. Empty trees are left out.
func problemTrees(doc *model.Model, pathway string) ([]string, error) {
	candidates := []string{pathway}
	if pathway == "" {
		for _, p := range doc.Pathways {
			if !p.Rejected() {
				candidates = append(candidates, p.ID)
			}
		}
	}
	var out []string
	for _, id := range candidates {
		t, err := problemtree.BuildFor(doc, id)
		if err != nil {
			return nil, err
		}
		if !t.Empty() {
			out = append(out, id)
		}
	}
	if pathway != "" && len(out) == 0 && !problemtree.Build(doc).Empty() {
		out = append(out, "")
	}
	return out, nil
}
