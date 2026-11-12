# MEL

A monitoring, evaluation, and learning (MEL) plan says how a team will know
whether its work is producing the results it intends, and what it will do with
what it learns. An outcome model already holds most of one: the outcomes and
how they lead to each other, the indicators that measure them, and the
assumptions they rest on. Add an evaluation plan, a learning agenda, and review
meetings, and the model is the whole plan.

```bash
outcome-model export mel outcome-model.yaml -O observations.yaml -o mel.json
```

The output is a MEL plan in the JSON format of
[prism-roadmap](https://github.com/grokify/prism-roadmap)'s `goals/mel`
package, which validates it and renders it as a Markdown document or as Pandoc
Markdown for a PDF:

```bash
splan goals mel validate mel.json
splan goals mel generate markdown mel.json -o mel-plan.md
```

## The mapping at a glance

| MEL plan | Comes from | In outcome model terms |
|----------|-----------|-----------------------|
| **Results framework** | `outcomes`, with `links` and `preconditionGroups` | The changes that should happen, at the same five levels, each contributing to the outcomes it leads to |
| **Activities** | `interventions` | The work, with the outcomes it targets, its owner, status, and dates |
| **Indicator reference sheets** | `indicators` | How each outcome, intervention, or assumption is measured |
| **Assumptions to monitor** | `assumptions` | What must hold, and its `test` as how it is watched |
| **Evaluation plan** | `evaluations` | When the organization asks whether and why results happened |
| **Learning agenda** | `learningQuestions` | What it most needs to find out, and by when |
| **Reviews** | `reviewMeetings` | When it stops to review what it has learned |
| **Reporting** | `audiences` | Who needs to know what, what they receive, and how often |

Two rules hold throughout:

- **Nothing is invented.** Every field in the plan is copied or translated from
  the outcome model. A field the model leaves empty stays empty: the plan's
  purpose, for one, has no source and is left for you to write.
- **What review rejected stays out.** A rejected outcome, intervention,
  assumption, indicator, evaluation, or learning question is not exported, and
  neither is any reference to it.

## One chain, end to end

From the fictional [intermediary example](examples.md), the Learning Exchange's
chain from coaching to practice exports as two results, one leading to the
other, shown trimmed with observations as of 2026-09-30:

```json
{
  "results": [
    {
      "id": "oc-im-leaders-practice",
      "statement": "Member leaders apply improved management practices",
      "level": "medium_term",
      "contributesTo": ["oc-lt-capable-members"]
    },
    {
      "id": "oc-st-leader-skills",
      "statement": "Leaders gain coaching-based skills and confidence",
      "level": "short_term",
      "contributesTo": ["oc-im-leaders-practice"]
    }
  ],
  "indicators": [
    {
      "id": "ind-practice-maturity",
      "name": "Leadership practice maturity",
      "measures": "oc-im-leaders-practice",
      "type": "qualitative",
      "baseline": { "value": "Emerging", "date": "2026-01-01", "source": "Rubric ratings of the first cohort" },
      "dataSource": "Coach and member rubric ratings",
      "collectionMethod": "Rubric rating at enrollment and at six and twelve months",
      "frequency": "Every six months",
      "dataQuality": {
        "checks": ["A second coach rates one in ten members, and ratings more than one level apart are discussed."],
        "limitations": ["Rated by coaches and members, so it reflects self-assessment as well as observed practice."],
        "lastAssessed": "2026-07-15"
      },
      "current": "Emerging"
    }
  ],
  "learningAgenda": [
    {
      "id": "lq-skills-applied",
      "question": "Do leaders who finish a cohort change how they run their organizations within six months?",
      "rationale": "Decides whether to keep investing in cohorts or shift to coaching.",
      "dueDate": "2027-03-31",
      "relatedIds": ["asm-skills-applied", "oc-st-leader-skills", "oc-im-leaders-practice"]
    }
  ]
}
```

- The two outcomes keep their levels, and the causal link between them becomes
  `contributesTo`.
- The indicator's quality checks become its data quality checks, and its
  caveats its known limitations.
- The learning question named the link `ln-skills-practice`. A MEL plan has no
  links, so the link is written as both of its ends.

## Results framework

Each outcome that review has not rejected is a result, at the same level:
`output`, `short_term`, `medium_term`, `long_term`, or `impact`. A result
contributes to the outcomes its causal links lead to, and to the target of any
precondition group it belongs to.

**Feedback links are left out.** An outcome model can mark a link as an
intentional loop, such as stronger collaboration feeding back into peer ties.
A results framework is a hierarchy, and a MEL plan treats a loop as an error,
so the export keeps the forward links only. The loop stays in the model, and
`report pathway` draws it.

## Activities

Each intervention is an activity that produces the outcomes it targets, with
its owner, status, `startDate`, and `endDate`. A paused intervention is a
paused activity.

## Indicators

Each indicator is a reference sheet. It measures what it measured in the
model: an outcome (a result), an intervention (an activity, whose delivery it
tracks), or an assumption (which it watches).

| Reference sheet field | From the indicator |
|-----------------------|--------------------|
| Definition, unit, direction, disaggregation | `definition`, `unit`, `direction`, `disaggregation` |
| Type | `measurementType`: quantitative, qualitative, and theme counts as quantitative; a hybrid indicator has no single type and is left blank |
| Baseline, targets | `baseline` and `targets`, with dates, labels, and the baseline's source |
| Data source | `dataSource`, or the collection `system` |
| Collection method | `method`, or the collection `mode` |
| Frequency | `frequency`, in words for a day, week, month, quarter, six months, or year |
| Responsible | the name of the actor in `ownerActorId` |
| Data quality | `qualityAssurance` checks and last assessment, and `caveats` as known limitations |
| Current | the latest observation, only when observations are given with `-O` |

## Assumptions

Each assumption is one to monitor. How it is monitored is its `test`, and the
indicators that watch it are its `indicatorIds` and any indicator that measures
it. The results that depend on it are what it applies to: an outcome itself, a
link's target, or an intervention's targets.

## Evaluation, learning, and review

Evaluations, learning questions, and review meetings carry over field by field,
with owners as actor names and frequencies in words. Statuses are written as
the MEL plan writes them, such as `planned` as Planned. A link that an
evaluation or learning question names is written as both of its ends. See
[Evaluation and learning](authoring.md#evaluation-and-learning) for how to add
them to the model.

## Reporting

Each audience is a report: its `needs` as what it needs to know, its
`channels` as what it receives, its `frequency`, who prepares it, and the
indicators it sees, its headline list first and then every other indicator
reported to it. The measurement budget's staff minutes a week become the plan's
budget.

## When the export is refused

A MEL plan needs at least one result. A model whose outcomes have all been
rejected is refused with exit status 2:

```
outcome-model: outcome-model.yaml: the outcome model cannot be drawn as a MEL plan: add an outcome
```

## Checking the result

The exported plan always passes prism-roadmap's default MEL validation. The
strict options also ask that every result has an indicator, every indicator a
full collection plan, every assumption a way to monitor it, and the plan an
evaluation, a learning question, and a review. Those are gaps in the model, so
`outcome-model lint` points at most of them first: OM-MEA-001 and OM-MEA-002 for
measurement, and OM-LRN-001 to 003 for evaluation and learning.

## Going the other way

There is no import from a MEL plan. A results framework says what leads to
what, but not why, and an outcome model built from one would have to invent the
rationale behind each link. Start the outcome model from the plan instead: each
result suggests an outcome, each activity an intervention, and each indicator,
evaluation, and learning question carries over as it is.
