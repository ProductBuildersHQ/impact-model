# Market adoption example (fictional)

Everything here is **fictional**: the Delegation Receipt standard, the Delegation
Standard Working Group, every organization and competitor, all numbers, and all
`example.org` URLs. It is invented to show the format and resembles no real
company, product, protocol, or standard.

The theory follows an open standard driven to adoption in a market with two
competing approaches (a proprietary single-vendor scheme and an extension of an
older authorization framework), using the Adopt-Adapt-Expand-Respond (AAER)
market-systems framing.

## Files

- `outcome-model.yaml`, `actors.yaml` (shared registry), `observations-2026-h2.yaml`

## Features shown

- AAER categories with adoption, practice, policy, and norm changes
- Five outcome levels; impact above the ceiling of accountability
- `at_least` logic, `all` and `any` groups, and a nested group: (A AND B) OR C
- Reinforcing and balancing feedback loops; a negative-polarity link; lags
- Interventions with capabilities and initiative references
- Assumptions spanning untested to invalidated, each tied to indicators
- Quantitative, qualitative (ordered levels), and maintain-direction indicators;
  shared metrics, a roll-up, data availability, and contribution claims
- Evidence that supports and contradicts; provenance with all five bases,
  including rejected elements
- Risks, external factors, decisions, consultations, revisions, alignments

## Commands

Run from this directory:

```bash
outcome-model validate outcome-model.yaml actors.yaml
outcome-model validate observations-2026-h2.yaml --model outcome-model.yaml
outcome-model graph outcome-model.yaml
outcome-model fmt --check outcome-model.yaml actors.yaml observations-2026-h2.yaml
```
