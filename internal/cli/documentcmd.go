package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ProductBuildersHQ/outcome-model/internal/zipfix"
	"github.com/ProductBuildersHQ/outcome-model/pkg/document"
	"github.com/ProductBuildersHQ/outcome-model/pkg/model"
	"github.com/ProductBuildersHQ/outcome-model/pkg/modelio"
	"github.com/ProductBuildersHQ/outcome-model/pkg/service"
)

func newExportDocumentCmd(e *env) *cobra.Command {
	var profiles []string
	for _, name := range document.ProfileNames() {
		var kinds []string
		for _, k := range document.Profiles[name] {
			kinds = append(kinds, string(k))
		}
		profiles = append(profiles, fmt.Sprintf("  %-12s %s", name, strings.Join(kinds, ", ")))
	}
	var sections []string
	for _, k := range document.SectionKinds() {
		sections = append(sections, fmt.Sprintf("  %-17s %s", k, k.Title()))
	}
	cmd := &cobra.Command{
		Use:   "document <model-file> --out <file>.md|.docx|.html",
		Short: "Export the outcome model as one document or presentation",
		Long: fmt.Sprintf(`Export the outcome model as one human-readable document: the sections you
choose, in the order you choose, each a projection of the same model.

The format follows the --out extension:

  .md     Markdown, with each diagram rendered to SVG in a folder next to it
          (<name>-diagrams/), or kept as D2 source with --diagrams d2
  .docx   a Word document, through Pandoc, with the diagrams as images
  .html   a single-file React presentation: one heading per slide and each
          diagram on its own slide; open it in a browser, or add ?print to
          lay out every slide as a page for printing to PDF

Choose the sections with --profile or --sections (not both). Profiles:

%s

Sections:

%s

Diagrams are rendered with the d2 tool, and Word documents with pandoc. If a
needed tool is not installed, nothing is written and the exit code is 3. A
strategy section the model cannot be drawn as, such as an OGSM without an
objective, is an error naming what to add; --skip-unavailable leaves it out
with a note instead.`, strings.Join(profiles, "\n"), strings.Join(sections, "\n")),
		Example: `  outcome-model export document outcome-model.yaml --out impact.docx
  outcome-model export document outcome-model.yaml --profile strategy --out strategy.html
  outcome-model export document outcome-model.yaml --sections theory-of-change,vsem,mel --out plan.md -O obs-q3.yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExportDocument(e, cmd, args[0])
		},
	}
	f := cmd.Flags()
	f.StringP("out", "o", "", "the .md, .docx, or .html file to write (required)")
	if err := cmd.MarkFlagRequired("out"); err != nil {
		panic(fmt.Sprintf("cli: mark --out required: %v", err))
	}
	f.String("profile", "", "a named list of sections (default "+document.DefaultProfile+")")
	f.String("sections", "", "the sections to include, comma-separated, in order")
	f.String("title", "", "the subtitle under the document's title")
	f.String("diagrams", "svg", "for .md: svg renders each diagram to a file; d2 keeps the source")
	f.Bool("skip-unavailable", false, "leave out a section the model cannot be drawn as, with a note")
	f.Bool("contents", true, "list the sections after the title (.md and .docx)")
	f.String("reference-doc", "", "a Pandoc reference document that styles a .docx")
	statusFlags(cmd)
	asOfFlags(cmd)
	return cmd
}

func runExportDocument(e *env, cmd *cobra.Command, path string) error {
	out := flagString(cmd, "out")
	format := strings.ToLower(filepath.Ext(out))
	switch format {
	case ".md", ".docx", ".html":
	default:
		return usageErr("--out must be a .md, .docx, or .html file, not %q", out)
	}
	diagrams := flagString(cmd, "diagrams")
	if diagrams != "svg" && diagrams != "d2" {
		return usageErr("--diagrams must be svg or d2, not %q", diagrams)
	}
	if diagrams == "d2" && format != ".md" {
		return usageErr("--diagrams d2 applies only to .md output")
	}
	sections, err := documentSections(cmd)
	if err != nil {
		return err
	}
	if ref := flagString(cmd, "reference-doc"); ref != "" {
		if format != ".docx" {
			return usageErr("--reference-doc applies only to .docx output")
		}
		if _, err := os.Stat(ref); err != nil {
			return usageErr("--reference-doc: %v", err)
		}
	}

	doc, err := loadForRender(path)
	if err != nil {
		return err
	}
	var st *model.Status
	var sets []*model.ObservationSet
	if len(observationPaths(cmd)) > 0 {
		in, loaded, err := loadStatus(e, cmd, path)
		if err != nil {
			return err
		}
		st, sets = in.Status, loaded
	}
	opts := document.Options{Sections: sections, Subtitle: flagString(cmd, "title"), SkipUnavailable: flagBool(cmd, "skip-unavailable")}
	if !flagBool(cmd, "no-stamp") {
		s, err := stampFor(cmd, doc)
		if err != nil {
			return err
		}
		opts.Stamp = s.Line()
	}
	d, err := service.Document(doc, st, sets, opts)
	if errors.Is(err, document.ErrSectionUnavailable) {
		return &ExitError{Code: ExitUsageErr, Err: fmt.Errorf("%s: %w (or pass --skip-unavailable)", path, err)}
	}
	if err != nil {
		return renderError(path, err)
	}
	for _, s := range d.Skipped {
		fmt.Fprintf(e.errOut, "note: section %s is left out: %s\n", s.Kind, s.Reason)
	}

	ctx := cmd.Context()
	if format == ".md" && diagrams == "d2" {
		return writeOut(out, d.Markdown(document.MarkdownOptions{Contents: flagBool(cmd, "contents")}))
	}
	d2, err := exec.LookPath("d2")
	if err != nil && len(d.Diagrams()) > 0 {
		return &ExitError{Code: ExitToolMissing, Err: errors.New("d2 is not installed or not on PATH; install it to render the diagrams, or write .md with --diagrams d2 to keep their source")}
	}
	switch format {
	case ".md":
		dir := strings.TrimSuffix(out, filepath.Ext(out)) + "-diagrams"
		svgs, err := renderD2(ctx, d2, d.Diagrams())
		if err != nil {
			return err
		}
		if len(svgs) > 0 {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return &ExitError{Code: ExitUsageErr, Err: err}
			}
		}
		for _, g := range d.Diagrams() {
			if err := modelio.WriteFile(filepath.Join(dir, g.ID+".svg"), svgs[g.ID]); err != nil {
				return &ExitError{Code: ExitUsageErr, Err: err}
			}
		}
		rel := filepath.Base(dir)
		md := d.Markdown(document.MarkdownOptions{
			Contents:    flagBool(cmd, "contents"),
			DiagramPath: func(g *document.Diagram) string { return rel + "/" + g.ID + ".svg" },
		})
		return writeOut(out, md)
	case ".html":
		svgs, err := renderD2(ctx, d2, d.Diagrams())
		if err != nil {
			return err
		}
		html, err := d.Deck(document.DeckOptions{SVG: func(g *document.Diagram) ([]byte, bool) {
			svg, ok := svgs[g.ID]
			return svg, ok
		}})
		if err != nil {
			return &ExitError{Code: ExitUsageErr, Err: err}
		}
		return writeOut(out, string(html))
	default:
		return exportDocumentDOCX(cmd, d, d2)
	}
}

// documentSections reads --sections or --profile.
func documentSections(cmd *cobra.Command) ([]document.SectionKind, error) {
	list, profile := flagString(cmd, "sections"), flagString(cmd, "profile")
	switch {
	case list != "" && profile != "":
		return nil, usageErr("give --sections or --profile, not both")
	case list != "":
		s, err := document.ParseSections(list)
		if err != nil {
			return nil, usageErr("--sections: %v", err)
		}
		return s, nil
	case profile != "":
		s, err := document.ProfileSections(profile)
		if err != nil {
			return nil, usageErr("--profile: %v", err)
		}
		return s, nil
	}
	return nil, nil
}

// exportDocumentDOCX renders the diagrams into a scratch directory and has
// pandoc convert the Markdown, with a table of contents, to Word.
func exportDocumentDOCX(cmd *cobra.Command, d *document.Document, d2 string) error {
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return &ExitError{Code: ExitToolMissing, Err: errors.New("pandoc is not installed or not on PATH; install it, or write .md or .html instead")}
	}
	dir, err := os.MkdirTemp("", "outcome-model-document-")
	if err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "outcome-model: remove %s: %v\n", dir, err)
		}
	}()
	ctx := cmd.Context()
	svgs, err := renderD2(ctx, d2, d.Diagrams())
	if err != nil {
		return err
	}
	for _, g := range d.Diagrams() {
		if err := modelio.WriteFile(filepath.Join(dir, g.ID+".svg"), svgs[g.ID]); err != nil {
			return &ExitError{Code: ExitUsageErr, Err: err}
		}
	}
	md := d.Markdown(document.MarkdownOptions{
		Contents:    flagBool(cmd, "contents"),
		DiagramPath: func(g *document.Diagram) string { return g.ID + ".svg" }, // found through --resource-path
		ImageAttrs:  "width=100%",
	})
	src := filepath.Join(dir, "document.md")
	if err := modelio.WriteFile(src, []byte(md)); err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	abs, err := filepath.Abs(flagString(cmd, "out"))
	if err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	args := []string{"-f", "markdown", "-t", "docx", "-o", abs, "--resource-path", dir}
	if ref := flagString(cmd, "reference-doc"); ref != "" {
		args = append(args, "--reference-doc", ref)
	}
	args = append(args, src)
	if err := runToolEnv(ctx, pandocEnv(cmd), pandoc, args...); err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	if err := zipfix.NormalizeFile(abs); err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	return nil
}

// renderD2 renders each diagram to SVG with the d2 tool, keyed by diagram ID.
func renderD2(ctx context.Context, d2 string, diagrams []*document.Diagram) (map[string][]byte, error) {
	out := map[string][]byte{}
	if len(diagrams) == 0 {
		return out, nil
	}
	dir, err := os.MkdirTemp("", "outcome-model-d2-")
	if err != nil {
		return nil, &ExitError{Code: ExitUsageErr, Err: err}
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Fprintf(os.Stderr, "outcome-model: remove %s: %v\n", dir, err)
		}
	}()
	for _, g := range diagrams {
		src := filepath.Join(dir, g.ID+".d2")
		svg := filepath.Join(dir, g.ID+".svg")
		if err := modelio.WriteFile(src, []byte(g.D2)); err != nil {
			return nil, &ExitError{Code: ExitUsageErr, Err: err}
		}
		if err := runTool(ctx, d2, src, svg); err != nil {
			return nil, &ExitError{Code: ExitUsageErr, Err: fmt.Errorf("diagram %s: %w", g.ID, err)}
		}
		data, err := os.ReadFile(svg)
		if err != nil {
			return nil, &ExitError{Code: ExitUsageErr, Err: err}
		}
		out[g.ID] = data
	}
	return out, nil
}

func writeOut(path, text string) error {
	if err := modelio.WriteFile(path, []byte(text)); err != nil {
		return &ExitError{Code: ExitUsageErr, Err: err}
	}
	return nil
}
