# Wallet Strategy

## Document Purpose

This document is the Phase 5.2 wallet strategy note for SignalArc before any wallet implementation or Circle integration.

## Inputs

- AGENTS.md
- PROJECT_STATE.md
- docs/arc-circle-research.md
- official Arc documentation
- official Circle documentation, only where exact facts are available

## Current Verified Constraints

- SignalArc frontend is UI-only.
- SignalArc backend owns business logic.
- Arc/Circle behavior must come from official documentation.
- Unsupported behavior must be marked unknown / not documented.
- Phase 5 is documentation and architecture planning only.
- Phase 6 starts contracts / settlement prototype work.
- Arc Testnet facts are already captured in docs/arc-circle-research.md.
- Circle product support for Arc Testnet remains unknown / not documented unless exact official documentation is found.

## Wallet Options Considered

### Option A — Browser Wallet Path

The browser wallet path uses the already selected frontend stack:

- Wagmi
- Viem
- RainbowKit

This path is suitable for user-controlled wallet UX. This path does not require SignalArc to custody user funds. Arc chain configuration must not be implemented in Phase 5.2. Actual runtime configuration is deferred until a later implementation phase.

### Option B — Circle Wallets Path

Circle Wallets are considered only at a planning level in this document.

Arc Testnet support for Circle Wallets: unknown / not documented.

This document does not define custody, key management, or transaction execution behavior for Circle Wallets beyond official documentation.

### Option C — Circle Agent Wallets / Agent Stack Path

Circle Agent Stack / Agent Wallets are considered only at a planning level in this document.

Arc Testnet support for Circle Agent Stack / Agent Wallets: unknown / not documented.

Agent execution is not allowed in the MVP unless later documentation and architecture explicitly approve it.

### Option D — Hybrid Path

A possible hybrid strategy is:

- browser wallet path for human users
- Circle Wallets only if officially documented and required for embedded or managed wallet UX
- Circle Agent Wallets only if officially documented and required for future agent execution

This is a planning option, not implementation approval.

## MVP Wallet Decision

Use browser wallet path as the default MVP direction because the repo stack already includes Wagmi, Viem, and RainbowKit, while Circle Wallets / Agent Wallets Arc Testnet support remains unknown / not documented.

- This is a planning decision, not implementation approval.
- No wallet code is added in Phase 5.2.
- Circle Wallets may be reconsidered after exact official Circle documentation snapshots are collected.
- Agent wallets are future-scope unless explicitly verified and approved.

## SignalArc Wallet Boundary

- Frontend may provide wallet connection UI in a later implementation phase.
- Backend owns business validation and intent handling.
- Backend must not receive private keys.
- Frontend must not contain Circle API secrets.
- No production custody assumptions are approved by this document.
- Trade execution, settlement, and claim flow remain out of scope for Phase 5.2.

## Open Questions For Phase 5.3

Phase 5.3 must create docs/usdc-integration-plan.md and answer:

- How USDC collateral is represented conceptually.
- How trade intent maps to future transaction flow.
- Which parts remain offchain in the backend.
- Which parts move to Phase 6 contracts.
- Whether CCTP, Gateway, Paymaster / Gas Station, or Circle webhooks are relevant to MVP.
- Which Circle behaviors remain unknown / not documented.

## Non-Goals

- This document does not implement wallet connect.
- This document does not configure Arc RPC.
- This document does not integrate Circle SDKs.
- This document does not approve custodial wallet behavior.
- This document does not approve agent trading.
- This document does not define settlement execution.
- This document does not create contracts.
