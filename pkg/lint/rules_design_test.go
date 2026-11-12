package lint

import (
	"slices"
	"testing"

	core "github.com/grokify/prism-core"

	"github.com/ProductBuildersHQ/outcome-model/pkg/graph"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
)

// addIndicator appends an indicator to the document and lets the caller
// adjust it.
func addIndicator(d *model.Model, id, measures string, adjust ...func(*model.Indicator)) {
	ind := model.Indicator{ID: id, Measures: measures, Name: id}
	for _, f := range adjust {
		f(&ind)
	}
	d.Indicators = append(d.Indicators, ind)
}

func selected(i *model.Indicator) { i.Selection = &model.Selection{Status: model.SelectionSelected} }

func outcome(d *model.Model, id string) *model.Outcome {
	for i := range d.Outcomes {
		if d.Outcomes[i].ID == id {
			return &d.Outcomes[i]
		}
	}
	panic("no outcome " + id)
}

func indicator(d *model.Model, id string) *model.Indicator {
	for i := range d.Indicators {
		if d.Indicators[i].ID == id {
			return &d.Indicators[i]
		}
	}
	panic("no indicator " + id)
}

// fires returns the elements a rule reports on the document.
func fires(t *testing.T, d *model.Model, rule string) []string {
	t.Helper()
	g, err := graph.Build(d)
	if err != nil {
		t.Fatalf("document does not build: %v", err)
	}
	res, err := Run(d, g, Options{Profile: ProfileFull})
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, f := range res.Findings {
		if f.Rule == rule {
			out = append(out, f.Element)
		}
	}
	return out
}

// Each case seeds one defect and checks the rule finds it, then repairs it
// and checks the rule falls silent about it. An empty element means a
// finding about the whole document.
func TestRulesFindSeededDefectsAndAcceptRepairs(t *testing.T) {
	horizon := func(d *model.Model) {
		d.Context.TimeHorizon = &model.TimeHorizon{Start: "2026-01-01", End: "2027-12-31"}
	}
	pathway := func(d *model.Model, outcomes ...string) {
		d.Pathways = append(d.Pathways, model.Pathway{ID: "pw-1", Title: "Pathway", OutcomeIDs: outcomes})
	}
	tests := []struct {
		rule, element string
		defect, fix   func(d *model.Model)
	}{
		{"OM-MEA-001", "oc-s1", func(d *model.Model) { d.Ceiling = &model.Ceiling{Level: model.LevelLongTerm} },
			func(d *model.Model) { addIndicator(d, "ind-1", "oc-s1") }},
		{"OM-MEA-001", "oc-s2", func(d *model.Model) {
			d.Ceiling = &model.Ceiling{Level: model.LevelLongTerm}
			addIndicator(d, "ind-1", "oc-s1")
		},
			func(d *model.Model) {
				addIndicator(d, "ind-2", "oc-b", func(i *model.Indicator) { i.ProxyFor = []string{"oc-s2"} })
			}},
		{"OM-MEA-002", "ind-x", func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1") },
			func(d *model.Model) {
				i := indicator(d, "ind-x")
				i.Baseline = &model.Baseline{Level: "low"}
				i.Targets = []model.Target{{Date: "2027-01-01", Level: "high"}}
				i.DataSource, i.Frequency = "Sign-in sheets", "P3M"
			}},
		{"OM-MEA-003", "ind-x", func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1") },
			func(d *model.Model) {
				indicator(d, "ind-x").Unit, indicator(d, "ind-x").Direction = "people", model.DirectionIncrease
			}},
		{"OM-MEA-004", "ind-x", func(d *model.Model) {
			horizon(d)
			addIndicator(d, "ind-x", "oc-s1", func(i *model.Indicator) { i.Targets = []model.Target{{Date: "2030-01-01", Level: "x"}} })
		}, func(d *model.Model) { indicator(d, "ind-x").Targets[0].Date = "2027-06-30" }},
		{"OM-MEA-005", "oc-a", func(d *model.Model) {
			addIndicator(d, "ind-s", "oc-s1", func(i *model.Indicator) { i.Targets = []model.Target{{Date: "2027-12-31", Level: "x"}} })
			addIndicator(d, "ind-a", "oc-a", func(i *model.Indicator) { i.Targets = []model.Target{{Date: "2026-12-31", Level: "x"}} })
		}, func(d *model.Model) { indicator(d, "ind-a").Targets[0].Date = "2028-12-31" }},
		{"OM-ACC-001", "", func(d *model.Model) {},
			func(d *model.Model) { d.Ceiling = &model.Ceiling{Level: model.LevelLongTerm} }},
		{"OM-ACC-002", "oc-impact", func(d *model.Model) {
			d.Ceiling = &model.Ceiling{Level: model.LevelLongTerm}
			outcome(d, "oc-impact").Responsibility = model.ResponsibilityOwn
		}, func(d *model.Model) { outcome(d, "oc-impact").Responsibility = model.ResponsibilityShared }},
		{"OM-CLM-001", "ind-x", func(d *model.Model) {
			d.Ceiling = &model.Ceiling{Level: model.LevelLongTerm}
			addIndicator(d, "ind-x", "oc-impact", func(i *model.Indicator) { i.Claim = model.ClaimAttribution })
		}, func(d *model.Model) { indicator(d, "ind-x").Claim = model.ClaimContribution }},
		{"OM-BUD-001", "", func(d *model.Model) {
			d.Budget = &model.MeasurementBudget{StaffMinutesPerWeek: 10}
			addIndicator(d, "ind-x", "oc-s1", selected, func(i *model.Indicator) {
				i.Collection = &model.Collection{Mode: model.ModeManual, MinutesPerPeriod: 100, Period: "P1W"}
			})
		}, func(d *model.Model) { d.Budget.StaffMinutesPerWeek = 200 }},
		{"OM-BUD-002", "", func(d *model.Model) {
			d.Budget = &model.MeasurementBudget{MaxQuestionsPerRespondentPerYear: 5}
			addIndicator(d, "ind-x", "oc-s1", selected, func(i *model.Indicator) {
				i.Collection = &model.Collection{Mode: model.ModeSurvey, Period: "P1M", QuestionsPerRespondent: 4}
			})
		}, func(d *model.Model) { d.Budget.MaxQuestionsPerRespondentPerYear = 100 }},
		{"OM-BUD-003", "", func(d *model.Model) {
			d.Budget = &model.MeasurementBudget{TargetSelectedMetrics: &model.Range{Min: 3, Max: 5}}
			addIndicator(d, "ind-x", "oc-s1", selected)
		}, func(d *model.Model) { d.Budget.TargetSelectedMetrics = &model.Range{Min: 1, Max: 2} }},
		{"OM-BUD-004", "ind-x", func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1", selected) },
			func(d *model.Model) {
				indicator(d, "ind-x").Collection = &model.Collection{Mode: model.ModePassive}
			}},
		{"OM-ASM-001", "asm-c", func(d *model.Model) {
			d.Assumptions = append(d.Assumptions, model.Assumption{ID: "asm-c", Statement: "S", Criticality: model.CriticalityCritical})
		}, func(d *model.Model) { d.Assumptions[1].Test = "Ask partners twice a year." }},
		{"OM-ASM-002", "asm-c", func(d *model.Model) {
			d.Assumptions = append(d.Assumptions, model.Assumption{ID: "asm-c", Statement: "S", Criticality: model.CriticalityHigh})
		}, func(d *model.Model) {
			d.Evidence = append(d.Evidence, model.Evidence{ID: "ev-1", Title: "Study"})
			d.Assumptions[1].EvidenceIDs = []string{"ev-1"}
		}},
		{"OM-ASM-003", "", func(d *model.Model) {},
			func(d *model.Model) { d.Assumptions[0].Type = model.AssumptionContextual }},
		{"OM-ASM-004", "asm-c", func(d *model.Model) {
			d.Assumptions = append(d.Assumptions, model.Assumption{ID: "asm-c", Statement: "S", NeedsTesting: true})
		}, func(d *model.Model) {
			d.Assumptions[1].Test = "Pilot it."
			d.Decisions = append(d.Decisions, model.Decision{ID: "dec-1", Question: "Q", RelatedIDs: []string{"asm-c"}})
		}},
		{"OM-ASM-005", "", func(d *model.Model) {
			outcome(d, "oc-s1").Responsibility = model.ResponsibilityPartner
		}, func(d *model.Model) {
			d.Assumptions = append(d.Assumptions, model.Assumption{ID: "asm-p", Statement: "Partners show up.", Type: model.AssumptionPartnership})
		}},
		{"OM-RSK-001", "asm-c", func(d *model.Model) {
			d.Assumptions = append(d.Assumptions, model.Assumption{ID: "asm-c", Statement: "S", Criticality: model.CriticalityCritical})
		}, func(d *model.Model) {
			d.Risks = append(d.Risks, model.Risk{Risk: core.Risk{ID: "rsk-1", Description: "D", Mitigation: "Plan B"}, RelatedAssumptionIDs: []string{"asm-c"}})
		}},
		{"OM-INT-002", "int-1", func(d *model.Model) {},
			func(d *model.Model) {
				d.Actors = append(d.Actors, model.Actor{ID: "act-1", Name: "Network"})
				d.Interventions[0].OwnerID = "act-1"
			}},
		{"OM-ACT-001", "oc-s1", func(d *model.Model) { outcome(d, "oc-s1").Sphere = model.SphereInfluence },
			func(d *model.Model) {
				d.Actors = append(d.Actors, model.Actor{ID: "act-1", Name: "Partner"})
				outcome(d, "oc-s1").ActorIDs = []string{"act-1"}
			}},
		{"OM-EVD-001", "ln-a-i", func(d *model.Model) {},
			func(d *model.Model) {
				d.Evidence = append(d.Evidence, model.Evidence{ID: "ev-1", Title: "Study"})
				d.Links[0].EvidenceIDs = []string{"ev-1"}
			}},
		{"OM-PTH-001", "pw-1", func(d *model.Model) {
			d.Outcomes = append(d.Outcomes, model.Outcome{ID: "oc-iso", Title: "Isolated", Level: model.LevelShortTerm})
			pathway(d, "oc-iso")
		}, func(d *model.Model) { d.Pathways[0].OutcomeIDs = []string{"oc-s1"} }},
		{"OM-PTH-002", "pw-1", func(d *model.Model) { pathway(d, "oc-s1") },
			func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1", selected) }},
		{"OM-PTH-003", "ind-x", func(d *model.Model) {
			pathway(d, "oc-s1")
			d.Pathways = append(d.Pathways, model.Pathway{ID: "pw-2", Title: "Other"})
			addIndicator(d, "ind-x", "oc-s1", func(i *model.Indicator) { i.PathwayIDs = []string{"pw-2"} })
		}, func(d *model.Model) { indicator(d, "ind-x").PathwayIDs = []string{"pw-1"} }},
		{"OM-MET-001", "met-1", func(d *model.Model) {
			pathway(d, "oc-s1")
			d.Metrics = append(d.Metrics, model.Metric{ID: "met-1", Name: "M", Definition: "D", Formula: "F", Layer: model.LayerCommon})
		}, func(d *model.Model) {
			addIndicator(d, "ind-x", "oc-s1", func(i *model.Indicator) { i.MetricID, i.PathwayIDs = "met-1", []string{"pw-1"} })
		}},
		{"OM-MET-003", "met-w", func(d *model.Model) {
			pathway(d, "oc-s1")
			d.Metrics = append(d.Metrics,
				model.Metric{ID: "met-wt", Name: "Weight", Definition: "D", Formula: "F"},
				model.Metric{ID: "met-w", Name: "W", Definition: "D", Formula: "F", Aggregation: model.AggregationWeightedMean, WeightMetricID: "met-wt"})
			addIndicator(d, "ind-w", "oc-s1", func(i *model.Indicator) { i.MetricID, i.PathwayIDs = "met-w", []string{"pw-1"} })
		}, func(d *model.Model) {
			addIndicator(d, "ind-wt", "oc-s2", func(i *model.Indicator) { i.MetricID, i.PathwayIDs = "met-wt", []string{"pw-1"} })
		}},
		{"OM-AUD-001", "aud-1", func(d *model.Model) {
			d.Audiences = append(d.Audiences, model.Audience{ID: "aud-1", Name: "Board"})
		},
			func(d *model.Model) { outcome(d, "oc-impact").AudienceIDs = []string{"aud-1"} }},
		{"OM-CAU-001", "cause-1", func(d *model.Model) {
			d.Causes = append(d.Causes, model.Cause{ID: "cause-1", Statement: "S", Level: model.CauseStructural})
		}, func(d *model.Model) { outcome(d, "oc-a").AddressesCauses = []string{"cause-1"} }},
		{"OM-CAU-002", "", func(d *model.Model) {
			d.Causes = append(d.Causes, model.Cause{ID: "cause-i", Statement: "S", Level: model.CauseImmediate})
			outcome(d, "oc-a").AddressesCauses = []string{"cause-i"}
		}, func(d *model.Model) {
			d.Causes = append(d.Causes, model.Cause{ID: "cause-u", Statement: "S", Level: model.CauseUnderlying})
			outcome(d, "oc-b").AddressesCauses = []string{"cause-u"}
		}},
		{"OM-RSP-001", "oc-s1", func(d *model.Model) { outcome(d, "oc-s1").Responsibility = model.ResponsibilityPartner },
			func(d *model.Model) {
				d.Actors = append(d.Actors, model.Actor{ID: "act-1", Name: "Partner"})
				outcome(d, "oc-s1").DeliveredBy = []string{"act-1"}
			}},
		{"OM-RSP-002", "oc-s1", func(d *model.Model) { outcome(d, "oc-s1").Responsibility = model.ResponsibilityPartner },
			func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1") }},
		{"OM-EQU-001", "oc-s1", func(d *model.Model) {
			d.TargetGroups = append(d.TargetGroups, model.TargetGroup{ID: "tg-1", Name: "Small members"})
			outcome(d, "oc-s1").TargetGroupIDs = []string{"tg-1"}
			addIndicator(d, "ind-x", "oc-s1")
		}, func(d *model.Model) { indicator(d, "ind-x").TargetGroupIDs = []string{"tg-1"} }},
		{"OM-VAL-001", "", func(d *model.Model) {},
			func(d *model.Model) {
				d.Consultations = append(d.Consultations, model.Consultation{ID: "con-1", Date: "2026-01-01", Step: model.StepWhole})
			}},
		{"OM-PRV-001", "oc-impact", func(d *model.Model) {},
			func(d *model.Model) {
				outcome(d, "oc-impact").Provenance = &model.Provenance{Basis: model.BasisInferred}
			}},
		{"OM-PRV-002", "oc-s1", func(d *model.Model) {
			outcome(d, "oc-s1").Provenance = &model.Provenance{Basis: model.BasisStated}
		}, func(d *model.Model) {
			d.Evidence = append(d.Evidence, model.Evidence{ID: "ev-1", Title: "Report", URL: "https://example.org/report"})
			outcome(d, "oc-s1").Provenance.EvidenceIDs = []string{"ev-1"}
		}},
		{"OM-PRV-003", "oc-s1", func(d *model.Model) {
			d.Metadata.Status = "approved"
			outcome(d, "oc-s1").Provenance = &model.Provenance{Basis: model.BasisHypothesis}
		}, func(d *model.Model) { outcome(d, "oc-s1").Provenance.Basis = model.BasisConfirmed }},
		{"OM-PRV-004", "ind-x", func(d *model.Model) { addIndicator(d, "ind-x", "oc-s1", selected) },
			func(d *model.Model) { indicator(d, "ind-x").DataAvailability = model.DataPublic }},
	}
	for _, tt := range tests {
		t.Run(tt.rule+"/"+tt.element, func(t *testing.T) {
			d := clean()
			tt.defect(d)
			if !slices.Contains(fires(t, d, tt.rule), tt.element) {
				t.Fatalf("%s did not report %q after the defect; reported %v", tt.rule, tt.element, fires(t, d, tt.rule))
			}
			tt.fix(d)
			if slices.Contains(fires(t, d, tt.rule), tt.element) {
				t.Errorf("%s still reports %q after the repair", tt.rule, tt.element)
			}
		})
	}
}

func TestEveryDocumentedRuleIsRegisteredOnce(t *testing.T) {
	// Rules the TRD lists that validation already enforces as errors, so
	// the document cannot reach lint in a state that breaks them.
	byValidation := []string{"OM-INT-001", "OM-MET-002"}
	want := []string{
		"OM-STR-001", "OM-STR-002", "OM-STR-003", "OM-STR-004", "OM-STR-005", "OM-STR-006", "OM-STR-007",
		"OM-LOG-001", "OM-LOG-002", "OM-LOG-003", "OM-LOG-004",
		"OM-ASM-001", "OM-ASM-002", "OM-ASM-003", "OM-ASM-004", "OM-ASM-005", "OM-RSK-001",
		"OM-MEA-001", "OM-MEA-002", "OM-MEA-003", "OM-MEA-004", "OM-MEA-005",
		"OM-ACC-001", "OM-ACC-002", "OM-CLM-001", "OM-INT-002", "OM-ACT-001", "OM-EVD-001",
		"OM-PTH-001", "OM-PTH-002", "OM-PTH-003", "OM-MET-001", "OM-MET-003",
		"OM-BUD-001", "OM-BUD-002", "OM-BUD-003", "OM-BUD-004", "OM-AUD-001",
		"OM-CAU-001", "OM-CAU-002", "OM-RSP-001", "OM-RSP-002", "OM-EQU-001", "OM-VAL-001",
		"OM-PRV-001", "OM-PRV-002", "OM-PRV-003", "OM-PRV-004",
		"OM-INT-003", "OM-INT-004", "OM-INT-005", "OM-INT-006", "OM-INT-007", "OM-INT-008",
		"OM-AUD-002", "OM-HLN-001", "OM-HLN-002", "OM-HLN-003", "OM-HLN-004",
		"OM-NSM-001", "OM-NSM-002", "OM-NSM-003", "OM-NSM-004", "OM-NSM-005", "OM-NSM-006", "OM-NSM-007",
		"OM-CTX-001", "OM-CTX-002",
		"OM-DMA-001", "OM-DMA-002", "OM-DMA-003",
		"OM-EFF-001", "OM-EFF-002", "OM-CAU-003",
		"OM-LRN-001", "OM-LRN-002", "OM-LRN-003",
		"OM-BSC-001", "OM-BSC-002", "OM-BSC-003", "OM-BSC-004",
	}
	slices.Sort(want)
	var got []string
	for _, r := range Default.Rules() {
		got = append(got, r.ID())
		if r.Description() == "" {
			t.Errorf("%s has no description", r.ID())
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("registered rules differ from the documented set\n got: %v\nwant: %v", got, want)
	}
	if len(byValidation) != 2 {
		t.Fatal("update the list of rules enforced by validation")
	}
}

func TestPRV001RunsOnlyInTheFullProfile(t *testing.T) {
	r, ok := Default.Rule("OM-PRV-001")
	if !ok {
		t.Fatal("OM-PRV-001 is not registered")
	}
	if Runs(r, ProfileStandard) || !Runs(r, ProfileFull) {
		t.Error("OM-PRV-001 should run under full only")
	}
}
