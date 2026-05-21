# Agent API Smoke Examples

These examples validate the local backend Agent API preview and confirm flow.

Run the backend locally first:

```bash
go run ./cmd/api
```

The examples assume:

```bash
API_BASE="http://localhost:8080"
AGENT_FACTORY_ADDRESS="0x69aE770e8b2F96297101FeC4dc123B3801dA7d80"
USER_WALLET="0x1111111111111111111111111111111111111111"
MARKET_CONTRACT_ADDRESS="0x2222222222222222222222222222222222222222"
RESOLVER_ADDRESS="0x3333333333333333333333333333333333333333"
COLLATERAL_TOKEN="0x4444444444444444444444444444444444444444"
```

The Agent API does not sign, broadcast, call RPC, call Circle APIs, or mutate onchain state in this phase.

## Create Market Intent

Create an intent preview:

```bash
curl -sS -X POST "$API_BASE/agent/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "create_market",
    "user_wallet": "'"$USER_WALLET"'",
    "market_id": "agent-market-smoke-1",
    "question": "Will SignalArc complete the agent API smoke test?",
    "close_timestamp": "1767225600",
    "resolver": "'"$RESOLVER_ADDRESS"'",
    "collateral_token": "'"$COLLATERAL_TOKEN"'"
  }'
```

Expected response shape:

```json
{
  "intent": {
    "intent_id": "agent_intent_...",
    "action": "create_market",
    "status": "preview",
    "requires_confirmation": true,
    "user_wallet": "0x1111111111111111111111111111111111111111",
    "address": "0x1111111111111111111111111111111111111111",
    "market_id": "agent-market-smoke-1",
    "resolver": "0x3333333333333333333333333333333333333333",
    "collateral_token": "0x4444444444444444444444444444444444444444",
    "close_timestamp": "1767225600",
    "question": "Will SignalArc complete the agent API smoke test?",
    "validation_result": {
      "valid": true,
      "errors": []
    },
    "warnings": [
      "preview only; no transaction has been executed"
    ],
    "created_at": "..."
  }
}
```

Get the intent:

```bash
CREATE_INTENT_ID="agent_intent_replace_me"

curl -sS "$API_BASE/agent/intents/$CREATE_INTENT_ID"
```

Expected response shape:

```json
{
  "intent": {
    "intent_id": "agent_intent_replace_me",
    "action": "create_market",
    "status": "preview",
    "requires_confirmation": true,
    "validation_result": {
      "valid": true,
      "errors": []
    }
  }
}
```

Confirm the intent:

```bash
curl -sS -X POST "$API_BASE/agent/intents/$CREATE_INTENT_ID/confirm"
```

Expected `transaction_request` shape:

```json
{
  "execution_plan": {
    "intent_id": "agent_intent_replace_me",
    "action": "create_market",
    "status": "confirmed",
    "execution_mode": "agent_contract",
    "network": "arc_testnet",
    "agent_factory_address": "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80",
    "requires_signature": true,
    "broadcast_performed": false,
    "transaction_hash": null,
    "transaction_request": {
      "to": "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80",
      "contract": "SignalArcAgentMarketFactory",
      "function": "createMarket",
      "args": [
        "agent-market-smoke-1",
        "Will SignalArc complete the agent API smoke test?",
        "1767225600",
        "0x3333333333333333333333333333333333333333",
        "0x4444444444444444444444444444444444444444"
      ],
      "value": "0",
      "chain": "arc_testnet",
      "broadcast_performed": false
    },
    "warnings": [
      "confirmation produced an execution plan only; no transaction has been executed"
    ]
  }
}
```

## Buy YES Intent

Create an intent preview:

```bash
curl -sS -X POST "$API_BASE/agent/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "buy_yes",
    "user_wallet": "'"$USER_WALLET"'",
    "market_id": "agent-market-smoke-1",
    "market_contract_address": "'"$MARKET_CONTRACT_ADDRESS"'",
    "amount": "25.5"
  }'
```

Expected response shape:

```json
{
  "intent": {
    "intent_id": "agent_intent_...",
    "action": "buy_yes",
    "status": "preview",
    "requires_confirmation": true,
    "user_wallet": "0x1111111111111111111111111111111111111111",
    "market_id": "agent-market-smoke-1",
    "market_contract_address": "0x2222222222222222222222222222222222222222",
    "amount": "25.5",
    "validation_result": {
      "valid": true,
      "errors": []
    },
    "warnings": [
      "preview only; no transaction has been executed"
    ],
    "created_at": "..."
  }
}
```

Get the intent:

```bash
BUY_YES_INTENT_ID="agent_intent_replace_me"

curl -sS "$API_BASE/agent/intents/$BUY_YES_INTENT_ID"
```

Confirm the intent:

```bash
curl -sS -X POST "$API_BASE/agent/intents/$BUY_YES_INTENT_ID/confirm"
```

Expected `transaction_request` shape:

```json
{
  "execution_plan": {
    "intent_id": "agent_intent_replace_me",
    "action": "buy_yes",
    "status": "confirmed",
    "execution_mode": "agent_contract",
    "network": "arc_testnet",
    "agent_factory_address": "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80",
    "requires_signature": true,
    "broadcast_performed": false,
    "transaction_hash": null,
    "transaction_request": {
      "to": "0x2222222222222222222222222222222222222222",
      "contract": "SignalArcAgentMarket",
      "function": "buyYes",
      "args": ["25.5"],
      "value": "0",
      "chain": "arc_testnet",
      "broadcast_performed": false
    }
  }
}
```

## Claim Refund Intent

Create an intent preview:

```bash
curl -sS -X POST "$API_BASE/agent/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "claim_refund",
    "user_wallet": "'"$USER_WALLET"'",
    "market_id": "agent-market-smoke-1",
    "market_contract_address": "'"$MARKET_CONTRACT_ADDRESS"'"
  }'
```

Expected response shape:

```json
{
  "intent": {
    "intent_id": "agent_intent_...",
    "action": "claim_refund",
    "status": "preview",
    "requires_confirmation": true,
    "user_wallet": "0x1111111111111111111111111111111111111111",
    "market_id": "agent-market-smoke-1",
    "market_contract_address": "0x2222222222222222222222222222222222222222",
    "validation_result": {
      "valid": true,
      "errors": []
    },
    "warnings": [
      "preview only; no transaction has been executed"
    ],
    "created_at": "..."
  }
}
```

Get the intent:

```bash
CLAIM_REFUND_INTENT_ID="agent_intent_replace_me"

curl -sS "$API_BASE/agent/intents/$CLAIM_REFUND_INTENT_ID"
```

Confirm the intent:

```bash
curl -sS -X POST "$API_BASE/agent/intents/$CLAIM_REFUND_INTENT_ID/confirm"
```

Expected `transaction_request` shape:

```json
{
  "execution_plan": {
    "intent_id": "agent_intent_replace_me",
    "action": "claim_refund",
    "status": "confirmed",
    "execution_mode": "agent_contract",
    "network": "arc_testnet",
    "agent_factory_address": "0x69aE770e8b2F96297101FeC4dc123B3801dA7d80",
    "requires_signature": true,
    "broadcast_performed": false,
    "transaction_hash": null,
    "transaction_request": {
      "to": "0x2222222222222222222222222222222222222222",
      "contract": "SignalArcAgentMarket",
      "function": "claimRefund",
      "args": [],
      "value": "0",
      "chain": "arc_testnet",
      "broadcast_performed": false
    }
  }
}
```

## Inspect Transaction Request

After confirming any intent, inspect only the transaction request:

```bash
curl -sS -X POST "$API_BASE/agent/intents/$BUY_YES_INTENT_ID/confirm" \
  | jq '.execution_plan.transaction_request'
```

Expected shape:

```json
{
  "to": "0x2222222222222222222222222222222222222222",
  "contract": "SignalArcAgentMarket",
  "function": "buyYes",
  "args": ["25.5"],
  "value": "0",
  "chain": "arc_testnet",
  "broadcast_performed": false
}
```

The enclosing execution plan must continue to report:

```json
{
  "broadcast_performed": false,
  "transaction_hash": null
}
```
