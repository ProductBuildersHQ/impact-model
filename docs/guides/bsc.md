# Balanced Scorecard

If you plan with a Balanced Scorecard, this page shows where each part comes
from in an impact model, and what the impact model adds that a scorecard
written by hand usually leaves out.

```bash
impact-model export bsc impact-model.yaml -O observations.yaml -o bsc.json
```

The output is a Balanced Scorecard in the JSON format of
[prism-roadmap](https://github.com/grokify/prism-roadmap)'s `goals/bsc`
package, which validates it and renders it as a Markdown document, with the
strategy map as a diagram, or as slides.

## The mapping at a glance

| Scorecard | Comes from | In impact model terms |
|-----------|-----------|-----------------------|
| **Mission** and **vision** | `context.mission`, `context.vision` | Why the organization exists, and the future it works toward |
| **Perspectives** | the `scorecardPerspective` values the outcomes use | Where an outcome sits: mission, financial, customer, internal process, or learning and growth |
| **Strategic objectives** | outcomes with a `scorecardPerspective` | A change the organization sets out to bring about |
| The **strategy map** | causal links and precondition groups between those outcomes | Why achieving one objective helps achieve another |
| **Strategic themes** | `pathways` that hold an objective | A program, product, or strategy: a group of outcomes and the interventions behind them |
| **Measures** | the objectives' indicators, with `leadLag` | How the organization will know each objective is being achieved |
| **Initiatives** | `interventions` | Something the organization does to change an outcome, with an owner and dates |
| Every **owner** | `ownerActorId` on the outcome, indicator, pathway, or intervention | The actor who answers for it |

Two rules hold throughout:

- **Nothing is invented.** Every field in the scorecard is copied or
  translated from something in the impact model. A field the model leaves
  empty stays empty, and a model with no outcome in a scorecard perspective is
  refused instead of written.
- **Order is kept.** Objectives, measures, themes, and initiatives appear in
  the order the impact model lists them, so the same model always gives the
  same scorecard.

## Placing outcomes in perspectives

A Balanced Scorecard sorts objectives into perspectives. An impact model sorts
outcomes by level instead, so you say which outcomes are scorecard objectives,
and where they sit:

```bash
impact-model outcome set impact-model.yaml oc-im-lower-overhead \
  --scorecard-perspective customer --owner act-network
impact-model indicator set impact-model.yaml ind-overhead-share --lead-lag lead
```

| Perspective | `scorecardPerspective` | Usually holds |
|-------------|------------------------|---------------|
| Mission | `mission` | The change the organization exists for; public and nonprofit scorecards put it on top |
| Financial | `financial` | How the organization succeeds for those who fund it |
| Customer | `customer` | What the people it serves experience |
| Internal Process | `internal_process` | The processes it must excel at |
| Learning and Growth | `learning_growth` | The people, systems, and culture it must build |

Only the perspectives the outcomes use are drawn, top first, with the mission
above the financial perspective. Not every outcome needs a perspective: a
scorecard keeps to the objectives that matter most, and the outcomes between
them still shape the strategy map.

## One chain, end to end

The fictional [intermediary example](examples.md) places twelve of its
outcomes. Its Learning Exchange program delivers peer-learning cohorts, which
build leaders' skills, which change how they manage their organizations, which
helps residents get good services. Trimmed to the fields that matter here:

```yaml
outcomes:
  - id: oc-resident-access
    title: Residents get timely, good-quality services from local nonprofits
    level: impact
    ownerActorId: act-network
    scorecardPerspective: mission
  - id: oc-im-leaders-practice
    title: Member leaders apply improved management practices
    level: medium_term
    ownerActorId: act-network
    scorecardPerspective: customer
  - id: oc-st-leader-skills
    title: Leaders gain coaching-based skills and confidence
    level: short_term
  - id: oc-out-cohorts
    title: Peer-learning cohorts and coaching sessions delivered
    level: output
    ownerActorId: act-network
    scorecardPerspective: internal_process
```

Skills are not placed, so the strategy map passes through them. The chain
becomes these objectives, shown trimmed, with observations as of 2026-09-30:

```json
{
  "objectives": [
    { "id": "oc-resident-access", "name": "Residents get timely, good-quality services from local nonprofits",
      "perspectiveId": "mission", "owner": "Riverbend Nonprofit Network" },
    { "id": "oc-im-leaders-practice", "name": "Member leaders apply improved management practices",
      "perspectiveId": "customer", "themeIds": ["pw-learning"], "owner": "Riverbend Nonprofit Network",
      "contributesTo": ["oc-resident-access"] },
    { "id": "oc-out-cohorts", "name": "Peer-learning cohorts and coaching sessions delivered",
      "perspectiveId": "internal-process", "themeIds": ["pw-learning"], "owner": "Riverbend Nonprofit Network",
      "contributesTo": ["oc-lt-collab", "oc-lt-collective-voice", "oc-im-leaders-practice"] }
  ],
  "measures": [
    { "id": "ind-practice-maturity", "name": "Leadership practice maturity",
      "objectiveId": "oc-im-leaders-practice", "type": "lead", "owner": "Riverbend Nonprofit Network",
      "baseline": "Emerging", "target": "Established", "current": "Emerging",
      "frequency": "Every six months", "status": "On Track" }
  ]
}
```

- Each placed outcome becomes an objective in its perspective, with its owner.
- The cohorts lead to the leaders' practice through skills, which are not an
  objective, so the arrow skips them. They also lead, through peer ties and
  joint positions, to two customer objectives.
- The leaders' practice leads, through capable members, to residents getting
  services.
- Both the cohorts and the leaders' practice are in the Learning Exchange
  pathway, so they belong to its theme.

## Strategy map

An objective contributes to the nearest objectives its outcome leads to. The
export follows causal links that are not rejected or feedback, and
precondition groups, through outcomes that are not objectives, and stops at
the first objective on each path. Feedback links are left out, since a
strategy map runs one way.

A strategy map runs from cause at the bottom to effect at the top, so lint
warns when an objective leads to one in a lower perspective (IMP-BSC-001):
one of the two perspectives may be wrong.

## Objectives and themes

| Objective field | From the outcome |
|-----------------|------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Perspective | `scorecardPerspective` |
| Owner | the name of the actor in `ownerActorId` |
| Themes | the pathways that list the outcome |
| Contributes to | the strategy map, above |

Each pathway that holds an objective becomes a strategic theme, with the
pathway's title, description, and owner. A rejected pathway is not a theme.

## Measures

An objective's measures are the indicators that measure its outcome.

| Measure field | From the indicator |
|---------------|--------------------|
| ID, Name, Unit | `id`, `name`, `unit` |
| Description | `definition` |
| Type | `leadLag`: `lead` moves before the results, `lag` reports one |
| Owner | the name of the actor in `ownerActorId` |
| Baseline | `baseline` |
| Target, Due | the latest-dated target and its date |
| Direction, Source | `direction`, `dataSource` |
| Frequency | `frequency`, in words for a day, week, month, quarter, six months, or year |
| Current, Progress, Status | the indicator's status, only when observations are given with `-O` |

The scorecard's stoplight thresholds are left empty: an indicator's status
comes from measurement.

## Initiatives

Initiatives are the **interventions**. An initiative serves the objectives
among its targets, or else the nearest objectives its targets lead to, as on
the strategy map. An intervention that leads to no objective is left out.

| Initiative field | From the intervention |
|------------------|-----------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Objectives | as above |
| Owner | the name of the actor in `ownerActorId` |
| Start, End | `startDate`, `endDate` |
| Budget | `resources` |
| Status | `planned` → Planned, `in_progress` → In Progress, `completed` → Completed, `cancelled` → Cancelled; `paused` is left blank |
| External links | `initiativeRefs`, keyed by framework |

## What a scorecard leaves out

| Left out | Where to see it |
|----------|-----------------|
| Outcomes without a perspective, and the links' reasons | `report pathway`, the theory of change |
| Assumptions and risks | `report assumptions`, or `export v2mom`, whose obstacles carry them |
| A one-sentence strategy | The model's pathways are its strategy; they become the themes |
| Indicators of outcomes that are not objectives | `report metric-tree` |

## What the impact model adds

A scorecard drawn by hand has arrows on its strategy map, but rarely says why
each one holds. Exported from an impact model:

- each arrow stands for a causal chain whose links carry a rationale,
  assumptions, and evidence
- each objective's place in the chain is checked by lint, so an arrow that
  runs down the perspectives is caught
- with `-O`, the measures carry the latest observed values and status

## When the export is refused

A scorecard needs at least one objective. If no outcome that review has not
rejected has a `scorecardPerspective`, `export bsc` writes nothing and exits
with status 2:

```
impact-model: impact-model.yaml: the impact model cannot be drawn as a Balanced Scorecard: add a scorecardPerspective to an outcome
```

## Checking the result

The exported scorecard always passes prism-roadmap's default validation. Its
strict options also require an owner on every objective and measure, a target
and a lead or lag type on every measure, both lead and lag measures, an
objective in every perspective drawn, and a connected strategy map. Lint rules
IMP-BSC-001 to 004 check the model for the same things before you export. The
intermediary example meets them.

## Going the other way

There is no import from a Balanced Scorecard. A scorecard's arrows say that one
objective helps another, not how or why, so an impact model built from one
would have to invent the reasoning. Start the impact model from the scorecard
instead: each objective suggests an outcome in its perspective, each arrow a
causal link to explain, each measure an indicator, and each initiative an
intervention.
