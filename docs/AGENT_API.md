# Agent API

SignalArc exposes agent-driven prediction market flows through a structured HTTP API designed for AI agents, institutional dashboards, monitoring systems, and Custom GPT actions. The same surface powers the SignalArc GPT Agent that judges and end users interact with directly.

Live production API base:

```text
https://api.signalarc.fun
```

The production API is live on GCP Cloud Run service `signalarc-backend-api` behind the custom domain `api.signalarc.fun`. The Cloud Run image bundles Node/npm and the Circle CLI (`@circle-fin/cli`) so ARC-TESTNET agent flows can run inside the container. The production database is GCP Cloud SQL, migrated to version 18.

ngrok URLs are local development conveniences only. They are not production endpoints and are not used by the published GPT Agent.

## Custom GPT Integration

The published SignalArc GPT Agent is preconfigured to call `https://api.signalarc.fun`. End users and judges do not need to import OpenAPI manually.

For maintainers who need to re-import the schema:

```text
https://raw.githubusercontent.com/wahyu241205/SignalArc/main/project-roadmap/signalarc-custom-gpt-openapi.json
```

## Live Architecture

| Surface | Target |
| --- | --- |
| `https://signalarc.fun` | Vercel frontend. |
| `https://api.signalarc.fun` | GCP Cloud Run service `signalarc-backend-api`. |
| Production database | GCP Cloud SQL PostgreSQL, schema version 18. |
| Backend container image | Includes `@circle-fin/cli` global on PATH so `circle wallet ...` commands run inside the container. |

## Agent Flow Overview

The end-to-end agent flow looks like this:

```text
onboarding start
    -> OTP verify
        -> active session
            -> wallet
                -> balance
                    -> faucet
                        -> create intent
                            -> confirm intent
                                -> execute intent
```

Each step maps to a backend endpoint described below.

## Capability Status

Available now:

| Capability | Endpoint |
| --- | --- |
| Health | `GET /health` |
| Onboarding start | `POST /agent/onboarding/start` |
| OTP verify | `POST /agent/onboarding/verify` |
| Onboarding lookup | `GET /agent/onboarding/{onboarding_id}` |
| Active session | `GET /agent/sessions/{agent_id}` |
| Registered wallet | `GET /agent/wallets/{agent_id}` |
| Balance (read-only) | `GET /agent/wallets/{agent_id}/balance` |
| ARC-TESTNET faucet | `POST /agent/wallets/{agent_id}/faucet` |
| Market intent preview | `POST /agent/intents` |
| Confirm intent | `POST /agent/intents/{intent_id}/confirm` |
| Execute intent | `POST /agent/intents/{intent_id}/execute` |
| Agent-readable market list | `GET /agent/markets` |

Not available / out of scope:

- Arbitrary transfer.
- Withdraw / deposit.
- Logout / agent session management endpoints.
- Mainnet funding.
- Arc mainnet contract deployment.
- API key enforcement, paid access, autonomous trading, or production SLA.

## Onboarding and Session Endpoints

### POST /agent/onboarding/start

Starts SignalArc Circle Agent Wallet onboarding and dispatches the OTP email through Circle.

```bash
curl -X POST https://api.signalarc.fun/agent/onboarding/start \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent_demo_001",
    "user_email": "you@example.com",
    "source_client": "chatgpt_custom_action",
    "channel": "chatgpt"
  }'
```

Response includes the onboarding record and a `next_step` value such as `circle_otp_required`. Raw Circle request IDs are not exposed.

### POST /agent/onboarding/verify

Verifies the OTP, registers the resolved Circle Agent Wallet address in the SignalArc database, and activates an agent session.

```bash
curl -X POST https://api.signalarc.fun/agent/onboarding/verify \
  -H "Content-Type: application/json" \
  -d '{
    "onboarding_id": "agent_onboarding_xxx",
    "otp": "123456"
  }'
```

Successful verification returns the onboarding, agent wallet, and agent session records, with `next_step` set to `agent_session_active`.

### GET /agent/onboarding/{onboarding_id}

Returns the current onboarding record. Useful for poll-style flows.

### GET /agent/sessions/{agent_id}

Returns the active agent session, including `agent_wallet_address`, `wallet_provider`, `chain`, and `allowed_actions`.

### GET /agent/wallets/{agent_id}

Returns the registered agent wallet record managed by SignalArc.

### GET /agent/wallets/{agent_id}/balance

Returns a read-only Circle Agent Wallet balance snapshot.

```json
{
  "agent_wallet_balance": {
    "agent_id": "agent_demo_001",
    "agent_wallet_address": "0x...",
    "chain": "ARC-TESTNET",
    "balances": []
  }
}
```

Errors map to `404 agent_wallet_not_found`, `501 circle_agent_wallet_balance_not_configured`, or `502 circle_agent_wallet_balance_failed`.

## ARC-TESTNET Faucet

### POST /agent/wallets/{agent_id}/faucet

Requests ARC-TESTNET faucet funding for the registered agent wallet on behalf of an active SignalArc agent. SignalArc forwards the request to the documented Circle CLI testnet faucet using only the registered wallet address.

Endpoint properties:

- ARC-TESTNET only.
- Token is fixed to `usdc`.
- No request body.
- Uses only the registered `agent_wallet_address` from the SignalArc database.
- Does not accept arbitrary recipient addresses; any `address`, `chain`, or `token` in the request body is ignored.
- Does not transfer, swap, execute contracts, create markets, or use mainnet funding.

Example:

```bash
curl -X POST https://api.signalarc.fun/agent/wallets/agent_demo_001/faucet
```

Underlying Circle CLI command shape (run inside the Cloud Run container, never exposed to the caller):

```text
circle wallet fund --address <registered_agent_wallet_address> --chain ARC-TESTNET --token usdc --output json
```

Success response (200):

```json
{
  "agent_wallet_faucet": {
    "agent_id": "agent_demo_001",
    "agent_wallet_address": "0x9999999999999999999999999999999999999999",
    "chain": "ARC-TESTNET",
    "token": "usdc",
    "status": "requested",
    "result": {}
  }
}
```

Response fields:

| Field | Meaning |
| --- | --- |
| `agent_id` | The SignalArc agent identifier. |
| `agent_wallet_address` | The registered Circle Agent Wallet address SignalArc actually targeted. |
| `chain` | Always `ARC-TESTNET`. |
| `token` | Always `usdc`. |
| `status` | Always `requested` when SignalArc successfully forwarded the call. SignalArc does not claim provider-side success. |
| `result` | Parsed Circle CLI JSON output when the CLI returned JSON, or `{ "message": "..." }` when the CLI returned successful text-only output. Sensitive content is sanitized server-side. |

Failure codes:

| Status | Code | Meaning |
| --- | --- | --- |
| 404 | `agent_wallet_not_found` | No registered agent wallet exists for the supplied `agent_id`. |
| 409 | `agent_wallet_status_invalid` | Registered agent wallet exists but is not `active`. |
| 409 | `agent_wallet_chain_invalid` | Registered agent wallet exists but is not on `ARC-TESTNET`. |
| 501 | `circle_agent_wallet_faucet_not_configured` | Faucet helper is not enabled in this runtime. |
| 502 | `circle_agent_wallet_faucet_failed` | Circle CLI faucet command failed; SignalArc returns a generic error and never echoes raw CLI output, credential paths, tokens, request IDs, session material, or emails. |

## Market Intent Lifecycle

Agents preview, confirm, and execute intents through these endpoints. They are documented at the surface level here; see the OpenAPI schema for full request/response shapes.

### POST /agent/intents

Creates a market intent preview. The request must include `agent_id`, `source_client`, `client_request_id`, `action`, and `user_wallet`. Supported actions:

- `create_market`
- `buy_yes`
- `buy_no`
- `close_market`
- `resolve_market`
- `claim_payout`
- `cancel_market`
- `claim_refund`

### POST /agent/intents/{intent_id}/confirm

Confirms a previewed intent and produces an execution plan. SignalArc validates allowed actions, market state, and ARC-TESTNET chain.

### POST /agent/intents/{intent_id}/execute

Executes a confirmed intent through the registered Circle Agent Wallet on ARC-TESTNET. Execution mode is `circle_agent_wallet_cli`.

## Agent-Readable Market Discovery

### GET /agent/markets

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

## Judge / User Testing Guide

This section is the recommended path for grant judges and end users testing SignalArc through the published GPT Agent.

1. Open the SignalArc GPT Agent (already wired to `https://api.signalarc.fun`).
2. Connect your SignalArc account by starting agent onboarding from the GPT Agent.
3. Provide a real, reachable email address when prompted. SignalArc dispatches the Circle OTP email to that address.
4. Enter the OTP from your email back into the GPT Agent. The agent will verify, register the resolved Circle Agent Wallet, and activate an agent session.
5. Check the registered wallet via the GPT Agent. The response includes the `agent_wallet_address` SignalArc will use for ARC-TESTNET actions.
6. Check the wallet balance via the GPT Agent. This is read-only.
7. Request the ARC-TESTNET faucet via the GPT Agent. SignalArc only funds the registered agent wallet address; you cannot supply a different recipient.
8. Create a draft market intent through the GPT Agent for one of the supported actions (for example, `create_market` or `buy_yes`).
9. Confirm the intent. SignalArc returns an execution plan but does not broadcast yet.
10. Execute the intent only after explicit approval. Execution runs through the Circle Agent Wallet on ARC-TESTNET only.

You should never be asked to provide a private key, seed phrase, Circle API key, or session token through the agent. SignalArc does not store any of those values in the repository, container image, or database.

## Example Use Cases

- Custom GPT action surface for prediction market workflows on Arc Testnet.
- Agent-driven onboarding, funding, and market lifecycle without requiring a private key in the agent.
- Market discovery for dashboards and automated reports.
- Probability-signal reading from structured market metadata.

## Current Limitations

- Autonomous unattended trading is not enabled by default.
- Mainnet funding, transfer, swap, and contract execution are out of scope.
- Logout/session management endpoints are not exposed.
- API key enforcement, paid access, and rate limits are not implemented.
- Production SLA is not claimed.
- Behavior beyond the documented endpoints is unknown / not documented.
