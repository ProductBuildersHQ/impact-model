# OGSM

If you plan with an OGSM (Objective, Goals, Strategies, Measures), this page
shows where each part comes from in an outcome model, and what the outcome model
adds that a one-page plan written by hand usually leaves out.

```bash
outcome-model export ogsm outcome-model.yaml -O observations.yaml -o ogsm.json
```

The output is an OGSM document in the JSON format of
[prism-roadmap](https://github.com/grokify/prism-roadmap)'s `goals/ogsm`
package, which validates it and renders it as a Markdown document, opening
with the classic one-page table, or as slides.

## The mapping at a glance

| OGSM | Comes from | In outcome model terms |
|------|-----------|-----------------------|
| **Objective** | `context.objective`, or the model's one impact outcome | The one thing the model sets out to achieve |
| The objective's **due date** | `context.timeHorizon.end` | The end of the period the model covers |
| **Goals** | indicators of the outcomes the objective stands for | How the organization will know the objective has been reached |
| **Strategies** | `pathways`, one strategy each | A pathway is a program, product, or strategy: a group of outcomes and the interventions that bring them about |
| The goals a strategy **serves** | the objective outcomes the pathway's outcomes are or lead to | The causal links and precondition groups between them |
| **Measures** | the pathway's `focus`: North Star, one metric that matters, guardrails | The few indicators the program steers by |
| **Actions** | the pathway's `interventionIds` | An intervention is something the organization does to change an outcome, with an owner and dates |
| Every **owner** | `ownerActorId` on the pathway, intervention, or indicator | The actor who answers for it |

Two rules hold throughout:

- **Nothing is invented.** Every field in the OGSM is copied or translated from
  something in the outcome model. A field the model leaves empty stays empty, and
  an OGSM without a clear objective or a strategy is refused instead of written.
- **Order is kept.** Goals, strategies, measures, and actions appear in the
  order the outcome model lists them, so the same model always gives the same
  OGSM.

## One pathway, end to end

The fictional [intermediary example](examples.md) names its objective and the
two resident-level outcomes that say it has been reached:

```yaml
context:
  timeHorizon:
    start: "2026-01-01"
    end: "2028-12-31"
  objective:
    statement: Member nonprofits are strong and connected enough that residents get good services and are heard
    outcomeIds:
      - oc-resident-access
      - oc-resident-voice
```

Its Learning Exchange program is a pathway:

```yaml
pathways:
  - id: pw-learning
    title: Learning Exchange
    description: Peer learning and leadership coaching for member leaders.
    ownerActorId: act-network
    interventionIds: [int-peer-cohorts, int-leader-coaching, int-training-series]
    focus:
      northStarId: ind-practice-maturity
      guardrailIds: [ind-retention-learning]
```

Together they become this OGSM, shown trimmed, with observations as of
2026-09-30:

```json
{
  "objective": {
    "statement": "Member nonprofits are strong and connected enough that residents get good services and are heard",
    "dueDate": "2028-12-31"
  },
  "goals": [
    {
      "id": "ind-residents-served",
      "name": "Residents served by member organizations",
      "owner": "Riverbend Nonprofit Network",
      "baseline": "61000",
      "target": "70000",
      "current": "62000",
      "unit": "residents",
      "dueDate": "2028-12-31",
      "progress": 0.11,
      "status": "On Track"
    }
  ],
  "strategies": [
    {
      "id": "pw-learning",
      "name": "Learning Exchange",
      "description": "Peer learning and leadership coaching for member leaders.",
      "owner": "Riverbend Nonprofit Network",
      "goalIds": ["ind-residents-served", "ind-residents-voice"]
    }
  ],
  "measures": [
    {
      "id": "ind-practice-maturity",
      "name": "Leadership practice maturity",
      "description": "North Star: Median level of participating members on the four-step leadership practice rubric, rated by a coach and the member together.",
      "strategyId": "pw-learning",
      "owner": "Riverbend Nonprofit Network",
      "baseline": "Emerging",
      "target": "Established",
      "current": "Emerging",
      "frequency": "Every six months",
      "status": "On Track"
    }
  ],
  "actions": [
    {
      "id": "int-training-series",
      "name": "Deliver a management skills training series",
      "strategyId": "pw-learning",
      "owner": "Northfield Learning Studio",
      "startDate": "2026-09-01",
      "endDate": "2027-05-31",
      "status": "Planned"
    }
  ]
}
```

- The objective's statement is copied, and it is due at the end of the time
  horizon.
- The indicators of the two resident outcomes are the goals (one shown).
- The pathway becomes a strategy. Its outcomes lead, through member-level
  outcomes, to both resident outcomes, so it serves both goals.
- Its North Star leads its measures, and the guardrail follows it (trimmed
  here). Each measure moves within months; the goals move over years.
- Its three interventions become actions (one shown).

## Objective

An OGSM has one objective. The export takes it from `context.objective`: the
statement, and the outcomes that say it has been reached. Name it with:

```bash
outcome-model patch outcome-model.yaml - <<'EOF'
[{"op":"add","path":"/context/objective","value":{
  "statement":"Member nonprofits are strong and connected enough that residents get good services and are heard",
  "outcomeIds":["oc-resident-access","oc-resident-voice"]}}]
EOF
```

Without `context.objective`, a model with a single impact-level outcome has
that outcome as its objective, its title as the statement. A model with
several impact outcomes and no `context.objective` has no one objective, and
lint says so (OM-CTX-002).

The objective's due date is `context.timeHorizon.end`, and the OGSM's period
is the whole horizon.

## Goals

The goals are the indicators of the objective's outcomes, in the order the
model lists them. They measure the objective, so they belong to no strategy.

| Goal field | From the indicator |
|------------|--------------------|
| ID, Name, Unit | `id`, `name`, `unit` |
| Description | `definition` |
| Owner | the name of the actor in `ownerActorId` |
| Baseline | `baseline` |
| Target, Due | the latest-dated target and its date |
| Direction | `direction` |
| Current, Progress, Status | the indicator's status, only when observations are given with `-O` |

## Strategies

Each **pathway** becomes one strategy, in the order the model lists them. A
pathway marked rejected, such as a program the organization decided against,
is not exported.

| Strategy field | From the pathway |
|----------------|------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Owner | the name of the actor in `ownerActorId` |
| Goals it serves | the goals of each objective outcome that the pathway's outcomes are, or lead to |

"Lead to" follows causal links and precondition groups, leaving out feedback
links and anything rejected. A pathway whose outcomes reach no objective
outcome serves no goal, which prism-roadmap's validation reports: either the
pathway does not belong in this plan, or a link is missing.

Pathways are not ranked, so a strategy has no priority.

## Measures

A strategy's measures are the indicators its pathway steers by, from its
`focus`: the **North Star**, then the **one metric that matters** if the
pathway names one, then the **guardrails**. Each description starts with its
role, such as "North Star:". An indicator that is already a goal is not
repeated as a measure.

| Measure field | From the indicator |
|---------------|--------------------|
| ID, Name, Unit | `id`, `name`, `unit` |
| Description | the role, then `definition` |
| Owner | the name of the actor in `ownerActorId` |
| Baseline | `baseline` |
| Target, Due | the latest-dated target and its date |
| Direction, Source | `direction`, `dataSource` |
| Frequency | `frequency`, in words for a day, week, month, quarter, six months, or year |
| Current, Progress, Status | the indicator's status, only when observations are given with `-O` |

An indicator that is a measure of two pathways appears once, under the first.

## Actions

Actions are each strategy's **interventions**, the things the organization
does.

| Action field | From the intervention |
|--------------|-----------------------|
| ID, Name, Description | `id`, `title`, `description` |
| Strategy | the pathway it is listed in |
| Owner | the name of the actor in `ownerActorId` |
| Status | `planned` → Planned, `in_progress` → In Progress, `completed` → Completed, `cancelled` → Cancelled; `paused` is left blank |
| Start, End | `startDate`, `endDate` |
| External links | `initiativeRefs`, keyed by framework, such as a Jira or VisionStudio item |

An intervention listed in more than one pathway appears once, under the first.
An intervention in no pathway has no strategy to carry out and is left out.

## What an OGSM leaves out

An OGSM is a one-page plan, so the export keeps only what fits it:

| Left out | Where to see it |
|----------|-----------------|
| The outcomes between the strategies and the objective, and the links between them | `report pathway`, the theory of change |
| Assumptions and risks | `report assumptions`, or `export v2mom`, whose obstacles carry them |
| Vision, mission, and values | `export vsem`, `export v2mom` |
| What a pathway commits | `export vsem`, whose strategies carry resources |
| Indicators other than the goals and each pathway's focus | `report metric-tree` |

## What the outcome model adds

An OGSM written by hand says which strategies serve which goals, but not why
they should. Exported from an outcome model:

- each strategy serves the goals its pathway's causal chain reaches, so a
  strategy linked to a goal it cannot move shows up as a missing link
- each measure is a North Star or guardrail that `outcome-model lint` has
  checked sits at or below the organization's ceiling of accountability
- the goals measure the outcomes the objective names, not a separate list
- with `-O`, the goals and measures carry the latest observed values and status

## When the export is refused

An OGSM needs one objective and at least one strategy. If the model has
several impact outcomes and no `context.objective`, or no pathway, `export
ogsm` writes nothing and exits with status 2, naming what to add:

```
outcome-model: outcome-model.yaml: the outcome model cannot be drawn as an OGSM: add context.objective and a pathway
```

## Checking the result

The exported OGSM always passes prism-roadmap's default validation. Its strict
options also require every goal and measure to have a target, every goal,
measure, and action to name an owner, the objective to have a due date, and
every strategy to serve a goal and have a measure and an action. The
intermediary example meets them.

## Going the other way

There is no import from an OGSM. An OGSM has no outcomes or causal links to
read, so an outcome model built from one would have to invent the reasoning
that makes it an outcome model. Start the outcome model from the OGSM instead:
the objective suggests the impact outcomes, each goal an indicator of them,
each strategy a pathway, each measure its North Star or a guardrail, and each
action an intervention.
