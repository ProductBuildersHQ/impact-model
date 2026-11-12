package edit

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// FocusOn returns a patch that makes m the pathway's one metric that matters
// from m.From. The focus it replaces, the entry still open on that day, ends
// the day before, so the periods stay apart and the earlier focus stays on
// record. It is refused when the new period would overlap one that is
// already closed, or start on or before the day the open one started.
func FocusOn(doc *model.Model, pathwayID string, m model.OneMetric) (Patch, error) {
	idx := -1
	for i, p := range doc.Pathways {
		if p.ID == pathwayID {
			idx = i
		}
	}
	if idx < 0 {
		return nil, fmt.Errorf("%w: no pathway %q", ErrPatch, pathwayID)
	}
	from, err := time.Parse(time.DateOnly, m.From)
	if err != nil {
		return nil, fmt.Errorf("%w: from %q is not a date (want YYYY-MM-DD)", ErrPatch, m.From)
	}
	if m.IndicatorID == "" {
		return nil, fmt.Errorf("%w: name the indicator to focus on", ErrPatch)
	}
	pw := strconv.Itoa(idx)
	f := doc.Pathways[idx].Focus
	if f == nil {
		op, err := Add(Pointer("pathways", pw, "focus"), model.Focus{OneMetrics: []model.OneMetric{m}})
		return Patch{op}, err
	}

	var p Patch
	dayBefore := from.AddDate(0, 0, -1).Format(time.DateOnly)
	for i, old := range f.OneMetrics {
		switch {
		case old.To == "" && old.From < m.From:
			op, err := Add(Pointer("pathways", pw, "focus", "oneMetrics", strconv.Itoa(i), "to"), dayBefore)
			if err != nil {
				return nil, err
			}
			p = append(p, op)
		case old.To == "":
			return nil, fmt.Errorf("%w: the current one metric that matters, %q, started on %s; a new focus must start after it", ErrPatch, old.IndicatorID, old.From)
		case old.To >= m.From && (m.To == "" || old.From <= m.To):
			return nil, fmt.Errorf("%w: the period overlaps the one metric that matters %q, from %s to %s", ErrPatch, old.IndicatorID, old.From, old.To)
		}
	}
	var op Op
	if len(f.OneMetrics) == 0 {
		op, err = Add(Pointer("pathways", pw, "focus", "oneMetrics"), []model.OneMetric{m})
	} else {
		op, err = Add(Pointer("pathways", pw, "focus", "oneMetrics", "-"), m)
	}
	if err != nil {
		return nil, err
	}
	return append(p, op), nil
}
