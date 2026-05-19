# Arc / Circle Research

## Document Purpose

This document is the Phase 5.1 research note for SignalArc before any wallet, Arc runtime, Circle, contract, settlement, webhook, or agent wallet implementation.

## Official Sources

- Arc docs: https://docs.arc.network/
- Arc docs alternate domain: https://docs.arc.io/
- Arc prediction market blueprint: https://www.arc.network/blog/build-institutional-grade-prediction-markets-on-arc-arc-blueprints
- Circle developer docs: https://developers.circle.com/
- Circle Agents docs: https://agents.circle.com/
- Circle Grants: https://www.circle.com/grant

## Verified Arc Facts From Current Project State

- Arc Testnet exists.
- Network: Arc Testnet.
- Chain ID: 5042002.
- Primary RPC: https://rpc.testnet.arc.network.
- Explorer: https://testnet.arcscan.app.
- Faucet: https://faucet.circle.com.
- Currency / native gas token: USDC.
- Arc is currently documented as testnet phase.
- Arc mainnet status: unknown / not documented from reviewed pages.
- Arc is EVM-compatible.
- Solidity smart contract deployment to Arc Testnet with Foundry is documented.
- Arc docs mention standard Ethereum tooling including Solidity, Foundry, Hardhat, and Viem.
- Arc uses USDC as native gas.
- Arc deterministic finality is documented.
- Arc USDC has dual interface behavior:
  - native balance uses 18 decimals.
  - ERC-20 interface uses 6 decimals.
  - both interfaces share one underlying USDC balance.

## Circle Documentation Status

| Area | Documentation Status | SignalArc Notes |
| --- | --- | --- |
| Circle Wallets | unknown / not documented in this repo snapshot | Do not assume Circle Wallets support Arc Testnet until exact official documentation snapshots are collected. |
| CCTP | unknown / not documented in this repo snapshot | Do not assume CCTP relevance or Arc Testnet support until exact official documentation snapshots are collected. |
| Gateway | unknown / not documented in this repo snapshot | Do not assume Gateway relevance or Arc Testnet support until exact official documentation snapshots are collected. |
| Paymaster / Gas Station | unknown / not documented in this repo snapshot | Do not assume Paymaster / Gas Station relevance or Arc Testnet support until exact official documentation snapshots are collected. |
| Circle Agent Stack | unknown / not documented in this repo snapshot | Do not assume Circle Agent Stack or agent wallet relevance until exact official documentation snapshots are collected. |
| Arc Testnet support | unknown / not documented in this repo snapshot | Do not claim that Circle products support Arc Testnet until exact official documentation evidence is available. |
| Circle webhooks / event handling | unknown / not documented in this repo snapshot | Do not design webhook behavior until exact official documentation snapshots are collected. |
| USDC flow for SignalArc MVP | unknown / not documented in this repo snapshot | Do not define deposit, trade, settlement, claim, or payout flow until exact official documentation snapshots are collected. |

## SignalArc Phase 5.1 Boundaries

- Documentation research only.
- No wallet implementation.
- No Circle SDK/API integration.
- No Arc RPC runtime config.
- No smart contracts.
- No contracts folder.
- No settlement execution.
- No claim flow.
- No backend code changes.
- No frontend code changes.
- No database migration changes.
- No dependency changes.
- No environment variable changes.

## Open Questions For Phase 5.2

Phase 5.2 must collect exact official Circle documentation snapshots before deciding:

- MVP wallet strategy.
- Circle Wallets relevance.
- Circle Agent Wallets relevance.
- CCTP relevance.
- Gateway relevance.
- Paymaster / Gas Station relevance.
- Circle webhook or event handling plan.
- USDC deposit, trade, settlement, and payout flow.
- Arc transaction flow for Phase 6/7.
- Contract requirements for Phase 6.

## Non-Goals

- This document does not define final architecture.
- This document does not approve implementation.
- This document does not create production custody assumptions.
- This document does not claim Circle or Arc support beyond official documentation evidence.
