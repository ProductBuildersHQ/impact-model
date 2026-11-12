package model

import "slices"

// Evaluation is a planned evaluation: a point at which the organization asks
// whether results happened and why. What an evaluation finds is recorded as
// evidence of type evaluation.
type Evaluation struct {
	ID        string         `json:"id" jsonschema:"required"`
	Title     string         `json:"title" jsonschema:"required"`
	Type      EvaluationType `json:"type,omitempty"`
	Questions []string       `json:"questions,omitempty"`
	// Method is the design and methods, e.g. "comparison with a holdout
	// cohort, plus interviews".
	Method    string           `json:"method,omitempty"`
	Evaluator Evaluator        `json:"evaluator,omitempty"`
	OwnerID   string           `json:"ownerActorId,omitempty"`
	StartDate string           `json:"startDate,omitempty" jsonschema:"format=date"`
	EndDate   string           `json:"endDate,omitempty" jsonschema:"format=date"`
	Status    EvaluationStatus `json:"status,omitempty"`
	Budget    string           `json:"budget,omitempty"`
	// RelatedIDs are the outcomes, links, assumptions, interventions,
	// indicators, and learning questions the evaluation addresses.
	RelatedIDs []string    `json:"relatedIds,omitempty"`
	Provenance *Provenance `json:"provenance,omitempty"`
}

// LearningQuestion is a question on the learning agenda: something the
// organization most needs to find out, and when. The questions a period's
// review took up are recorded in its observation set's learning reviews.
type LearningQuestion struct {
	ID       string `json:"id" jsonschema:"required"`
	Question string `json:"question" jsonschema:"required"`
	// Rationale says why the answer matters: the decision it informs.
	Rationale string `json:"rationale,omitempty"`
	// Methods say how it will be answered, e.g. "cohort interviews".
	Methods []string       `json:"methods,omitempty"`
	OwnerID string         `json:"ownerActorId,omitempty"`
	DueDate string         `json:"dueDate,omitempty" jsonschema:"format=date"`
	Status  QuestionStatus `json:"status,omitempty"`
	Answer  string         `json:"answer,omitempty"`
	// RelatedIDs are the outcomes, links, assumptions, interventions, and
	// indicators the question tests or bears on.
	RelatedIDs []string    `json:"relatedIds,omitempty"`
	Provenance *Provenance `json:"provenance,omitempty"`
}

// ReviewMeeting is a recurring meeting at which the organization reviews what
// it has measured and learned. Each meeting held is recorded as a learning
// review in that period's observation set.
type ReviewMeeting struct {
	ID      string `json:"id" jsonschema:"required"`
	Name    string `json:"name" jsonschema:"required"`
	Purpose string `json:"purpose,omitempty"`
	// Frequency is an ISO 8601 duration, e.g. "P3M".
	Frequency    string   `json:"frequency,omitempty"`
	Participants []string `json:"participants,omitempty"`
	OwnerID      string   `json:"ownerActorId,omitempty"`
	NextDate     string   `json:"nextDate,omitempty" jsonschema:"format=date"`
	// QuestionIDs are the learning questions the meeting takes up.
	QuestionIDs []string `json:"questionIds,omitempty"`
}

// QualityAssurance records how an indicator's data is checked. What the data
// cannot show belongs in the indicator's caveats.
type QualityAssurance struct {
	Checks       []string `json:"checks,omitempty"`
	LastAssessed string   `json:"lastAssessed,omitempty" jsonschema:"format=date"`
}

// Rejected reports whether review rejected the evaluation.
func (e Evaluation) Rejected() bool { return e.Provenance.IsRejected() }

// Rejected reports whether review rejected the learning question.
func (q LearningQuestion) Rejected() bool { return q.Provenance.IsRejected() }

// EvaluationType says what an evaluation is for.
type EvaluationType string

const (
	EvaluationBaseline      EvaluationType = "baseline"
	EvaluationProcess       EvaluationType = "process"
	EvaluationFormative     EvaluationType = "formative"
	EvaluationMidterm       EvaluationType = "midterm"
	EvaluationFinal         EvaluationType = "final"
	EvaluationImpact        EvaluationType = "impact"
	EvaluationDevelopmental EvaluationType = "developmental"
)

// AllEvaluationTypes returns all evaluation types.
func AllEvaluationTypes() []EvaluationType {
	return []EvaluationType{
		EvaluationBaseline, EvaluationProcess, EvaluationFormative, EvaluationMidterm,
		EvaluationFinal, EvaluationImpact, EvaluationDevelopmental,
	}
}

// Valid reports whether v is a known EvaluationType.
func (v EvaluationType) Valid() bool { return slices.Contains(AllEvaluationTypes(), v) }

// Evaluator says who carries out an evaluation.
type Evaluator string

const (
	EvaluatorInternal Evaluator = "internal"
	EvaluatorExternal Evaluator = "external"
	EvaluatorMixed    Evaluator = "mixed"
)

// AllEvaluators returns all evaluators.
func AllEvaluators() []Evaluator {
	return []Evaluator{EvaluatorInternal, EvaluatorExternal, EvaluatorMixed}
}

// Valid reports whether v is a known Evaluator.
func (v Evaluator) Valid() bool { return slices.Contains(AllEvaluators(), v) }

// EvaluationStatus is where an evaluation stands.
type EvaluationStatus string

const (
	EvaluationPlanned    EvaluationStatus = "planned"
	EvaluationInProgress EvaluationStatus = "in_progress"
	EvaluationCompleted  EvaluationStatus = "completed"
	EvaluationCancelled  EvaluationStatus = "cancelled"
)

// AllEvaluationStatuses returns all evaluation statuses.
func AllEvaluationStatuses() []EvaluationStatus {
	return []EvaluationStatus{EvaluationPlanned, EvaluationInProgress, EvaluationCompleted, EvaluationCancelled}
}

// Valid reports whether v is a known EvaluationStatus.
func (v EvaluationStatus) Valid() bool { return slices.Contains(AllEvaluationStatuses(), v) }

// QuestionStatus is where a learning question stands.
type QuestionStatus string

const (
	QuestionOpen     QuestionStatus = "open"
	QuestionAnswered QuestionStatus = "answered"
	QuestionRetired  QuestionStatus = "retired"
)

// AllQuestionStatuses returns all learning question statuses.
func AllQuestionStatuses() []QuestionStatus {
	return []QuestionStatus{QuestionOpen, QuestionAnswered, QuestionRetired}
}

// Valid reports whether v is a known QuestionStatus.
func (v QuestionStatus) Valid() bool { return slices.Contains(AllQuestionStatuses(), v) }
