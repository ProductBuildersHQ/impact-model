# ADR 0002: Rename the project to outcome-model

**Status:** Accepted
**Date:** 2026-11-12

## Context

[ADR 0001](0001-rename-to-impact-model.md) named the project the impact model,
with the theory of change, the logic model, and OKRs as projections of it. Since
then the model has grown to hold what a strategy, a measurement plan, and an
improvement effort need as well: V2MOM, VSEM, OGSM, and the Balanced Scorecard;
a monitoring, evaluation, and learning plan; and DMAIC. One document now serves
four views of the same work:

| View | Question | Projections |
|------|----------|-------------|
| **Causality** | Why should what we do lead to the change we want? | Problem tree, theory of change, outcome map, logic model, logframe |
| **Strategy** | What do we commit to, and who owns it? | OKRs, V2MOM, VSEM, OGSM, Balanced Scorecard |
| **Measurement and learning** | How will we know, and what must we learn? | Status, metric tree, headline metrics, MEL plan |
| **Improvement** | How do we make a measured process better and keep the gain? | DMAIC |

"Impact model" no longer fits:

- **Scope.** "Impact" names the causal view, and the top outcome level within it.
  A strategy, a scorecard, or an improvement effort is not an impact model, yet
  each is drawn from this one.
- **Audience.** "Impact" reads as social impact to a company, and as financial or
  environmental impact assessment to others. The format is meant for companies
  and nonprofits alike.

What every view shares is the outcome: the causal view links outcomes, a strategy
commits to them, measurement tracks them, and improvement moves them. Releases up
to v0.5.0 had no known users outside this project, so renaming still costs little.

## Decision

The project is the **outcome model**: one document whose causality, strategy,
measurement and learning, and improvement views are drawn as projections.

| Was | Is |
|-----|----|
| `github.com/ProductBuildersHQ/impact-model` | `github.com/ProductBuildersHQ/outcome-model` |
| command `impact-model` | command `outcome-model` |
| packages `pkg/impact`, `pkg/impactio` | `pkg/model`, `pkg/modelio` |
| type `impact.Model` | `model.Model` |
| root package `impactmodel` | `outcomemodel` |
| `kind: ImpactModel`, `ImpactModelStatus` | `kind: OutcomeModel`, `OutcomeModelStatus` |
| schemas `impact-model`, `impact-model-status` | `outcome-model`, `outcome-model-status` |
| `impact-model.yaml` | `outcome-model.yaml` |
| lint rule IDs `IMP-STR-001` | `OM-STR-001` |
| MCP tools `impact_*`, skills `impact-*` | `outcome_*`, `outcome-*` |

The package is `model` rather than `outcome`, so that the root type reads
`model.Model` and the element type `model.Outcome`, not `outcome.Outcome`.

## What did not change

- **The impact level.** `level: impact`, impact outcomes, and the reports that
  speak of the route to impact keep the word, since there it names the top of the
  causal chain.
- **The method names.** The theory of change, the logic model, the outcome map, the
  logframe, and every strategy and improvement format keep their own names.
- **`specVersion`**, which stays `0.1`. The reader accepts any `0.x`.
- **The cross-document reference** `model:<id>#<element>` and the `--model` flag,
  which already used the neutral word.
- **Identifiers of the work.** The initiative `INIT-THEORYOFCHANGESPEC-001` and
  the roadmap items `RMI-THEORYOFCHANGESPEC-*` and `RMI-IMPACTMODEL-*` keep their
  names. New roadmap items take the repository's new slug, `RMI-OUTCOMEMODEL-001`
  onward.
- **Records.** Earlier release notes, ADR 0001, and the judges' recorded verdicts
  describe what was true when they were written.

## Consequences

- **Releases.** v0.3.0 to v0.5.0 were published under the old module path and stay
  there. The first release under the new name is **v0.6.0**.
- **Migration.** A document written for v0.5.0 needs `kind: OutcomeModel` (or
  `OutcomeModelStatus` for a status snapshot), and a `$schema` that names the new
  schema. Lint suppressions and scripts that name rule IDs change `IMP-` to `OM-`.
  Nothing else in the format moves.
- **Links.** GitHub redirects the old repository address. The documentation site
  moves to a new path.
- **No alias.** The old command, kinds, and rule IDs are not kept.

## Alternatives considered

| Name | Why not |
|------|---------|
| Keep `impact-model` | Names one view, and reads as social or financial impact to many readers |
| `strategy-model` | Names another single view, and suggests a plan without the reasons behind it |
| `causal-model` | Accurate for one view, and a term of art in statistics that promises more than a graph of beliefs |
| `results-model` | Close, but "results" leans toward reporting what happened over explaining why |
| `goal-model` | Suggests targets without the causes, evidence, and measures behind them |

`outcome-model` has one drawback: "Outcome Mapping" (IDRC) and the outcome map
projection here share the word. The outcome map is one projection of the model,
and the documentation says it is not IDRC's method.
