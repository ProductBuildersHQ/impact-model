## DMAIC plan

DMAIC (Define, Measure, Analyze, Improve, Control) is a cycle for improving a program that is measured often. One metric can serve several phases. A metric serves **Define** when it sizes the problem (a long-term or impact outcome, or one that relieves an effect), **Measure** always, **Analyze** when it is broken down by group, **Improve** when it moves soon after an intervention changes, and **Control** when it is measured at least once a year. `export dmaic` includes a metric once it has a numeric baseline and target, and draws control limits once it has 20 measurements.

### Supply side: providers and platforms

**Define.** Core problem: Providers and platforms do not ship a common delegated-access standard.

Root causes to analyze, deepest first:

- Providers see no customer demand for a common standard. (`cause-no-demand-signal`, underlying)

**Improve.** The interventions an improvement would change:

- Run interoperability events (`int-interop`), in progress
- Publish reference implementations (`int-reference`), planned

**Metrics.** 6 metrics: 1 have a baseline, 1 a numeric target, and 0 have everything control limits need.

| Metric | D | M | A | I | C | Cadence | Data | Still needs |
|--------|---|---|---|---|---|---------|------|-------------|
| Identity providers with generally available support (`ind-idp-count`) | ✓ | ✓ | ✓ | ✓ | ✓ | quarterly | public | 20+ measurements a year for control limits |
| Reference implementation downloads (`ind-reference-downloads`) |  | ✓ |  | ✓ |  | — | internal_new | baseline; numeric target; data source; frequency |
| Providers with GA support (North America) (`ind-idp-count-na`) | ✓ | ✓ |  | ✓ |  | — | partner_held | baseline; numeric target; data source; frequency |
| Providers with GA support (Europe) (`ind-idp-count-eu`) | ✓ | ✓ |  | ✓ |  | — | unknown | baseline; numeric target; data source; frequency |
| Providers with GA support (total) (`ind-idp-total`) | ✓ | ✓ |  | ✓ |  | — | — | baseline; numeric target; data source; frequency |
| Stories of provider adoption (`ind-adoption-stories`) | ✓ | ✓ |  |  |  | — | internal_new | baseline; numeric target; data source; frequency |

