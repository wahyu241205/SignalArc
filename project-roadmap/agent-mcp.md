# SignalArc Agent / MCP Roadmap

## Scope

Backend/API agent only. No frontend chat UI.

## Goal

Build a separate agent execution path for SignalArc.

Agent can later:
- read markets
- create market
- trade YES / NO
- claim refund
- claim payout
- resolve/cancel only if authorized

## Core Rule

Backend remains source of truth.

Agent must not bypass backend.

## Contract Separation

Current user/frontend contracts stay untouched:

SignalArcMarketFactory -> SignalArcMarket

Future agent contracts:

SignalArcAgentMarketFactory -> SignalArcAgentMarket

Planned files:
- contracts/src/agent/SignalArcAgentMarket.sol
- contracts/src/agent/SignalArcAgentMarketFactory.sol

Rules:
- Do not modify current live contracts.
- Do not reuse current factory address.
- Deploy new agent factory later.
- Agent contracts stay testnet-only until validated.
- Add Foundry tests before deployment.

## Planned Backend Agent API

- POST /agent/intents
- POST /agent/intents/{id}/confirm
- GET /agent/intents/{id}
- GET /agent/activity
- GET /agent/markets
- GET /agent/positions
- GET /agent/claimable

## Planned Actions

- create_market
- buy_yes
- buy_no
- cancel_market
- close_market
- resolve_market
- claim_refund
- claim_payout

## MCP Boundary

Arc MCP is only for developer tooling and official docs lookup.

Arc MCP is not the runtime trading agent.

## Circle Agent Wallet Boundary

Circle Agent Wallet is planned, not integrated.

Before implementation:
- verify supported blockchains
- verify Arc Testnet support
- verify SDK/CLI path
- verify faucet/funding path

If not documented, mark as unknown / not documented.

## Implementation Order

1. Separate agent contract prototype — DONE
2. Foundry tests — DONE
3. Agent factory deploy on Arc Testnet — DONE
4. Backend Agent API intent model — DONE
5. Backend execution path to agent contract — NEXT
6. Circle Agent Wallet proof of concept
7. Policy-limited agent wallet execution



## Agent Factory Deployment

Status: DONE.

- Network: Arc Testnet
- Agent factory: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
- Deploy transaction: `0x8e4dfa481a2863a08a749fe4add30c4f030f178a4d6ba4658419df6730ebf10b`
- Read validation: `marketCount() == 0`


## Real Arc Testnet Agent Market Smoke

Status: DONE.

This is a real Arc Testnet transaction, not a mock or local simulation.

- Action: `SignalArcAgentMarketFactory.createMarket`
- Factory: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
- Transaction: `0xff4d6eb644792a1c064992704ba767b6712b7cc02c1b44635859e199efdfc69d`
- Created agent market: `0x4e26143A63457cf06A34112b8B9044F3760d3007`
- Read validation:
  - `marketCount() == 1`
  - `allMarkets(0) == 0x4e26143A63457cf06A34112b8B9044F3760d3007`
  - `isMarket(0x4e26143A63457cf06A34112b8B9044F3760d3007) == true`

Note: this validates real factory lifecycle only. Trading validation still requires a valid Arc Testnet collateral token.

## Real Arc Testnet Agent USDC Collateral / Trading Smoke

Status: DONE.

This is a real Arc Testnet transaction sequence, not a mock, fake smoke, dry run, or local simulation.

Official docs finding:
- Arc docs document Arc Testnet USDC ERC-20 interface: `0x3600000000000000000000000000000000000000`.
- Circle docs document Arc Testnet USDC token address: `0x3600000000000000000000000000000000000000`.
- Token readback on Arc Testnet:
  - `name() == "USDC"`
  - `symbol() == "USDC"`
  - `decimals() == 6`

Previous agent market collateral read:
- Market: `0x4e26143A63457cf06A34112b8B9044F3760d3007`
- `collateralToken() == 0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
- Result: not usable as ERC-20 collateral; this is the deployer/resolver EOA, not the documented Arc Testnet USDC contract.

USDC-backed agent market:
- Action: `SignalArcAgentMarketFactory.createMarket`
- Factory: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
- Transaction: `0x9e1dc4b2b65ea6220605f8960dd13ad1be3907b69b72560083e39d3f2c77f579`
- Created agent market: `0xd76c5633c3D8C1761F7edae46506B44cDeEe43a7`
- Collateral token: `0x3600000000000000000000000000000000000000`
- Read validation:
  - `marketCount() == 2`
  - `collateralToken() == 0x3600000000000000000000000000000000000000`
  - `admin() == 0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
  - `resolver() == 0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
  - `isOpen() == true`

Real buyYes / buyNo validation:
- USDC approve transaction: `0xc5963b4ad20b8c66dfb7787c260fcf04deb25fcf53667349958109a10f7584c9`
- `buyYes(1000000)` transaction: `0x1a23107407b17a4dced81d8fd79c2ead8acb8148190a03aaae2f34662621e79f`
- `buyNo(1000000)` transaction: `0x41e0a39271940265894c45893733f794936a6e49cf3f475ae8881d5f9a7bb073`
- Read validation:
  - `yesPositions(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
  - `noPositions(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
  - `totalYes() == 1000000`
  - `totalNo() == 1000000`
  - `totalCollateral() == 2000000`
  - `USDC.balanceOf(0xd76c5633c3D8C1761F7edae46506B44cDeEe43a7) == 2000000`
  - `USDC.allowance(deployer, market) == 0`

## Non-Claims

Not implemented yet:
- frontend chat UI
- Circle Agent Wallet
- autonomous trading
- mainnet settlement
- real funds
- smart contract audit
