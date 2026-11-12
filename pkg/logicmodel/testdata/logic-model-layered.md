## Layered logic model: Learning Exchange

```mermaid
flowchart LR
  classDef activities fill:#e4f3e5,stroke:#2e7d32
  classDef impact fill:#c5cfe8,stroke:#1f3a7a
  classDef inputs fill:#fff1dc,stroke:#e08a1e
  classDef long_term fill:#d3dcee,stroke:#2c4a8a
  classDef medium_term fill:#dfe8f6,stroke:#3f5fa8
  classDef outputs fill:#ece4fb,stroke:#7e57c2
  classDef short_term fill:#e1f3fd,stroke:#2d9cdb
  classDef stage fill:#f4f6fa,stroke:#dde3ec,stroke-width:1px
  subgraph n_inputs["Inputs"]
    direction TB
    n_inputs__input1["Two program staff<br/>and meeting space<br/>donated by the<br/>county library<br/>system"]
    n_inputs__input2["One program staff<br/>member and a pool of<br/>volunteer coaches"]
    n_inputs__input3["Fee-for-service<br/>agreement with the<br/>training partner"]
  end
  subgraph n_activities["Activities"]
    direction TB
    n_activities_int_peer_cohorts["Run mixed-size<br/>peer-learning<br/>cohorts"]
    n_activities_int_leader_coaching["Offer one-to-one<br/>leadership coaching"]
    n_activities_int_training_series["Deliver a management<br/>skills training<br/>series (inferred)"]
  end
  subgraph n_outputs["Outputs"]
    direction TB
    n_outputs__output1["Cohort sessions held"]
    n_outputs__output2["Peer-pairing<br/>introductions"]
    n_outputs_oc_out_cohorts["Peer-learning<br/>cohorts and coaching<br/>sessions delivered"]
    n_outputs__output3["Coaching hours<br/>delivered"]
    n_outputs__output4["Workshops delivered"]
    n_outputs__output5["Participant<br/>materials"]
  end
  subgraph n_short_term["Short-term outcomes"]
    direction TB
    n_short_term_oc_st_peer_ties["Members form peer<br/>relationships across<br/>sizes (inferred)"]
    n_short_term_oc_st_leader_skills["Leaders gain<br/>coaching-based<br/>skills and<br/>confidence"]
  end
  subgraph n_medium_term["Medium-term outcomes"]
    direction TB
    n_medium_term_oc_im_peer_network["Members turn to each<br/>other for advice and<br/>joint work<br/>(inferred)"]
    n_medium_term_oc_im_joint_positions["Members adopt and<br/>share joint<br/>positions on local<br/>issues (hypothesis)"]
    n_medium_term_oc_im_leaders_practice["Member leaders apply<br/>improved management<br/>practices"]
  end
  subgraph n_long_term["Long-term outcomes"]
    direction TB
    n_long_term_oc_lt_collab["Members routinely<br/>collaborate instead<br/>of duplicating<br/>effort"]
    n_long_term_oc_lt_collective_voice["Local government and<br/>funders consult<br/>members through a<br/>coordinated channel<br/>(hypothesis)"]
    n_long_term_oc_lt_capable_members["Member nonprofits<br/>run strong,<br/>well-managed<br/>organizations"]
  end
  subgraph n_impact["Impact"]
    direction TB
    n_impact_oc_resident_voice["Residents'<br/>priorities shape<br/>local public<br/>decisions<br/>(hypothesis)"]
    n_impact_oc_resident_access["Residents get<br/>timely, good-quality<br/>services from local<br/>nonprofits<br/>(inferred)"]
  end
  n_inputs__input1 --> n_activities_int_peer_cohorts
  n_activities_int_peer_cohorts --> n_outputs__output1
  n_activities_int_peer_cohorts --> n_outputs__output2
  n_activities_int_peer_cohorts --> n_outputs_oc_out_cohorts
  n_outputs__output1 --> n_short_term_oc_st_peer_ties
  n_outputs__output2 --> n_short_term_oc_st_peer_ties
  n_inputs__input2 --> n_activities_int_leader_coaching
  n_activities_int_leader_coaching --> n_outputs__output3
  n_activities_int_leader_coaching --> n_outputs_oc_out_cohorts
  n_outputs__output3 --> n_short_term_oc_st_leader_skills
  n_inputs__input3 --> n_activities_int_training_series
  n_activities_int_training_series --> n_outputs__output4
  n_activities_int_training_series --> n_outputs__output5
  n_outputs__output4 --> n_short_term_oc_st_leader_skills
  n_outputs__output5 --> n_short_term_oc_st_leader_skills
  n_outputs_oc_out_cohorts --> n_short_term_oc_st_leader_skills
  n_outputs_oc_out_cohorts --> n_short_term_oc_st_peer_ties
  n_short_term_oc_st_leader_skills --> n_medium_term_oc_im_leaders_practice
  n_short_term_oc_st_peer_ties --> n_medium_term_oc_im_peer_network
  n_medium_term_oc_im_leaders_practice --> n_long_term_oc_lt_capable_members
  n_medium_term_oc_im_peer_network --> n_long_term_oc_lt_collab
  n_short_term_oc_st_peer_ties --> n_medium_term_oc_im_joint_positions
  n_medium_term_oc_im_joint_positions --> n_long_term_oc_lt_collab
  n_medium_term_oc_im_joint_positions --> n_long_term_oc_lt_collective_voice
  n_long_term_oc_lt_capable_members --> n_impact_oc_resident_access
  n_long_term_oc_lt_collab --> n_impact_oc_resident_access
  n_long_term_oc_lt_collective_voice --> n_impact_oc_resident_voice
  n_long_term_oc_lt_collab --> n_impact_oc_resident_voice
  class n_inputs__input1 inputs
  class n_inputs__input2 inputs
  class n_inputs__input3 inputs
  class n_inputs stage
  class n_activities_int_peer_cohorts activities
  class n_activities_int_leader_coaching activities
  class n_activities_int_training_series activities
  class n_activities stage
  class n_outputs__output1 outputs
  class n_outputs__output2 outputs
  class n_outputs_oc_out_cohorts outputs
  class n_outputs__output3 outputs
  class n_outputs__output4 outputs
  class n_outputs__output5 outputs
  class n_outputs stage
  class n_short_term_oc_st_peer_ties short_term
  class n_short_term_oc_st_leader_skills short_term
  class n_short_term stage
  class n_medium_term_oc_im_peer_network medium_term
  class n_medium_term_oc_im_joint_positions medium_term
  class n_medium_term_oc_im_leaders_practice medium_term
  class n_medium_term stage
  class n_long_term_oc_lt_collab long_term
  class n_long_term_oc_lt_collective_voice long_term
  class n_long_term_oc_lt_capable_members long_term
  class n_long_term stage
  class n_impact_oc_resident_voice impact
  class n_impact_oc_resident_access impact
  class n_impact stage
```

Feedback, not drawn:

- Members routinely collaborate instead of duplicating effort leads back to Members form peer relationships across sizes (inferred)

## Layered logic model: Shared Services

```mermaid
flowchart LR
  classDef activities fill:#e4f3e5,stroke:#2e7d32
  classDef impact fill:#c5cfe8,stroke:#1f3a7a
  classDef inputs fill:#fff1dc,stroke:#e08a1e
  classDef long_term fill:#d3dcee,stroke:#2c4a8a
  classDef medium_term fill:#dfe8f6,stroke:#3f5fa8
  classDef outputs fill:#ece4fb,stroke:#7e57c2
  classDef short_term fill:#e1f3fd,stroke:#2d9cdb
  classDef stage fill:#f4f6fa,stroke:#dde3ec,stroke-width:1px
  subgraph n_inputs["Inputs"]
    direction TB
    n_inputs__input1["One shared-services<br/>coordinator and<br/>contracted<br/>bookkeepers"]
    n_inputs__input2["Part of one staff<br/>role"]
  end
  subgraph n_activities["Activities"]
    direction TB
    n_activities_int_pooled_backoffice["Run pooled<br/>bookkeeping and<br/>payroll help"]
    n_activities_int_group_purchasing["Negotiate group<br/>purchasing<br/>agreements"]
    n_activities_int_grants_calendar["Maintain a shared<br/>grants calendar"]
  end
  subgraph n_outputs["Outputs"]
    direction TB
    n_outputs__output1["Member accounts<br/>onboarded"]
    n_outputs__output2["Shared service desk"]
    n_outputs_oc_out_services["Pooled service<br/>contracts and group<br/>purchasing deals in<br/>place"]
    n_outputs__output3["Vendor agreements<br/>signed"]
    n_outputs_oc_out_calendar["Shared grants<br/>calendar published<br/>and kept current<br/>(inferred)"]
    n_outputs__output4["Monthly calendar<br/>update"]
    n_outputs__output5["Deadline reminders"]
  end
  subgraph n_short_term["Short-term outcomes"]
    direction TB
    n_short_term_oc_st_pooled_use["Members use pooled<br/>back-office services<br/>(inferred)"]
    n_short_term_oc_st_calendar_use["Members consult the<br/>shared grants<br/>calendar when<br/>planning<br/>(hypothesis)"]
  end
  subgraph n_medium_term["Medium-term outcomes"]
    direction TB
    n_medium_term_oc_im_lower_overhead["Members spend less<br/>staff time and money<br/>on back-office work"]
    n_medium_term_oc_im_grant_pipeline["Members apply for<br/>more grants that fit<br/>their mission<br/>(inferred)"]
  end
  subgraph n_long_term["Long-term outcomes"]
    direction TB
    n_long_term_oc_lt_capable_members["Member nonprofits<br/>run strong,<br/>well-managed<br/>organizations"]
    n_long_term_oc_lt_collab["Members routinely<br/>collaborate instead<br/>of duplicating<br/>effort"]
    n_long_term_oc_lt_sustainable["Members have<br/>diversified, stable<br/>funding (inferred)"]
  end
  subgraph n_impact["Impact"]
    direction TB
    n_impact_oc_resident_access["Residents get<br/>timely, good-quality<br/>services from local<br/>nonprofits<br/>(inferred)"]
    n_impact_oc_resident_voice["Residents'<br/>priorities shape<br/>local public<br/>decisions<br/>(hypothesis)"]
  end
  n_inputs__input1 --> n_activities_int_pooled_backoffice
  n_activities_int_pooled_backoffice --> n_outputs__output1
  n_activities_int_pooled_backoffice --> n_outputs__output2
  n_activities_int_pooled_backoffice --> n_outputs_oc_out_services
  n_outputs__output1 --> n_short_term_oc_st_pooled_use
  n_outputs__output2 --> n_short_term_oc_st_pooled_use
  n_activities_int_group_purchasing --> n_outputs__output3
  n_activities_int_group_purchasing --> n_outputs_oc_out_services
  n_inputs__input2 --> n_activities_int_grants_calendar
  n_activities_int_grants_calendar --> n_outputs__output4
  n_activities_int_grants_calendar --> n_outputs__output5
  n_activities_int_grants_calendar --> n_outputs_oc_out_calendar
  n_outputs__output4 --> n_short_term_oc_st_calendar_use
  n_outputs__output5 --> n_short_term_oc_st_calendar_use
  n_outputs_oc_out_services --> n_short_term_oc_st_pooled_use
  n_outputs_oc_out_calendar --> n_short_term_oc_st_calendar_use
  n_short_term_oc_st_pooled_use --> n_medium_term_oc_im_lower_overhead
  n_short_term_oc_st_calendar_use --> n_medium_term_oc_im_grant_pipeline
  n_medium_term_oc_im_lower_overhead --> n_long_term_oc_lt_capable_members
  n_medium_term_oc_im_lower_overhead --> n_long_term_oc_lt_collab
  n_medium_term_oc_im_lower_overhead --> n_long_term_oc_lt_sustainable
  n_medium_term_oc_im_grant_pipeline --> n_long_term_oc_lt_sustainable
  n_long_term_oc_lt_capable_members --> n_impact_oc_resident_access
  n_long_term_oc_lt_collab --> n_impact_oc_resident_access
  n_long_term_oc_lt_sustainable --> n_impact_oc_resident_access
  n_long_term_oc_lt_collab --> n_impact_oc_resident_voice
  class n_inputs__input1 inputs
  class n_inputs__input2 inputs
  class n_inputs stage
  class n_activities_int_pooled_backoffice activities
  class n_activities_int_group_purchasing activities
  class n_activities_int_grants_calendar activities
  class n_activities stage
  class n_outputs__output1 outputs
  class n_outputs__output2 outputs
  class n_outputs_oc_out_services outputs
  class n_outputs__output3 outputs
  class n_outputs_oc_out_calendar outputs
  class n_outputs__output4 outputs
  class n_outputs__output5 outputs
  class n_outputs stage
  class n_short_term_oc_st_pooled_use short_term
  class n_short_term_oc_st_calendar_use short_term
  class n_short_term stage
  class n_medium_term_oc_im_lower_overhead medium_term
  class n_medium_term_oc_im_grant_pipeline medium_term
  class n_medium_term stage
  class n_long_term_oc_lt_capable_members long_term
  class n_long_term_oc_lt_collab long_term
  class n_long_term_oc_lt_sustainable long_term
  class n_long_term stage
  class n_impact_oc_resident_access impact
  class n_impact_oc_resident_voice impact
  class n_impact stage
```

## Layered logic model: Civic Voice

```mermaid
flowchart LR
  classDef activities fill:#e4f3e5,stroke:#2e7d32
  classDef impact fill:#c5cfe8,stroke:#1f3a7a
  classDef long_term fill:#d3dcee,stroke:#2c4a8a
  classDef medium_term fill:#dfe8f6,stroke:#3f5fa8
  classDef outputs fill:#ece4fb,stroke:#7e57c2
  classDef short_term fill:#e1f3fd,stroke:#2d9cdb
  classDef stage fill:#f4f6fa,stroke:#dde3ec,stroke-width:1px
  subgraph n_activities["Activities"]
    direction TB
    n_activities_int_policy_briefings["Publish<br/>plain-language<br/>policy briefings<br/>(hypothesis)"]
    n_activities_int_joint_forums["Convene member<br/>forums to draft<br/>joint positions"]
    n_activities_int_liaison_channel["Set up a standing<br/>liaison channel with<br/>the county planning<br/>office"]
    n_activities_int_funder_roundtable["Host a funder and<br/>member roundtable"]
  end
  subgraph n_outputs["Outputs"]
    direction TB
    n_outputs_oc_out_briefings["Policy briefings and<br/>joint-position<br/>drafts produced<br/>(hypothesis)"]
    n_outputs__output1["Briefings published"]
    n_outputs__output2["Forums held"]
    n_outputs__output3["Joint position<br/>drafts"]
    n_outputs__output4["Quarterly liaison<br/>meeting"]
    n_outputs__output5["Named contacts on<br/>both sides"]
    n_outputs__output6["Roundtable held<br/>twice a year"]
  end
  subgraph n_short_term["Short-term outcomes"]
    direction TB
    n_short_term_oc_st_policy_know["Members understand<br/>local policy and<br/>funding processes<br/>(inferred)"]
    n_short_term_oc_st_channel["A standing channel<br/>exists between<br/>members and the<br/>county planning<br/>office"]
  end
  subgraph n_medium_term["Medium-term outcomes"]
    direction TB
    n_medium_term_oc_im_joint_positions["Members adopt and<br/>share joint<br/>positions on local<br/>issues (hypothesis)"]
    n_medium_term_oc_im_officials_engage["Local officials and<br/>funders engage<br/>members in planning"]
  end
  subgraph n_long_term["Long-term outcomes"]
    direction TB
    n_long_term_oc_lt_collab["Members routinely<br/>collaborate instead<br/>of duplicating<br/>effort"]
    n_long_term_oc_lt_collective_voice["Local government and<br/>funders consult<br/>members through a<br/>coordinated channel<br/>(hypothesis)"]
  end
  subgraph n_impact["Impact"]
    direction TB
    n_impact_oc_resident_access["Residents get<br/>timely, good-quality<br/>services from local<br/>nonprofits<br/>(inferred)"]
    n_impact_oc_resident_voice["Residents'<br/>priorities shape<br/>local public<br/>decisions<br/>(hypothesis)"]
  end
  n_activities_int_policy_briefings --> n_outputs__output1
  n_activities_int_policy_briefings --> n_outputs_oc_out_briefings
  n_outputs__output1 --> n_short_term_oc_st_policy_know
  n_activities_int_joint_forums --> n_outputs__output2
  n_activities_int_joint_forums --> n_outputs__output3
  n_outputs__output2 --> n_short_term_oc_st_policy_know
  n_outputs__output3 --> n_short_term_oc_st_policy_know
  n_outputs__output2 --> n_medium_term_oc_im_joint_positions
  n_outputs__output3 --> n_medium_term_oc_im_joint_positions
  n_activities_int_liaison_channel --> n_outputs__output4
  n_activities_int_liaison_channel --> n_outputs__output5
  n_outputs__output4 --> n_short_term_oc_st_channel
  n_outputs__output5 --> n_short_term_oc_st_channel
  n_activities_int_funder_roundtable --> n_outputs__output6
  n_outputs__output6 --> n_short_term_oc_st_channel
  n_outputs__output6 --> n_medium_term_oc_im_officials_engage
  n_outputs_oc_out_briefings --> n_short_term_oc_st_policy_know
  n_short_term_oc_st_policy_know --> n_medium_term_oc_im_joint_positions
  n_medium_term_oc_im_joint_positions --> n_long_term_oc_lt_collab
  n_medium_term_oc_im_joint_positions --> n_long_term_oc_lt_collective_voice
  n_medium_term_oc_im_officials_engage --> n_long_term_oc_lt_collective_voice
  n_short_term_oc_st_channel --> n_medium_term_oc_im_officials_engage
  n_long_term_oc_lt_collab --> n_impact_oc_resident_access
  n_long_term_oc_lt_collective_voice --> n_impact_oc_resident_voice
  n_long_term_oc_lt_collab --> n_impact_oc_resident_voice
  class n_activities_int_policy_briefings activities
  class n_activities_int_joint_forums activities
  class n_activities_int_liaison_channel activities
  class n_activities_int_funder_roundtable activities
  class n_activities stage
  class n_outputs_oc_out_briefings outputs
  class n_outputs__output1 outputs
  class n_outputs__output2 outputs
  class n_outputs__output3 outputs
  class n_outputs__output4 outputs
  class n_outputs__output5 outputs
  class n_outputs__output6 outputs
  class n_outputs stage
  class n_short_term_oc_st_policy_know short_term
  class n_short_term_oc_st_channel short_term
  class n_short_term stage
  class n_medium_term_oc_im_joint_positions medium_term
  class n_medium_term_oc_im_officials_engage medium_term
  class n_medium_term stage
  class n_long_term_oc_lt_collab long_term
  class n_long_term_oc_lt_collective_voice long_term
  class n_long_term stage
  class n_impact_oc_resident_access impact
  class n_impact_oc_resident_voice impact
  class n_impact stage
```

