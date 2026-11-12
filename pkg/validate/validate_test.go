package validate

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
)

const fixtures = "../model/testdata/"

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(fixtures + name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return data
}

// fullDoc returns a fresh copy of the valid full fixture.
func fullDoc(t *testing.T) *model.Model {
	t.Helper()
	var doc model.Model
	if err := json.Unmarshal(readFixture(t, "full.json"), &doc); err != nil {
		t.Fatalf("decode full.json: %v", err)
	}
	return &doc
}

func obsSet(t *testing.T) *model.ObservationSet {
	t.Helper()
	var s model.ObservationSet
	if err := json.Unmarshal(readFixture(t, "observations.json"), &s); err != nil {
		t.Fatalf("decode observations.json: %v", err)
	}
	return &s
}

func describe(issues []Issue) string {
	var b strings.Builder
	for _, i := range issues {
		b.WriteString("\n  " + i.String())
	}
	return b.String()
}

func TestValidFixturesHaveNoIssues(t *testing.T) {
	for _, name := range []string{"minimal.json", "full.json"} {
		var doc model.Model
		if err := json.Unmarshal(readFixture(t, name), &doc); err != nil {
			t.Fatal(err)
		}
		if issues := Document(&doc); len(issues) != 0 {
			t.Errorf("%s: unexpected issues:%s", name, describe(issues))
		}
	}
	if issues := Observations(fullDoc(t), obsSet(t)); len(issues) != 0 {
		t.Errorf("observations: unexpected issues:%s", describe(issues))
	}
}

func TestDocumentIssues(t *testing.T) {
	f := func(v float64) *float64 { return &v }
	_ = f
	tests := []struct {
		name     string
		mutate   func(d *model.Model)
		code     string
		path     string
		severity Severity
	}{
		{"duplicate id", func(d *model.Model) { d.Outcomes[1].ID = "oc-safer-agents" }, "VAL-ID-002", "/outcomes/1/id", SeverityError},
		{"empty id", func(d *model.Model) { d.Links[0].ID = "" }, "VAL-ID-001", "/links/0/id", SeverityError},
		{"bad id format", func(d *model.Model) { d.Actors[0].ID = "Bad_ID" }, "VAL-ID-003", "/actors/0/id", SeverityError},
		{"id shared across kinds", func(d *model.Model) { d.Risks[0].ID = "ev-survey" }, "VAL-ID-002", "/risks/0/id", SeverityError},
		{"link to unknown", func(d *model.Model) { d.Links[0].To = "oc-missing" }, "VAL-REF-001", "/links/0/to", SeverityError},
		{"link from assumption", func(d *model.Model) { d.Links[0].From = "asm-no-fork" }, "VAL-REF-002", "/links/0/from", SeverityError},
		{"link to intervention", func(d *model.Model) { d.Links[0].To = "int-interop" }, "VAL-REF-002", "/links/0/to", SeverityError},
		{"self link", func(d *model.Model) { d.Links[0].To = d.Links[0].From }, "VAL-LINK-001", "/links/0", SeverityError},
		{"loop type without feedback", func(d *model.Model) { d.Links[0].LoopType = model.LoopBalancing }, "VAL-LINK-002", "/links/0/loopType", SeverityError},
		{"bad lag", func(d *model.Model) { d.Links[0].Lag = "6 months" }, "VAL-DUR-001", "/links/0/lag", SeverityError},
		{"bare P duration", func(d *model.Model) { d.Links[0].Lag = "P" }, "VAL-DUR-001", "/links/0/lag", SeverityError},
		{"link assumption of wrong kind", func(d *model.Model) { d.Links[0].AssumptionIDs = []string{"ev-survey"} }, "VAL-REF-002", "/links/0/assumptionIds/0", SeverityError},
		{"group too small", func(d *model.Model) { d.Groups[0].Members = []string{"oc-idp-support"} }, "VAL-GRP-001", "/preconditionGroups/0/members", SeverityError},
		{"group contains target", func(d *model.Model) { d.Groups[0].Members[0] = "oc-enterprise-require" }, "VAL-GRP-002", "/preconditionGroups/0/members/0", SeverityError},
		{"group member unknown", func(d *model.Model) { d.Groups[0].Members[1] = "nope" }, "VAL-REF-001", "/preconditionGroups/0/members/1", SeverityError},
		{"group target not outcome", func(d *model.Model) { d.Groups[0].Target = "ln-spec-idp" }, "VAL-REF-002", "/preconditionGroups/0/target", SeverityError},
		{"group k too large", func(d *model.Model) { d.Groups[0].Logic = model.Logic{Type: model.LogicAtLeast, K: 3} }, "VAL-LOGIC-003", "/preconditionGroups/0/logic/k", SeverityError},
		{"at_least without k", func(d *model.Model) { d.Outcomes[1].Logic = &model.Logic{Type: model.LogicAtLeast} }, "VAL-LOGIC-001", "/outcomes/1/logic/k", SeverityError},
		{"k with all", func(d *model.Model) { d.Outcomes[1].Logic = &model.Logic{Type: model.LogicAll, K: 2} }, "VAL-LOGIC-002", "/outcomes/1/logic/k", SeverityError},
		{"intervention without targets", func(d *model.Model) { d.Interventions[0].Targets = nil }, "VAL-INT-001", "/interventions/0/targets", SeverityError},
		{"intervention start not a date", func(d *model.Model) { d.Interventions[0].StartDate = "2027-02-30" }, "VAL-DATE-001", "/interventions/0/startDate", SeverityError},
		{"intervention ends before it starts", func(d *model.Model) {
			d.Interventions[0].StartDate, d.Interventions[0].EndDate = "2027-03-01", "2027-02-01"
		}, "VAL-DATE-002", "/interventions/0", SeverityError},
		{"objective without statement or one outcome", func(d *model.Model) { d.Context.Objective = &model.Objective{} }, "VAL-OBJ-001", "/context/objective", SeverityError},
		{"objective names an unknown outcome", func(d *model.Model) { d.Context.Objective.OutcomeIDs = []string{"oc-ghost"} }, "VAL-REF-001", "/context/objective/outcomeIds/0", SeverityError},
		{"evaluation without title", func(d *model.Model) { d.Evaluations[0].Title = " " }, "VAL-EVL-001", "/evaluations/0/title", SeverityError},
		{"evaluation with an empty question", func(d *model.Model) { d.Evaluations[0].Questions = []string{""} }, "VAL-EVL-002", "/evaluations/0/questions/0", SeverityError},
		{"evaluation type unknown", func(d *model.Model) { d.Evaluations[0].Type = "summative" }, "VAL-ENUM-001", "/evaluations/0/type", SeverityError},
		{"evaluation evaluator unknown", func(d *model.Model) { d.Evaluations[0].Evaluator = "vendor" }, "VAL-ENUM-001", "/evaluations/0/evaluator", SeverityError},
		{"evaluation ends before it starts", func(d *model.Model) { d.Evaluations[0].EndDate = "2026-12-01" }, "VAL-DATE-002", "/evaluations/0", SeverityError},
		{"evaluation related to an actor", func(d *model.Model) { d.Evaluations[0].RelatedIDs = []string{"act-idp"} }, "VAL-REF-002", "/evaluations/0/relatedIds/0", SeverityError},
		{"evaluation owner unknown", func(d *model.Model) { d.Evaluations[0].OwnerID = "act-ghost" }, "VAL-REF-001", "/evaluations/0/ownerActorId", SeverityError},
		{"learning question empty", func(d *model.Model) { d.LearningQuestions[0].Question = "" }, "VAL-LRN-001", "/learningQuestions/0/question", SeverityError},
		{"answered without an answer", func(d *model.Model) { d.LearningQuestions[0].Answer = "" }, "VAL-LRN-002", "/learningQuestions/0/answer", SeverityError},
		{"question status unknown", func(d *model.Model) { d.LearningQuestions[0].Status = "done" }, "VAL-ENUM-001", "/learningQuestions/0/status", SeverityError},
		{"question about another question", func(d *model.Model) { d.LearningQuestions[0].RelatedIDs = []string{"lq-incentive"} }, "VAL-REF-002", "/learningQuestions/0/relatedIds/0", SeverityError},
		{"question due date bad", func(d *model.Model) { d.LearningQuestions[0].DueDate = "soon" }, "VAL-DATE-001", "/learningQuestions/0/dueDate", SeverityError},
		{"review meeting without name", func(d *model.Model) { d.ReviewMeetings[0].Name = "" }, "VAL-MTG-001", "/reviewMeetings/0/name", SeverityError},
		{"review meeting frequency bad", func(d *model.Model) { d.ReviewMeetings[0].Frequency = "quarterly" }, "VAL-DUR-001", "/reviewMeetings/0/frequency", SeverityError},
		{"review meeting takes up an evaluation", func(d *model.Model) { d.ReviewMeetings[0].QuestionIDs = []string{"eval-midterm"} }, "VAL-REF-002", "/reviewMeetings/0/questionIds/0", SeverityError},
		{"quality assessed date bad", func(d *model.Model) {
			d.Indicators[0].QualityAssurance = &model.QualityAssurance{LastAssessed: "2026"}
		}, "VAL-DATE-001", "/indicators/0/qualityAssurance/lastAssessed", SeverityError},
		{"audience frequency bad", func(d *model.Model) { d.Audiences[0].Frequency = "often" }, "VAL-DUR-001", "/audiences/0/frequency", SeverityError},
		{"audience owner unknown", func(d *model.Model) { d.Audiences[0].OwnerID = "act-ghost" }, "VAL-REF-001", "/audiences/0/ownerActorId", SeverityError},
		{"outcome owner unknown", func(d *model.Model) { d.Outcomes[0].OwnerID = "act-ghost" }, "VAL-REF-001", "/outcomes/0/ownerActorId", SeverityError},
		{"scorecard perspective bad", func(d *model.Model) { d.Outcomes[0].ScorecardPerspective = "people" }, "VAL-ENUM-001", "/outcomes/0/scorecardPerspective", SeverityError},
		{"lead or lag bad", func(d *model.Model) { d.Indicators[0].LeadLag = "leading" }, "VAL-ENUM-001", "/indicators/0/leadLag", SeverityError},
		{"intervention targets non-outcome", func(d *model.Model) { d.Interventions[0].Targets = []string{"asm-no-fork"} }, "VAL-REF-002", "/interventions/0/targets/0", SeverityError},
		{"indicator measures unknown", func(d *model.Model) { d.Indicators[0].Measures = "oc-missing" }, "VAL-REF-001", "/indicators/0/measures", SeverityError},
		{"indicator measures link", func(d *model.Model) { d.Indicators[0].Measures = "ln-spec-idp" }, "VAL-REF-002", "/indicators/0/measures", SeverityError},
		{"baseline level not in levels", func(d *model.Model) { d.Indicators[1].Baseline.Level = "gigantic" }, "VAL-IND-001", "/indicators/1/baseline/level", SeverityError},
		{"target level not in levels", func(d *model.Model) { d.Indicators[1].Targets[0].Level = "huge" }, "VAL-IND-001", "/indicators/1/targets/0/level", SeverityError},
		{"bad baseline date", func(d *model.Model) { d.Indicators[0].Baseline.Date = "01/02/2026" }, "VAL-DATE-001", "/indicators/0/baseline/date", SeverityError},
		{"impossible date", func(d *model.Model) { d.Indicators[0].Targets[0].Date = "2026-02-30" }, "VAL-DATE-001", "/indicators/0/targets/0/date", SeverityError},
		{"target without date", func(d *model.Model) { d.Indicators[0].Targets[0].Date = "" }, "VAL-DATE-003", "/indicators/0/targets/0/date", SeverityError},
		{"bad frequency", func(d *model.Model) { d.Indicators[0].Frequency = "quarterly" }, "VAL-DUR-001", "/indicators/0/frequency", SeverityError},
		{"horizon reversed", func(d *model.Model) { d.Context.TimeHorizon.End = "2025-01-01" }, "VAL-DATE-002", "/context/timeHorizon", SeverityError},
		{"invalid level", func(d *model.Model) { d.Outcomes[2].Level = "huge" }, "VAL-ENUM-001", "/outcomes/2/level", SeverityError},
		{"missing level", func(d *model.Model) { d.Outcomes[2].Level = "" }, "VAL-OUT-002", "/outcomes/2/level", SeverityError},
		{"invalid role in list", func(d *model.Model) { d.Actors[0].Roles[1] = "overlord" }, "VAL-ENUM-001", "/actors/0/roles/1", SeverityError},
		{"invalid assumption status", func(d *model.Model) { d.Assumptions[0].Status = "maybe" }, "VAL-ENUM-001", "/assumptions/0/status", SeverityError},
		{"empty outcome title", func(d *model.Model) { d.Outcomes[0].Title = "  " }, "VAL-OUT-001", "/outcomes/0/title", SeverityError},
		{"wrong kind", func(d *model.Model) { d.Kind = model.KindStatus }, "VAL-DOC-001", "/kind", SeverityError},
		{"unsupported spec version", func(d *model.Model) { d.SpecVersion = "1.0" }, "VAL-DOC-002", "/specVersion", SeverityError},
		{"bad document id", func(d *model.Model) { d.ID = "Not A Slug" }, "VAL-DOC-003", "/id", SeverityError},
		{"no outcomes", func(d *model.Model) {
			d.Outcomes, d.Links, d.Groups, d.Indicators = nil, nil, nil, nil
			d.Interventions, d.Alignments, d.Decisions = nil, nil, nil
			d.Assumptions, d.ExternalFactors, d.Risks, d.Actors = nil, nil, nil, nil
			d.Context.Market = nil
		}, "VAL-DOC-004", "/outcomes", SeverityError},
		{"empty problem statement", func(d *model.Model) { d.Context.ProblemStatement = "" }, "VAL-DOC-005", "/context/problemStatement", SeverityError},
		{"value without name", func(d *model.Model) { d.Context.Values = []model.Value{{Name: " "}} }, "VAL-DOC-006", "/context/values/0/name", SeverityError},
		{"unknown projection", func(d *model.Model) { d.Pathways[0].Projections = []model.Projection{"six_sigma"} }, "VAL-ENUM-001", "/pathways/0/projections/0", SeverityError},
		{"duplicate projection", func(d *model.Model) {
			d.Pathways[0].Projections = []model.Projection{model.ProjectionDMAIC, model.ProjectionDMAIC}
		}, "VAL-PTH-007", "/pathways/0/projections/1", SeverityError},
		{"duplicate value", func(d *model.Model) {
			d.Context.Values = []model.Value{{Name: "Candor"}, {Name: "Candor"}}
		}, "VAL-DOC-007", "/context/values/1/name", SeverityError},
		{"empty ceiling", func(d *model.Model) { d.Ceiling = &model.Ceiling{} }, "VAL-CEIL-001", "/ceilingOfAccountability", SeverityError},
		{"ceiling outcome unknown", func(d *model.Model) { d.Ceiling = &model.Ceiling{OutcomeIDs: []string{"zzz"}} }, "VAL-REF-001", "/ceilingOfAccountability/outcomeIds/0", SeverityError},
		{"unknown actor", func(d *model.Model) { d.Outcomes[1].ActorIDs = []string{"act-ghost"} }, "VAL-REF-001", "/outcomes/1/actorIds/0", SeverityError},
		{"actor ref to non-actor", func(d *model.Model) { d.Interventions[0].OwnerID = "oc-safer-agents" }, "VAL-REF-002", "/interventions/0/ownerActorId", SeverityError},
		{"perspective actor unknown", func(d *model.Model) { d.Perspective.ActorID = "act-ghost" }, "VAL-REF-001", "/perspective/actorId", SeverityError},
		{"assumption applies to evidence", func(d *model.Model) { d.Assumptions[0].AppliesTo = []string{"ev-survey"} }, "VAL-REF-002", "/assumptions/0/appliesTo/0", SeverityError},
		{"risk relates to non-assumption", func(d *model.Model) { d.Risks[0].RelatedAssumptionIDs = []string{"oc-safer-agents"} }, "VAL-REF-002", "/risks/0/relatedAssumptionIds/0", SeverityError},
		{"alignment to unknown element", func(d *model.Model) { d.Alignments[0].ElementID = "ghost" }, "VAL-REF-001", "/alignments/0/elementId", SeverityError},
		{"external ref without ref", func(d *model.Model) { d.Actors[1].ExternalRefs[0].Ref = "" }, "VAL-EXT-001", "/actors/1/externalRefs/0/ref", SeverityError},
		{"capability without id", func(d *model.Model) { d.Interventions[0].Capabilities[0].CapabilityID = "" }, "VAL-INT-003", "/interventions/0/capabilities/0/capabilityId", SeverityError},
		{"bad decision due date", func(d *model.Model) { d.Decisions[0].DueDate = "next week" }, "VAL-DATE-001", "/decisions/0/dueDate", SeverityError},
		{"revision trigger evidence unknown", func(d *model.Model) { d.Revisions[0].TriggerEvidenceIDs = []string{"ev-ghost"} }, "VAL-REF-001", "/revisions/0/triggerEvidenceIds/0", SeverityError},
		{"revision changed id removed", func(d *model.Model) { d.Revisions[0].ChangedIDs = []string{"oc-removed"} }, "VAL-REF-003", "/revisions/0/changedIds/0", SeverityWarning},
		{"actor missing from local list but registry named", func(d *model.Model) {
			d.ActorRegistry = "shared-actors"
			d.Outcomes[1].ActorIDs = []string{"act-elsewhere"}
		}, "VAL-REF-004", "/outcomes/1/actorIds/0", SeverityWarning},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := fullDoc(t)
			tt.mutate(doc)
			issues := Document(doc)
			var got *Issue
			for i := range issues {
				if issues[i].Code == tt.code && issues[i].Path == tt.path {
					got = &issues[i]
					break
				}
			}
			if got == nil {
				t.Fatalf("expected %s at %s; got:%s", tt.code, tt.path, describe(issues))
			}
			if got.Severity != tt.severity {
				t.Fatalf("%s severity = %s, want %s", tt.code, got.Severity, tt.severity)
			}
			if got.Message == "" {
				t.Fatalf("%s has an empty message", tt.code)
			}
		})
	}
}

func TestActorRegistryResolution(t *testing.T) {
	doc := fullDoc(t)
	doc.ActorRegistry = "shared-actors"
	doc.Outcomes[1].ActorIDs = []string{"act-elsewhere"}

	reg := &model.ActorRegistry{Actors: []model.Actor{{ID: "act-elsewhere", Name: "Elsewhere"}}}
	if issues := Document(doc, WithActorRegistry(reg)); len(issues) != 0 {
		t.Errorf("registry actor should resolve:%s", describe(issues))
	}

	empty := &model.ActorRegistry{Actors: []model.Actor{{ID: "someone-else", Name: "Else"}}}
	issues := Document(doc, WithActorRegistry(empty))
	if len(issues) != 1 || issues[0].Code != "VAL-REF-001" || issues[0].Severity != SeverityError {
		t.Errorf("actor missing from supplied registry should be an error:%s", describe(issues))
	}

	if issues := Document(doc, WithActorRegistry(nil)); len(issues) != 1 || issues[0].Code != "VAL-REF-004" {
		t.Errorf("nil registry should leave the unverifiable warning:%s", describe(issues))
	}
}

func TestDocumentIsDeterministic(t *testing.T) {
	doc := fullDoc(t)
	doc.Links[0].To = "oc-missing"
	doc.Outcomes[2].Level = "huge"
	doc.Actors[0].ID = "Bad_ID"
	a, b := Document(doc), Document(doc)
	if !reflect.DeepEqual(a, b) || len(a) < 3 {
		t.Fatalf("issues differ between runs or are missing:%s", describe(a))
	}
}

func TestObservationIssues(t *testing.T) {
	v := func(f float64) *float64 { return &f }
	tests := []struct {
		name     string
		mutate   func(s *model.ObservationSet)
		code     string
		path     string
		severity Severity
	}{
		{"wrong kind", func(s *model.ObservationSet) { s.Kind = model.KindStatus }, "VAL-OBS-001", "/kind", SeverityError},
		{"different document", func(s *model.ObservationSet) { s.ModelID = "other-doc" }, "VAL-OBS-002", "/modelId", SeverityError},
		{"unknown indicator", func(s *model.ObservationSet) { s.Observations[0].IndicatorID = "ind-ghost" }, "VAL-REF-001", "/observations/0/indicatorId", SeverityError},
		{"indicator is an outcome", func(s *model.ObservationSet) { s.Observations[0].IndicatorID = "oc-idp-support" }, "VAL-REF-002", "/observations/0/indicatorId", SeverityError},
		{"empty indicator id", func(s *model.ObservationSet) { s.Observations[0].IndicatorID = "" }, "VAL-OBS-005", "/observations/0/indicatorId", SeverityError},
		{"bad date", func(s *model.ObservationSet) { s.Observations[0].Date = "December 15" }, "VAL-DATE-001", "/observations/0/date", SeverityError},
		{"missing date", func(s *model.ObservationSet) { s.Observations[0].Date = "" }, "VAL-DATE-003", "/observations/0/date", SeverityError},
		{"neither value nor level", func(s *model.ObservationSet) { s.Observations[0].Value = nil }, "VAL-OBS-006", "/observations/0", SeverityError},
		{"level not allowed", func(s *model.ObservationSet) { s.Observations[1].Level = "galactic" }, "VAL-OBS-008", "/observations/1/level", SeverityError},
		{"quantitative with only a level", func(s *model.ObservationSet) {
			s.Observations[0].Value = nil
			s.Observations[0].Level = "some"
		}, "VAL-OBS-006", "/observations/0/value", SeverityError},
		{"qualitative with only a value", func(s *model.ObservationSet) {
			s.Observations[1].Level = ""
			s.Observations[1].Value = v(3)
		}, "VAL-OBS-006", "/observations/1/level", SeverityError},
		{"qualitative with value and level", func(s *model.ObservationSet) { s.Observations[1].Value = v(3) }, "VAL-OBS-007", "/observations/1/value", SeverityWarning},
		{"invalid quality", func(s *model.ObservationSet) { s.Observations[0].Quality = "vibes" }, "VAL-ENUM-001", "/observations/0/quality", SeverityError},
		{"unknown evidence", func(s *model.ObservationSet) { s.Observations[0].EvidenceIDs = []string{"ev-ghost"} }, "VAL-REF-001", "/observations/0/evidenceIds/0", SeverityError},
		{"unknown assumption", func(s *model.ObservationSet) { s.Assessments[0].AssumptionID = "asm-ghost" }, "VAL-REF-001", "/assumptionAssessments/0/assumptionId", SeverityError},
		{"assessment of non-assumption", func(s *model.ObservationSet) { s.Assessments[0].AssumptionID = "oc-idp-support" }, "VAL-REF-002", "/assumptionAssessments/0/assumptionId", SeverityError},
		{"assessment without status", func(s *model.ObservationSet) { s.Assessments[0].Status = "" }, "VAL-OBS-004", "/assumptionAssessments/0/status", SeverityError},
		{"assessment bad status", func(s *model.ObservationSet) { s.Assessments[0].Status = "fine" }, "VAL-ENUM-001", "/assumptionAssessments/0/status", SeverityError},
		{"note related to unknown", func(s *model.ObservationSet) { s.Notes[0].RelatedIDs = []string{"ghost"} }, "VAL-REF-003", "/notes/0/relatedIds/0", SeverityWarning},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := obsSet(t)
			tt.mutate(set)
			issues := Observations(fullDoc(t), set)
			for _, i := range issues {
				if i.Code == tt.code && i.Path == tt.path {
					if i.Severity != tt.severity {
						t.Fatalf("%s severity = %s, want %s", tt.code, i.Severity, tt.severity)
					}
					return
				}
			}
			t.Fatalf("expected %s at %s; got:%s", tt.code, tt.path, describe(issues))
		})
	}
}

func TestObservationSourcesWhenMultipleSets(t *testing.T) {
	bad := obsSet(t)
	bad.Observations[0].IndicatorID = "ind-ghost"
	issues := Observations(fullDoc(t), obsSet(t), bad)
	if len(issues) != 1 || issues[0].Source != "observations[1]" {
		t.Fatalf("expected one issue sourced to observations[1]:%s", describe(issues))
	}
	single := Observations(fullDoc(t), bad)
	if len(single) != 1 || single[0].Source != "" {
		t.Fatalf("a single set should carry no source label:%s", describe(single))
	}
}

func TestSchemaIssues(t *testing.T) {
	if issues := Schema(readFixture(t, "full.json"), modelio.FormatJSON); issues != nil {
		t.Fatalf("valid document produced schema issues:%s", describe(issues))
	}
	bad := `{"kind":"OutcomeModel","specVersion":"0.1","id":"a","metadata":{},"context":{},"outcomes":[{"id":"o","level":"huge"}],"bogus":1}`
	issues := Schema([]byte(bad), modelio.FormatJSON)
	want := map[string]string{
		"/context":          "problemStatement",
		"/outcomes/0":       "title",
		"/outcomes/0/level": "one of",
		"":                  "bogus",
	}
	for path, fragment := range want {
		found := false
		for _, i := range issues {
			if i.Path == path && i.Code == "VAL-SCH-001" && strings.Contains(i.Message, fragment) {
				found = true
			}
		}
		if !found {
			t.Errorf("missing schema issue at %q mentioning %q:%s", path, fragment, describe(issues))
		}
	}
	seen := map[string]bool{}
	for _, i := range issues {
		k := i.Path + i.Message
		if seen[k] {
			t.Errorf("duplicate schema issue %s", i)
		}
		seen[k] = true
	}

	for name, data := range map[string]string{
		"unparseable": `{"kind":`,
		"no kind":     `{"id":"x"}`,
		"unknown":     `{"kind":"Nope"}`,
	} {
		got := Schema([]byte(data), modelio.FormatJSON)
		if len(got) != 1 || got[0].Severity != SeverityError {
			t.Errorf("%s: want exactly one error issue, got:%s", name, describe(got))
		}
	}
	if got := Schema([]byte("kind: ["), modelio.FormatYAML); len(got) != 1 || got[0].Code != "VAL-SCH-002" {
		t.Errorf("bad YAML:%s", describe(got))
	}
}

func TestSchemaOtherKinds(t *testing.T) {
	for _, name := range []string{"observations.json", "status.json", "registry.json"} {
		if issues := Schema(readFixture(t, name), modelio.FormatJSON); issues != nil {
			t.Errorf("%s: unexpected schema issues:%s", name, describe(issues))
		}
	}
}

func TestSource(t *testing.T) {
	doc, issues := Source(readFixture(t, "full.json"), modelio.FormatJSON)
	if doc == nil || len(issues) != 0 {
		t.Fatalf("valid source: doc=%v issues:%s", doc != nil, describe(issues))
	}

	// Conforms to the schema but has a dangling reference.
	dangling := strings.Replace(string(readFixture(t, "minimal.json")), `"to": "oc-funding"`, `"to": "oc-nowhere"`, 1)
	doc, issues = Source([]byte(dangling), modelio.FormatJSON)
	if doc == nil || len(issues) != 1 || issues[0].Code != "VAL-REF-001" {
		t.Fatalf("dangling reference: doc=%v issues:%s", doc != nil, describe(issues))
	}

	// Violates the schema: no document is returned.
	doc, issues = Source([]byte(`{"kind":"OutcomeModel"}`), modelio.FormatJSON)
	if doc != nil || len(issues) == 0 {
		t.Fatalf("schema failure: doc=%v issues:%s", doc != nil, describe(issues))
	}

	// Other kinds pass through with schema checking only.
	doc, issues = Source(readFixture(t, "observations.json"), modelio.FormatJSON)
	if doc != nil || issues != nil {
		t.Fatalf("observation set: doc=%v issues:%s", doc != nil, describe(issues))
	}
}

func TestSeverityHelpers(t *testing.T) {
	issues := []Issue{
		{Code: "A", Severity: SeverityError},
		{Code: "B", Severity: SeverityWarning},
		{Code: "C", Severity: SeverityError},
	}
	if !HasErrors(issues) || HasErrors(issues[1:2]) || HasErrors(nil) {
		t.Fatal("HasErrors mismatch")
	}
	if CountSeverity(issues, SeverityError) != 2 || CountSeverity(issues, SeverityInfo) != 0 {
		t.Fatal("CountSeverity mismatch")
	}
	if errs := Errors(issues); len(errs) != 2 || errs[1].Code != "C" {
		t.Fatalf("Errors = %+v", errs)
	}
}

func TestIssueString(t *testing.T) {
	got := Issue{Code: "VAL-REF-001", Severity: SeverityError, Path: "/links/0/to", Source: "observations[1]", Message: "bad"}.String()
	if got != "error VAL-REF-001 observations[1] /links/0/to: bad" {
		t.Fatalf("String() = %q", got)
	}
	if got := (Issue{Code: "X", Severity: SeverityInfo, Message: "m"}).String(); got != "info X: m" {
		t.Fatalf("String() = %q", got)
	}
}

func TestDelta(t *testing.T) {
	a := Issue{Code: "VAL-REF-001", Element: "ln-1", Message: "references unknown element \"x\"", Path: "/links/0/to"}
	b := Issue{Code: "VAL-ENUM-001", Element: "oc-1", Message: "bad", Path: "/outcomes/0/level"}
	c := Issue{Code: "VAL-DOC-005", Message: "empty", Path: "/context/problemStatement"}

	// An unchanged problem whose array index shifted is neither new nor resolved.
	shifted := a
	shifted.Path = "/links/3/to"
	intro, res := Delta([]Issue{a, b}, []Issue{shifted, c})
	if len(intro) != 1 || intro[0].Code != "VAL-DOC-005" || len(res) != 1 || res[0].Code != "VAL-ENUM-001" {
		t.Fatalf("Delta = introduced %+v resolved %+v", intro, res)
	}

	// Duplicates match one-to-one.
	intro, res = Delta([]Issue{a}, []Issue{a, a})
	if len(intro) != 1 || len(res) != 0 {
		t.Fatalf("duplicate Delta = introduced %+v resolved %+v", intro, res)
	}
	if intro, res = Delta(nil, nil); intro != nil || res != nil {
		t.Fatal("empty Delta should be nil")
	}
}

func TestPointerEscapes(t *testing.T) {
	if got := pointer("a/b", "c~d", "0"); got != "/a~1b/c~0d/0" {
		t.Fatalf("pointer = %q", got)
	}
	if got := pointer(); got != "" {
		t.Fatalf("empty pointer = %q", got)
	}
}
