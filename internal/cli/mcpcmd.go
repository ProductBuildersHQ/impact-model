package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/ProductBuildersHQ/outcome-model/pkg/mcp"
)

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Serve the tools to an AI agent over MCP (stdio)",
		Long: `Start a Model Context Protocol server on standard input and output, so an AI
assistant can validate, lint, query, measure, report on, and edit an outcome model.

Tools: outcome_schema, outcome_validate, outcome_lint, outcome_graph, outcome_status, outcome_report,
outcome_patch, outcome_deck, outcome_judge_bundle, outcome_diff. A tool takes an outcome model inline
(document) or by path. A path must lie inside --root, which defaults to the
current directory: ".." is refused, and so is a symbolic link that leaves the
root. outcome_patch rejects any patch that adds validation errors.

The server makes no network calls and no model calls.`,
		Example: `  outcome-model mcp --root ./outcome-models`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := flagString(cmd, "root")
			srv, err := mcp.New(root)
			if err != nil {
				return &ExitError{Code: ExitUsageErr, Err: err}
			}
			parent := cmd.Context()
			if parent == nil {
				parent = context.Background()
			}
			ctx, stop := signal.NotifyContext(parent, os.Interrupt)
			defer stop()
			if err := srv.Run(ctx); err != nil && ctx.Err() == nil {
				return &ExitError{Code: ExitUsageErr, Err: err}
			}
			return nil
		},
	}
	cmd.Flags().String("root", ".", "the directory tool paths are confined to")
	return cmd
}
