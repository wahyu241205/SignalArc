# Agent API

SignalArc exposes agent-readable market intelligence as structured API data that can be consumed by AI agents, institutional dashboards, monitoring systems, and reporting tools. In the current implementation, this means compact market discovery data, not autonomous trading or agent-managed wallet execution.

Primary planned API base:

```text
https://api.signalarc.fun
```

The production API is not live yet. DNS and deployment are not completed.

## Current Agent-Readable Endpoint

| Method | Path | Status |
| --- | --- | --- |
| GET | `/agent/markets` | Implemented locally. |

## GET /agent/markets

Returns up to 50 markets in a compact shape suitable for market discovery and signal-reading workflows.

```bash
curl https://api.signalarc.fun/agent/markets
```

Response shape:

```json
{
  "markets": [
    {
      "id": "10000000-0000-4000-8000-000000000003",
      "title": "Will SignalArc complete its public docs?",
      "status": "OPEN",
      "category": "product",
      "collateral_asset": "USDC",
      "chain": "Arc Testnet",
      "closes_at": "2026-06-01T00:00:00Z",
      "resolution_source": "Project repository evidence"
    }
  ]
}
```

Fields:

| Field | Meaning |
| --- | --- |
| `id` | Market UUID. |
| `title` | Human-readable market question/title. |
| `status` | Current backend market status. |
| `category` | Optional market category. |
| `collateral_asset` | Stored collateral asset label, typically `USDC` in current data. |
| `chain` | Stored chain label. |
| `closes_at` | RFC3339 close timestamp. |
| `resolution_source` | Optional source or method for resolution evidence. |

Known error:

```json
{
  "error": {
    "code": "markets_list_failed",
    "message": "failed to list markets"
  }
}
```

## Example Use Cases

- Market discovery for dashboards and automated reports.
- Event-risk monitoring across open or closing markets.
- Probability-signal reading from structured market metadata.
- Institutional intelligence dashboards.
- AI agent reporting that summarizes market status and resolution sources.

## Current Limitations

- Autonomous trading is not implemented.
- Circle Agents integration is not implemented.
- API key enforcement is not implemented.
- Paid access and billing are not implemented.
- Rate limits are not implemented.
- No production SLA is available.
- Agent wallet policy is not implemented.
- Agent behavior beyond reading the implemented API is unknown / not documented.

## Planned Roadmap

The following items are planned concepts, not implemented behavior:

- API keys.
- Scoped agent access.
- Rate limits.
- Agent wallet policy.
- Circle Agents integration if later approved and implemented.
- Expanded market intelligence endpoints.
- Production monitoring and SLA documentation.
