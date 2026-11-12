## Layered logic model: Learning Exchange

```d2
vars: {
  d2-config: {
    layout-engine: elk
  }
}

direction: up

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  inputs: {
    style.fill: "#fff1dc"
    style.stroke: "#e08a1e"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  stage: {
    style.fill: "#f4f6fa"
    style.stroke: "#dde3ec"
    style.stroke-width: 1
    style.font-size: 28
    style.bold: true
  }
}

inputs: Inputs {
  direction: right
  class: stage

  _input1: "Two program staff\nand meeting space\ndonated by the\ncounty library\nsystem" {
    class: inputs
  }
  _input2: "One program staff\nmember and a pool of\nvolunteer coaches" {
    class: inputs
  }
  _input3: "Fee-for-service\nagreement with the\ntraining partner" {
    class: inputs
  }
}

activities: Activities {
  direction: right
  class: stage

  "int-peer-cohorts": "Run mixed-size\npeer-learning\ncohorts" {
    class: activities
  }
  "int-leader-coaching": "Offer one-to-one\nleadership coaching" {
    class: activities
  }
  "int-training-series": "Deliver a management\nskills training\nseries (inferred)" {
    class: activities
  }
}

outputs: Outputs {
  direction: right
  class: stage

  _output1: Cohort sessions held {
    class: outputs
  }
  _output2: "Peer-pairing\nintroductions" {
    class: outputs
  }
  "oc-out-cohorts": "Peer-learning\ncohorts and coaching\nsessions delivered" {
    class: outputs
  }
  _output3: "Coaching hours\ndelivered" {
    class: outputs
  }
  _output4: Workshops delivered {
    class: outputs
  }
  _output5: "Participant\nmaterials" {
    class: outputs
  }
}

short_term: Short-term outcomes {
  direction: right
  class: stage

  "oc-st-peer-ties": "Members form peer\nrelationships across\nsizes (inferred)" {
    class: short_term
  }
  "oc-st-leader-skills": "Leaders gain\ncoaching-based\nskills and\nconfidence" {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  direction: right
  class: stage

  "oc-im-peer-network": "Members turn to each\nother for advice and\njoint work\n(inferred)" {
    class: medium_term
  }
  "oc-im-joint-positions": "Members adopt and\nshare joint\npositions on local\nissues (hypothesis)" {
    class: medium_term
  }
  "oc-im-leaders-practice": "Member leaders apply\nimproved management\npractices" {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  direction: right
  class: stage

  "oc-lt-collab": "Members routinely\ncollaborate instead\nof duplicating\neffort" {
    class: long_term
  }
  "oc-lt-collective-voice": "Local government and\nfunders consult\nmembers through a\ncoordinated channel\n(hypothesis)" {
    class: long_term
  }
  "oc-lt-capable-members": "Member nonprofits\nrun strong,\nwell-managed\norganizations" {
    class: long_term
  }
}

impact: Impact {
  direction: right
  class: stage

  "oc-resident-voice": "Residents'\npriorities shape\nlocal public\ndecisions\n(hypothesis)" {
    class: impact
  }
  "oc-resident-access": "Residents get\ntimely, good-quality\nservices from local\nnonprofits\n(inferred)" {
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

## Layered logic model: Shared Services

```d2
vars: {
  d2-config: {
    layout-engine: elk
  }
}

direction: up

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  inputs: {
    style.fill: "#fff1dc"
    style.stroke: "#e08a1e"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  stage: {
    style.fill: "#f4f6fa"
    style.stroke: "#dde3ec"
    style.stroke-width: 1
    style.font-size: 28
    style.bold: true
  }
}

inputs: Inputs {
  direction: right
  class: stage

  _input1: "One shared-services\ncoordinator and\ncontracted\nbookkeepers" {
    class: inputs
  }
  _input2: "Part of one staff\nrole" {
    class: inputs
  }
}

activities: Activities {
  direction: right
  class: stage

  "int-pooled-backoffice": "Run pooled\nbookkeeping and\npayroll help" {
    class: activities
  }
  "int-group-purchasing": "Negotiate group\npurchasing\nagreements" {
    class: activities
  }
  "int-grants-calendar": "Maintain a shared\ngrants calendar" {
    class: activities
  }
}

outputs: Outputs {
  direction: right
  class: stage

  _output1: "Member accounts\nonboarded" {
    class: outputs
  }
  _output2: Shared service desk {
    class: outputs
  }
  "oc-out-services": "Pooled service\ncontracts and group\npurchasing deals in\nplace" {
    class: outputs
  }
  _output3: "Vendor agreements\nsigned" {
    class: outputs
  }
  "oc-out-calendar": "Shared grants\ncalendar published\nand kept current\n(inferred)" {
    class: outputs
  }
  _output4: "Monthly calendar\nupdate" {
    class: outputs
  }
  _output5: Deadline reminders {
    class: outputs
  }
}

short_term: Short-term outcomes {
  direction: right
  class: stage

  "oc-st-pooled-use": "Members use pooled\nback-office services\n(inferred)" {
    class: short_term
  }
  "oc-st-calendar-use": "Members consult the\nshared grants\ncalendar when\nplanning\n(hypothesis)" {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  direction: right
  class: stage

  "oc-im-lower-overhead": "Members spend less\nstaff time and money\non back-office work" {
    class: medium_term
  }
  "oc-im-grant-pipeline": "Members apply for\nmore grants that fit\ntheir mission\n(inferred)" {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  direction: right
  class: stage

  "oc-lt-capable-members": "Member nonprofits\nrun strong,\nwell-managed\norganizations" {
    class: long_term
  }
  "oc-lt-collab": "Members routinely\ncollaborate instead\nof duplicating\neffort" {
    class: long_term
  }
  "oc-lt-sustainable": "Members have\ndiversified, stable\nfunding (inferred)" {
    class: long_term
  }
}

impact: Impact {
  direction: right
  class: stage

  "oc-resident-access": "Residents get\ntimely, good-quality\nservices from local\nnonprofits\n(inferred)" {
    class: impact
  }
  "oc-resident-voice": "Residents'\npriorities shape\nlocal public\ndecisions\n(hypothesis)" {
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

## Layered logic model: Civic Voice

```d2
vars: {
  d2-config: {
    layout-engine: elk
  }
}

direction: up

classes: {
  activities: {
    style.fill: "#e4f3e5"
    style.stroke: "#2e7d32"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  impact: {
    style.fill: "#c5cfe8"
    style.stroke: "#1f3a7a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  long_term: {
    style.fill: "#d3dcee"
    style.stroke: "#2c4a8a"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  medium_term: {
    style.fill: "#dfe8f6"
    style.stroke: "#3f5fa8"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  outputs: {
    style.fill: "#ece4fb"
    style.stroke: "#7e57c2"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  short_term: {
    style.fill: "#e1f3fd"
    style.stroke: "#2d9cdb"
    style.border-radius: 12
    style.font-size: 24
    width: 270
  }
  stage: {
    style.fill: "#f4f6fa"
    style.stroke: "#dde3ec"
    style.stroke-width: 1
    style.font-size: 28
    style.bold: true
  }
}

activities: Activities {
  direction: right
  class: stage

  "int-policy-briefings": "Publish\nplain-language\npolicy briefings\n(hypothesis)" {
    class: activities
  }
  "int-joint-forums": "Convene member\nforums to draft\njoint positions" {
    class: activities
  }
  "int-liaison-channel": "Set up a standing\nliaison channel with\nthe county planning\noffice" {
    class: activities
  }
  "int-funder-roundtable": "Host a funder and\nmember roundtable" {
    class: activities
  }
}

outputs: Outputs {
  direction: right
  class: stage

  "oc-out-briefings": "Policy briefings and\njoint-position\ndrafts produced\n(hypothesis)" {
    class: outputs
  }
  _output1: Briefings published {
    class: outputs
  }
  _output2: Forums held {
    class: outputs
  }
  _output3: "Joint position\ndrafts" {
    class: outputs
  }
  _output4: "Quarterly liaison\nmeeting" {
    class: outputs
  }
  _output5: "Named contacts on\nboth sides" {
    class: outputs
  }
  _output6: "Roundtable held\ntwice a year" {
    class: outputs
  }
}

short_term: Short-term outcomes {
  direction: right
  class: stage

  "oc-st-policy-know": "Members understand\nlocal policy and\nfunding processes\n(inferred)" {
    class: short_term
  }
  "oc-st-channel": "A standing channel\nexists between\nmembers and the\ncounty planning\noffice" {
    class: short_term
  }
}

medium_term: Medium-term outcomes {
  direction: right
  class: stage

  "oc-im-joint-positions": "Members adopt and\nshare joint\npositions on local\nissues (hypothesis)" {
    class: medium_term
  }
  "oc-im-officials-engage": "Local officials and\nfunders engage\nmembers in planning" {
    class: medium_term
  }
}

long_term: Long-term outcomes {
  direction: right
  class: stage

  "oc-lt-collab": "Members routinely\ncollaborate instead\nof duplicating\neffort" {
    class: long_term
  }
  "oc-lt-collective-voice": "Local government and\nfunders consult\nmembers through a\ncoordinated channel\n(hypothesis)" {
    class: long_term
  }
}

impact: Impact {
  direction: right
  class: stage

  "oc-resident-access": "Residents get\ntimely, good-quality\nservices from local\nnonprofits\n(inferred)" {
    class: impact
  }
  "oc-resident-voice": "Residents'\npriorities shape\nlocal public\ndecisions\n(hypothesis)" {
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

