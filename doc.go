// Package outcomemodel is a Go-first specification for representing,
// measuring, and reporting on an impact: why a program expects to bring about
// change, what it does, and what it measures. The theory of change, the logic
// model, the outcome map, the logframe, and OKRs are projections of one outcome
// model.
//
// Go structs are the source of truth; JSON Schemas are generated from them
// and embedded. Tooling built on the model validates documents, analyzes
// their causal graphs, edits them through validated patches, rolls up
// indicator observations into pathway status, and renders reports and
// slide decks.
//
// Key packages:
//   - pkg/ref: typed references to outcome models and their elements
//   - pkg/model: the outcome model
//   - schema: JSON Schemas generated from pkg/model, embedded for runtime use
//
// See docs/specs/initiatives/INIT-THEORYOFCHANGESPEC-001 for the PRD, TRD,
// plan, and roadmap, and docs/research for the methodology background.
package outcomemodel
