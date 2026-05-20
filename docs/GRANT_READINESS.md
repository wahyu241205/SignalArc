# Grant Readiness

SignalArc is positioned as Arc-native prediction market infrastructure for USDC-settled event markets, market intelligence, resolver workflow, and agent-readable APIs.

This document separates implemented work from testnet prototype behavior, planned work, and unknown / not documented behavior. It does not claim grant readiness is complete.

## Implemented

| Area | Current state |
| --- | --- |
| Frontend MVP | Local Next.js application exists and runs manually. |
| Backend API | Go/Chi API implements health, readiness, schema validation, markets, trade intents, positions, resolutions, settlements, agent markets, and Arc contract metadata. |
| Local backend stack | Docker Compose runs PostgreSQL and backend. |
| Database schema | Local schema is complete through migration version 13 according to project state. |
| Agent-readable API | `GET /agent/markets` is implemented. |
| Wallet frontend | Current working tree includes RainbowKit/Wagmi/Viem external wallet UI. |
| Arc Testnet contract | `SignalArcMarket` prototype is deployed on Arc Testnet. |

## Testnet Prototype

| Area | Current state |
| --- | --- |
| Contract | Arc Testnet `SignalArcMarket` prototype at `0xf4ccc11A9e24fb996679F946C23C04AFd2797F26`. |
| Browser transaction flow | Current frontend working tree includes Arc Testnet USDC approval and `openPosition`. |
| Explorer links | Current frontend generates Arcscan transaction links. |
| Settlement | Prototype contract supports claim/refund paths in Solidity, but production settlement is not approved. |

## Circle Relevance

| Area | Status |
| --- | --- |
| USDC collateral | Prototype contract uses a USDC-like ERC20 collateral model and Arc Testnet USDC interface address. |
| Circle Developer Platform | Planned/possible only. Not implemented in the current repository. |
| Circle Agents | Planned/possible only. Not implemented in the current repository. |
| Circle API keys | Must not be committed or exposed. |

## Planned

- Production frontend deployment at `https://signalarc.fun`.
- Production API deployment at `https://api.signalarc.fun`.
- Documentation hosting at `https://docs.signalarc.fun`.
- Production database.
- Production CORS configuration.
- API keys and scoped agent access.
- Rate limits.
- Expanded agent-readable intelligence endpoints.
- Circle integration only after official documentation review and explicit implementation approval.

## Grant Gaps

- Live app URL.
- Verified deployed frontend.
- Verified deployed backend API.
- Configured DNS.
- Demo video.
- Product README and technical README completeness.
- Production deployment plan execution.
- Risk/compliance disclaimer.
- Circle integration depth.
- Production monitoring/logging baseline.
- Clear live demo flow.
- Any required security review or audit status.

## Unknown / Not Documented

- Arc mainnet deployment path for SignalArc: unknown / not documented in this repository.
- Production custody model: unknown / not documented.
- Production settlement approval: unknown / not documented.
- Circle SDK/API behavior inside SignalArc: not implemented.
- Circle Agents behavior inside SignalArc: not implemented.
- Compliance approval: unknown / not documented.

## Readiness Assessment

SignalArc has a credible local MVP, an Arc Testnet contract prototype, and a clear API-first product direction. It is not yet grant-submission complete because live deployment, DNS, demo materials, production documentation, and deeper Circle integration evidence are still gaps.
