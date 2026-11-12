# Getting Started

## Install

```bash
go install github.com/ProductBuildersHQ/outcome-model/cmd/outcome-model@latest
outcome-model --version
```

The tool is a single Go binary. Three optional programs add output formats; the
command that needs one tells you if it is missing and exits with code 3
without writing anything:

| Program | Needed for |
|---------|-----------|
| [Marp](https://marp.app) (`marp`) | Rendering decks to PDF, PPTX, or HTML |
| [Mermaid CLI](https://github.com/mermaid-js/mermaid-cli) (`mmdc`) | Diagrams as images in decks and Word documents |
| [Pandoc](https://pandoc.org) (`pandoc`) | Word and PowerPoint export |

Markdown reports, CSV tables, Marp Markdown, and everything that checks or
measures an outcome model need none of them.

## Create an outcome model

```bash
outcome-model init --id my-program --name "My program" --out outcome-model.yaml
```

This writes a small starter outcome model that already passes validation, with `TODO`
placeholders to replace. An outcome model is a YAML or JSON file; the format follows the
file extension. Open it in any editor, or change it through the commands in
[Authoring an outcome model](authoring.md).

## Check it

```bash
outcome-model validate outcome-model.yaml   # schema and structure
outcome-model lint outcome-model.yaml       # is the outcome model sound?
outcome-model fmt outcome-model.yaml        # write it in canonical form
```

`validate` says the document is well formed: the schema, duplicate IDs,
references that do not resolve, invalid values, bad dates. `lint` says whether
the outcome model is sound. A new outcome model is GO under `--profile draft`; the default
`standard` profile also reports warnings.

`fmt` rewrites a file in canonical form: a stable field order, two-space
indentation, and a trailing newline, so diffs stay readable. Formatting is
idempotent. YAML comments are not kept, because tool-written files are
canonical; put notes in `description` and `narrative` fields.

## Ask questions of the causal graph

```bash
outcome-model graph outcome-model.yaml                      # a summary
outcome-model graph outcome-model.yaml --orphans            # outcomes nothing leads to
outcome-model graph outcome-model.yaml --dead-ends          # outcomes that lead nowhere
outcome-model graph outcome-model.yaml --upstream --from oc-adoption
outcome-model graph outcome-model.yaml --format dot | dot -Tsvg > graph.svg
```

## Get a first report

```bash
outcome-model report definitions outcome-model.yaml --out definitions.md
```

This renders the problem tree (once the model records causes or effects), the
pathways, logic models, outcome maps, logframes, metric definitions, and
glossary as one Markdown document. It needs no observations, so a workshop can review it before
any data exists.

To start from the problem rather than the solution, draw the problem tree: the
core problem, its causes below, and its effects above. See the
[problem tree guide](problem-tree.md).

```bash
outcome-model report problem-tree outcome-model.yaml --out problem-tree.md
```

## Measure progress

Once you have numbers, put them in an observation set (see
[Measuring progress](measurement.md)) and ask how far the outcome model has come:

```bash
outcome-model status outcome-model.yaml -O observations.yaml
outcome-model report exec outcome-model.yaml -O observations.yaml -o exec.md
```

## Working in a Git repository

Specifications are meant to be diffed and reviewed. Two habits help:

- Run `outcome-model fmt --check *.yaml` in CI. It exits with code 1 if any
  file is not in canonical form.
- Add a `.gitattributes` line so Windows checkouts keep LF line endings. The
  canonical form is byte-exact LF, and a checkout that converts to CRLF will
  look non-canonical.

  ```text
  * text=auto eol=lf
  *.csv -text
  ```

  CSV files exported by the tool follow RFC 4180 and keep their CRLF, so they
  are left alone.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | The input is invalid, or a change was rejected, or a check did not pass |
| 2 | A usage or I/O error |
| 3 | An optional external tool (`marp`, `mmdc`, `pandoc`) is missing |

## Next

Read [Authoring an outcome model](authoring.md), or open a complete outcome model in
[Examples](examples.md).
