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

## Live AI Agent Transaction MVP

Status: IN PROGRESS - Circle Agent Wallet createMarket, buyYes, and buyNo are validated by user-provided Circle CLI evidence; payout/refund lifecycle, backend Circle execution automation, and external client trigger are pending.

Primary objective:
- A user owns or controls an agent.
- The agent has its own wallet address.
- The agent wallet is not the deployer wallet.
- The agent wallet is not the user personal wallet unless an explicitly documented user-controlled custody link is implemented.
- ChatGPT, Claude, Telegram, or other external clients are command interfaces only.
- SignalArc Backend Agent API remains the source of truth for intent validation, policy, permissions, balances, market status, execution state, and audit state.
- Every milestone transaction must include a real Arc Testnet tx hash, receipt success, onchain readback, and proof that the signer/wallet is the agent wallet.

### Agent Wallet Boundary

Do not use `contracts/.env` `PRIVATE_KEY` as an agent wallet.

Do not map `PRIVATE_KEY` to `AGENT_EXECUTOR_PRIVATE_KEY`.

Do not use the deployer, resolver, or user wallet as the agent wallet.

Do not treat `AGENT_EXECUTOR_PRIVATE_KEY` as the final Agent Wallet design. The existing backend EOA executor is legacy execution plumbing only and is not valid proof of the Live AI Agent Transaction MVP.

Deployer/resolver wallet explicitly forbidden as an agent wallet:
- `0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`

Preferred wallet provider:
- `circle_agent_wallet`, if the live Circle Agent Wallet CLI path is authenticated and proven on `ARC-TESTNET`.

Documented fallback provider name only:
- `temporary_testnet_agent_eoa`

The fallback must never use the deployer private key, resolver wallet, user wallet, or any key copied from `contracts/.env`.

### Official Documentation Findings

Circle Agent Wallets:
- Designed for agents to hold, spend, trade, and earn USDC/tokens within spending controls.
- Built on Circle user-controlled wallets with 2-of-2 MPC.
- Agent key shares are not exposed to the agent.
- Users retain custody.
- Operated through Circle CLI.

Circle Agent Wallet ARC-TESTNET support:
- Circle documents `ARC-TESTNET` as the Arc Testnet chain identifier for Agent Wallet CLI commands.
- Circle documents `circle wallet login you@example.com --testnet`.
- Circle documents `circle wallet list --type agent --chain ARC-TESTNET`.
- Circle documents `circle wallet fund --address 0xYourWalletAddress --chain ARC-TESTNET` for testnet funding, with testnet wallets auto-funded from the Circle faucet.
- Circle documents `circle wallet balance --address 0xYourWalletAddress --chain ARC-TESTNET`.
- Circle documents `circle wallet execute <abiFunctionSignature> ... --address <agent_wallet> --contract <contract> --chain ARC-TESTNET`.
- Circle documents `circle contract query <abiFunctionSignature> ... --contract <contract> --chain ARC-TESTNET`.

These are official-doc and CLI help command shapes only. They are not authenticated local validation results unless this document records the exact command output, wallet address, tx hash or Circle transaction id, receipt status, and onchain readback.

Arc:
- Arc MCP is documentation lookup tooling only.
- Arc MCP is not the runtime trading agent.
- Arc Agentic Economy is relevant to the long-term SignalArc design.

Unknown / not documented:
- Whether SignalArc should store Circle policy IDs, detailed spending caps, or Circle-side wallet policy state in the backend before a real Circle Agent Wallet is authenticated and inspected.
- Whether Circle Agent Wallet contract execution for this exact SignalArc contract call shape succeeds on this local account before CLI login and OTP authentication.

### Checkpoint 2 - Circle Agent Wallet POC Discovery

Local discovery result:
- `circle` CLI was installed globally with `npm install -g @circle-fin/cli` because it is an external/global tool and not a SignalArc project dependency.
- `circle --version` returned `0.0.3`.
- Safe help discovery was run for `circle wallet`, `wallet login`, `wallet list`, `wallet create`, `wallet fund`, `wallet balance`, `wallet execute`, and `contract query`.
- Additional safe help discovery was run for `circle --help` and `circle terms --help`.
- Safe auth/chain discovery commands were blocked before login by Circle CLI Terms acceptance:
  - `circle wallet status --type agent --output json`
  - `circle blockchain list -q`
  - `circle blockchain list --output json`
- Circle CLI testnet authentication could not be tested because accepting Terms of Use and Privacy Policy is a user-controlled manual step.
- No Circle Agent Wallet address was created or listed.
- No Circle Agent Wallet was funded.
- No Circle Agent Wallet contract execution was performed.

Manual setup commands from official docs:

```bash
npm install -g @circle-fin/cli
circle wallet login you@example.com --testnet
circle wallet list --type agent --chain ARC-TESTNET
circle wallet fund --address 0xYourWalletAddress --chain ARC-TESTNET
circle wallet balance --address 0xYourWalletAddress --chain ARC-TESTNET
circle wallet execute "createMarket(string,string,uint256,address,address)" <marketId> <question> <closeTimestamp> <resolver> <collateralToken> --address 0xAgentWallet --contract 0x69aE770e8b2F96297101FeC4dc123B3801dA7d80 --chain ARC-TESTNET
circle contract query "marketCount()" --contract 0x69aE770e8b2F96297101FeC4dc123B3801dA7d80 --chain ARC-TESTNET
```

If OTP/manual login is required, stop and let the user complete it. Do not bypass user control.

Current manual blocker:

```bash
# User-controlled step; do not run automatically.
circle terms accept
circle wallet login you@example.com --testnet
```

Circle CLI help says login emails an OTP. The user must accept Circle Terms/Privacy Policy and complete OTP manually before SignalArc can list, fund, or execute from an agent wallet.

Note: official Circle docs show `circle wallet login you@example.com --testnet`; the installed `circle wallet login --help` output describes separate mainnet/testnet sessions but did not list a `--testnet` option. Treat the exact testnet login flag as blocked until the user completes the CLI Terms step and the command can be checked interactively.

### Backend Agent Wallet Model

Persistent backend model added for:
- `agent_id`
- `user_wallet`
- `user_email`
- `agent_wallet_address`
- `wallet_provider`
- `chain`
- `allowed_actions`
- `status`
- `policy_metadata`
- `source_client`
- `created_at`
- `updated_at`

Current implementation uses the DB-backed `agent_wallets` table for production routing. Test code uses an in-memory registry with the same interface for handler coverage only.

No private keys, Circle OTPs, secret tokens, or undocumented Circle session material are stored.

### Agent Wallet Onboarding API

`POST /agent/wallets`
- Registers an already-created user-owned Circle Agent Wallet with SignalArc.
- Does not create the Circle Agent Wallet.
- Does not authenticate Circle.
- Does not store secrets.
- Rejects deployer/resolver wallet reuse.
- Rejects agent wallet reuse of `user_wallet` until a documented custody-link model exists.
- Requires `wallet_provider == circle_agent_wallet`.
- Requires `chain == ARC-TESTNET`.
- Requires non-empty `allowed_actions`.

`GET /agent/wallets/{agent_id}`
- Returns registered wallet metadata only.
- Does not return secrets.

`POST /agent/wallets/{agent_id}/disable`
- Marks a registered agent wallet disabled.
- Disabled wallets cannot execute.

Onboarding flow:
1. User controls email and Circle Agent Wallet setup.
2. User completes Circle CLI login/OTP outside SignalArc backend.
3. User obtains Circle Agent Wallet address on `ARC-TESTNET`.
4. Client registers wallet metadata with SignalArc through `POST /agent/wallets`.
5. Any external client can submit intents by `agent_id`.
6. Backend validates wallet status, chain, provider, action allowlist, and wallet boundary before execution readiness.

### Backend Execution Guardrails

`POST /agent/intents` now accepts:
- `agent_id`
- `agent_wallet_address`
- `source_client`
- `client_request_id`
- existing action, market, and amount fields

`POST /agent/intents/{id}/execute` rejects execution when:
- `agent_id` has no registered agent wallet.
- agent wallet is inactive.
- agent wallet chain is not `ARC-TESTNET`.
- action is not in `allowed_actions`.
- agent wallet equals the deployer/resolver wallet.
- agent wallet equals the registered user wallet.
- agent wallet matches `SIGNALARC_FORBIDDEN_AGENT_WALLETS`.
- request `agent_wallet_address` does not match the registered agent wallet.
- provider is `circle_agent_wallet` but Circle CLI authentication and live wallet proof are not complete.
- provider is `temporary_testnet_agent_eoa`; this is only a documented fallback name and is not enabled.

### Channel-Agnostic Client Contract

External clients can include:
- WhatsApp agent
- Telegram bot
- ChatGPT custom action
- Claude tool
- Web client

Client responsibilities:
- Call SignalArc Backend Agent API.
- Reference a registered `agent_id`.
- Provide `source_client` and optional `client_request_id` for auditability.
- Never hold private keys.
- Never receive OTPs or Circle session tokens from SignalArc.
- Never execute contracts directly.

Backend responsibilities:
- Maintain the user/agent wallet registry.
- Validate intent shape and registered wallet state.
- Enforce provider, chain, status, and action allowlist checks.
- Return an execution plan or execution readiness failure.
- Keep Circle Agent Wallet execution fail-closed until a safe server-side Circle auth/session strategy is designed and proven.

Chosen MVP client remains: ChatGPT custom action / API-call style client.

Live external agent validation is not complete until a real client submits an intent and the registered agent wallet executes the onchain transaction through the backend-approved path.

### Backend Circle Provider Decision

Do not automate Circle CLI execution in backend by relying on a user's local interactive CLI session in production.

Current backend behavior:
- `circle_agent_wallet` execution uses a guarded Circle CLI adapter only when explicitly enabled.
- Registered wallet metadata and intent validation are DB-backed.
- Onchain Circle Agent Wallet transactions are proven manually through Circle CLI evidence only.
- By default, `CIRCLE_AGENT_WALLET_EXECUTION_ENABLED=false`, so execution fails closed.
- The backend never calls Circle login.
- The backend never accepts OTP.
- The backend never stores or logs Circle session files, OTPs, Circle tokens, private keys, or deployer keys.

Backend Circle CLI Provider:
- Provider: `circle_agent_wallet`
- Execution mode: `circle_agent_wallet_cli`
- Chain: `ARC-TESTNET`
- CLI path config: `CIRCLE_CLI_PATH`, default `circle`
- Enable config: `CIRCLE_AGENT_WALLET_EXECUTION_ENABLED`, default `false`
- Chain config: `CIRCLE_AGENT_WALLET_CHAIN`, default `ARC-TESTNET`
- Timeout config: `CIRCLE_AGENT_WALLET_TIMEOUT_SECONDS`, default `120`

Supported backend Circle CLI actions in this MVP:
- `create_market`
- `buy_yes`
- `buy_no`

Not implemented in this provider checkpoint:
- `close_market`
- `resolve_market`
- `claim_payout`
- `cancel_market`
- `claim_refund`

Circle CLI command mapping:
- `create_market`: `circle wallet execute "createMarket(string,string,uint256,address,address)" ... --address <agent_wallet> --contract <factory> --chain ARC-TESTNET --output json`
- `buy_yes`: `approve(address,uint256)` against Arc Testnet USDC, then `buyYes(uint256)` against the market contract.
- `buy_no`: `approve(address,uint256)` against Arc Testnet USDC, then `buyNo(uint256)` against the market contract.

Readback mapping:
- After `create_market`: `marketCount()`, `allMarkets(lastIndex)`, `isMarket(created_market)`.
- After `buy_yes`: `yesPositions(agent_wallet)`, `totalYes()`, `totalCollateral()`, `USDC.balanceOf(market)`.
- After `buy_no`: `noPositions(agent_wallet)`, `totalNo()`, `totalCollateral()`, `USDC.balanceOf(market)`.

Runtime requirement:
- The backend runtime must have Circle CLI installed and an operator-authenticated Circle Agent Wallet session available.
- If official docs do not document a non-CLI server API for this exact execution path, production auth/session strategy remains a separate checkpoint.
- This remains a testnet/dev execution path until production auth, policy, audit, and lifecycle behavior are designed and proven.

### Current Non-Claims

Not complete:
- The payout lifecycle has not been completed from this agent wallet in this Codex shell.
- The cancel/refund lifecycle has not been completed from this agent wallet in this Codex shell.
- No live external ChatGPT/Claude/Telegram client triggered execution.
- SignalArc is not mainnet ready.
- SignalArc does not claim production policy-limited execution on ARC-TESTNET. Circle CLI help says `wallet limit` is mainnet only.
- SignalArc is not claiming autonomous trading.

### User-Proven Circle Agent Wallet Evidence

Status: PARTIALLY VALIDATED BY USER-PROVIDED CIRCLE CLI EVIDENCE.

This evidence supersedes the earlier deployer-signed smoke tests for the create-market and buy-position portions of the Live AI Agent Transaction MVP. It does not complete payout, refund, external-agent-client, or production policy-limit validation.

Agent wallet:
- Chain: `ARC-TESTNET`
- Wallet provider: `circle_agent_wallet`
- Circle CLI version observed: `0.0.3`
- Agent wallet address: `0x96d5051a005547eba149f71604ccf58ae1a7c950`
- Deployer/resolver wallet: `0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
- Boundary proof: agent wallet is different from deployer/resolver wallet.
- Funding proof: Circle Agent Wallet balance showed `20 USDC` before the recorded execution sequence.

Factory:
- `SignalArcAgentMarketFactory`: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
- Initial readback: `marketCount() == 7`

Circle Agent Wallet `createMarket`:
- Transaction: `0x7142dbd7eebe7cbfb19199d9984efa5cef814d0e6038c17b98f2e98cc731cacf`
- Circle source address: `0x96d5051a005547eba149f71604ccf58ae1a7c950`
- Circle state: `COMPLETE`
- Post-create factory readback:
  - `marketCount() == 8`
  - `allMarkets(7) == 0xAbCf081E456C1a11106deF590666A07B76D456f8`

Created market readback:
- Market: `0xAbCf081E456C1a11106deF590666A07B76D456f8`
- `collateralToken() == 0x3600000000000000000000000000000000000000`
- `admin() == 0x96d5051a005547eba149f71604ccf58ae1a7c950`
- `resolver() == 0x96d5051a005547eba149f71604ccf58ae1a7c950`
- `isOpen() == true`

Circle Agent Wallet YES position:
- USDC approve transaction: `0xeb7304b0a1be9f5dc575f62fb705dfaf384bc720da13f7e4ffe9563442c036ca`
- `buyYes(1000000)` transaction: `0xe311d999e15e6f34fa6f623a8f27bc724c665d7c3296632460339326b6094b16`
- Readback:
  - `yesPositions(0x96d5051a005547eba149f71604ccf58ae1a7c950) == 1000000`
  - `totalYes() == 1000000`
  - `totalCollateral() == 1000000`

Circle Agent Wallet NO position:
- USDC approve transaction: `0x6ea6a10293a4df5d7ed50e077821115571787d8e9d6b9507a984ebf33fc52a9b`
- `buyNo(1000000)` transaction: `0xaefe8bcdcec794c811d615517f0dfa800b9e263631200a74c85d000374aa8f24`
- Readback:
  - `noPositions(0x96d5051a005547eba149f71604ccf58ae1a7c950) == 1000000`
  - `totalNo() == 1000000`
  - `totalCollateral() == 2000000`
  - `USDC.balanceOf(0xAbCf081E456C1a11106deF590666A07B76D456f8) == 2000000`

Continuation attempt from this Codex shell:
- `CIRCLE_ACCEPT_TERMS=1` was set only for Circle CLI processes so the non-interactive shell could use the existing session.
- `circle wallet status --type agent --output json` returned `AUTH_REQUIRED`.
- `circle wallet list --type agent --chain ARC-TESTNET --output json` returned not logged in or session expired.
- `circle wallet balance --address 0x96d5051a005547eba149f71604ccf58ae1a7c950 --chain ARC-TESTNET` returned no active agent session.
- `circle contract query` readbacks attempted from this shell returned `fetch failed`.

Current blocker:
- The authenticated Circle Agent Wallet session is not available to this non-interactive Codex shell.
- Required manual step: authenticate the Circle CLI in the same shell/context that Codex can access, then rerun status/list/balance before any further transaction.
- Do not ask for or print OTP, tokens, private keys, or Circle credentials.

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

Real cancel / refund lifecycle validation:
- Pre-cancel read validation:
  - `status() == 0`
  - `claimableRefund(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 0`
  - `USDC.balanceOf(0xd76c5633c3D8C1761F7edae46506B44cDeEe43a7) == 2000000`
  - `USDC.balanceOf(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 32622253`
- `cancelMarket()` transaction: `0x93df5450fcfc50054b7cbc4f260fafc876a29d339d3923807cdb3d4f6323274d`
- Cancel receipt status: `1 (success)`
- Post-cancel read validation:
  - `status() == 3`
  - `claimableRefund(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 2000000`
  - `USDC.balanceOf(0xd76c5633c3D8C1761F7edae46506B44cDeEe43a7) == 2000000`
- `claimRefund()` transaction: `0x54bb3939995f613212531e45345d42c31d57e0ce61eccb26eb5d526921ce4453`
- Claim refund receipt status: `1 (success)`
- Post-refund read validation:
  - `hasClaimed(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == true`
  - `USDC.balanceOf(0xd76c5633c3D8C1761F7edae46506B44cDeEe43a7) == 0`
  - `USDC.balanceOf(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 34619899`
  - `claimableRefund(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 2000000`
- Note: `claimableRefund` reports the cancelled position amount and does not check `hasClaimed`; refund completion is proven by `hasClaimed == true`, market USDC balance `0`, and the successful USDC transfer in the claim transaction.

## Real Arc Testnet Agent YES Payout Lifecycle Smoke

Status: DONE.

This is a real Arc Testnet transaction sequence, not a mock, fake smoke, dry run, or local simulation.

Fresh USDC-backed payout market:
- Action: `SignalArcAgentMarketFactory.createMarket`
- Factory: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
- Transaction: `0x4ac895622bc802ac6639095707675269d5b3de8a08e60991fe8a002c794aa75d`
- Created agent market: `0xcCE012A74Cdf7d17138cd6A514394c79b092B6E7`
- Close timestamp: `1779357303`
- Collateral token: `0x3600000000000000000000000000000000000000`
- Read validation:
  - `marketCount() == 4`
  - `allMarkets(3) == 0xcCE012A74Cdf7d17138cd6A514394c79b092B6E7`
  - `collateralToken() == 0x3600000000000000000000000000000000000000`
  - `status() == 0`
  - `isOpen() == true`

Real YES position:
- USDC approve transaction: `0x73d9ac7f00cda86c02acc9d8082dafd8432874818696fc7588bb705305ce31ad`
- `buyYes(1000000)` transaction: `0x584d9d762b4e5279e6db475ffe2a0ab43b0c06220c86ed799ba7e0dbc8d15311`
- Read validation:
  - `yesPositions(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
  - `totalYes() == 1000000`
  - `totalCollateral() == 1000000`
  - `USDC.balanceOf(0xcCE012A74Cdf7d17138cd6A514394c79b092B6E7) == 1000000`

Real close / resolve / payout lifecycle validation:
- Wait validation:
  - latest block timestamp before close: `1779357317`
  - close timestamp: `1779357303`
- `closeMarket()` transaction: `0x0fb4194bbb9b75097180bcae078ecaf4ce20e3c9425a75de7358ac2d2f8f34e4`
- Close receipt status: `1 (success)`
- `resolve(YES)` transaction: `0x961153bc68239ffcdfccfe7472e7498305d98d1765fe90fb3b4a7013f8ba7afb`
- Resolve receipt status: `1 (success)`
- Pre-claim read validation:
  - `status() == 2`
  - `winningOutcome() == 1`
  - `claimablePayout(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
  - `USDC.balanceOf(0xcCE012A74Cdf7d17138cd6A514394c79b092B6E7) == 1000000`
- `claimPayout()` transaction: `0xca6e837d455a43e99136a6f2c50bdeb2ba76aef41c5235d5d9994baadc83b631`
- Claim payout receipt status: `1 (success)`
- Post-payout read validation:
  - `status() == 2`
  - `winningOutcome() == 1`
  - `hasClaimed(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == true`
  - `USDC.balanceOf(0xcCE012A74Cdf7d17138cd6A514394c79b092B6E7) == 0`
  - `USDC.balanceOf(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 34570580`
  - `claimablePayout(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
  - `USDC.allowance(deployer, market) == 0`
- Note: `claimablePayout` reports the winning position amount and does not check `hasClaimed`; payout completion is proven by `hasClaimed == true`, market USDC balance `0`, and the successful USDC transfer in the claim transaction.

## Backend Agent create_market Execution Endpoint

Status: DONE.

This validates the backend `POST /agent/intents/{id}/execute` path for `action=create_market` with a real Arc Testnet transaction, not a mock, fake smoke, dry run, or local simulation.

Implementation boundary:
- Implemented action: `create_market`
- Other actions: return `501 not_implemented`
- Runtime env names:
  - `ARC_TESTNET_RPC_URL`
  - `AGENT_EXECUTOR_PRIVATE_KEY`
  - `AGENT_FACTORY_ADDRESS`
- Private key is read from env only and is not logged or returned.

Backend local validation:
- Backend route: `POST /agent/intents/{id}/execute`
- Local validation port: `4107`
- Intent id: `agent_intent_4e6c6445c3d542cd6b6cf2dff24611d9`
- Market id: `backend-create-market-1779358755`
- Question: `Will SignalArc backend execute create_market on Arc Testnet?`
- Close timestamp: `1779445152`
- Resolver: `0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
- Collateral token: `0x3600000000000000000000000000000000000000`

Execution result:
- Transaction: `0x0f89ff6e50e31769c0708d3f725c86c4612ba553f7e2bdab7721f8e6dd2674c4`
- Created agent market: `0x0d4aF15Bee1Caf6FB61F55668c2Cd8CB7a051e81`
- Backend response readback:
  - `status == executed`
  - `broadcast_performed == true`
  - `transaction_hash == 0x0f89ff6e50e31769c0708d3f725c86c4612ba553f7e2bdab7721f8e6dd2674c4`
  - `network == arc_testnet`
  - `readback.market_count == 5`
  - `readback.created_market == 0x0d4aF15Bee1Caf6FB61F55668c2Cd8CB7a051e81`
  - `readback.is_market == true`

Independent onchain readback:
- Receipt status: `1 (success)`
- Factory `marketCount() == 5`
- Factory `isMarket(0x0d4aF15Bee1Caf6FB61F55668c2Cd8CB7a051e81) == true`
- Created market `collateralToken() == 0x3600000000000000000000000000000000000000`

## Backend Agent buy_yes Execution Endpoint

Status: DONE.

This validates the backend `POST /agent/intents/{id}/execute` path for `action=buy_yes` with real Arc Testnet transactions, not a mock, fake smoke, dry run, or local simulation.

Implementation boundary:
- Implemented additional action: `buy_yes`
- Still not implemented in this phase: `buy_no`, `cancel_market`, `close_market`, `resolve_market`, `claim_refund`, `claim_payout`
- Runtime env names:
  - `ARC_TESTNET_RPC_URL`
  - `AGENT_EXECUTOR_PRIVATE_KEY`
  - `AGENT_FACTORY_ADDRESS`
- Private key is read from env only and is not logged or returned.
- Execution uses Arc Testnet USDC: `0x3600000000000000000000000000000000000000`

Fresh market for backend buy_yes validation:
- Created through backend `create_market` execution endpoint.
- Intent id: `agent_intent_3c7ccaed030100b61eb3818aec277883`
- Market id: `backend-buy-yes-market-1779359535`
- Create transaction: `0xb7acd12fe390c1f3cc59543938710cd549855d4636682f9a841fd16fa500f8f3`
- Created agent market: `0xCB8D34fFdA32a3b58f355b92c3e720deCCF1C437`
- Backend create readback:
  - `readback.market_count == 6`
  - `readback.created_market == 0xCB8D34fFdA32a3b58f355b92c3e720deCCF1C437`
  - `readback.is_market == true`

Backend buy_yes validation:
- Intent id: `agent_intent_296671404fb958240b433ea31c79f0b5`
- Market contract address: `0xCB8D34fFdA32a3b58f355b92c3e720deCCF1C437`
- Amount: `1000000`
- USDC approve transaction: `0x5b3668792b391a7ec5fa892168569c26cf690e939ee8afa333ce54396871c1aa`
- `buyYes(1000000)` transaction: `0xc5d455a6f31bee34c443045575c2dedc5e62dfa53576d52b05916206a5fa8825`
- Backend response readback:
  - `status == executed`
  - `action == buy_yes`
  - `network == arc_testnet`
  - `broadcast_performed == true`
  - `approve_transaction_hash == 0x5b3668792b391a7ec5fa892168569c26cf690e939ee8afa333ce54396871c1aa`
  - `transaction_hash == 0xc5d455a6f31bee34c443045575c2dedc5e62dfa53576d52b05916206a5fa8825`
  - `readback.yes_positions == 1000000`
  - `readback.total_yes == 1000000`
  - `readback.total_collateral == 1000000`
  - `readback.usdc_balance == 1000000`

Independent onchain readback:
- Approve receipt status: `1 (success)`
- Buy YES receipt status: `1 (success)`
- `yesPositions(0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC) == 1000000`
- `totalYes() == 1000000`
- `totalCollateral() == 1000000`
- `USDC.balanceOf(0xCB8D34fFdA32a3b58f355b92c3e720deCCF1C437) == 1000000`
- `USDC.allowance(deployer, market) == 0`

## Non-Claims

Not implemented yet:
- frontend chat UI
- Circle Agent Wallet
- autonomous trading
- mainnet settlement
- real funds
- smart contract audit
