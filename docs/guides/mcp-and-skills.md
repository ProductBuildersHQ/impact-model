# MCP Server and Agent Skills

The same validated operations the command line offers are available to an AI
agent over the [Model Context Protocol](https://modelcontextprotocol.io), and
five agent skills describe the common workflows.

## The MCP server

```bash
outcome-model mcp --root ./outcome-models
```

The server speaks MCP over standard input and output. Configure an MCP client to
run that command; for example, in a client that takes a JSON configuration:

```json
{
  "mcpServers": {
    "outcome-model": {
      "command": "outcome-model",
      "args": ["mcp", "--root", "./outcome-models"]
    }
  }
}
```

It makes no network calls and no model calls.

### Tools

| Tool | What it does |
|------|--------------|
| `outcome_schema` | List the embedded JSON Schemas or return one |
| `outcome_validate` | Validate an outcome model |
| `outcome_lint` | Lint an outcome model with a profile |
| `outcome_graph` | Query the causal graph |
| `outcome_status` | Compute a status snapshot from observations |
| `outcome_report` | Render a definition or status report |
| `outcome_deck` | Generate a deck as Marp Markdown |
| `outcome_export` | Export for the PRISM tools as JSON: `okr`, `journey`, `v2mom`, `vsem`, `ogsm`, `bsc`, `dmaic`, or `mel`. With observations, the export carries the latest values and status. A model that lacks what the format needs is refused, naming what to add |
| `outcome_document` | One Markdown document with the sections of a profile, or those listed, in order. Diagrams are D2 blocks; sections it skipped are listed |
| `outcome_diff` | Compare two versions of an outcome model |
| `outcome_judge_bundle` | Build the judge bundle |
| `outcome_patch` | Apply a JSON Patch, atomically |

A tool takes an outcome model **inline** (the `document` argument) or **by path**.
`outcome_patch` rejects any patch that adds validation errors, exactly as the
`patch` command does, so an agent cannot leave an outcome model worse than it found it.

### Where an agent can read and write

A path must lie inside `--root`, which defaults to the current directory. These
are refused with the same answer on every platform: a path with a `..` element,
a rooted or drive path outside the root, and a symbolic link that leads outside
the root. An inline document is limited in size.

## Agent skills

A skill is a Markdown file that tells an agent how to carry out one task with the
commands above. They live in the repository's `skills/` directory.

| Skill | Use it to |
|-------|-----------|
| `outcome-draft-from-sources` | Draft an outcome model for an organization from its public pages and reports. Every source is recorded as evidence, and every element is labelled stated, inferred, or hypothesis |
| `outcome-author` | Add, change, or remove outcomes, links, assumptions, indicators, actors, and other elements with validated edits |
| `outcome-status-refresh` | Load a new period's numbers from a spreadsheet, then report where the outcome model stands and what changed |
| `outcome-learning-review` | Prepare and record a learning review for a period |
| `outcome-audience-story` | Draft a short update for one audience from the outcome model's numbers and consented stories |

Each skill is tested: the commands and flags it names must exist, and its
sequence of steps is run against the examples to check it works and repeats.

## Next

[Examples](examples.md) · [CLI reference](cli.md)
