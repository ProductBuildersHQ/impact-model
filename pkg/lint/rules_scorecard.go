package lint

import (
	"fmt"
	"strconv"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
	"github.com/ProductBuildersHQ/outcome-model/pkg/validate"
)

// The scorecard rules apply once an outcome names a Balanced Scorecard
// perspective. They check what a scorecard drawn from the model needs: a
// strategy map that runs upward, and objectives that are owned and measured
// by lead or lag indicators.
func init() {
	register("IMP-BSC-001", validate.SeverityWarning, "The strategy map runs upward: no scorecard objective leads to one in a lower perspective.", func(c *Context) []Finding {
		var out []Finding
		forObjectives(c.Doc, func(i int, o impact.Outcome) {
			for _, id := range c.Doc.ScorecardContributesTo(o.ID) {
				t, _ := c.Doc.Outcome(id)
				if t.ScorecardPerspective.Rank() > o.ScorecardPerspective.Rank() {
					out = append(out, Finding{Path: outcomePath(i) + "/scorecardPerspective", Element: o.ID,
						Message: fmt.Sprintf("leads to %s in the %s perspective, below its own %s perspective; a strategy map runs from learning and growth up to the results, so one of the two perspectives may be wrong",
							id, t.ScorecardPerspective, o.ScorecardPerspective)})
				}
			}
		})
		return out
	})

	register("IMP-BSC-002", validate.SeverityInfo, "Every scorecard objective names the actor who answers for it.", func(c *Context) []Finding {
		var out []Finding
		forObjectives(c.Doc, func(i int, o impact.Outcome) {
			if o.OwnerID == "" {
				out = append(out, Finding{Path: outcomePath(i), Element: o.ID,
					Message: "a scorecard objective with no ownerActorId; a Balanced Scorecard names who answers for each objective"})
			}
		})
		return out
	})

	register("IMP-BSC-003", validate.SeverityInfo, "Every scorecard objective has an indicator.", func(c *Context) []Finding {
		var out []Finding
		forObjectives(c.Doc, func(i int, o impact.Outcome) {
			for _, ind := range c.Doc.Indicators {
				if ind.Measures == o.ID && !ind.Rejected() {
					return
				}
			}
			out = append(out, Finding{Path: outcomePath(i), Element: o.ID,
				Message: "a scorecard objective that no indicator measures"})
		})
		return out
	})

	register("IMP-BSC-004", validate.SeverityInfo, "Every indicator of a scorecard objective says whether it is a lead or a lag measure.", func(c *Context) []Finding {
		objective := map[string]bool{}
		for _, o := range c.Doc.ScorecardObjectives() {
			objective[o.ID] = true
		}
		var out []Finding
		for i, ind := range c.Doc.Indicators {
			if objective[ind.Measures] && !ind.Rejected() && ind.LeadLag == "" {
				out = append(out, Finding{Path: "/indicators/" + strconv.Itoa(i) + "/leadLag", Element: ind.ID,
					Message: "an indicator of a scorecard objective that is neither lead nor lag; a balanced scorecard mixes lead measures, which move first, with lag results"})
			}
		}
		return out
	})
}

// forObjectives calls fn for each scorecard objective with its index among
// the outcomes.
func forObjectives(d *impact.Model, fn func(i int, o impact.Outcome)) {
	for i, o := range d.Outcomes {
		if o.ScorecardPerspective != "" && !o.Rejected() {
			fn(i, o)
		}
	}
}

func outcomePath(i int) string { return "/outcomes/" + strconv.Itoa(i) }
