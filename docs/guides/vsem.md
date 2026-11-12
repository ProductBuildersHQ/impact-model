# VSEM

If you plan with a VSEM (Vision, Strategy, Execution, Metrics), this page shows
where each part comes from in an outcome model, and what the outcome model adds
that a one-page plan written by hand usually leaves out.

```bash
outcome-model export vsem outcome-model.yaml -O observations.yaml -o vsem.json
```

The output is a VSEM document in the JSON format of
[prism-roadmap](https://github.com/grokify/prism-roadmap)'s `goals/vsem`
package, which validates it and renders it as a Markdown document or slides.

## The mapping at a glance

| VSEM | Comes from | In outcome model terms |
|------|-----------|-----------------------|
| **Vision** | `context.vision`, and `context.mission` | The future the organization works toward, and why it exists |
| **Strategy** | `pathways`, one strategy each | A pathway is a program, product, or strategy: a group of outcomes and the interventions that bring them about |
| A strategy's **resources** | the pathway's `resources` | What the pathway commits: people, money, tooling |
| **Execution** | the pathway's `interventionIds` | An intervention is something the organization does to change an outcome, with an owner and dates |
| A strategy's **metrics** | the pathway's `focus`: North Star, one metric that matters, guardrails | The few indicators the program steers by |
| **Metrics** for the vision | indicators of `impact`-level outcomes | How the organization will know the vision is getting closer |
| Every **owner** | `ownerActorId` on the pathway, intervention, or indicator | The actor who answers for it |

Two rules hold throughout:

- **Nothing is invented.** Every field in the VSEM is copied or translated from
  something in the outcome model. A field the model leaves empty stays empty, and
  a VSEM without a vision or a strategy is refused instead of written.
- **Order is kept.** Strategies, execution, and metrics appear in the order the
  outcome model lists them, so the same model always gives the same VSEM.

## One pathway, end to end

This pathway from the fictional [intermediary example](examples.md) is a program
called Learning Exchange:

```yaml
pathways:
  - id: pw-learning
    title: Learning Exchange
    kind: program
    description: Peer learning and leadership coaching for member leaders.
    ownerActorId: act-network
    interventionIds: [int-peer-cohorts, int-leader-coaching, int-training-series]
    focus:
      northStarId: ind-practice-maturity
      guardrailIds: [ind-retention-learning]
    resources:
      - name: Program staff
        allocation: 3 staff
      - name: Volunteer coaches
        allocation: A pool of about 15 coaches
      - name: Training partner
        allocation: Fee-for-service agreement
```

It becomes this strategy, with its execution and metrics, shown trimmed, with
observations as of 2026-09-30:

```json
{
  "strategies": [
    {
      "id": "pw-learning",
      "name": "Learning Exchange",
      "description": "Peer learning and leadership coaching for member leaders.",
      "owner": "Riverbend Nonprofit Network",
      "resources": [
        { "name": "Program staff", "allocation": "3 staff" },
        { "name": "Volunteer coaches", "allocation": "A pool of about 15 coaches" },
        { "name": "Training partner", "allocation": "Fee-for-service agreement" }
      ]
    }
  ],
  "execution": [
    {
      "id": "int-training-series",
      "name": "Deliver a management skills training series",
      "strategyId": "pw-learning",
      "kind": "initiative",
      "owner": "Northfield Learning Studio",
      "status": "Planned",
      "startDate": "2026-09-01",
      "endDate": "2027-05-31"
    }
  ],
  "metrics": [
    {
      "id": "ind-practice-maturity",
      "name": "Leadership practice maturity",
      "description": "North Star: Median level of participating members on the four-step leadership practice rubric, rated by a coach and the member together.",
      "strategyId": "pw-learning",
      "owner": "Riverbend Nonprofit Network",
      "baseline": "Emerging",
      "target": "Established",
      "current": "Emerging",
      "dueDate": "2028-12-31",
      "frequency": "Every six months",
      "status": "On Track"
    }
  ]
}
```

- The pathway's title, description, owner, and resources become the strategy's.
- Its three interventions become execution items (one shown), each naming the
  strategy it carries out.
- The North Star leads the strategy's metrics, and the guardrail follows it
  (trimmed here).
- The outcomes and the links between them are not in the VSEM. They are the
  reasoning behind the strategy, and `report pathway` draws them.

## Vision

The VSEM's vision is `context.vision`, word for word, and its mission is
`context.mission`. A VSEM needs one of them; some teams state a mission in
place of a vision.

## Strategy

Each **pathway** becomes one strategy, in the order the model lists them. A
pathway marked rejected, such as a program the organization decided against, is
not exported.

| Strategy field | From the pathway |
|----------------|------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Owner | the name of the actor in `ownerActorId` |
| Resources | `resources`, each with its name, allocation, and description |

Pathways are not ranked, so a strategy has no priority.

### Resources

A pathway's `resources` say what it commits as a whole, which is what a VSEM's
strategy section asks for:

```bash
outcome-model pathway set outcome-model.yaml pw-learning \
  --set 'resources=[{"name":"Program staff","allocation":"3 staff"}]'
```

An intervention also has a `resources` field. That one says what a single piece
of work uses, and the logic model shows it as an input; it is not exported to
the VSEM.

## Execution

Execution is each strategy's **interventions**, the things the organization
does.

| Execution field | From the intervention |
|-----------------|-----------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Strategy | the pathway it is listed in |
| Kind | `program` when the intervention's type is `program`; otherwise `initiative` |
| Owner | the name of the actor in `ownerActorId` |
| Status | `planned` → Planned, `in_progress` → In Progress, `completed` → Completed, `cancelled` → Cancelled; `paused` is left blank |
| Start, End | `startDate`, `endDate` |
| External links | `initiativeRefs`, keyed by framework, such as a Jira or VisionStudio item |

Set an intervention's dates with:

```bash
outcome-model intervention set outcome-model.yaml int-training-series \
  --start-date 2026-09-01 --end-date 2027-05-31
```

Ongoing work has a start and no end. Lint warns when either date falls outside
the model's time horizon (OM-INT-008).

An intervention listed in more than one pathway appears once, under the first,
because each item in a VSEM has its own ID. An intervention in no pathway has no
strategy to carry out and is left out.

## Metrics

A strategy's metrics are the indicators its pathway steers by, from its `focus`:
the **North Star**, then the **one metric that matters** if the pathway names
one, then the **guardrails**. Each description starts with its role, such as
"North Star:".

The **vision's metrics** are the indicators of `impact`-level outcomes that are
not already a strategy's metric. They have no strategy.

| Metric field | From the indicator |
|--------------|--------------------|
| ID, Name, Unit | `id`, `name`, `unit` |
| Description | the role, then `definition` |
| Owner | the name of the actor in `ownerActorId`: the person accountable for the number |
| Baseline | `baseline` |
| Target, Due | the latest-dated target and its date |
| Direction, Source | `direction`, `dataSource` |
| Frequency | `frequency`, in words for a day, week, month, quarter, six months, or year |
| Current, Progress, Status | the indicator's status, only when observations are given with `-O` |

Status is On Track, At Risk, Behind, or Achieved. An indicator that is a metric
of two pathways appears once, under the first.

## What a VSEM leaves out

A VSEM is a one-page plan, so the export keeps only what fits it:

| Left out | Where to see it |
|----------|-----------------|
| The outcomes and the causal links between them | `report pathway`, the theory of change |
| Assumptions and risks | `report assumptions`, or `export v2mom`, whose obstacles carry them |
| Values | `export v2mom` |
| Indicators other than a pathway's focus and the impact outcomes | `report metric-tree` |
| What each audience is told | `report headline` |

## What the outcome model adds

A VSEM written by hand says where effort goes and who answers for each number,
but not why that effort should move those numbers. Exported from an outcome
model:

- each strategy is backed by a pathway of outcomes, each link with its reason
- each North Star is one that `outcome-model lint` has checked sits at or below
  the organization's ceiling of accountability
- the vision's metrics are the measures of the impact the strategies serve
- with `-O`, the metrics carry the latest observed values and status

## When the export is refused

A VSEM needs a vision or mission and at least one strategy. If the outcome model
has neither `context.vision` nor `context.mission`, or no pathway, `export
vsem` writes nothing and exits with status 2, naming what to add:

```
outcome-model: outcome-model.yaml: the outcome model cannot be drawn as a VSEM: add context.vision or context.mission and a pathway
```

## Checking the result

The exported VSEM always passes prism-roadmap's default validation. Its strict
options also require every metric and execution item to name an owner, every
metric to have a target, and every strategy to commit resources and have work
and a metric. A model that sets those fields passes them too.

## Going the other way

There is no import from a VSEM. A VSEM has no outcomes or causal links to read,
so an outcome model built from one would have to invent the reasoning that makes
it an outcome model. Start the outcome model from the VSEM instead: each strategy
suggests a pathway, each execution item an intervention, and each metric an
indicator.
