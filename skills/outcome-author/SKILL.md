---
name: outcome-author
description: Create or revise a Theory of Change with validated edits. Use when asked to add, change, or remove outcomes, links, assumptions, indicators, actors, or other elements of a Theory of Change, or to check and tidy one.
---

# Author and revise a Theory of Change

## Purpose

Change a theory without ever leaving the document in a broken state. Every
edit is checked before it is written, and every edit is small enough to read.

## Rules

- Judgment is limited to deciding what the theory should say. Every check
  comes from a `outcome-model` command.
- Change the document only with the element commands (`outcome-model
  outcome add`, `link set`, `assumption rm --cascade`, and so on) or with
  `outcome-model patch`. Never edit the file by hand.
- Run an edit with `--dry-run` first when it is not obviously safe. A
  rejected edit writes nothing, and the command says which issues it would
  have added.
- Outcomes are states, not activities: "Leaders gain coaching skills", not
  "Run coaching". Interventions are the activities.
- Give every link a reason (`--rationale`) and name the assumptions it rests
  on. A link into a medium-term, long-term, or impact outcome without one
  is a lint warning.
- Do not mark an element `confirmed` or `rejected` unless the people it
  concerns said so. Use `--basis stated|inferred|hypothesis`, and say why in
  `--basis-note`.
- Removing an element that others refer to needs `--cascade`; read what it
  would remove with `--dry-run` first.

## Command sequence

1. Check the starting point:

   ```bash
   outcome-model validate outcome-model.yaml
   ```

2. Preview each edit, then make it:

   ```bash
   outcome-model outcome set outcome-model.yaml <id> --title "<new title>" --dry-run
   outcome-model outcome set outcome-model.yaml <id> --title "<new title>"
   outcome-model link add outcome-model.yaml --id <id> --from <outcome> --to <outcome> --rationale "<why>"
   outcome-model assumption add outcome-model.yaml --id <id> --statement "<what must be true>" --applies-to <link-id>
   ```

   For a change several elements share, write one patch and apply it as a unit:

   ```bash
   outcome-model patch outcome-model.yaml change.patch.json --dry-run
   outcome-model patch outcome-model.yaml change.patch.json
   ```

3. Check the result and the quality of the theory:

   ```bash
   outcome-model validate outcome-model.yaml
   outcome-model lint outcome-model.yaml
   ```

4. Show what changed since the last version, if one is saved:

   ```bash
   outcome-model diff previous.yaml outcome-model.yaml
   ```

## Checks to report

- Validation result.
- Lint: errors, warnings, and notes, with the rule IDs of anything new.
- The diff: elements added, changed, and removed.

## Outputs

The edited `outcome-model.yaml`. Tell the user which warnings remain and offer to
fix them.
