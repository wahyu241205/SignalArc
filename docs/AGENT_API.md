# Agent API

SignalArc exposes a framework-neutral Agent API for Arc Testnet YES/NO prediction market workflows. Any external agent runtime can integrate through HTTP: Custom GPTs, Claude Code, LangChain, ElizaOS, Telegram or WhatsApp bots, Hermes, OpenClaw, internal dashboards, or a plain script.

SignalArc owns the API, wallet/session records, intent lifecycle, execution coordination, and read models. It does not own the user's agent runtime, prompt loop, messaging channel, or autonomous policy engine.

## Base URLs

Production:

```text
https://api.signalarc.fun
```

Local development:

```text
http://localhost:4000
```

Useful local health checks:

```bash
curl http://localhost:4000/health
curl http://localhost:4000/readyz
curl http://localhost:4000/schema/validate
```

## Non-Goals

- No smart contract changes.
- No Smart Contract V2.
- No channel-specific bot design.
- No autonomous unattended trading by default.
- No mainnet funding.
- No arbitrary transfers, swaps, withdrawals, or deposits.
- No private keys, seed phrases, Circle API keys, Circle session tokens, or credential paths in API responses.

## Agent Identity

`agent_id` is the stable caller-supplied identity for an external agent/user pair. Reuse the same `agent_id` across onboarding, wallet/session reads, intent creation, confirmation, execution, portfolio, and activity.

Validation rules:

- Required on onboarding, wallet registration, wallet/session reads, portfolio/activity reads, and intent requests where an agent identity is present.
- Must start with `agent_`.
- Must be at least 10 characters long.
- May contain ASCII letters, digits, underscores, and hyphens.
- Must not be a generic placeholder such as `signalarc-gpt-agent`, `agent_desi_001`, `default`, `test`, `demo`, `user`, `agent`, or `chatgpt`.

Safe example used in this guide:

```text
agent_demo_custom_001
```

## Wallet And Session Model

An agent wallet is a SignalArc database record keyed by `agent_id`. Current production flows resolve and register a Circle Agent Wallet on `ARC-TESTNET`.

Important fields:

| Field | Meaning |
| --- | --- |
| `agent_wallet_address` | Registered wallet address SignalArc will target. |
| `wallet_provider` | Current production provider is `circle_agent_wallet`. |
| `chain` | Current execution and faucet flows require `ARC-TESTNET`. |
| `allowed_actions` | Wallet/session scoped action allowlist. |
| `policy_metadata` | Optional policy object; `max_trade_amount` is enforced for buy intents when present. |
| `status` | Wallet must be `active` for execution and faucet flows. |

SignalArc never asks the external agent to provide a private key. Circle OTPs, request IDs, session files, tokens, and credentials are not exposed in public responses.

## Supported Actions

Intent actions:

- `create_market`
- `buy_yes`
- `buy_no`
- `close_market`
- `resolve_market`
- `claim_payout`
- `cancel_market`
- `claim_refund`

`buy_yes` and `buy_no` require a positive finite decimal `amount`. `create_market` requires `question`, `resolver`, `collateral_token`, and `close_timestamp` as UTC RFC3339 or unix-seconds string.

## Lifecycle

Transaction actions must follow this lifecycle:

```text
create intent preview -> confirm intent -> execute confirmed intent
```

`POST /agent/intents` validates shape and returns a preview. It does not broadcast.

`POST /agent/intents/{intent_id}/confirm` confirms the preview and returns an execution plan. It does not broadcast.

`POST /agent/intents/{intent_id}/execute` executes only a confirmed intent through the configured wallet provider. Unconfirmed execution returns `409 agent_intent_not_confirmed`.

Configured backends persist durable intent and execution records in Postgres. The in-process store remains a fallback for tests and unconfigured runtimes.

## Idempotency

Intent creation supports idempotency with:

- `agent_id`
- `source_client`
- `client_request_id`

If the same agent sends the same `source_client + client_request_id` again, configured durable storage returns the existing intent instead of creating a duplicate.

Example values:

```json
{
  "agent_id": "agent_demo_custom_001",
  "source_client": "custom-agent",
  "client_request_id": "demo-001"
}
```

## Endpoint Summary

| Capability | Endpoint |
| --- | --- |
| Health | `GET /health` |
| Readiness | `GET /readyz` |
| Schema validation | `GET /schema/validate` |
| Agent market discovery | `GET /agent/markets` |
| Full market list | `GET /markets` |
| Full market detail | `GET /markets/{market_id}` |
| Start onboarding | `POST /agent/onboarding/start` |
| Verify onboarding OTP | `POST /agent/onboarding/verify` |
| Onboarding lookup | `GET /agent/onboarding/{onboarding_id}` |
| Active session | `GET /agent/sessions/{agent_id}` |
| Registered wallet | `GET /agent/wallets/{agent_id}` |
| Wallet balance | `GET /agent/wallets/{agent_id}/balance` |
| Testnet faucet | `POST /agent/wallets/{agent_id}/faucet` |
| Create intent | `POST /agent/intents` |
| Get intent | `GET /agent/intents/{intent_id}` |
| Confirm intent | `POST /agent/intents/{intent_id}/confirm` |
| Execute intent | `POST /agent/intents/{intent_id}/execute` |
| List intent executions | `GET /agent/intents/{intent_id}/executions` |
| Agent portfolio | `GET /agent/portfolio/{agent_id}` |
| Agent activity | `GET /agent/activity/{agent_id}` |

## Curl Examples

Set the base URL once:

```bash
BASE_URL="https://api.signalarc.fun"
AGENT_ID="agent_demo_custom_001"
```

For local development, use:

```bash
BASE_URL="http://localhost:4000"
```

### List Agent Markets

```bash
curl "$BASE_URL/agent/markets"
```

### Start Onboarding

```bash
curl -X POST "$BASE_URL/agent/onboarding/start" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent_demo_custom_001",
    "user_email": "user@example.com",
    "source_client": "custom-agent",
    "channel": "custom-agent"
  }'
```

### Get Session And Wallet

```bash
curl "$BASE_URL/agent/sessions/$AGENT_ID"
curl "$BASE_URL/agent/wallets/$AGENT_ID"
curl "$BASE_URL/agent/wallets/$AGENT_ID/balance"
```

### Create A buy_yes Intent

```bash
curl -X POST "$BASE_URL/agent/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent_demo_custom_001",
    "source_client": "custom-agent",
    "client_request_id": "demo-001",
    "action": "buy_yes",
    "user_wallet": "0x1111111111111111111111111111111111111111",
    "market_id": "example-market-id",
    "market_contract_address": "0x3333333333333333333333333333333333333333",
    "amount": "1"
  }'
```

Save the returned `intent.intent_id`.

### Confirm Intent

```bash
INTENT_ID="agent_intent_example"

curl -X POST "$BASE_URL/agent/intents/$INTENT_ID/confirm"
```

### Execute Intent

```bash
curl -X POST "$BASE_URL/agent/intents/$INTENT_ID/execute"
```

### Get Intent

```bash
curl "$BASE_URL/agent/intents/$INTENT_ID"
```

### List Intent Executions

```bash
curl "$BASE_URL/agent/intents/$INTENT_ID/executions"
```

### Get Portfolio

```bash
curl "$BASE_URL/agent/portfolio/$AGENT_ID"
```

### Get Activity

```bash
curl "$BASE_URL/agent/activity/$AGENT_ID"
```

## Read Models

### Portfolio

`GET /agent/portfolio/{agent_id}` returns a compact read-only portfolio summary for a registered agent wallet.

Current fields include:

- `agent_id`
- `agent_wallet_address`
- `chain`
- `wallet_provider`
- `active_positions_count`
- `resolved_or_closed_positions_count`
- `claimable_refundable_count`
- `total_exposure`
- `positions`
- `settlements`
- `unavailable_fields`

Current limitation: positions are derived from executed Agent API buy intents and executions, not live onchain wallet-indexed balance reads. Existing settlement rows are internal-user keyed, so agent settlement fields remain unavailable until a wallet-indexed model exists.

### Activity

`GET /agent/activity/{agent_id}` returns recent framework-neutral activity from durable intent and execution records. Items include intent/execution type, action, status, market id, market contract address, amount, outcome/side, tx hashes, sanitized error fields, readback, and timestamps when available.

Agents with a registered wallet but no activity return an empty `items` array. Unknown agents return `404 agent_wallet_not_found`.

### Intent Executions

`GET /agent/intents/{intent_id}/executions` returns execution attempts for one intent, including pending/executed/failed status, tx hashes, readback JSON, and sanitized error code/message when available.

## Safety Model

- `agent_id` validation rejects malformed or generic IDs with `agent_id_invalid`.
- `allowed_actions` gates executable intent preview, confirmation when a wallet is available, and execution.
- `execute` requires a confirmed intent.
- `policy_metadata.max_trade_amount` is optional. When present on the registered wallet, `buy_yes` and `buy_no` amounts above the cap are rejected.
- Circle CLI/provider failures are sanitized in public responses and durable execution records.
- Deployer/resolver wallet reuse is rejected for agent wallets.
- Agent wallet reuse of `user_wallet` is rejected until a documented custody-link model exists.
- Current execution, faucet, and wallet policies are pinned to `ARC-TESTNET`.

## Stable Error Catalog

Common Agent API errors:

| Status | Code | Meaning |
| --- | --- | --- |
| 400 | `invalid_json` | Request body is not valid JSON. |
| 400 | `agent_id_invalid` | `agent_id` is missing, malformed, too short, or generic. |
| 400 | `agent_intent_invalid` | Intent validation failed; details may include invalid fields. |
| 400 | `agent_wallet_missing` | Execution was requested without a registered agent wallet. |
| 400 | `agent_onboarding_invalid` | Onboarding request validation failed. |
| 400 | `agent_wallet_invalid` | Wallet registration validation failed. |
| 403 | `agent_action_forbidden` | Action is outside wallet/session `allowed_actions`. |
| 403 | `agent_policy_violation` | Optional wallet policy blocks the request, such as `max_trade_amount`. |
| 404 | `agent_intent_not_found` | No intent exists for the supplied `intent_id`. |
| 404 | `agent_wallet_not_found` | No registered wallet exists for the supplied `agent_id`. |
| 404 | `agent_session_not_found` | No active session exists for the supplied `agent_id`. |
| 409 | `agent_intent_not_confirmed` | Execute was called before confirm. |
| 409 | `agent_wallet_status_invalid` | Wallet is not active. |
| 409 | `agent_wallet_chain_invalid` | Wallet is not on `ARC-TESTNET`. |
| 500 | `agent_activity_get_failed` | Activity read failed. |
| 500 | `agent_portfolio_get_failed` | Portfolio read failed. |
| 501 | `circle_agent_wallet_execution_not_enabled` | Circle Agent Wallet execution is not enabled in this runtime. |
| 501 | `circle_agent_wallet_balance_not_configured` | Balance lookup is not configured. |
| 501 | `circle_agent_wallet_faucet_not_configured` | Faucet helper is not configured. |
| 502 | `agent_execution_failed` | Provider execution failed; response is sanitized. |
| 502 | `circle_agent_wallet_balance_failed` | Circle balance lookup failed; response is sanitized. |
| 502 | `circle_agent_wallet_faucet_failed` | Circle faucet request failed; response is sanitized. |
| 503 | `agent_execution_config_invalid` | Execution environment is not configured. |

## OpenAPI

The Custom GPT action schema lives at:

```text
project-roadmap/signalarc-custom-gpt-openapi.json
```

The schema is useful for any HTTP client, not only Custom GPTs. It documents the current Agent API endpoints and stable response shapes. Custom GPT-specific behavior belongs in integration instructions, not in the core API contract.

## Known Limitations

- Circle Agent Wallet execution depends on runtime CLI/session availability.
- The in-process intent store remains a fallback for tests or unconfigured runtimes.
- API key enforcement, paid access, rate limits, and production SLA are not implemented.
- Portfolio positions are intent/execution-derived, not live onchain position indexing.
- Claimable/refundable eligibility remains limited until wallet-indexed claim/refund state is indexed.
- Mainnet funding and Arc mainnet contract execution are not supported.
- Behavior beyond the documented endpoints is unknown or not documented.
