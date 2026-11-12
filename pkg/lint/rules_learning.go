package lint

import (
	"fmt"
	"slices"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/validate"
)

// takenUp reports whether a learning question or an evaluation that review
// has not rejected names the element in its related IDs.
func takenUp(d *model.Model, id string) bool {
	for _, q := range d.LearningQuestions {
		if !q.Rejected() && slices.Contains(q.RelatedIDs, id) {
			return true
		}
	}
	for _, e := range d.Evaluations {
		if !e.Rejected() && slices.Contains(e.RelatedIDs, id) {
			return true
		}
	}
	return false
}

func init() {
	register("OM-LRN-001", validate.SeverityWarning, "Once the impact model plans evaluations or learning questions, every high or critical assumption not yet shown to hold is taken up by one.", func(c *Context) []Finding {
		if len(c.Doc.LearningQuestions) == 0 && len(c.Doc.Evaluations) == 0 {
			return nil // no evaluation plan or learning agenda yet; OM-ASM-004 covers testing
		}
		var out []Finding
		for i, a := range c.Doc.Assumptions {
			if a.Provenance.IsRejected() || (a.Criticality != model.CriticalityHigh && a.Criticality != model.CriticalityCritical) {
				continue
			}
			switch a.Status {
			case "", model.AssumptionUntested, model.AssumptionAtRisk, model.AssumptionInvalidated:
			default:
				continue
			}
			if !takenUp(c.Doc, a.ID) {
				out = append(out, Finding{Path: at("assumptions", i), Element: a.ID,
					Message: fmt.Sprintf("a %s assumption that has not been shown to hold, but no learning question or evaluation takes it up", a.Criticality)})
			}
		}
		return out
	})

	register("OM-LRN-002", validate.SeverityInfo, "Every open learning question is taken up by a review meeting.", func(c *Context) []Finding {
		reviewed := map[string]bool{}
		for _, m := range c.Doc.ReviewMeetings {
			for _, id := range m.QuestionIDs {
				reviewed[id] = true
			}
		}
		var out []Finding
		for i, q := range c.Doc.LearningQuestions {
			if q.Rejected() || q.Status == model.QuestionAnswered || q.Status == model.QuestionRetired || reviewed[q.ID] {
				continue
			}
			out = append(out, Finding{Path: at("learningQuestions", i), Element: q.ID,
				Message: "no review meeting takes this question up, so nothing says when the team will look at the answer"})
		}
		return out
	})

	register("OM-LRN-003", validate.SeverityWarning, "Evaluation dates and learning question due dates fall within the time horizon of the impact model.", func(c *Context) []Finding {
		h := c.Doc.Context.TimeHorizon
		if h == nil {
			return nil
		}
		outside := func(date string) bool { return date != "" && (date < h.Start || date > h.End) }
		var out []Finding
		for i, e := range c.Doc.Evaluations {
			for _, d := range []struct{ field, date string }{{"startDate", e.StartDate}, {"endDate", e.EndDate}} {
				if !e.Rejected() && outside(d.date) {
					out = append(out, Finding{Path: at("evaluations", i) + "/" + d.field, Element: e.ID,
						Message: fmt.Sprintf("%s %s is outside the time horizon %s to %s", d.field, d.date, h.Start, h.End)})
				}
			}
		}
		for i, q := range c.Doc.LearningQuestions {
			if !q.Rejected() && outside(q.DueDate) {
				out = append(out, Finding{Path: at("learningQuestions", i) + "/dueDate", Element: q.ID,
					Message: fmt.Sprintf("due date %s is outside the time horizon %s to %s", q.DueDate, h.Start, h.End)})
			}
		}
		return out
	})
}
