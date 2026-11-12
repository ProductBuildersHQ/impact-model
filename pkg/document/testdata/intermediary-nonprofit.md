# Riverbend Nonprofit Network theory of change

Outcome model: Theory of Change, Logic Models, VSEM, Monitoring, Evaluation, and Learning Plan, Glossary

*Example stamp*

## Contents

1. Theory of Change
2. Logic Models
3. VSEM: Vision, Strategy, Execution, Metrics
4. Monitoring, Evaluation, and Learning Plan
5. Glossary

## Theory of Change

### Theory of change on one page

```d2
grid-rows: 5
vertical-gap: 24

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 28
    style.bold: true
  }
}

impact: Impact {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-resident-access": "Residents get timely,\ngood-quality services from\nlocal nonprofits\n(inferred)" {
    class: impact
  }
  "oc-resident-voice": "Residents' priorities\nshape local public\ndecisions (hypothesis)" {
    class: impact
  }
}

long_term: Long-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-lt-capable-members": "Member nonprofits run\nstrong, well-managed\norganizations" {
    class: long_term
  }
  "oc-lt-collab": "Members routinely\ncollaborate instead of\nduplicating effort" {
    class: long_term
  }
  "oc-lt-sustainable": "Members have diversified,\nstable funding (inferred)" {
    class: long_term
  }
  "oc-lt-collective-voice": "Local government and\nfunders consult members\nthrough a coordinated\nchannel (hypothesis)" {
    class: long_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-im-peer-network": "Members turn to each other\nfor advice and joint work\n(inferred)" {
    class: medium_term
  }
  "oc-im-leaders-practice": "Member leaders apply\nimproved management\npractices" {
    class: medium_term
  }
  "oc-im-lower-overhead": "Members spend less staff\ntime and money on\nback-office work" {
    class: medium_term
  }
  "oc-im-grant-pipeline": "Members apply for more\ngrants that fit their\nmission (inferred)" {
    class: medium_term
  }
  "oc-im-joint-positions": "Members adopt and share\njoint positions on local\nissues (hypothesis)" {
    class: medium_term
  }
  "oc-im-officials-engage": "Local officials and\nfunders engage members in\nplanning" {
    class: medium_term
  }
}

short_term: Short-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-st-peer-ties": "Members form peer\nrelationships across sizes\n(inferred)" {
    class: short_term
  }
  "oc-st-leader-skills": "Leaders gain\ncoaching-based skills and\nconfidence" {
    class: short_term
  }
  "oc-st-pooled-use": "Members use pooled\nback-office services\n(inferred)" {
    class: short_term
  }
  "oc-st-calendar-use": "Members consult the shared\ngrants calendar when\nplanning (hypothesis)" {
    class: short_term
  }
  "oc-st-policy-know": "Members understand local\npolicy and funding\nprocesses (inferred)" {
    class: short_term
  }
  "oc-st-channel": "A standing channel exists\nbetween members and the\ncounty planning office" {
    class: short_term
  }
}

activities: Activities {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "int-peer-cohorts": "Run mixed-size\npeer-learning cohorts" {
    class: activities
  }
  "int-leader-coaching": "Offer one-to-one\nleadership coaching" {
    class: activities
  }
  "int-training-series": "Deliver a management\nskills training series\n(inferred)" {
    class: activities
  }
  "int-pooled-backoffice": "Run pooled bookkeeping and\npayroll help" {
    class: activities
  }
  "int-group-purchasing": "Negotiate group purchasing\nagreements" {
    class: activities
  }
  "int-grants-calendar": "Maintain a shared grants\ncalendar" {
    class: activities
  }
  "int-policy-briefings": "Publish plain-language\npolicy briefings\n(hypothesis)" {
    class: activities
  }
  "int-joint-forums": "Convene member forums to\ndraft joint positions" {
    class: activities
  }
  "int-liaison-channel": "Set up a standing liaison\nchannel with the county\nplanning office" {
    class: activities
  }
  "int-funder-roundtable": "Host a funder and member\nroundtable" {
    class: activities
  }
}

activities."int-peer-cohorts" -> short_term."oc-st-peer-ties"
activities."int-peer-cohorts" -> short_term."oc-st-leader-skills"
activities."int-leader-coaching" -> short_term."oc-st-leader-skills"
activities."int-leader-coaching" -> short_term."oc-st-peer-ties"
activities."int-training-series" -> short_term."oc-st-leader-skills"
activities."int-pooled-backoffice" -> short_term."oc-st-pooled-use"
activities."int-group-purchasing" -> short_term."oc-st-pooled-use"
activities."int-grants-calendar" -> short_term."oc-st-calendar-use"
activities."int-policy-briefings" -> short_term."oc-st-policy-know"
activities."int-joint-forums" -> short_term."oc-st-policy-know"
activities."int-joint-forums" -> medium_term."oc-im-joint-positions"
activities."int-liaison-channel" -> short_term."oc-st-channel"
activities."int-funder-roundtable" -> short_term."oc-st-channel"
activities."int-funder-roundtable" -> medium_term."oc-im-officials-engage"
short_term."oc-st-leader-skills" -> medium_term."oc-im-leaders-practice"
short_term."oc-st-peer-ties" -> medium_term."oc-im-peer-network"
medium_term."oc-im-leaders-practice" -> long_term."oc-lt-capable-members"
medium_term."oc-im-peer-network" -> long_term."oc-lt-collab"
short_term."oc-st-pooled-use" -> medium_term."oc-im-lower-overhead"
short_term."oc-st-calendar-use" -> medium_term."oc-im-grant-pipeline"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-capable-members"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-collab"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-sustainable"
medium_term."oc-im-grant-pipeline" -> long_term."oc-lt-sustainable"
short_term."oc-st-policy-know" -> medium_term."oc-im-joint-positions"
short_term."oc-st-peer-ties" -> medium_term."oc-im-joint-positions"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collab"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collective-voice"
medium_term."oc-im-officials-engage" -> long_term."oc-lt-collective-voice"
short_term."oc-st-channel" -> medium_term."oc-im-officials-engage"
long_term."oc-lt-capable-members" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-sustainable" -> impact."oc-resident-access"
long_term."oc-lt-collective-voice" -> impact."oc-resident-voice"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

Feedback, not drawn:

- Members routinely collaborate instead of duplicating effort leads back to Members form peer relationships across sizes (inferred)

### Theory of change: Learning Exchange

```d2
grid-rows: 5
vertical-gap: 24

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 28
    style.bold: true
  }
}

impact: Impact {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-resident-voice": "Residents' priorities\nshape local public\ndecisions (hypothesis)" {
    class: impact
  }
  "oc-resident-access": "Residents get timely,\ngood-quality services from\nlocal nonprofits\n(inferred)" {
    class: impact
  }
}

long_term: Long-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-lt-collab": "Members routinely\ncollaborate instead of\nduplicating effort" {
    class: long_term
  }
  "oc-lt-collective-voice": "Local government and\nfunders consult members\nthrough a coordinated\nchannel (hypothesis)" {
    class: long_term
  }
  "oc-lt-capable-members": "Member nonprofits run\nstrong, well-managed\norganizations" {
    class: long_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-im-peer-network": "Members turn to each other\nfor advice and joint work\n(inferred)" {
    class: medium_term
  }
  "oc-im-joint-positions": "Members adopt and share\njoint positions on local\nissues (hypothesis)" {
    class: medium_term
  }
  "oc-im-leaders-practice": "Member leaders apply\nimproved management\npractices" {
    class: medium_term
  }
}

short_term: Short-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-st-peer-ties": "Members form peer\nrelationships across sizes\n(inferred)" {
    class: short_term
  }
  "oc-st-leader-skills": "Leaders gain\ncoaching-based skills and\nconfidence" {
    class: short_term
  }
}

activities: Activities {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "int-peer-cohorts": "Run mixed-size\npeer-learning cohorts" {
    class: activities
  }
  "int-leader-coaching": "Offer one-to-one\nleadership coaching" {
    class: activities
  }
  "int-training-series": "Deliver a management\nskills training series\n(inferred)" {
    class: activities
  }
}

activities."int-peer-cohorts" -> short_term."oc-st-peer-ties"
activities."int-peer-cohorts" -> short_term."oc-st-leader-skills"
activities."int-leader-coaching" -> short_term."oc-st-leader-skills"
activities."int-leader-coaching" -> short_term."oc-st-peer-ties"
activities."int-training-series" -> short_term."oc-st-leader-skills"
short_term."oc-st-leader-skills" -> medium_term."oc-im-leaders-practice"
short_term."oc-st-peer-ties" -> medium_term."oc-im-peer-network"
medium_term."oc-im-leaders-practice" -> long_term."oc-lt-capable-members"
medium_term."oc-im-peer-network" -> long_term."oc-lt-collab"
short_term."oc-st-peer-ties" -> medium_term."oc-im-joint-positions"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collab"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collective-voice"
long_term."oc-lt-capable-members" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-collective-voice" -> impact."oc-resident-voice"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

Feedback, not drawn:

- Members routinely collaborate instead of duplicating effort leads back to Members form peer relationships across sizes (inferred)

### Theory of change: Shared Services

```d2
grid-rows: 5
vertical-gap: 24

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 28
    style.bold: true
  }
}

impact: Impact {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-resident-access": "Residents get timely,\ngood-quality services from\nlocal nonprofits\n(inferred)" {
    class: impact
  }
  "oc-resident-voice": "Residents' priorities\nshape local public\ndecisions (hypothesis)" {
    class: impact
  }
}

long_term: Long-term outcomes {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-lt-capable-members": "Member nonprofits run\nstrong, well-managed\norganizations" {
    class: long_term
  }
  "oc-lt-collab": "Members routinely\ncollaborate instead of\nduplicating effort" {
    class: long_term
  }
  "oc-lt-sustainable": "Members have diversified,\nstable funding (inferred)" {
    class: long_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-im-lower-overhead": "Members spend less staff\ntime and money on\nback-office work" {
    class: medium_term
  }
  "oc-im-grant-pipeline": "Members apply for more\ngrants that fit their\nmission (inferred)" {
    class: medium_term
  }
}

short_term: Short-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-st-pooled-use": "Members use pooled\nback-office services\n(inferred)" {
    class: short_term
  }
  "oc-st-calendar-use": "Members consult the shared\ngrants calendar when\nplanning (hypothesis)" {
    class: short_term
  }
}

activities: Activities {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "int-pooled-backoffice": "Run pooled bookkeeping and\npayroll help" {
    class: activities
  }
  "int-group-purchasing": "Negotiate group purchasing\nagreements" {
    class: activities
  }
  "int-grants-calendar": "Maintain a shared grants\ncalendar" {
    class: activities
  }
}

activities."int-pooled-backoffice" -> short_term."oc-st-pooled-use"
activities."int-group-purchasing" -> short_term."oc-st-pooled-use"
activities."int-grants-calendar" -> short_term."oc-st-calendar-use"
short_term."oc-st-pooled-use" -> medium_term."oc-im-lower-overhead"
short_term."oc-st-calendar-use" -> medium_term."oc-im-grant-pipeline"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-capable-members"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-collab"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-sustainable"
medium_term."oc-im-grant-pipeline" -> long_term."oc-lt-sustainable"
long_term."oc-lt-capable-members" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-sustainable" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

### Theory of change: Civic Voice

```d2
grid-rows: 5
vertical-gap: 24

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 28
    style.bold: true
  }
}

impact: Impact {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-resident-access": "Residents get timely,\ngood-quality services from\nlocal nonprofits\n(inferred)" {
    class: impact
  }
  "oc-resident-voice": "Residents' priorities\nshape local public\ndecisions (hypothesis)" {
    class: impact
  }
}

long_term: Long-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-lt-collab": "Members routinely\ncollaborate instead of\nduplicating effort" {
    class: long_term
  }
  "oc-lt-collective-voice": "Local government and\nfunders consult members\nthrough a coordinated\nchannel (hypothesis)" {
    class: long_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-im-joint-positions": "Members adopt and share\njoint positions on local\nissues (hypothesis)" {
    class: medium_term
  }
  "oc-im-officials-engage": "Local officials and\nfunders engage members in\nplanning" {
    class: medium_term
  }
}

short_term: Short-term outcomes {
  grid-columns: 2
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "oc-st-policy-know": "Members understand local\npolicy and funding\nprocesses (inferred)" {
    class: short_term
  }
  "oc-st-channel": "A standing channel exists\nbetween members and the\ncounty planning office" {
    class: short_term
  }
}

activities: Activities {
  grid-columns: 3
  horizontal-gap: 24
  vertical-gap: 24
  class: stage

  "int-policy-briefings": "Publish plain-language\npolicy briefings\n(hypothesis)" {
    class: activities
  }
  "int-joint-forums": "Convene member forums to\ndraft joint positions" {
    class: activities
  }
  "int-liaison-channel": "Set up a standing liaison\nchannel with the county\nplanning office" {
    class: activities
  }
  "int-funder-roundtable": "Host a funder and member\nroundtable" {
    class: activities
  }
}

activities."int-policy-briefings" -> short_term."oc-st-policy-know"
activities."int-joint-forums" -> short_term."oc-st-policy-know"
activities."int-joint-forums" -> medium_term."oc-im-joint-positions"
activities."int-liaison-channel" -> short_term."oc-st-channel"
activities."int-funder-roundtable" -> short_term."oc-st-channel"
activities."int-funder-roundtable" -> medium_term."oc-im-officials-engage"
short_term."oc-st-policy-know" -> medium_term."oc-im-joint-positions"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collab"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collective-voice"
medium_term."oc-im-officials-engage" -> long_term."oc-lt-collective-voice"
short_term."oc-st-channel" -> medium_term."oc-im-officials-engage"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-collective-voice" -> impact."oc-resident-voice"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

### Narrative summary

Small and mid-sized nonprofits in Riverbend County spend scarce staff time on back-office work and grant hunting, work alone on shared problems, and are rarely heard in local planning, so residents get fewer and less reliable services than the need requires.

**Who is affected.** About 150 member nonprofits, and through them the residents of Riverbend County who use their services.

Riverbend Nonprofit Network is a membership network with about 8 staff. It works through three program areas, each a pathway in this theory. Learning Exchange builds leadership and peer ties, Shared Services lowers overhead and widens the grant pipeline, and Civic Voice helps members take joint positions and be heard by local government and funders. The three pathways converge on shared long-term outcomes about member collaboration, then on resident-level impact. Layer 1 is change in member nonprofits, which the network influences and measures directly. Layer 2 is change for residents, which the network tracks and contributes to but does not own.

### Causal analysis

#### Structural causes

| Cause | Who is affected | What addresses it |
|-------|-----------------|-------------------|
| Members have no standing channel to speak with one voice to local government and funders. (`cause-no-standing-voice`) | Small member nonprofits | Local government and funders consult members through a coordinated channel (`oc-lt-collective-voice`); Local officials and funders engage members in planning (`oc-im-officials-engage`) |

#### Underlying causes

| Cause | Who is affected | What addresses it |
|-------|-----------------|-------------------|
| Each organization works alone, so effort and costs are duplicated. (`cause-isolation`) | — | Local government and funders consult members through a coordinated channel (`oc-lt-collective-voice`) |
| Grant funding is unpredictable and concentrated in a few funders. (`cause-funding-volatility`) | — | Set aside: Funder behavior is outside what the network can change; it only widens the pipeline. |

#### Immediate causes

| Cause | Who is affected | What addresses it |
|-------|-----------------|-------------------|
| Leaders get little management support or peer exchange. (`cause-no-peer-exchange`) | Small member nonprofits | Members turn to each other for advice and joint work (`oc-im-peer-network`) |
| Back-office work takes scarce staff time. (`cause-high-overhead`) | Small member nonprofits | Members spend less staff time and money on back-office work (`oc-im-lower-overhead`) |

#### Priority groups

- **Small member nonprofits**: Members with budgets under 500,000 and no dedicated finance or development staff.
- **Residents of low-access areas**: People who live far from services or cannot easily travel to them.

### Priorities and why

#### Outcomes

| Outcome | Equity | Comparative advantage | Evidence of what works | Weighted score | Rationale |
|---|---|---|---|---|---|
| Members turn to each other for advice and joint work (`oc-im-peer-network`) | 4 | 5 | 4 | 4.40 | Cheapest way to reach small members; strong evidence from other networks. |
| Residents get timely, good-quality services from local nonprofits (`oc-resident-access`) | 5 | 2 | — | 3.50 | The reason the network exists, but only partly in its hands. |
| Members spend less staff time and money on back-office work (`oc-im-lower-overhead`) | 4 | 3 | 3 | 3.40 | — |

#### Interventions

| Intervention | Feasibility | Reach | Weighted score | Rationale |
|---|---|---|---|---|
| Run mixed-size peer-learning cohorts (`int-peer-cohorts`) | 5 | 4 | 4.67 | Eight staff can run it with volunteer hosts. |
| Run pooled bookkeeping and payroll help (`int-pooled-backoffice`) | 3 | 3 | 3.00 | Needs a vendor partner before it can scale. |

### Assumptions and risks

| Assumption | Type | Criticality | Plausibility | Needs testing |
|------------|------|-------------|--------------|---------------|
| Member leaders can free enough time to attend cohort sessions and coaching. (`asm-time-to-attend`) | Implementation | High | Medium | Yes |
| Leaders put coaching-based skills into practice back at their organizations. (`asm-skills-applied`) _(inferred)_ | Behavioral | Critical | Medium | Yes |
| Small and large members come to see each other as peers rather than as competitors for the same funding. (`asm-peer-trust`) _(hypothesis)_ | Contextual | Medium | Medium | No |
| Vendors offer meaningful discounts and service terms at the volume the members can pool. (`asm-vendor-terms`) | External | High | Medium | Yes |
| The grants calendar stays accurate enough that members trust it. (`asm-calendar-reliable`) _(hypothesis)_ | Implementation | Medium | High | No |
| Time and money saved on overhead are reinvested in programs, leadership, and joint work rather than absorbed elsewhere. (`asm-savings-reinvested`) _(inferred)_ | Causal | Medium | Medium | No |
| Regional funders keep funding priorities open enough that better-matched applications succeed. (`asm-funders-fit`) _(inferred)_ | External | Medium | Medium | No |
| The county planning office and municipal councils stay open to a standing channel with the network. (`asm-officials-receptive`) | Partnership | High | Medium | Yes |
| Once members adopt a joint position, they keep to it in public rather than splintering. (`asm-positions-hold`) _(hypothesis)_ | Behavioral | High | Low | Yes |
| Stronger, more connected, better-funded members deliver more and better services to residents. (`asm-members-deliver`) _(inferred)_ | Causal | Critical | Medium | Yes |
| Members who benefit from the network recruit their peers. (`asm-peers-recruit`) | Behavioral | Low | High | No |
| Leadership turnover at members stays low enough for new practices to take root. (`asm-leadership-turnover`) _(hypothesis)_ | Contextual | Medium | Medium | No |

#### Risks and how they are managed

- `rsk-time-burden` Leaders drop out of cohorts because they cannot spare the time. Mitigation: Hold sessions at off-peak hours, offer remote options, and review attendance each quarter.
- `rsk-large-dominate` Larger members dominate shared agendas and smaller members withdraw. Mitigation: Rotate facilitation, cap vote weights, and give small members a seat on each steering group.
- `rsk-vendor-lockin` A single pooled-service vendor becomes a critical dependency and raises prices after the first term. Mitigation: Sign short initial terms, keep two qualified vendors, and review prices each year.
- `rsk-funder-shift` The regional funder changes priorities and withdraws support for the network or its members. Mitigation: Diversify the network's own funding and keep funders informed of results.
- `rsk-staff-capacity` With about 8 staff, the network cannot sustain three program areas if one grows faster than planned. Mitigation: Stage launches, review workload each half year, and use partners for delivery.

### What partners deliver

#### Delivered by partners

These are part of the solution and are monitored, but others answer for them.

| Outcome | Delivered by | Measured by | Status |
|---------|--------------|-------------|--------|
| Local officials and funders engage members in planning (`oc-im-officials-engage`) | Riverbend County Office of Community Planning, Millbrook Town Council | ind-officials-meetings | — |

#### Above the ceiling of accountability

These are tracked and contributed to, but not owned.

| Outcome | Delivered by | Measured by | Status |
|---------|--------------|-------------|--------|
| Residents get timely, good-quality services from local nonprofits (`oc-resident-access`) | Riverbend Nonprofit Network, Small member nonprofits, Larger member nonprofits | ind-residents-served | — |
| Residents' priorities shape local public decisions (`oc-resident-voice`) | — | ind-residents-voice | — |

### Who does what

| Intervention | Owner | Status |
|--------------|-------|--------|
| Run mixed-size peer-learning cohorts (`int-peer-cohorts`) | Riverbend Nonprofit Network | In progress |
| Offer one-to-one leadership coaching (`int-leader-coaching`) | Riverbend Nonprofit Network | In progress |
| Deliver a management skills training series (`int-training-series`) | Northfield Learning Studio | Planned |
| Run pooled bookkeeping and payroll help (`int-pooled-backoffice`) | Riverbend Nonprofit Network | In progress |
| Negotiate group purchasing agreements (`int-group-purchasing`) | Riverbend Nonprofit Network | In progress |
| Maintain a shared grants calendar (`int-grants-calendar`) | Riverbend Nonprofit Network | Completed |
| Publish plain-language policy briefings (`int-policy-briefings`) | Riverbend Nonprofit Network | Planned |
| Convene member forums to draft joint positions (`int-joint-forums`) | Riverbend Nonprofit Network | In progress |
| Set up a standing liaison channel with the county planning office (`int-liaison-channel`) | Riverbend County Office of Community Planning | In progress |
| Host a funder and member roundtable (`int-funder-roundtable`) | Lowland Regional Fund | Planned |

The division of labour across all actors is in the actors report.

### Monitoring and evaluation

#### Assumptions to test

- `asm-time-to-attend` Member leaders can free enough time to attend cohort sessions and coaching. Test: Track the share of enrolled leaders who attend at least three quarters of sessions, and ask drop-outs why they left.
- `asm-skills-applied` Leaders put coaching-based skills into practice back at their organizations. Test: Interview a sample of cohort graduates six months later and score their practice with the maturity rubric.
- `asm-vendor-terms` Vendors offer meaningful discounts and service terms at the volume the members can pool. Test: Compare pooled vendor quotes with the prices a sample of members pay now, before signing each agreement.
- `asm-officials-receptive` The county planning office and municipal councils stay open to a standing channel with the network. Test: Review the county planning calendar each quarter for processes that invited members, and ask planning staff whether the channel was used.
- `asm-positions-hold` Once members adopt a joint position, they keep to it in public rather than splintering. Test: After each joint position, record whether any member publicly took a conflicting stance within the following quarter.
- `asm-members-deliver` Stronger, more connected, better-funded members deliver more and better services to residents. Test: Compare service volumes and waits at members with high and low maturity ratings, using the annual member service-count report.

#### How often it is measured

- P1Y: 5 indicator(s)
- P3M: 17 indicator(s)
- P6M: 5 indicator(s)

#### Questions an evaluation should answer

- Does it hold that "Peer-learning cohorts and coaching sessions delivered" leads to "Members form peer relationships across sizes"? (`ln-cohorts-ties` rests on no evidence yet.)
- Does it hold that "Members form peer relationships across sizes" leads to "Members turn to each other for advice and joint work"? (`ln-ties-network` rests on no evidence yet.)
- Does it hold that "Member leaders apply improved management practices" leads to "Member nonprofits run strong, well-managed organizations"? (`ln-practice-capable` rests on no evidence yet.)
- Does it hold that "Pooled service contracts and group purchasing deals in place" leads to "Members use pooled back-office services"? (`ln-services-use` rests on no evidence yet.)
- Does it hold that "Shared grants calendar published and kept current" leads to "Members consult the shared grants calendar when planning"? (`ln-calendar-use` rests on no evidence yet.)

### Who has been consulted

| Date | Step | Method | Who | What was learned |
|------|------|--------|-----|------------------|
| 2026-02-05 | Whole | Workshop | Network staff | Staff agreed on the three program areas and on the two impact layers. Staff doubted that members would hold a joint position when funders compete. |
| 2026-02-24 | Change analysis | Interview | Small member nonprofits, Larger member nonprofits | Small members asked for help with bookkeeping before any other service. Larger members questioned the value of cohorts for them and asked for a say in the agenda. |
| 2026-03-02 | Assumptions risks | Interview | Regional funders | The funder is open to a roundtable but may change its priorities each cycle. |
| 2026-03-20 | Focus | Meeting | Member assembly | The assembly rejected a certification program and asked the network to focus on the three program areas. |
| 2026-04-14 | Measurement | Workshop | Network staff, Member representatives | Participants chose indicators that are already on sign-in sheets over new surveys. Participants asked for one resident-level number even though it is only contributed to. |
| 2026-03-05 | Change analysis | Document review | Network board | The board agreed that funding volatility is real but outside the network's reach. |
| 2026-03-12 | Focus | Interview | Residents of low-access areas, Neighborhood association volunteers | Residents named travel time and unclear eligibility as the main barriers to services. Residents asked that members share intake forms so people tell their story once. |
| 2026-03-19 | Actors | Meeting | County planning staff | Staff would use a named contact and a quarterly meeting. Staff cannot commit to inviting members into every planning process. |

### Key messages

- The change this work is for: Residents get timely, good-quality services from local nonprofits; Residents' priorities shape local public decisions.
- The organization answers for the outcome model up to: Long-term outcome and below.
- 19 outcome(s) are the organization's to deliver or share; 3 are monitored but delivered by others or beyond the ceiling.

## Logic Models

### Learning Exchange

```d2
grid-columns: 7
horizontal-gap: 90
direction: right

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
  }
  inputs: {
    style.fill: "#fff1dc"
    style.stroke: "#e08a1e"
    style.border-radius: 12
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 22
    style.bold: true
  }
}

inputs: Inputs {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  _input1: Two program staff and meeting space donated by the county library system {
    class: inputs
  }
  _input2: One program staff member and a pool of volunteer coaches {
    class: inputs
  }
  _input3: Fee-for-service agreement with the training partner {
    class: inputs
  }
}

activities: Activities {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "int-peer-cohorts": Run mixed-size peer-learning cohorts {
    class: activities
  }
  "int-leader-coaching": Offer one-to-one leadership coaching {
    class: activities
  }
  "int-training-series": Deliver a management skills training series (inferred) {
    class: activities
  }
}

outputs: Outputs {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  _output1: Cohort sessions held {
    class: outputs
  }
  _output2: Peer-pairing introductions {
    class: outputs
  }
  "oc-out-cohorts": Peer-learning cohorts and coaching sessions delivered {
    class: outputs
  }
  _output3: Coaching hours delivered {
    class: outputs
  }
  _output4: Workshops delivered {
    class: outputs
  }
  _output5: Participant materials {
    class: outputs
  }
}

short_term: Short-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-st-peer-ties": Members form peer relationships across sizes (inferred) {
    class: short_term
  }
  "oc-st-leader-skills": Leaders gain coaching-based skills and confidence {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-im-peer-network": Members turn to each other for advice and joint work (inferred) {
    class: medium_term
  }
  "oc-im-joint-positions": Members adopt and share joint positions on local issues (hypothesis) {
    class: medium_term
  }
  "oc-im-leaders-practice": Member leaders apply improved management practices {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-lt-collab": Members routinely collaborate instead of duplicating effort {
    class: long_term
  }
  "oc-lt-collective-voice": Local government and funders consult members through a coordinated channel (hypothesis) {
    class: long_term
  }
  "oc-lt-capable-members": Member nonprofits run strong, well-managed organizations {
    class: long_term
  }
}

impact: Impact {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-resident-voice": 'Residents'' priorities shape local public decisions (hypothesis)' {
    class: impact
  }
  "oc-resident-access": Residents get timely, good-quality services from local nonprofits (inferred) {
    class: impact
  }
}

inputs._input1 -> activities."int-peer-cohorts"
activities."int-peer-cohorts" -> outputs._output1
activities."int-peer-cohorts" -> outputs._output2
activities."int-peer-cohorts" -> outputs."oc-out-cohorts"
outputs._output1 -> short_term."oc-st-peer-ties"
outputs._output2 -> short_term."oc-st-peer-ties"
inputs._input2 -> activities."int-leader-coaching"
activities."int-leader-coaching" -> outputs._output3
activities."int-leader-coaching" -> outputs."oc-out-cohorts"
outputs._output3 -> short_term."oc-st-leader-skills"
inputs._input3 -> activities."int-training-series"
activities."int-training-series" -> outputs._output4
activities."int-training-series" -> outputs._output5
outputs._output4 -> short_term."oc-st-leader-skills"
outputs._output5 -> short_term."oc-st-leader-skills"
outputs."oc-out-cohorts" -> short_term."oc-st-leader-skills"
outputs."oc-out-cohorts" -> short_term."oc-st-peer-ties"
short_term."oc-st-leader-skills" -> medium_term."oc-im-leaders-practice"
short_term."oc-st-peer-ties" -> medium_term."oc-im-peer-network"
medium_term."oc-im-leaders-practice" -> long_term."oc-lt-capable-members"
medium_term."oc-im-peer-network" -> long_term."oc-lt-collab"
short_term."oc-st-peer-ties" -> medium_term."oc-im-joint-positions"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collab"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collective-voice"
long_term."oc-lt-capable-members" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-collective-voice" -> impact."oc-resident-voice"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

Feedback, not drawn:

- Members routinely collaborate instead of duplicating effort leads back to Members form peer relationships across sizes (inferred)


Peer learning and leadership coaching for member leaders.

| Stage | Elements |
|-------|----------|
| Inputs | Two program staff and meeting space donated by the county library system; One program staff member and a pool of volunteer coaches; Fee-for-service agreement with the training partner |
| Activities | Run mixed-size peer-learning cohorts; Offer one-to-one leadership coaching; Deliver a management skills training series _(inferred)_ |
| Outputs | Cohort sessions held; Peer-pairing introductions; Coaching hours delivered; Workshops delivered; Participant materials; Peer-learning cohorts and coaching sessions delivered |
| Short-term outcomes | Leaders gain coaching-based skills and confidence; Members form peer relationships across sizes _(inferred)_ |
| Medium-term outcomes | Member leaders apply improved management practices; Members turn to each other for advice and joint work _(inferred)_; Members adopt and share joint positions on local issues _(hypothesis)_ |
| Long-term outcomes | Member nonprofits run strong, well-managed organizations; Members routinely collaborate instead of duplicating effort; Local government and funders consult members through a coordinated channel _(hypothesis)_ |
| Impact | Residents get timely, good-quality services from local nonprofits _(inferred)_; Residents' priorities shape local public decisions _(hypothesis)_ |


### Shared Services

```d2
grid-columns: 7
horizontal-gap: 90
direction: right

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
  }
  inputs: {
    style.fill: "#fff1dc"
    style.stroke: "#e08a1e"
    style.border-radius: 12
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 22
    style.bold: true
  }
}

inputs: Inputs {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  _input1: One shared-services coordinator and contracted bookkeepers {
    class: inputs
  }
  _input2: Part of one staff role {
    class: inputs
  }
}

activities: Activities {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "int-pooled-backoffice": Run pooled bookkeeping and payroll help {
    class: activities
  }
  "int-group-purchasing": Negotiate group purchasing agreements {
    class: activities
  }
  "int-grants-calendar": Maintain a shared grants calendar {
    class: activities
  }
}

outputs: Outputs {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  _output1: Member accounts onboarded {
    class: outputs
  }
  _output2: Shared service desk {
    class: outputs
  }
  "oc-out-services": Pooled service contracts and group purchasing deals in place {
    class: outputs
  }
  _output3: Vendor agreements signed {
    class: outputs
  }
  "oc-out-calendar": Shared grants calendar published and kept current (inferred) {
    class: outputs
  }
  _output4: Monthly calendar update {
    class: outputs
  }
  _output5: Deadline reminders {
    class: outputs
  }
}

short_term: Short-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-st-pooled-use": Members use pooled back-office services (inferred) {
    class: short_term
  }
  "oc-st-calendar-use": Members consult the shared grants calendar when planning (hypothesis) {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-im-lower-overhead": Members spend less staff time and money on back-office work {
    class: medium_term
  }
  "oc-im-grant-pipeline": Members apply for more grants that fit their mission (inferred) {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-lt-capable-members": Member nonprofits run strong, well-managed organizations {
    class: long_term
  }
  "oc-lt-collab": Members routinely collaborate instead of duplicating effort {
    class: long_term
  }
  "oc-lt-sustainable": Members have diversified, stable funding (inferred) {
    class: long_term
  }
}

impact: Impact {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-resident-access": Residents get timely, good-quality services from local nonprofits (inferred) {
    class: impact
  }
  "oc-resident-voice": 'Residents'' priorities shape local public decisions (hypothesis)' {
    class: impact
  }
}

inputs._input1 -> activities."int-pooled-backoffice"
activities."int-pooled-backoffice" -> outputs._output1
activities."int-pooled-backoffice" -> outputs._output2
activities."int-pooled-backoffice" -> outputs."oc-out-services"
outputs._output1 -> short_term."oc-st-pooled-use"
outputs._output2 -> short_term."oc-st-pooled-use"
activities."int-group-purchasing" -> outputs._output3
activities."int-group-purchasing" -> outputs."oc-out-services"
inputs._input2 -> activities."int-grants-calendar"
activities."int-grants-calendar" -> outputs._output4
activities."int-grants-calendar" -> outputs._output5
activities."int-grants-calendar" -> outputs."oc-out-calendar"
outputs._output4 -> short_term."oc-st-calendar-use"
outputs._output5 -> short_term."oc-st-calendar-use"
outputs."oc-out-services" -> short_term."oc-st-pooled-use"
outputs."oc-out-calendar" -> short_term."oc-st-calendar-use"
short_term."oc-st-pooled-use" -> medium_term."oc-im-lower-overhead"
short_term."oc-st-calendar-use" -> medium_term."oc-im-grant-pipeline"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-capable-members"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-collab"
medium_term."oc-im-lower-overhead" -> long_term."oc-lt-sustainable"
medium_term."oc-im-grant-pipeline" -> long_term."oc-lt-sustainable"
long_term."oc-lt-capable-members" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-sustainable" -> impact."oc-resident-access"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

Pooled back-office help, group purchasing, and a shared grants calendar.

| Stage | Elements |
|-------|----------|
| Inputs | One shared-services coordinator and contracted bookkeepers; Part of one staff role |
| Activities | Run pooled bookkeeping and payroll help; Negotiate group purchasing agreements; Maintain a shared grants calendar |
| Outputs | Member accounts onboarded; Shared service desk; Vendor agreements signed; Monthly calendar update; Deadline reminders; Pooled service contracts and group purchasing deals in place; Shared grants calendar published and kept current _(inferred)_ |
| Short-term outcomes | Members use pooled back-office services _(inferred)_; Members consult the shared grants calendar when planning _(hypothesis)_ |
| Medium-term outcomes | Members spend less staff time and money on back-office work; Members apply for more grants that fit their mission _(inferred)_ |
| Long-term outcomes | Member nonprofits run strong, well-managed organizations; Members routinely collaborate instead of duplicating effort; Members have diversified, stable funding _(inferred)_ |
| Impact | Residents get timely, good-quality services from local nonprofits _(inferred)_; Residents' priorities shape local public decisions _(hypothesis)_ |


### Civic Voice

```d2
grid-columns: 6
horizontal-gap: 90
direction: right

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
  }
  stage: {
    style.fill: "transparent"
    style.stroke-width: 0
    style.font-size: 22
    style.bold: true
  }
}

activities: Activities {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "int-policy-briefings": Publish plain-language policy briefings (hypothesis) {
    class: activities
  }
  "int-joint-forums": Convene member forums to draft joint positions {
    class: activities
  }
  "int-liaison-channel": Set up a standing liaison channel with the county planning office {
    class: activities
  }
  "int-funder-roundtable": Host a funder and member roundtable {
    class: activities
  }
}

outputs: Outputs {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-out-briefings": Policy briefings and joint-position drafts produced (hypothesis) {
    class: outputs
  }
  _output1: Briefings published {
    class: outputs
  }
  _output2: Forums held {
    class: outputs
  }
  _output3: Joint position drafts {
    class: outputs
  }
  _output4: Quarterly liaison meeting {
    class: outputs
  }
  _output5: Named contacts on both sides {
    class: outputs
  }
  _output6: Roundtable held twice a year {
    class: outputs
  }
}

short_term: Short-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-st-policy-know": Members understand local policy and funding processes (inferred) {
    class: short_term
  }
  "oc-st-channel": A standing channel exists between members and the county planning office {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-im-joint-positions": Members adopt and share joint positions on local issues (hypothesis) {
    class: medium_term
  }
  "oc-im-officials-engage": Local officials and funders engage members in planning {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-lt-collab": Members routinely collaborate instead of duplicating effort {
    class: long_term
  }
  "oc-lt-collective-voice": Local government and funders consult members through a coordinated channel (hypothesis) {
    class: long_term
  }
}

impact: Impact {
  grid-columns: 1
  vertical-gap: 24
  class: stage

  "oc-resident-access": Residents get timely, good-quality services from local nonprofits (inferred) {
    class: impact
  }
  "oc-resident-voice": 'Residents'' priorities shape local public decisions (hypothesis)' {
    class: impact
  }
}

activities."int-policy-briefings" -> outputs._output1
activities."int-policy-briefings" -> outputs."oc-out-briefings"
outputs._output1 -> short_term."oc-st-policy-know"
activities."int-joint-forums" -> outputs._output2
activities."int-joint-forums" -> outputs._output3
outputs._output2 -> short_term."oc-st-policy-know"
outputs._output3 -> short_term."oc-st-policy-know"
outputs._output2 -> medium_term."oc-im-joint-positions"
outputs._output3 -> medium_term."oc-im-joint-positions"
activities."int-liaison-channel" -> outputs._output4
activities."int-liaison-channel" -> outputs._output5
outputs._output4 -> short_term."oc-st-channel"
outputs._output5 -> short_term."oc-st-channel"
activities."int-funder-roundtable" -> outputs._output6
outputs._output6 -> short_term."oc-st-channel"
outputs._output6 -> medium_term."oc-im-officials-engage"
outputs."oc-out-briefings" -> short_term."oc-st-policy-know"
short_term."oc-st-policy-know" -> medium_term."oc-im-joint-positions"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collab"
medium_term."oc-im-joint-positions" -> long_term."oc-lt-collective-voice"
medium_term."oc-im-officials-engage" -> long_term."oc-lt-collective-voice"
short_term."oc-st-channel" -> medium_term."oc-im-officials-engage"
long_term."oc-lt-collab" -> impact."oc-resident-access"
long_term."oc-lt-collective-voice" -> impact."oc-resident-voice"
long_term."oc-lt-collab" -> impact."oc-resident-voice"
```

Helping members speak together to local government and funders.

| Stage | Elements |
|-------|----------|
| Inputs | — |
| Activities | Publish plain-language policy briefings _(hypothesis)_; Convene member forums to draft joint positions; Set up a standing liaison channel with the county planning office; Host a funder and member roundtable |
| Outputs | Briefings published; Forums held; Joint position drafts; Quarterly liaison meeting; Named contacts on both sides; Roundtable held twice a year; Policy briefings and joint-position drafts produced _(hypothesis)_ |
| Short-term outcomes | Members understand local policy and funding processes _(inferred)_; A standing channel exists between members and the county planning office |
| Medium-term outcomes | Members adopt and share joint positions on local issues _(hypothesis)_; Local officials and funders engage members in planning |
| Long-term outcomes | Members routinely collaborate instead of duplicating effort; Local government and funders consult members through a coordinated channel _(hypothesis)_ |
| Impact | Residents get timely, good-quality services from local nonprofits _(inferred)_; Residents' priorities shape local public decisions _(hypothesis)_ |

## VSEM: Vision, Strategy, Execution, Metrics

### Vision

A county where every resident can count on a strong local nonprofit when they need one.

**Mission:** Strengthen the nonprofits of Riverbend County so that residents get timely, good-quality services.

### Strategy

#### 1. Learning Exchange

**Owner:** Riverbend Nonprofit Network

Peer learning and leadership coaching for member leaders.

**Resources**

- Program staff: 3 staff
- Volunteer coaches: A pool of about 15 coaches
- Training partner: Fee-for-service agreement

**Execution**

| Work | Kind | Owner | Timing | Status |
|------|------|-------|--------|--------|
| Run mixed-size peer-learning cohorts<br>Cohorts of 10 to 12 leaders meet monthly for a year, mixing small and large members. | program | Riverbend Nonprofit Network | from 2026-02-01 | In Progress |
| Offer one-to-one leadership coaching | program | Riverbend Nonprofit Network |  | In Progress |
| Deliver a management skills training series | initiative | Northfield Learning Studio | 2026-09-01 to 2027-05-31 | Planned |

**Metrics**

| Metric | Accountable | Baseline | Target | Current | Due | Status |
|--------|-------------|----------|--------|---------|-----|--------|
| Leadership practice maturity<br>North Star: Median level of participating members on the four-step leadership practice rubric, rated by a coach and the member together. | Riverbend Nonprofit Network | Emerging | Established |  | 2028-12-31 |  |
| Participant retention: Learning Exchange<br>Guardrail: Share of enrolled cohort and coaching participants still active at the end of the program year. | Riverbend Nonprofit Network | 78 percent | 88 percent |  | 2028-12-31 |  |

#### 2. Shared Services

**Owner:** Riverbend Nonprofit Network

Pooled back-office help, group purchasing, and a shared grants calendar.

**Resources**

- Shared-services coordinator: 1 staff
- Contracted bookkeepers: Paid from members' service fees

**Execution**

| Work | Kind | Owner | Timing | Status |
|------|------|-------|--------|--------|
| Run pooled bookkeeping and payroll help | initiative | Riverbend Nonprofit Network |  | In Progress |
| Negotiate group purchasing agreements | initiative | Riverbend Nonprofit Network |  | In Progress |
| Maintain a shared grants calendar | initiative | Riverbend Nonprofit Network | 2026-02-01 to 2026-06-30 | Completed |

**Metrics**

| Metric | Accountable | Baseline | Target | Current | Due | Status |
|--------|-------------|----------|--------|---------|-----|--------|
| Back-office share of operating budget<br>North Star: Share of a member's operating budget spent on bookkeeping, payroll, insurance, and purchased supplies, averaged across members using pooled services. | Riverbend Nonprofit Network | 21.5 percent | 17 percent |  | 2028-12-31 |  |
| Participant retention: Shared Services<br>Guardrail: Share of member organizations on a pooled service at the start of the year that still use it at year end. | Riverbend Nonprofit Network | 85 percent | 92 percent |  | 2028-12-31 |  |

#### 3. Civic Voice

**Owner:** Riverbend Nonprofit Network

Helping members speak together to local government and funders.

**Resources**

- Policy and organizing staff: 1.5 staff
- Forum costs: Venue and childcare for four forums a year

**Execution**

| Work | Kind | Owner | Timing | Status |
|------|------|-------|--------|--------|
| Publish plain-language policy briefings | initiative | Riverbend Nonprofit Network |  | Planned |
| Convene member forums to draft joint positions | initiative | Riverbend Nonprofit Network |  | In Progress |
| Set up a standing liaison channel with the county planning office | initiative | Riverbend County Office of Community Planning |  | In Progress |
| Host a funder and member roundtable | initiative | Lowland Regional Fund | 2027-03-01 to 2027-03-31 | Planned |

**Metrics**

| Metric | Accountable | Baseline | Target | Current | Due | Status |
|--------|-------------|----------|--------|---------|-----|--------|
| Joint positions adopted by members<br>North Star: Count of positions on local issues formally adopted by at least ten member organizations in the last twelve months. | Riverbend Nonprofit Network | 1 positions | 6 positions |  | 2028-12-31 |  |
| Member organizations reached: Civic Voice<br>Guardrail: Distinct member organizations attending a forum or receiving a briefing in the last twelve months. | Riverbend Nonprofit Network | 40 member organizations | 80 member organizations |  | 2028-12-31 |  |

### Metrics for the Vision

| Metric | Accountable | Baseline | Target | Current | Due | Status |
|--------|-------------|----------|--------|---------|-----|--------|
| Residents served by member organizations<br>Residents who used at least one service of a member organization in the year, as reported by members. | Riverbend Nonprofit Network | 61000 residents | 70000 residents |  | 2028-12-31 |  |
| Residents who say local organizations help them be heard<br>Share of surveyed residents who agree that local organizations help them influence local decisions. | Riverbend Nonprofit Network | 31 percent | 42 percent |  | 2028-12-31 |  |

## Monitoring, Evaluation, and Learning Plan

### Results Framework

- **Impact:** Residents get timely, good-quality services from local nonprofits (`oc-resident-access`) — Layer 2 impact. People in the county who need community services can reach them without long waits.
  - **Long-term outcome:** Member nonprofits run strong, well-managed organizations (`oc-lt-capable-members`) — Layer 1. Members show sound leadership practice with low overhead.
    - **Medium-term outcome:** Member leaders apply improved management practices (`oc-im-leaders-practice`)
      - **Short-term outcome:** Leaders gain coaching-based skills and confidence (`oc-st-leader-skills`)
        - **Output:** Peer-learning cohorts and coaching sessions delivered (`oc-out-cohorts`)
    - **Medium-term outcome:** Members spend less staff time and money on back-office work (`oc-im-lower-overhead`)
      - **Short-term outcome:** Members use pooled back-office services (`oc-st-pooled-use`)
        - **Output:** Pooled service contracts and group purchasing deals in place (`oc-out-services`)
  - **Long-term outcome:** Members routinely collaborate instead of duplicating effort (`oc-lt-collab`) — Layer 1. Members share contracts, referrals, and joint work. Three pathways converge here.
    - **Medium-term outcome:** Members turn to each other for advice and joint work (`oc-im-peer-network`)
      - **Short-term outcome:** Members form peer relationships across sizes (`oc-st-peer-ties`)
        - **Output:** Peer-learning cohorts and coaching sessions delivered (`oc-out-cohorts`), shown above
    - **Medium-term outcome:** Members spend less staff time and money on back-office work (`oc-im-lower-overhead`), shown above
    - **Medium-term outcome:** Members adopt and share joint positions on local issues (`oc-im-joint-positions`)
      - **Short-term outcome:** Members form peer relationships across sizes (`oc-st-peer-ties`), shown above
      - **Short-term outcome:** Members understand local policy and funding processes (`oc-st-policy-know`)
        - **Output:** Policy briefings and joint-position drafts produced (`oc-out-briefings`)
  - **Long-term outcome:** Members have diversified, stable funding (`oc-lt-sustainable`) — Layer 1. Members hold several funding sources and are less exposed to the loss of one.
    - **Medium-term outcome:** Members spend less staff time and money on back-office work (`oc-im-lower-overhead`), shown above
    - **Medium-term outcome:** Members apply for more grants that fit their mission (`oc-im-grant-pipeline`)
      - **Short-term outcome:** Members consult the shared grants calendar when planning (`oc-st-calendar-use`)
        - **Output:** Shared grants calendar published and kept current (`oc-out-calendar`)
- **Impact:** Residents' priorities shape local public decisions (`oc-resident-voice`) — Layer 2 impact. Planning and funding decisions reflect what residents say they need.
  - **Long-term outcome:** Members routinely collaborate instead of duplicating effort (`oc-lt-collab`), shown above
  - **Long-term outcome:** Local government and funders consult members through a coordinated channel (`oc-lt-collective-voice`) — Layer 1. Officials and funders ask the network, not each member, for input on local issues.
    - **Medium-term outcome:** Members adopt and share joint positions on local issues (`oc-im-joint-positions`), shown above
    - **Medium-term outcome:** Local officials and funders engage members in planning (`oc-im-officials-engage`)
      - **Short-term outcome:** A standing channel exists between members and the county planning office (`oc-st-channel`)

### Activities

| Activity | Produces | Lead | Timing | Status |
|----------|----------|------|--------|--------|
| Run mixed-size peer-learning cohorts | Peer-learning cohorts and coaching sessions delivered; Members form peer relationships across sizes | Riverbend Nonprofit Network | from 2026-02-01 | In Progress |
| Offer one-to-one leadership coaching | Peer-learning cohorts and coaching sessions delivered; Leaders gain coaching-based skills and confidence | Riverbend Nonprofit Network |  | In Progress |
| Deliver a management skills training series | Leaders gain coaching-based skills and confidence | Northfield Learning Studio | 2026-09-01 to 2027-05-31 | Planned |
| Run pooled bookkeeping and payroll help | Pooled service contracts and group purchasing deals in place; Members use pooled back-office services | Riverbend Nonprofit Network |  | In Progress |
| Negotiate group purchasing agreements | Pooled service contracts and group purchasing deals in place | Riverbend Nonprofit Network |  | In Progress |
| Maintain a shared grants calendar | Shared grants calendar published and kept current; Members consult the shared grants calendar when planning | Riverbend Nonprofit Network | 2026-02-01 to 2026-06-30 | Completed |
| Publish plain-language policy briefings | Policy briefings and joint-position drafts produced; Members understand local policy and funding processes | Riverbend Nonprofit Network |  | Planned |
| Convene member forums to draft joint positions | Members understand local policy and funding processes; Members adopt and share joint positions on local issues | Riverbend Nonprofit Network |  | In Progress |
| Set up a standing liaison channel with the county planning office | A standing channel exists between members and the county planning office | Riverbend County Office of Community Planning |  | In Progress |
| Host a funder and member roundtable | A standing channel exists between members and the county planning office; Local officials and funders engage members in planning | Lowland Regional Fund | 2027-03-01 to 2027-03-31 | Planned |

### Indicators

| Indicator | Measures | Baseline | Target | Frequency | Source | Responsible |
|-----------|----------|----------|--------|-----------|--------|-------------|
| Attendee hours: Learning Exchange | Run mixed-size peer-learning cohorts | 3200 hours | 5400 hours by 2028-12-31 | Quarterly | Attendance sheets | Riverbend Nonprofit Network |
| Attendee hours: Shared Services | Run pooled bookkeeping and payroll help | 600 hours | 1600 hours by 2028-12-31 | Quarterly | Attendance sheets | Riverbend Nonprofit Network |
| Attendee hours: Civic Voice | Convene member forums to draft joint positions | 350 hours | 1100 hours by 2028-12-31 | Quarterly | Attendance sheets | Riverbend Nonprofit Network |
| Attendee hours across all three program areas | Members form peer relationships across sizes | 4150 hours | 8100 hours by 2028-12-31 | Quarterly | Roll-up of program-area indicators | Riverbend Nonprofit Network |
| Member organizations reached: Learning Exchange | Peer-learning cohorts and coaching sessions delivered | 48 member organizations | 80 member organizations by 2028-12-31 | Quarterly | Attendance sheets matched to the member list | Riverbend Nonprofit Network |
| Member organizations reached: Shared Services | Members use pooled back-office services | 22 member organizations | 60 member organizations by 2028-12-31 | Quarterly | Shared service desk records | Riverbend Nonprofit Network |
| Member organizations reached: Civic Voice | Members understand local policy and funding processes | 40 member organizations | 80 member organizations by 2028-12-31 | Quarterly | Attendance sheets and briefing mailing list | Riverbend Nonprofit Network |
| Member organizations reached across the whole network | Members form peer relationships across sizes | 71 member organizations | 110 member organizations by 2028-12-31 | Quarterly | Attendance sheets, service desk records, briefing list | Riverbend Nonprofit Network |
| Participant retention: Learning Exchange | Leaders gain coaching-based skills and confidence | 78 percent | 88 percent by 2028-12-31 | Every six months | Cohort rosters | Riverbend Nonprofit Network |
| Participant retention: Shared Services | Members use pooled back-office services | 85 percent | 92 percent by 2028-12-31 | Every six months | Service desk records | Riverbend Nonprofit Network |
| Leadership practice maturity | Member leaders apply improved management practices | Emerging | Established by 2028-12-31 | Every six months | Coach and member rubric ratings | Riverbend Nonprofit Network |
| Member collaboration maturity | Members routinely collaborate instead of duplicating effort | Parallel | Cooperating by 2028-12-31 | Yearly | Annual member interviews | Riverbend Nonprofit Network |
| Level of members' voice with local government and funders | Local government and funders consult members through a coordinated channel | Not heard | Co-designing by 2028-12-31 | Every six months | Staff assessment checked with county office and regional funder | Riverbend Nonprofit Network |
| Back-office share of operating budget | Members spend less staff time and money on back-office work | 21.5 percent | 17 percent by 2028-12-31 | Yearly | Annual member cost survey | Riverbend Nonprofit Network |
| Grant applications submitted by members using the calendar | Members apply for more grants that fit their mission | 120 applications | 200 applications by 2028-12-31 | Quarterly | Quarterly member check-in | Riverbend Nonprofit Network |
| Members with three or more distinct funding sources | Members have diversified, stable funding | 38 percent | 55 percent by 2028-12-31 | Yearly | Annual member survey | Riverbend Nonprofit Network |
| Joint positions adopted by members | Members adopt and share joint positions on local issues | 1 positions | 6 positions by 2028-12-31 | Quarterly | Forum minutes and sign-on lists | Riverbend Nonprofit Network |
| Public planning processes that invited members | Local officials and funders engage members in planning | 2 processes | 8 processes by 2028-12-31 | Quarterly | County planning calendar and municipal agendas | Riverbend County Office of Community Planning |
| Residents served by member organizations | Residents get timely, good-quality services from local nonprofits | 61000 residents | 70000 residents by 2028-12-31 | Yearly | Annual member service-count report | Riverbend Nonprofit Network |
| Residents who say local organizations help them be heard | Residents' priorities shape local public decisions | 31 percent | 42 percent by 2028-12-31 | Yearly | County resident survey | Riverbend Nonprofit Network |
| Cohort attendance rate | Member leaders can free enough time to attend cohort sessions and coaching. | 71 percent | 80 percent by 2027-12-31 | Quarterly | Cohort rosters | Riverbend Nonprofit Network |
| Stories of members helping each other | Members form peer relationships across sizes | 0 | 12 by 2027-12-31 | Quarterly | Story log | Riverbend Nonprofit Network |
| Pooled contracts and purchasing deals in place | Pooled service contracts and group purchasing deals in place | 2 deals | 8 deals by 2027-12-31 | Every six months | Contract file | Riverbend Nonprofit Network |
| Months the grants calendar was updated on schedule | Shared grants calendar published and kept current | 6 months | 12 months by 2027-12-31 | Quarterly | Calendar edit log | Riverbend Nonprofit Network |
| Policy briefings and joint-position drafts published | Policy briefings and joint-position drafts produced | 3 documents | 10 documents by 2027-12-31 | Quarterly | Civic Voice publication list | Riverbend Nonprofit Network |
| Enrollments: Learning Exchange | Peer-learning cohorts and coaching sessions delivered | 120 participants | 170 participants by 2027-12-31 | Quarterly | Enrollment sheets | Riverbend Nonprofit Network |
| Enrollments: Shared Services | Run pooled bookkeeping and payroll help | 60 participants | 100 participants by 2027-12-31 | Quarterly | Enrollment sheets | Riverbend Nonprofit Network |

### Indicator Reference Sheets

#### Attendee hours: Learning Exchange

| | |
|---|---|
| **ID** | `ind-hours-learning` |
| **Measures** | Run mixed-size peer-learning cohorts |
| **Definition** | Trailing twelve-month attendee hours at cohort sessions, coaching, and trainings. |
| **Type** | quantitative |
| **Unit** | hours |
| **Direction** | increase |
| **Baseline** | 3200 hours (2026-01-01); source: 2025 attendance sheets |
| **Targets** | 4000 hours by 2026-12-31; 4800 hours by 2027-12-31; 5400 hours by 2028-12-31 |
| **Data source** | Attendance sheets |
| **Collection method** | Sum of attendees times session length, from sign-in sheets |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Attendee hours: Shared Services

| | |
|---|---|
| **ID** | `ind-hours-shared` |
| **Measures** | Run pooled bookkeeping and payroll help |
| **Definition** | Trailing twelve-month attendee hours at onboarding sessions, vendor briefings, and service-desk workshops. |
| **Type** | quantitative |
| **Unit** | hours |
| **Direction** | increase |
| **Baseline** | 600 hours (2026-01-01); source: 2025 attendance sheets |
| **Targets** | 900 hours by 2026-12-31; 1300 hours by 2027-12-31; 1600 hours by 2028-12-31 |
| **Data source** | Attendance sheets |
| **Collection method** | Sum of attendees times session length, from sign-in sheets |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Attendee hours: Civic Voice

| | |
|---|---|
| **ID** | `ind-hours-civic` |
| **Measures** | Convene member forums to draft joint positions |
| **Definition** | Trailing twelve-month attendee hours at member forums and policy briefings. |
| **Type** | quantitative |
| **Unit** | hours |
| **Direction** | increase |
| **Baseline** | 350 hours (2026-01-01); source: 2025 attendance sheets |
| **Targets** | 600 hours by 2026-12-31; 900 hours by 2027-12-31; 1100 hours by 2028-12-31 |
| **Data source** | Attendance sheets |
| **Collection method** | Sum of attendees times session length, from sign-in sheets |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Attendee hours across all three program areas

| | |
|---|---|
| **ID** | `ind-hours-total` |
| **Measures** | Members form peer relationships across sizes |
| **Definition** | Network-wide attendee hours, a proxy for the time members spend together. Rolls up the three program-area indicators because hours can be added. |
| **Type** | quantitative |
| **Unit** | hours |
| **Direction** | increase |
| **Baseline** | 4150 hours (2026-01-01); source: Sum of the 2025 program-area figures |
| **Targets** | 5500 hours by 2026-12-31; 7000 hours by 2027-12-31; 8100 hours by 2028-12-31 |
| **Data source** | Roll-up of program-area indicators |
| **Collection method** | Sum of the three program-area values |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Member organizations reached: Learning Exchange

| | |
|---|---|
| **ID** | `ind-reach-learning` |
| **Measures** | Peer-learning cohorts and coaching sessions delivered |
| **Definition** | Distinct member organizations with at least one leader in a cohort, coaching, or training in the last twelve months. |
| **Type** | quantitative |
| **Unit** | member organizations |
| **Direction** | increase |
| **Disaggregation** | member size |
| **Baseline** | 48 member organizations (2026-01-01); source: 2025 attendance sheets |
| **Targets** | 60 member organizations by 2026-12-31; 72 member organizations by 2027-12-31; 80 member organizations by 2028-12-31 |
| **Data source** | Attendance sheets matched to the member list |
| **Collection method** | Count distinct member organization ids |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Member organizations reached: Shared Services

| | |
|---|---|
| **ID** | `ind-reach-shared` |
| **Measures** | Members use pooled back-office services |
| **Definition** | Distinct member organizations using at least one pooled service in the last twelve months. |
| **Type** | quantitative |
| **Unit** | member organizations |
| **Direction** | increase |
| **Baseline** | 22 member organizations (2026-01-01); source: Service desk records |
| **Targets** | 35 member organizations by 2026-12-31; 50 member organizations by 2027-12-31; 60 member organizations by 2028-12-31 |
| **Data source** | Shared service desk records |
| **Collection method** | Count distinct member organization ids |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Member organizations reached: Civic Voice

| | |
|---|---|
| **ID** | `ind-reach-civic` |
| **Measures** | Members understand local policy and funding processes |
| **Definition** | Distinct member organizations attending a forum or receiving a briefing in the last twelve months. |
| **Type** | quantitative |
| **Unit** | member organizations |
| **Direction** | increase |
| **Baseline** | 40 member organizations (2026-01-01); source: 2025 attendance sheets |
| **Targets** | 55 member organizations by 2026-12-31; 70 member organizations by 2027-12-31; 80 member organizations by 2028-12-31 |
| **Data source** | Attendance sheets and briefing mailing list |
| **Collection method** | Count distinct member organization ids |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Member organizations reached across the whole network

| | |
|---|---|
| **ID** | `ind-reach-network` |
| **Measures** | Members form peer relationships across sizes |
| **Definition** | Distinct member organizations reached by any program area. This is counted directly, not summed from the program-area indicators, because the same member often appears in more than one of them. |
| **Type** | quantitative |
| **Unit** | member organizations |
| **Direction** | increase |
| **Disaggregation** | member size |
| **Baseline** | 71 member organizations (2026-01-01); source: 2025 attendance sheets and service desk records |
| **Targets** | 85 member organizations by 2026-12-31; 100 member organizations by 2027-12-31; 110 member organizations by 2028-12-31 |
| **Data source** | Attendance sheets, service desk records, briefing list |
| **Collection method** | Count distinct member organization ids across all sources |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Not equal to the sum of the program-area counts, which double count members in more than one program. |

#### Participant retention: Learning Exchange

| | |
|---|---|
| **ID** | `ind-retention-learning` |
| **Measures** | Leaders gain coaching-based skills and confidence |
| **Definition** | Share of enrolled cohort and coaching participants still active at the end of the program year. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | increase |
| **Baseline** | 78 percent (2026-01-01); source: 2025 cohort records |
| **Targets** | 82 percent by 2026-12-31; 85 percent by 2027-12-31; 88 percent by 2028-12-31 |
| **Data source** | Cohort rosters |
| **Collection method** | Active participants at year end divided by enrolled participants |
| **Frequency** | Every six months |
| **Responsible** | Riverbend Nonprofit Network |

#### Participant retention: Shared Services

| | |
|---|---|
| **ID** | `ind-retention-shared` |
| **Measures** | Members use pooled back-office services |
| **Definition** | Share of member organizations on a pooled service at the start of the year that still use it at year end. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | increase |
| **Baseline** | 85 percent (2026-01-01); source: Service desk records |
| **Targets** | 88 percent by 2026-12-31; 90 percent by 2027-12-31; 92 percent by 2028-12-31 |
| **Data source** | Service desk records |
| **Collection method** | Organizations still active at year end divided by organizations active at year start |
| **Frequency** | Every six months |
| **Responsible** | Riverbend Nonprofit Network |

#### Leadership practice maturity

| | |
|---|---|
| **ID** | `ind-practice-maturity` |
| **Measures** | Member leaders apply improved management practices |
| **Definition** | Median level of participating members on the four-step leadership practice rubric, rated by a coach and the member together. |
| **Type** | qualitative |
| **Direction** | increase |
| **Baseline** | Emerging (2026-01-01); source: Rubric ratings of the first cohort |
| **Targets** | Emerging by 2026-12-31; Established by 2027-12-31; Established by 2028-12-31 |
| **Data source** | Coach and member rubric ratings |
| **Collection method** | Rubric rating at enrollment and at six and twelve months |
| **Frequency** | Every six months |
| **Responsible** | Riverbend Nonprofit Network |
| **Data quality checks** | A second coach rates one in ten members, and ratings more than one level apart are discussed. |
| **Known limitations** | Rated by coaches and members, so it reflects self-assessment as well as observed practice. |
| **Quality last assessed** | 2026-07-15 |

#### Member collaboration maturity

| | |
|---|---|
| **ID** | `ind-collab-maturity` |
| **Measures** | Members routinely collaborate instead of duplicating effort |
| **Definition** | Typical level of member collaboration on a four-step scale, from members working in parallel to co-delivering services. |
| **Type** | qualitative |
| **Direction** | increase |
| **Baseline** | Parallel (2026-01-01); source: Member interviews |
| **Targets** | Coordinating by 2027-12-31; Cooperating by 2028-12-31 |
| **Data source** | Annual member interviews |
| **Collection method** | Interviewers rate each member against the scale and report the most common level |
| **Frequency** | Yearly |
| **Responsible** | Riverbend Nonprofit Network |

#### Level of members' voice with local government and funders

| | |
|---|---|
| **ID** | `ind-voice-level` |
| **Measures** | Local government and funders consult members through a coordinated channel |
| **Definition** | How officials and funders treat the members' coordinated channel, rated by network staff and checked with the officials concerned. |
| **Type** | qualitative |
| **Direction** | increase |
| **Baseline** | Not heard (2026-01-01); source: Staff assessment |
| **Targets** | Heard by 2026-12-31; Consulted by 2027-12-31; Co-designing by 2028-12-31 |
| **Data source** | Staff assessment checked with county office and regional funder |
| **Collection method** | Annual rating with a short call to each official concerned |
| **Frequency** | Every six months |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Depends on officials' own account, which may be generous or guarded. |

#### Back-office share of operating budget

| | |
|---|---|
| **ID** | `ind-overhead-share` |
| **Measures** | Members spend less staff time and money on back-office work |
| **Definition** | Share of a member's operating budget spent on bookkeeping, payroll, insurance, and purchased supplies, averaged across members using pooled services. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | decrease |
| **Baseline** | 21.5 percent (2026-01-01); source: 2025 member cost survey |
| **Targets** | 20 percent by 2026-12-31; 18.5 percent by 2027-12-31; 17 percent by 2028-12-31 |
| **Data source** | Annual member cost survey |
| **Collection method** | Members report back-office spending and total budget |
| **Frequency** | Yearly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Self-reported by members; other cost changes can also move the figure. |

#### Grant applications submitted by members using the calendar

| | |
|---|---|
| **ID** | `ind-grant-applications` |
| **Measures** | Members apply for more grants that fit their mission |
| **Definition** | Count of applications submitted in the last twelve months by members who consulted the calendar. |
| **Type** | quantitative |
| **Unit** | applications |
| **Direction** | increase |
| **Baseline** | 120 applications (2026-01-01); source: Member survey |
| **Targets** | 140 applications by 2026-12-31; 170 applications by 2027-12-31; 200 applications by 2028-12-31 |
| **Data source** | Quarterly member check-in |
| **Collection method** | Short form listing applications submitted |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Members with three or more distinct funding sources

| | |
|---|---|
| **ID** | `ind-funding-mix` |
| **Measures** | Members have diversified, stable funding |
| **Definition** | Share of responding members whose revenue comes from at least three separate funders or income types. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | increase |
| **Baseline** | 38 percent (2026-01-01); source: 2025 member survey |
| **Targets** | 42 percent by 2026-12-31; 48 percent by 2027-12-31; 55 percent by 2028-12-31 |
| **Data source** | Annual member survey |
| **Collection method** | Survey question on funding sources |
| **Frequency** | Yearly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Based on responding members only; response rates differ by member size. |

#### Joint positions adopted by members

| | |
|---|---|
| **ID** | `ind-joint-positions` |
| **Measures** | Members adopt and share joint positions on local issues |
| **Definition** | Count of positions on local issues formally adopted by at least ten member organizations in the last twelve months. |
| **Type** | quantitative |
| **Unit** | positions |
| **Direction** | increase |
| **Baseline** | 1 positions (2026-01-01); source: Forum minutes |
| **Targets** | 3 positions by 2026-12-31; 5 positions by 2027-12-31; 6 positions by 2028-12-31 |
| **Data source** | Forum minutes and sign-on lists |
| **Collection method** | Count adopted positions meeting the sign-on threshold |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Public planning processes that invited members

| | |
|---|---|
| **ID** | `ind-officials-meetings` |
| **Measures** | Local officials and funders engage members in planning |
| **Definition** | Count of county or municipal planning processes in the last twelve months that formally invited member representatives. |
| **Type** | quantitative |
| **Unit** | processes |
| **Direction** | increase |
| **Baseline** | 2 processes (2026-01-01); source: County planning calendar |
| **Targets** | 4 processes by 2026-12-31; 6 processes by 2027-12-31; 8 processes by 2028-12-31 |
| **Data source** | County planning calendar and municipal agendas |
| **Collection method** | Review of published calendars and agendas |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend County Office of Community Planning |

#### Residents served by member organizations

| | |
|---|---|
| **ID** | `ind-residents-served` |
| **Measures** | Residents get timely, good-quality services from local nonprofits |
| **Definition** | Residents who used at least one service of a member organization in the year, as reported by members. |
| **Type** | quantitative |
| **Unit** | residents |
| **Direction** | increase |
| **Disaggregation** | neighborhood, age group |
| **Baseline** | 61000 residents (2026-01-01); source: 2025 member service-count report |
| **Targets** | 63000 residents by 2026-12-31; 66000 residents by 2027-12-31; 70000 residents by 2028-12-31 |
| **Data source** | Annual member service-count report |
| **Collection method** | Members submit counts; the network removes obvious duplicates but cannot match individuals across members |
| **Frequency** | Yearly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Reflects members' reported reach; Riverbend Nonprofit Network cannot attribute resident outcomes to its own work.; Residents who use several members are counted more than once. |

#### Residents who say local organizations help them be heard

| | |
|---|---|
| **ID** | `ind-residents-voice` |
| **Measures** | Residents' priorities shape local public decisions |
| **Definition** | Share of surveyed residents who agree that local organizations help them influence local decisions. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | increase |
| **Disaggregation** | neighborhood |
| **Baseline** | 31 percent (2026-01-01); source: 2025 county resident survey |
| **Targets** | 34 percent by 2026-12-31; 38 percent by 2027-12-31; 42 percent by 2028-12-31 |
| **Data source** | County resident survey |
| **Collection method** | Published results of the annual county resident survey |
| **Frequency** | Yearly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Reflects residents' general view of local organizations; Riverbend Nonprofit Network cannot attribute changes to its own work. |

#### Cohort attendance rate

| | |
|---|---|
| **ID** | `ind-attendance-rate` |
| **Measures** | Member leaders can free enough time to attend cohort sessions and coaching. |
| **Definition** | Share of enrolled leaders who attend at least three quarters of their cohort sessions. |
| **Type** | quantitative |
| **Unit** | percent |
| **Direction** | increase |
| **Baseline** | 71 percent (2026-01-01); source: 2025 cohort rosters |
| **Targets** | 75 percent by 2026-12-31; 80 percent by 2027-12-31 |
| **Data source** | Cohort rosters |
| **Collection method** | Leaders meeting the attendance bar divided by enrolled leaders |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Stories of members helping each other

| | |
|---|---|
| **ID** | `ind-story-peer-help` |
| **Measures** | Members form peer relationships across sizes |
| **Definition** | Stories coded with the peer-help theme in the period. |
| **Type** | quantitative |
| **Baseline** | 0 (2026-01-01); source: No stories were collected before 2026 |
| **Targets** | 8 by 2026-12-31; 12 by 2027-12-31 |
| **Data source** | Story log |
| **Collection method** | interview |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |
| **Known limitations** | Counts stories collected, not every instance of peer help. |

#### Pooled contracts and purchasing deals in place

| | |
|---|---|
| **ID** | `ind-service-deals` |
| **Measures** | Pooled service contracts and group purchasing deals in place |
| **Type** | quantitative |
| **Unit** | deals |
| **Direction** | increase |
| **Baseline** | 2 deals (2026-01-01); source: 2025 contract file |
| **Targets** | 5 deals by 2026-12-31; 8 deals by 2027-12-31 |
| **Data source** | Contract file |
| **Frequency** | Every six months |
| **Responsible** | Riverbend Nonprofit Network |

#### Months the grants calendar was updated on schedule

| | |
|---|---|
| **ID** | `ind-calendar-current` |
| **Measures** | Shared grants calendar published and kept current |
| **Type** | quantitative |
| **Unit** | months |
| **Direction** | increase |
| **Baseline** | 6 months (2026-01-01); source: 2025 calendar edit log |
| **Targets** | 12 months by 2026-12-31; 12 months by 2027-12-31 |
| **Data source** | Calendar edit log |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Policy briefings and joint-position drafts published

| | |
|---|---|
| **ID** | `ind-briefings` |
| **Measures** | Policy briefings and joint-position drafts produced |
| **Type** | quantitative |
| **Unit** | documents |
| **Direction** | increase |
| **Baseline** | 3 documents (2026-01-01); source: 2025 publication list |
| **Targets** | 8 documents by 2026-12-31; 10 documents by 2027-12-31 |
| **Data source** | Civic Voice publication list |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Enrollments: Learning Exchange

| | |
|---|---|
| **ID** | `ind-enroll-learning` |
| **Measures** | Peer-learning cohorts and coaching sessions delivered |
| **Definition** | Participants enrolled in the period, the weight for the retention rate. |
| **Type** | quantitative |
| **Unit** | participants |
| **Direction** | increase |
| **Baseline** | 120 participants (2026-01-01); source: 2025 enrollment sheets |
| **Targets** | 150 participants by 2026-12-31; 170 participants by 2027-12-31 |
| **Data source** | Enrollment sheets |
| **Collection method** | passive |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

#### Enrollments: Shared Services

| | |
|---|---|
| **ID** | `ind-enroll-shared` |
| **Measures** | Run pooled bookkeeping and payroll help |
| **Definition** | Participants enrolled in the period, the weight for the retention rate. |
| **Type** | quantitative |
| **Unit** | participants |
| **Direction** | increase |
| **Baseline** | 60 participants (2026-01-01); source: 2025 enrollment sheets |
| **Targets** | 80 participants by 2026-12-31; 100 participants by 2027-12-31 |
| **Data source** | Enrollment sheets |
| **Collection method** | passive |
| **Frequency** | Quarterly |
| **Responsible** | Riverbend Nonprofit Network |

### Assumptions to Monitor

| Assumption | Results that depend on it | How it is monitored |
|------------|---------------------------|---------------------|
| Member leaders can free enough time to attend cohort sessions and coaching. | Leaders gain coaching-based skills and confidence; Members form peer relationships across sizes | Track the share of enrolled leaders who attend at least three quarters of sessions, and ask drop-outs why they left. Indicators: Cohort attendance rate |
| Leaders put coaching-based skills into practice back at their organizations. | Member leaders apply improved management practices | Interview a sample of cohort graduates six months later and score their practice with the maturity rubric. Indicators: Leadership practice maturity |
| Small and large members come to see each other as peers rather than as competitors for the same funding. | Members turn to each other for advice and joint work; Members routinely collaborate instead of duplicating effort; Members adopt and share joint positions on local issues |  |
| Vendors offer meaningful discounts and service terms at the volume the members can pool. | Members use pooled back-office services; Members spend less staff time and money on back-office work | Compare pooled vendor quotes with the prices a sample of members pay now, before signing each agreement. Indicators: Back-office share of operating budget |
| The grants calendar stays accurate enough that members trust it. | Members consult the shared grants calendar when planning; Members apply for more grants that fit their mission |  |
| Time and money saved on overhead are reinvested in programs, leadership, and joint work rather than absorbed elsewhere. | Member nonprofits run strong, well-managed organizations; Members routinely collaborate instead of duplicating effort; Members have diversified, stable funding |  |
| Regional funders keep funding priorities open enough that better-matched applications succeed. | Members have diversified, stable funding |  |
| The county planning office and municipal councils stay open to a standing channel with the network. | Local officials and funders engage members in planning; Local government and funders consult members through a coordinated channel; Residents' priorities shape local public decisions | Review the county planning calendar each quarter for processes that invited members, and ask planning staff whether the channel was used. Indicators: Public planning processes that invited members |
| Once members adopt a joint position, they keep to it in public rather than splintering. | Members adopt and share joint positions on local issues; Members routinely collaborate instead of duplicating effort; Local government and funders consult members through a coordinated channel; Residents' priorities shape local public decisions | After each joint position, record whether any member publicly took a conflicting stance within the following quarter. Indicators: Joint positions adopted by members |
| Stronger, more connected, better-funded members deliver more and better services to residents. | Residents get timely, good-quality services from local nonprofits | Compare service volumes and waits at members with high and low maturity ratings, using the annual member service-count report. Indicators: Residents served by member organizations |
| Members who benefit from the network recruit their peers. | Members form peer relationships across sizes |  |
| Leadership turnover at members stays low enough for new practices to take root. | Member nonprofits run strong, well-managed organizations |  |

### Evaluation Plan

| Evaluation | Type | Timing | Evaluator | Lead | Status |
|------------|------|--------|-----------|------|--------|
| Midterm review of the three programs | midterm | 2027-04-01 to 2027-06-30 | external | Riverbend Nonprofit Network | Planned |
| Final evaluation of the network strategy | final | 2028-07-01 to 2028-12-31 | external | Riverbend Nonprofit Network | Planned |

#### Midterm review of the three programs

**Method:** Interviews with 20 member leaders and county staff, with the rubric ratings and joint-position records

**Budget:** Funded by the regional funder's evaluation grant

1. Are member leaders applying what they learn, and what helps or stops them?
2. Do joint positions reach the county planning office, and does it act on them?

#### Final evaluation of the network strategy

**Method:** Contribution analysis using the rubric ratings, member survey, county service data, and outcome harvesting

1. Did members become more capable and more collaborative, and how much did the network contribute?
2. Did residents receive more and better services from member organizations?

### Learning Agenda

| Question | Why it matters | How it will be answered | Lead | Due | Status |
|----------|----------------|-------------------------|------|-----|--------|
| Do leaders who finish a cohort change how they run their organizations within six months? | Decides whether to keep investing in cohorts or shift to coaching. | Rubric ratings at six months; Follow-up interviews with a sample of leaders | Riverbend Nonprofit Network | 2027-03-31 | Open |
| Is the Learning Exchange behind plan because of attendance barriers or because of scheduling? | Decides whether to change session times or offer travel help. | Cohort exit form; Stories from small members | Riverbend Nonprofit Network | 2026-12-31 | Open |
| When members become more capable, do residents see better services? | Tests the link the whole strategy rests on. | County service data; Resident survey | Riverbend Nonprofit Network | 2028-06-30 | Open |

### Reviews and Reporting

| Review | Frequency | Participants | Lead | Next | Questions taken up |
|--------|-----------|--------------|------|------|--------------------|
| Quarterly learning review | Quarterly | Program leads, Two member representatives | Riverbend Nonprofit Network | 2026-12-10 | `lq-skills-applied`, `lq-attendance-barriers`, `lq-members-deliver` |
| Annual strategy review with the board | Yearly | Network board, Program leads | Riverbend Nonprofit Network | 2027-01-20 |  |

| Audience | Needs to know | Receives | Frequency | Prepared by | Indicators |
|----------|---------------|----------|-----------|-------------|------------|
| Network board | Whether each program area is on track and where it is stuck. | board deck, quarterly summary | Quarterly | Riverbend Nonprofit Network | Leadership practice maturity; Back-office share of operating budget; Joint positions adopted by members; Member organizations reached across the whole network; Attendee hours across all three program areas; Residents served by member organizations |
| Members | What the network did with their dues, and what other members learned. | newsletter, member assembly | Every six months | Riverbend Nonprofit Network | Stories of members helping each other; Leadership practice maturity; Member organizations reached across the whole network |
| Funders | Evidence that members are stronger and that residents are better served. | annual report, grant report | Yearly | Riverbend Nonprofit Network | Residents served by member organizations; Member organizations reached across the whole network; Back-office share of operating budget; Leadership practice maturity; Attendee hours across all three program areas; Stories of members helping each other |
| County planning office | What members jointly recommend. | policy briefing | Every six months | Riverbend Nonprofit Network | Joint positions adopted by members |

### Budget

90 staff minutes a week for measurement

## Glossary

| Term | Definition | Examples |
|------|------------|----------|
| Output | A direct product of the network's own work, within its control. | Twelve cohort sessions held in a year; A pooled bookkeeping contract signed |
| Outcome | A change in members or in the actors around them that must happen for the theory to work, stated as a condition and not as an activity. | Member leaders apply improved management practices; Local officials and funders engage members in planning |
| Impact | Change for residents of the county, which the network contributes to but does not own. | Residents get timely, good-quality services from local nonprofits |
| Member organization | A nonprofit in the county that has joined the network. The network counts members, not individual staff. | — |
| Joint position | A written position on a local issue that at least ten member organizations have formally signed. | A shared statement asking the county for earlier notice of planning meetings |
| Pooled services | Back-office services such as bookkeeping and payroll that several members buy together to lower cost and effort. | — |
| Attribution and contribution | Attribution means the network can credit a result to its own work, as with hours it delivered. Contribution means the network played a part but other forces matter, as with resident outcomes. | — |
