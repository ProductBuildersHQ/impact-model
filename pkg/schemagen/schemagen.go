// Package schemagen generates the JSON Schemas for the impact model
// document kinds from the Go types in pkg/model.
//
// Go structs are the source of truth. Enumeration types are emitted as
// string enums whose values come from the pkg/model All… functions, so the
// schema and the code cannot drift.
package schemagen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/invopop/jsonschema"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// BaseURL is the prefix of every schema $id.
const BaseURL = "https://github.com/ProductBuildersHQ/outcome-model/schema/v0/"

// Definition describes one generated schema.
type Definition struct {
	// Name is the file stem: <Name>.schema.json.
	Name        string
	Title       string
	Description string
	Value       any
}

// Definitions returns the schemas this module publishes, in a stable order.
func Definitions() []Definition {
	return []Definition{
		{"impact-model", "Impact Model", "An impact model document: context, outcomes, causal links, interventions, assumptions, indicators, actors, and evidence.", &model.Model{}},
		{"observation-set", "Impact Model Observation Set", "Indicator observations, assumption assessments, and notes for one impact model and period.", &model.ObservationSet{}},
		{"impact-model-status", "Impact Model Status", "A computed status snapshot of indicators, outcomes, pathways, links, and assumptions at a date.", &model.Status{}},
		{"status-diff", "Impact Model Status Diff", "Changes between two status snapshots.", &model.StatusDiff{}},
		{"actor-registry", "Impact Model Actor Registry", "A shared list of actors referenced by several impact models.", &model.ActorRegistry{}},
		{"csv-mapping", "Impact Model CSV Mapping", "How the columns of a spreadsheet map onto indicator observations.", &model.CSVMapping{}},
		{"review-response", "Impact Model Review Response", "What reviewers said about an impact model, by element, independent of the channel it was collected through.", &model.ReviewResponse{}},
	}
}

// ID returns the $id URL for a schema name.
func ID(name string) string {
	return BaseURL + name + ".schema.json"
}

// FileName returns the file name for a schema name.
func FileName(name string) string {
	return name + ".schema.json"
}

// File is one generated schema file.
type File struct {
	Name string
	Data []byte
}

// Generate returns every schema as indented JSON with a trailing newline,
// in Definitions order. Output is deterministic.
func Generate() ([]File, error) {
	defs := Definitions()
	files := make([]File, 0, len(defs))
	for _, d := range defs {
		data, err := GenerateOne(d)
		if err != nil {
			return nil, err
		}
		files = append(files, File{Name: FileName(d.Name), Data: data})
	}
	return files, nil
}

// GenerateOne returns one schema as indented JSON with a trailing newline.
func GenerateOne(d Definition) ([]byte, error) {
	r := &jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,
		Mapper:                     enumMapper,
	}
	s := r.Reflect(d.Value)
	if s == nil {
		return nil, fmt.Errorf("schemagen: reflect %s: nil schema", d.Name)
	}
	s.ID = jsonschema.ID(ID(d.Name))
	s.Title = d.Title
	s.Description = d.Description

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return nil, fmt.Errorf("schemagen: encode %s: %w", d.Name, err)
	}
	return buf.Bytes(), nil
}

// EnumValues maps each pkg/model enumeration type to its values.
func EnumValues() map[reflect.Type][]string {
	m := map[reflect.Type][]string{}
	add := func(t reflect.Type, values []string) { m[t] = values }
	add(typeOf[model.Level](), strs(model.AllLevels()))
	add(typeOf[model.Sphere](), strs(model.AllSpheres()))
	add(typeOf[model.ScorecardPerspective](), strs(model.AllScorecardPerspectives()))
	add(typeOf[model.LeadLag](), strs(model.AllLeadLags()))
	add(typeOf[model.LogicType](), strs(model.AllLogicTypes()))
	add(typeOf[model.ChangeType](), strs(model.AllChangeTypes()))
	add(typeOf[model.LinkType](), strs(model.AllLinkTypes()))
	add(typeOf[model.Polarity](), strs(model.AllPolarities()))
	add(typeOf[model.LoopType](), strs(model.AllLoopTypes()))
	add(typeOf[model.InterventionType](), strs(model.AllInterventionTypes()))
	add(typeOf[model.InterventionStatus](), strs(model.AllInterventionStatuses()))
	add(typeOf[model.AssumptionType](), strs(model.AllAssumptionTypes()))
	add(typeOf[model.AssumptionStatus](), strs(model.AllAssumptionStatuses()))
	add(typeOf[model.Criticality](), strs(model.AllCriticalities()))
	add(typeOf[model.Rating](), strs(model.AllRatings()))
	add(typeOf[model.MeasurementType](), strs(model.AllMeasurementTypes()))
	add(typeOf[model.Direction](), strs(model.AllDirections()))
	add(typeOf[model.ActorType](), strs(model.AllActorTypes()))
	add(typeOf[model.ActorRole](), strs(model.AllActorRoles()))
	add(typeOf[model.Stance](), strs(model.AllStances()))
	add(typeOf[model.PerspectiveRole](), strs(model.AllPerspectiveRoles()))
	add(typeOf[model.EvidenceType](), strs(model.AllEvidenceTypes()))
	add(typeOf[model.Strength](), strs(model.AllStrengths()))
	add(typeOf[model.EvidenceDirection](), strs(model.AllEvidenceDirections()))
	add(typeOf[model.ProcessTracingTest](), strs(model.AllProcessTracingTests()))
	add(typeOf[model.Sensitivity](), strs(model.AllSensitivities()))
	add(typeOf[model.FactorKind](), strs(model.AllFactorKinds()))
	add(typeOf[model.RiskCategory](), strs(model.AllRiskCategories()))
	add(typeOf[model.PathwayKind](), strs(model.AllPathwayKinds()))
	add(typeOf[model.Projection](), strs(model.AllProjections()))
	add(typeOf[model.Aggregation](), strs(model.AllAggregations()))
	add(typeOf[model.CompositionOp](), strs(model.AllCompositionOps()))
	add(typeOf[model.MetricLayer](), strs(model.AllMetricLayers()))
	add(typeOf[model.Claim](), strs(model.AllClaims()))
	add(typeOf[model.Basis](), strs(model.AllBases()))
	add(typeOf[model.DataAvailability](), strs(model.AllDataAvailabilities()))
	add(typeOf[model.ConsultationStep](), strs(model.AllConsultationSteps()))
	add(typeOf[model.Responsibility](), strs(model.AllResponsibilities()))
	add(typeOf[model.CauseLevel](), strs(model.AllCauseLevels()))
	add(typeOf[model.Consent](), strs(model.AllConsents()))
	add(typeOf[model.CollectionMode](), strs(model.AllCollectionModes()))
	add(typeOf[model.SelectionStatus](), strs(model.AllSelectionStatuses()))
	add(typeOf[model.SelectionSubject](), strs(model.AllSelectionSubjects()))
	add(typeOf[model.ReviewMethod](), strs(model.AllReviewMethods()))
	add(typeOf[model.ReviewDecision](), strs(model.AllReviewDecisions()))
	add(typeOf[model.ReviewOrigin](), strs(model.AllReviewOrigins()))
	add(typeOf[model.ConsultationMethod](), strs(model.AllConsultationMethods()))
	add(typeOf[model.Framework](), strs(model.AllFrameworks()))
	add(typeOf[model.Relation](), strs(model.AllRelations()))
	add(typeOf[model.RevisionTrigger](), strs(model.AllRevisionTriggers()))
	add(typeOf[model.QualityRating](), strs(model.AllQualityRatings()))
	add(typeOf[model.DecisionStatus](), strs(model.AllDecisionStatuses()))
	add(typeOf[model.DataQuality](), strs(model.AllDataQualities()))
	add(typeOf[model.State](), strs(model.AllStates()))
	add(typeOf[model.LinkState](), strs(model.AllLinkStates()))
	add(typeOf[model.RollupMethod](), strs(model.AllRollupMethods()))
	add(typeOf[model.AssumptionSource](), strs(model.AllAssumptionSources()))
	add(typeOf[model.ElementKind](), strs(model.AllElementKinds()))
	add(typeOf[model.EvaluationType](), strs(model.AllEvaluationTypes()))
	add(typeOf[model.Evaluator](), strs(model.AllEvaluators()))
	add(typeOf[model.EvaluationStatus](), strs(model.AllEvaluationStatuses()))
	add(typeOf[model.QuestionStatus](), strs(model.AllQuestionStatuses()))
	return m
}

// enumMapper emits enumeration types as string enums. The document Kind
// type is left to per-field enum tags, since each document kind allows one
// value.
func enumMapper(t reflect.Type) *jsonschema.Schema {
	values, ok := EnumValues()[t]
	if !ok {
		return nil
	}
	enum := make([]any, len(values))
	for i, v := range values {
		enum[i] = v
	}
	return &jsonschema.Schema{Type: "string", Enum: enum}
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func strs[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = string(v)
	}
	return slices.Clip(out)
}
