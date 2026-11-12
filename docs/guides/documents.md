# Documents and presentations

`outcome-model export document` assembles one document for people to read from
an outcome model: the projections you choose, in the order you choose. Each
section is drawn from the same model, so the theory of change, the logic
models, the VSEM, and the MEL plan never disagree.

```bash
outcome-model export document outcome-model.yaml --out impact.docx
outcome-model export document outcome-model.yaml --out impact.md
outcome-model export document outcome-model.yaml --out impact.html
```

## Formats

The format follows the `--out` extension.

| Extension | What you get |
|-----------|--------------|
| `.md` | Markdown, with each diagram rendered to SVG in a folder beside it (`impact-diagrams/`). With `--diagrams d2` the diagrams stay as D2 source and no tool is needed |
| `.docx` | A Word document, through Pandoc, with a contents list and the diagrams as images. `--reference-doc` styles it |
| `.html` | A single-file presentation built with React (see [Presentations](#presentations)) |

Diagrams are rendered with the [`d2`](https://d2lang.com) tool, and Word
documents with [Pandoc](https://pandoc.org). If a tool that is needed is not
installed, nothing is written and the exit code is 3.

## Choosing sections

Name a profile with `--profile`, or list sections with `--sections`, in the
order they should appear. The two cannot be combined. Without either, the
`nonprofit` profile is used.

| Profile | Sections | For |
|---------|----------|-----|
| `nonprofit` | theory-of-change, logic-models, vsem, mel, glossary | A nonprofit's board, funders, and staff |
| `strategy` | theory-of-change, vsem, ogsm | A leadership plan |
| `evaluation` | theory-of-change, logic-models, logframe, metrics, mel | Evaluators and funders |
| `improvement` | problem-trees, dmaic | Process improvement work |
| `full` | every section | A complete reference |

| Section | What it holds |
|---------|---------------|
| `theory-of-change` | The theory of change on one page, as a diagram, then the narrative: causes, priorities, assumptions and risks, who does what, and how it is monitored |
| `problem-trees` | The whole model's problem tree, then each pathway's own |
| `logic-models` | For each pathway, the layered logic model diagram and its table |
| `outcome-map` | Outcomes by level, for each pathway |
| `logframe` | The logframe matrix, for each pathway |
| `metrics` | The indicator definitions |
| `vsem` | The [VSEM](vsem.md) |
| `ogsm` | The [OGSM](ogsm.md) |
| `mel` | The [MEL plan](mel.md), with an indicator reference sheet for each indicator |
| `dmaic` | The [DMAIC](dmaic.md) improvement plan |
| `glossary` | The terms the model defines |

```bash
outcome-model export document outcome-model.yaml --sections theory-of-change,vsem,mel --out plan.docx
outcome-model export document outcome-model.yaml --profile strategy --out strategy.html
```

The VSEM, OGSM, MEL, and DMAIC sections are written by the
[prism-roadmap](https://github.com/grokify/prism-roadmap) renderers and placed
under the document's own headings.

## Sections a model cannot draw

Some sections need something the model may not have yet: an OGSM needs an
objective, a VSEM a vision, and a DMAIC plan a pathway marked
`projections: [dmaic]`. Asking for such a section is an error that names what
to add. With `--skip-unavailable` the section is left out instead, and a note
on standard error says why.

## Status

With observations (`-O`), the strategy and MEL sections carry the latest
values, progress, and status, as the matching `export` commands do. Without
them, the document describes the plan alone. Unless `--no-stamp` is given, the
document carries the export stamp (model ID, version, content hash, and
`--as-of` date).

## Presentations

The `.html` output is one file that opens in any browser. It starts with a
title slide and a contents slide, and then gives each section a divider
slide, each heading a slide of its own, and each diagram a full slide. Long
text and long tables run on to further slides, and a table's header repeats.

- **Moving around:** the arrow keys, Page Up and Page Down, Space, Home, and
  End. The address keeps the slide number (`deck.html#12`), so a link opens on
  that slide.
- **Printing to PDF:** open it with `?print` (`deck.html?print`) to lay out
  every slide as a 1280×720 page. The page follows the deck-pdf/v1 contract of
  [presentation-spec](https://github.com/grokify/presentation-spec)'s
  `htmldeck`, which converts it to PDF or one PNG per slide; a browser's
  Print to PDF works too.
- **What it loads:** React from cdnjs. Everything else, including the
  diagrams, is in the file.

## Determinism

The same model, observations, and options always produce the same Markdown,
the same presentation, and the same diagrams, so a document can be kept in
version control and rebuilt without a diff.

## From the MCP server

The `outcome_document` tool returns the same document as Markdown, with the
diagrams as D2 blocks, and lists any sections it skipped.
