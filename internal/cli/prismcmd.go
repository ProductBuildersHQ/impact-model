package cli

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ProductBuildersHQ/outcome-model/pkg/impact"
	"github.com/ProductBuildersHQ/outcome-model/pkg/impactio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/prismx"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

func newPrismExportCmd(e *env, kind string) *cobra.Command {
	short, long := "Export OKRs for the PRISM tools", `Export the outcomes the organization answers for as an OKR document: each
outcome with measured indicators is an objective, and each of its indicators a
key result with baseline, target, current value, score, status, and
confidence. With observations (-O), the key results carry the latest values.`
	switch kind {
	case "journey":
		short, long = "Export outcome journeys for the PRISM tools", `Export each numeric indicator as an outcome journey: where the measure stands
now, from the latest observation or the baseline, and the dated targets it is
meant to reach. With observations (-O), target confidence comes from the
status.`
	case "v2mom":
		short, long = "Export the impact model as a V2MOM", `Export the impact model as a V2MOM (Vision, Values, Methods, Obstacles,
Measures):

  Vision     context.vision
  Values     context.values, in order
  Methods    one per pathway; its interventions are the method's projects
  Measures   each pathway's focus: North Star, one metric, guardrails
  Obstacles  risks, and high or critical assumptions not yet shown to hold
             (their assessed status with -O); under the method of the
             pathway they apply to, otherwise for the whole V2MOM

With observations (-O), the measures carry current values, progress, and
status. A model without a vision, values, or a pathway is refused, with a
message naming what to add.`
	case "vsem":
		short, long = "Export the impact model as a VSEM", `Export the impact model as a VSEM (Vision, Strategy, Execution, Metrics):

  Vision     context.vision, and context.mission
  Strategy   one per pathway, with the resources it commits
  Execution  each pathway's interventions, with owners, status, and dates;
             type program is a program, any other an initiative
  Metrics    each pathway's focus: North Star, one metric, guardrails;
             then, for the vision, indicators of impact-level outcomes

An intervention or indicator in more than one pathway is listed under the
first. With observations (-O), the metrics carry current values, progress,
and status. A model without a vision or mission, or without a pathway, is
refused, with a message naming what to add.`
	case "dmaic":
		short, long = "Export DMAIC pathways as a DMAIC metrics document", `Export the pathways marked projections: [dmaic] as a DMAIC (Define, Measure,
Analyze, Improve, Control) metrics document:

  Category     one per marked pathway
  Metric       each of its indicators with a numeric baseline and target;
               others are listed as notes on standard error
  Data points  the indicator's observations (-O), without disaggregated slices
  Limits       an individuals chart, only from 20 or more observations
  Root causes  the causes the measured outcome addresses
  Initiatives  the pathway's interventions aimed at the measured outcome
  Phase        Control once achieved; Measure until there are enough points
               for limits; Improve while an intervention is in progress;
               otherwise Analyze

Process capability is not filled: it needs specification limits, which an
impact model does not record. A model with no marked pathway, or none with a
numeric indicator, is refused.`
	case "mel":
		short, long = "Export the impact model as a MEL plan", `Export the impact model as a monitoring, evaluation, and learning (MEL) plan:

  Results      the outcomes at their levels, each contributing to the
               outcomes its links and precondition groups lead to (feedback
               links, which are loops, are left out)
  Activities   the interventions, producing the outcomes they target
  Indicators   reference sheets: definition, baseline, targets, data source,
               method, frequency, owner, quality checks, and caveats
  Assumptions  each with its test as how it is monitored
  Evaluations  the evaluation plan
  Learning     the learning questions; a link one names becomes both ends
  Reviews      the review meetings
  Reports      the audiences: needs, channels, frequency, and indicators

With observations (-O), the indicators carry their latest values. A model
with no outcome that review has not rejected is refused.`
	case "ogsm":
		short, long = "Export the impact model as an OGSM", `Export the impact model as an OGSM (Objective, Goals, Strategies, Measures):

  Objective   context.objective, or the model's one impact outcome, due at
              the end of the time horizon
  Goals       the indicators of the outcomes the objective stands for
  Strategies  one per pathway, serving the goals its outcomes lead to
  Measures    each pathway's focus: North Star, one metric, guardrails,
              leaving out any that is already a goal
  Actions     each pathway's interventions, with owners, status, and dates

An intervention or indicator in more than one pathway is listed under the
first. With observations (-O), the goals and measures carry current values,
progress, and status. A model with several impact outcomes and no
context.objective, or without a pathway, is refused, with a message naming
what to add.`
	case "bsc":
		short, long = "Export the impact model as a Balanced Scorecard", `Export the impact model as a Balanced Scorecard:

  Mission, vision  context.mission and context.vision
  Perspectives     those the outcomes use, top first: mission, financial,
                   customer, internal process, learning and growth
  Objectives       each outcome with a scorecardPerspective, owned by its
                   ownerActorId
  Strategy map     each objective contributes to the nearest objectives its
                   outcome leads to, through causal links (not feedback) and
                   precondition groups
  Themes           each pathway that holds an objective
  Measures         each objective's indicators, lead or lag
  Initiatives      each intervention, serving the objectives among its
                   targets, or else the nearest ones they lead to

With observations (-O), the measures carry current values, progress, and
status. A model with no outcome in a scorecard perspective is refused, with a
message naming what to add.`
	}
	cmd := &cobra.Command{
		Use:     kind + " <model-file>",
		Short:   short,
		Long:    long + "\n\nThe output is JSON.",
		Example: fmt.Sprintf("  impact-model export %s impact-model.yaml -O obs-q3.yaml -o %s.json", kind, kind),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrismExport(e, cmd, kind, args[0])
		},
	}
	statusFlags(cmd)
	cmd.Flags().String("as-of", "", "the status date, YYYY-MM-DD; default: the latest date in the observations")
	cmd.Flags().StringP("out", "o", "", "write to this file instead of standard output")
	return cmd
}

func runPrismExport(e *env, cmd *cobra.Command, kind, path string) error {
	doc, err := loadForRender(path)
	if err != nil {
		return err
	}
	var st *impact.Status
	var sets []*impact.ObservationSet
	if len(observationPaths(cmd)) > 0 || flagString(cmd, "as-of") != "" {
		in, loaded, err := loadStatus(e, cmd, path)
		if err != nil {
			return err
		}
		st, sets = in.Status, loaded
	}
	exp, err := service.ExportPRISM(kind, doc, st, sets)
	for _, o := range exp.Omissions {
		fmt.Fprintf(e.errOut, "note: %s is left out: %s\n", o.IndicatorID, o.Reason)
	}
	if errors.Is(err, prismx.ErrNotV2MOMReady) || errors.Is(err, prismx.ErrNotVSEMReady) || errors.Is(err, prismx.ErrNotDMAICReady) || errors.Is(err, prismx.ErrNotMELReady) || errors.Is(err, prismx.ErrNotOGSMReady) || errors.Is(err, prismx.ErrNotBSCReady) {
		return &ExitError{Code: ExitUsageErr, Err: fmt.Errorf("%s: %w", path, err)}
	}
	if err != nil {
		return renderError(path, err)
	}
	v := exp.Document
	var buf bytes.Buffer
	if err := writeJSON(&buf, v); err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	if out := flagString(cmd, "out"); out != "" {
		if err := impactio.WriteFile(out, buf.Bytes()); err != nil {
			return &ExitError{Code: ExitUsageErr, Err: err}
		}
		return nil
	}
	fmt.Fprint(e.out, buf.String())
	return nil
}
