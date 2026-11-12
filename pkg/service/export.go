package service

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
	"github.com/ProductBuildersHQ/outcome-model/pkg/prismx"
)

// PRISMExportKinds are the documents an impact model exports for the PRISM tools.
var PRISMExportKinds = []string{"okr", "journey", "v2mom", "vsem", "dmaic", "mel", "ogsm", "bsc"}

// PRISMExport is one document exported for the PRISM tools, and what was left
// out of it.
type PRISMExport struct {
	// Kind is the export kind, one of PRISMExportKinds.
	Kind string `json:"kind"`
	// Document is the exported document, ready to write as JSON.
	Document any `json:"document"`
	// Omissions are the indicators the export left out, and why; only the
	// DMAIC export reports them.
	Omissions []prismx.Omission `json:"omissions,omitempty"`
}

// ExportPRISM draws the impact model as the export kind. On error the
// result carries no document, but may still list omissions. st and sets may be
// nil, for an export without observations. It is the one place the command
// line and the MCP server export from, so both refuse the same incomplete
// models with the same errors (prismx.ErrNotV2MOMReady, ErrNotVSEMReady,
// ErrNotDMAICReady, ErrNotMELReady, ErrNotOGSMReady, ErrNotBSCReady).
func ExportPRISM(kind string, doc *impact.Model, st *impact.Status, sets []*impact.ObservationSet) (PRISMExport, error) {
	out := PRISMExport{Kind: kind}
	var err error
	switch kind {
	case "okr":
		out.Document, err = prismx.ToOKR(doc, st)
	case "journey":
		out.Document, err = prismx.ToOutcomeJourneys(doc, st)
	case "v2mom":
		out.Document, err = prismx.ToV2MOM(doc, st)
	case "vsem":
		out.Document, err = prismx.ToVSEM(doc, st)
	case "dmaic":
		out.Document, out.Omissions, err = prismx.ToDMAIC(doc, st, sets)
	case "mel":
		out.Document, err = prismx.ToMEL(doc, st)
	case "ogsm":
		out.Document, err = prismx.ToOGSM(doc, st)
	case "bsc":
		out.Document, err = prismx.ToBSC(doc, st)
	default:
		return PRISMExport{}, fmt.Errorf("unknown export %q (want one of %s)", kind, strings.Join(PRISMExportKinds, ", "))
	}
	if err != nil {
		// A refused DMAIC export still says which indicators it left out.
		return PRISMExport{Kind: kind, Omissions: out.Omissions}, err
	}
	return out, nil
}

// IsPRISMExport reports whether kind is one of PRISMExportKinds.
func IsPRISMExport(kind string) bool { return slices.Contains(PRISMExportKinds, kind) }
