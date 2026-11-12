package cli

import (
	"github.com/spf13/cobra"

	"github.com/ProductBuildersHQ/outcome-model/pkg/edit"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
)

func newFocusCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "focus <model-file> <pathway-id> --one-metric <indicator-id> --from <date>",
		Short: "Change a pathway's one metric that matters from a date",
		Long: `Make an indicator the pathway's one metric that matters from a date. The
focus it replaces ends the day before, and stays on record, so the pathway
keeps a history of what the team steered by in each period.

A new focus must start after the current one started, and must not overlap a
period that has already ended.`,
		Example: `  outcome-model pathway focus outcome-model.yaml pw-onboarding --one-metric ind-first-week-recon \
    --from 2027-01-01 --stage activation --reason "Activation is the bottleneck"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doc, err := modelio.LoadDocument(args[0])
			if err != nil {
				return &ExitError{Code: ExitUsageErr, Err: err}
			}
			p, err := edit.FocusOn(doc, args[1], model.OneMetric{
				IndicatorID: flagString(cmd, "one-metric"),
				From:        flagString(cmd, "from"),
				To:          flagString(cmd, "to"),
				Stage:       flagString(cmd, "stage"),
				Reason:      flagString(cmd, "reason"),
			})
			if err != nil {
				return usageErr("%v", err)
			}
			return e.applyElementPatch(cmd, args[0], p)
		},
	}
	cmd.Flags().String("one-metric", "", "indicator ID to focus on (required)")
	cmd.Flags().String("from", "", "first day of the focus, YYYY-MM-DD (required)")
	cmd.Flags().String("to", "", "last day of the focus, YYYY-MM-DD (default: open)")
	cmd.Flags().String("stage", "", "the stage the focus fits, e.g. activation")
	cmd.Flags().String("reason", "", "why this metric matters most for the period")
	for _, f := range []string{"one-metric", "from"} {
		if err := cmd.MarkFlagRequired(f); err != nil {
			panic("cli: mark --" + f + " required: " + err.Error())
		}
	}
	addCommonFlags(cmd)
	return cmd
}
