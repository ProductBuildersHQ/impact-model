// Package service holds one function per user operation. The CLI and the
// MCP server are thin adapters over it, so both behave identically.
package service

import (
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/validate"
)

// Report is the validation result for one input.
type Report struct {
	// Source is the file path or another label for the input.
	Source string           `json:"source"`
	Format modelio.Format   `json:"format,omitempty"`
	Kind   model.Kind       `json:"kind,omitempty"`
	Issues []validate.Issue `json:"issues"`
}

// Valid reports whether the input has no error-severity issues.
func (r Report) Valid() bool { return !validate.HasErrors(r.Issues) }

// ValidateBytes validates one input of any document kind. A impact model of
// Change is checked against its schema and structurally; other kinds are
// checked against their schema. When against is given, an observation set
// is also checked against that impact model.
func ValidateBytes(source string, data []byte, f modelio.Format, against *model.Model, opts ...validate.Option) Report {
	r := Report{Source: source, Format: f, Issues: []validate.Issue{}}
	if kind, err := modelio.DetectKind(data, f); err == nil {
		r.Kind = kind
	}
	_, issues := validate.Source(data, f, opts...)
	r.Issues = append(r.Issues, issues...)
	if r.Kind == model.KindObservationSet && against != nil && !validate.HasErrors(issues) {
		set, err := modelio.ParseObservationSet(data, f)
		if err != nil {
			r.Issues = append(r.Issues, validate.Issue{Code: "VAL-SCH-002", Severity: validate.SeverityError, Message: err.Error()})
		} else {
			r.Issues = append(r.Issues, validate.Observations(against, set)...)
		}
	}
	return r
}

// ValidateFile reads and validates a file. The format comes from the
// extension.
func ValidateFile(path string, against *model.Model, opts ...validate.Option) (Report, error) {
	data, f, err := readInput(path)
	if err != nil {
		return Report{}, err
	}
	return ValidateBytes(path, data, f, against, opts...), nil
}

// ValidateFamily checks that child impact models resolve against their parent:
// each names the parent, and every outcome reference into the parent
// exists there. Only impact model files can be children.
func ValidateFamily(parentPath string, childPaths ...string) (Report, error) {
	parent, err := modelio.LoadDocument(parentPath)
	if err != nil {
		return Report{}, err
	}
	children := make([]*model.Model, 0, len(childPaths))
	for _, p := range childPaths {
		child, err := modelio.LoadDocument(p)
		if err != nil {
			return Report{}, err
		}
		children = append(children, child)
	}
	return Report{
		Source: parentPath,
		Kind:   model.KindOutcomeModel,
		Issues: append([]validate.Issue{}, validate.Family(parent, children...)...),
	}, nil
}

// ValidateFutures checks impact models that describe competing futures of one
// market against one another (see validate.Futures). Each must also be valid
// alone; use ValidateFile with the shared actor registry for that.
func ValidateFutures(paths ...string) (Report, error) {
	docs := make([]*model.Model, 0, len(paths))
	for _, p := range paths {
		d, err := modelio.LoadDocument(p)
		if err != nil {
			return Report{}, err
		}
		docs = append(docs, d)
	}
	return Report{Source: "competing futures", Kind: model.KindOutcomeModel, Issues: append([]validate.Issue{}, validate.Futures(docs...)...)}, nil
}

// LoadActorRegistry reads an actor registry file.
func LoadActorRegistry(path string) (*model.ActorRegistry, error) {
	data, f, err := readInput(path)
	if err != nil {
		return nil, err
	}
	return modelio.ParseActorRegistry(data, f)
}

func readInput(path string) ([]byte, modelio.Format, error) {
	f, err := modelio.FormatFromPath(path)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}
	if info.Size() > modelio.MaxFileBytes {
		return nil, "", fmt.Errorf("%w: %s is %d bytes (limit %d)", modelio.ErrTooLarge, path, info.Size(), modelio.MaxFileBytes)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return data, f, nil
}
