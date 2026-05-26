# SignalArc Project State

This file is the handoff source of truth for continuing SignalArc work in a new chat or coding session.

## Project

- Name: SignalArc
- Repository: https://github.com/wahyu241205/SignalArc
- Positioning: Arc-native prediction market infrastructure, not a Polymarket clone
- Core idea: Convert market opinions into structured probability signals, support USDC-settled event markets, expose market intelligence through APIs, and make the system usable by creators, institutions, developers, and AI agents.

## Working Style

- Use English for repository-facing guidance and project notes.
- Use controlled debugging.
- One step at a time.
- Do not jump phases without justification.
- Do not provide multiple unrelated solutions at once.
- For Arc and Circle behavior, use official documentation only.
- If behavior is not documented, state: unknown / not documented.
- Do not invent integrations or system behavior.
- Keep commits small and reviewable.

## Current Overall Status

Foundation complete. Backend and frontend local MVP integration verified. Product MVP is not live yet.

Estimated progress toward live grant submission: 30-35%.

## Phase 1 — Foundation Repo and Local Infra

Status: Mostly done, about 80%.

Done:

- Created SignalArc GitHub repository.
- Set correct Git remote: https://github.com/wahyu241205/SignalArc
- Configured Git identity locally for this repo.
- Converted project into modular monorepo.
- Added pnpm workspace.
- Moved Next.js frontend into `apps/web`.
- Removed npm `package-lock.json`.
- Added `pnpm-lock.yaml`.
- Added monorepo `.gitignore` rules.
- Added `AGENTS.md` project instructions.
- Set `CLAUDE.md` to reference `AGENTS.md`.
- Verified frontend runs with `pnpm dev:web`.
- Cleaned old WizPay Docker containers, images, and volumes.
- Added Docker Compose for SignalArc PostgreSQL.
- PostgreSQL runs locally on `127.0.0.1:15433`.
- Added `backend/.env.example`.
- Created local `backend/.env`, ignored by Git.
- Installed Go `1.26.3` in WSL.
- Installed `golang-migrate` CLI.
- Repo pushed and clean.

Not done:

- Add README architecture section.
- Add docs folder.
- Add contracts folder.
- Add CI workflow.
- Add Makefile or task runner commands.
- Add project roadmap file.

## Phase 2 — Core Database Schema

Status: COMPLETE.

Done:

- Created users table migration.
- Applied users migration locally.
- Committed users migration.
- Created wallets table migration.
- Applied wallets migration locally.
- Committed wallets migration.
- Created markets table migration.
- Applied markets migration locally.
- Committed markets migration.
- Created positions table migration.
- Applied positions migration locally.
- Committed positions migration.
- Created trades table migration.
- Applied trades migration locally.
- Created liquidity table migration.
- Applied liquidity migration locally.
- Created resolutions table migration.
- Applied resolutions migration locally.
- Created settlements table migration.
- Applied settlements migration locally.
- Created oracle_events table migration.
- Applied oracle_events migration locally.
- Created audit_logs table migration.
- Applied audit_logs migration locally.
- Created api_keys table migration.
- Applied api_keys migration locally.
- Created webhooks table migration.
- Applied webhooks migration locally.
- Created agent_access table migration.
- Applied agent_access migration locally.

Completed migrations:

- `000001_create_users_table`
- `000002_create_wallets_table`
- `000003_create_markets_table`
- `000004_create_positions_table`
- `000005_create_trades_table`
- `000006_create_liquidity_table`
- `000007_create_resolutions_table`
- `000008_create_settlements_table`
- `000009_create_oracle_events_table`
- `000010_create_audit_logs_table`
- `000011_create_api_keys_table`
- `000012_create_webhooks_table`
- `000013_create_agent_access_table`

Current local database tables:

- `users`
- `wallets`
- `markets`
- `positions`
- `trades`
- `liquidity`
- `resolutions`
- `settlements`
- `oracle_events`
- `audit_logs`
- `api_keys`
- `webhooks`
- `agent_access`
- `schema_migrations`

Current local database migration status:

- version: `13`
- dirty: `false`

Validation results:

- final schema review: `PASS`
- rollback/down migration test: `PASS`
- local demo seed validation: `PASS`
- migration status: `version=13`, `dirty=false`

Not done:

- None for Phase 2 core schema.

Next step:

Start Phase 5 Arc / Circle Integration planning by verifying current official Arc and Circle documentation before implementing any integration.

## Phase 3 — Backend API MVP

Status: Complete.

Done:

- Added Go backend module.
- Added backend dependencies:
  - `chi`
  - `pgx`
  - `zerolog`
  - `validator`
  - `godotenv`
- Added backend API entrypoint.
- Added health endpoint.
- Added backend config package.
- Added PostgreSQL database connection package using pgxpool.
- Added startup database ping.
- Added readiness endpoint.
- Added schema validation endpoint for Phase 2 tables and migration version.
- Verified `/health`, `/readyz`, and `/schema/validate` locally.
- Refactored backend API route registration into `internal/api`.
- Added consistent JSON response/error helper package in `internal/httpjson`.
- Added read-only repository layer for users, wallets, and markets.
- Added read-only market listing endpoint.
- Added read-only market detail endpoint.
- Added read-only repository layer for positions, trades, resolutions, and settlements.
- Added read-only position endpoints.
- Added read-only resolution status endpoint.
- Added read-only settlement status endpoints.
- Added market creation endpoint.
- Added market creation request validation.
- Added trade intent endpoint baseline.
- Added agent-readable market API baseline.
- Added request ID middleware.
- Added structured request logging middleware.
- Added panic recoverer middleware.
- Split oversized backend API router into focused handler, response, middleware, and validation files.
- Completed final backend smoke validation for Phase 3 MVP.
- Verified:

```bash
curl http://127.0.0.1:4000/health
```

Expected result:

```json
{"status":"ok"}
```

Current endpoints:

- `GET /health`
- `GET /readyz`
- `GET /schema/validate`
- `GET /markets`
- `GET /markets/{id}`
- `POST /markets`
- `POST /trade-intents`
- `GET /agent/markets`
- `GET /users/{user_id}/positions`
- `GET /markets/{market_id}/positions`
- `GET /markets/{market_id}/resolution`
- `GET /users/{user_id}/settlements`
- `GET /markets/{market_id}/settlements`

Not done:

- Resolver endpoint.
- Request validation.
- CORS middleware.
- Auth middleware.
- API key middleware.
- Unit tests.
- Integration tests.

## Phase 4 — Frontend MVP

Status: COMPLETE.

Done:

- Next.js app exists.
- Frontend moved to `apps/web`.
- shadcn/ui base components exist.
- Frontend runs locally.
- Added frontend API environment example at `apps/web/.env.example`.
- Added typed frontend API client at `apps/web/src/lib/api.ts`.
- Added frontend API helpers for Phase 3 backend endpoints:
  - `GET /health`
  - `GET /readyz`
  - `GET /schema/validate`
  - `GET /markets`
  - `GET /markets/{id}`
  - `GET /agent/markets`
- Added market list page at `apps/web/src/app/markets/page.tsx`.
- Added market list component at `apps/web/src/features/markets/market-list.tsx`.
- Wired market list UI to `GET /markets` through the frontend API client.
- Added market list loading, empty, and error states.
- Added market detail page at `apps/web/src/app/markets/[id]/page.tsx`.
- Added market detail component at `apps/web/src/features/markets/market-detail.tsx`.
- Wired market detail UI to `GET /markets/{id}` through the frontend API client.
- Added market detail loading and error states.
- Added frontend `createMarket()` API helper for `POST /markets`.
- Added create market page at `apps/web/src/app/markets/new/page.tsx`.
- Added create market form at `apps/web/src/features/markets/create-market-form.tsx`.
- Added create market idle, submitting, success, and error states.
- Added frontend `createTradeIntent()` API helper for `POST /trade-intents`.
- Added trade intent panel at `apps/web/src/features/markets/trade-intent-panel.tsx`.
- Wired trade intent panel into the market detail UI.
- Added trade intent idle, submitting, success, and error states.
- Kept trade intent UI explicitly intent-only: no wallet execution, no onchain settlement, no position update.
- Added frontend portfolio API helpers for `GET /users/{user_id}/positions` and `GET /users/{user_id}/settlements`.
- Added portfolio page at `apps/web/src/app/portfolio/page.tsx`.
- Added portfolio view at `apps/web/src/features/portfolio/portfolio-view.tsx`.
- Added read-only portfolio idle, loading, empty, error, and loaded states.
- Kept portfolio UI explicitly read-only: no wallet balance, no claim flow, no settlement mutation.
- Added frontend market resolution API helpers for `GET /markets/{market_id}/resolution` and `GET /markets/{market_id}/settlements`.
- Added market resolution panel at `apps/web/src/features/markets/market-resolution-panel.tsx`.
- Wired market resolution panel into the market detail UI.
- Added read-only resolution loading, empty/not-found, error, and loaded states.
- Kept resolution UI explicitly read-only: no resolver submission, no claim flow, no settlement execution, no eligibility inference.
- Added intelligence page at `apps/web/src/app/intelligence/page.tsx`.
- Added intelligence dashboard at `apps/web/src/features/intelligence/intelligence-dashboard.tsx`.
- Wired intelligence dashboard to `GET /agent/markets` through the frontend API client.
- Added intelligence loading, empty, error, and loaded states.
- Kept intelligence UI explicitly read-only: no agent execution, no paid access, no API key enforcement, no trading automation.
- Completed local Phase 4 smoke validation:
  - backend `/health` returned 200
  - backend `/readyz` returned 200
  - backend `/schema/validate` returned 200 with `migration_version=13`, `dirty=false`, and no missing tables
  - frontend `/markets` returned 200
  - frontend `/markets/new` returned 200
  - frontend `/markets/{id}` returned 200
  - frontend `/portfolio` returned 200
  - frontend `/intelligence` returned 200
  - `GET /agent/markets` returned 200
  - `POST /markets` returned 201
  - `POST /trade-intents` returned expected backend validation error `market_not_open` for a DRAFT market
  - `pnpm lint` passed
  - `pnpm exec tsc --noEmit` passed
  - `pnpm build` passed

Not done:

- SignalArc landing page.
- Wallet connect UI integration.
- Responsive layout polish.
- Demo-ready UI.

## Phase 5 — Arc / Circle Integration

Status: COMPLETE.

Done:

- Official documentation sources identified:
  - Arc docs: https://docs.arc.network/
  - Arc prediction market blueprint: https://www.arc.network/blog/build-institutional-grade-prediction-markets-on-arc-arc-blueprints
  - Circle developer docs: https://developers.circle.com/
  - Circle Agents docs: https://agents.circle.com/
  - Circle grants page: https://www.circle.com/grant
- Project rule added: Arc/Circle behavior must use official docs only.
- Reviewed Arc official documentation snapshots for:
  - Connect to Arc
  - Deploy on Arc
  - EVM compatibility
- Confirmed documented Arc Testnet network details:
  - Network: Arc Testnet
  - Chain ID: `5042002`
  - Primary RPC: `https://rpc.testnet.arc.network`
  - Explorer: `https://testnet.arcscan.app`
  - Faucet: `https://faucet.circle.com`
  - Currency / native gas token: USDC
- Confirmed Arc is currently documented as testnet phase; Arc mainnet remains unknown / not documented from reviewed Arc pages.
- Confirmed Arc EVM compatibility is documented.
- Confirmed Solidity smart contract deployment to Arc Testnet with Foundry is documented.
- Confirmed Arc docs mention standard Ethereum tooling including Solidity, Foundry, Hardhat, and Viem.
- Confirmed Arc uses USDC as native gas.
- Confirmed Arc deterministic finality is documented.
- Confirmed Arc USDC has dual interface behavior:
  - native balance uses 18 decimals
  - ERC-20 interface uses 6 decimals
  - both interfaces share one underlying USDC balance
- Collected preliminary Circle Docs AI response for Circle Wallets, CCTP, Gateway, Paymaster, Circle Agent Stack, and Arc Testnet support; this still needs exact official page snapshots before implementation decisions.
- Completed Phase 5.1 Arc / Circle research as an internal/local-only planning note under ignored `docs/internal/`.
- Completed Phase 5.2 wallet strategy as an internal/local-only planning note under ignored `docs/internal/`.
- Completed Phase 5.3 USDC integration planning as an internal/local-only note under ignored `docs/internal/`.
- Completed Phase 5.4 agent wallet boundary planning as an internal/local-only note under ignored `docs/internal/`.
- Completed Phase 5.5 contract requirements handoff as an internal/local-only note under ignored `docs/internal/`.

Not done:

- Do not implement contracts until Phase 6.

## Phase 6 — Contracts / Settlement Prototype

Status: COMPLETE.

Done:

- Stack decision: Solidity + Foundry + OpenZeppelin.
- Created `contracts` folder with `src`, `test`, and `script` placeholders.
- Added minimal Foundry workspace config at `contracts/foundry.toml`.
- Added `contracts/README.md` with prototype boundaries.
- Added minimal `SignalArcMarket` contract scope for binary market lifecycle without collateral, settlement, or claim logic.
- Added test-only `MockUSDC` collateral model and Foundry tests for market lifecycle and collateral assumptions.
- Added local binary position and trade prototype logic with test-only MockUSDC collateral accounting.
- Added manual resolution and local claim/refund prototype logic with Foundry coverage.
- Completed contract test coverage hardening and added `contracts/SECURITY_BOUNDARIES.md`.
- Completed final Phase 6 review and added `contracts/PHASE_6_REVIEW.md`.
- Completed manual Arc Testnet deployment, source verification, and onchain smoke tests; recorded results in `contracts/ARC_TESTNET_DEPLOYMENT.md`.

Not done:

- Backend/frontend contract integration is deferred to a later integration phase.
- Production deployment remains not approved.

## Phase 7 — Live Deployment

Status: Not started, 0%.

Target deployment:

- Frontend: Vercel.
- Backend: GCP Cloud Run.
- Database: hosted PostgreSQL.
- Domains:
  - `signalarc.xyz`
  - `app.signalarc.xyz`
  - `api.signalarc.xyz`
  - `docs.signalarc.xyz`

Not done:

- Buy/configure domain DNS.
- Deploy frontend.
- Deploy backend.
- Configure production database.
- Configure backend env vars.
- Configure CORS.
- Configure HTTPS.
- Add API health URL.
- Add frontend/backend integration.
- Add deployment README.
- Add monitoring/logging baseline.
- Verify live demo flow.

## Phase 8 — Grant Submission Package

Status: Not started, 0%.

Required before submission:

- Live app URL.
- GitHub repo public or shareable.
- Demo video.
- Pitch deck.
- Architecture diagram.
- Technical README.
- Product README.
- Roadmap.
- Clear Arc usage explanation.
- Clear Circle/USDC usage explanation.
- MVP screenshots.
- Demo user flow.
- Risk/compliance disclaimer.
- Grant form answers.

Do not submit yet until:

- Live MVP can be clicked.
- Market creation flow works.
- Market detail flow works.
- User/wallet flow is visible.
- Trade or simulated trade flow works.
- Resolution/settlement flow is visible.
- Arc + Circle relevance is clear.
- README explains why SignalArc exists.

## Frontend Wallet and Production-Facing UI Polish

Status: COMPLETE.

Done:

- Added external wallet connection UI using RainbowKit, Wagmi, and Viem.
- Added Arc Testnet frontend chain configuration.
- Added app-level Web3 providers for wallet connection.
- Added responsive dark Web3 layout and navigation.
- Added visible wallet connection control in the header.
- Added connected wallet address display through RainbowKit.
- Added wrong-network warning for non-Arc Testnet wallet connections.
- Reworked the landing page into a production-facing SignalArc product homepage.
- Reduced internal developer wording from user-facing pages.
- Moved contract metadata away from the primary landing hero into lower-priority developer/network context.
- Reworked markets, market detail, create market, portfolio, and intelligence pages for a cleaner user-facing experience.
- Kept trade flow intent-only.
- No Circle SDK added.
- No Circle API key added.
- No onchain write execution added.
- No USDC approve, open position, claim, refund, or resolver transaction UI added.
- No backend, Docker, contract, DNS, deployment, or mainnet changes added in this step.

Validation:

- `pnpm --dir apps/web lint` passed.
- `pnpm --dir apps/web exec tsc --noEmit` passed.
- `pnpm --dir apps/web build` passed.
- Browser wallet connection UI rendered.
- Wallet address displayed after connection.
- Existing frontend pages rendered with the polished dark Web3 UI.

## Local Docker Backend Step

Status: COMPLETE.

Done:

- Added `backend/Dockerfile` for local backend container builds.
- Added `backend/.dockerignore` to exclude local env and unnecessary files.
- Added `signalarc-backend` service to `docker-compose.yml`.
- Local Docker Compose now runs PostgreSQL and backend API.
- Backend container connects to PostgreSQL through Docker network using `signalarc-postgres:5432`.
- Backend is exposed only on host loopback at `127.0.0.1:4000`.
- Frontend remains a manual Next.js dev server.
- Smart contract is not Dockerized.
- No production deployment, DNS change, contract redeploy, or mainnet config added.

Validation:

- `docker compose build signalarc-backend` passed.
- `docker compose up -d signalarc-postgres signalarc-backend` passed.
- `GET /health` passed.
- `GET /readyz` passed.
- `GET /arc/contract` passed.
- `GET /markets` passed.

## Local MVP Integration Fix

Status: COMPLETE.

Done:

- Local frontend/backend MVP usability smoke test completed.
- Backend local CORS added for the local Next.js frontend.
- Backend read-only `GET /arc/contract` endpoint added.
- Frontend root page now shows SignalArc local MVP instead of the default Next.js template.
- Frontend can load markets from backend without "Failed to fetch".
- Local browser flow passed:
  - `/` — landing page renders
  - `/markets` — market list loads from backend
  - `/markets/new` — create market form renders
  - Created one test market successfully
  - Opened market detail page
  - Submitted one trade intent with status `not_executed`
  - `/portfolio` — portfolio page renders
  - `/intelligence` — intelligence dashboard renders
- Arc Testnet contract reference visible as prototype/testnet metadata only.
- No wallet/onchain write execution added.
- No Circle integration added.
- No production deployment approved.
- No contract redeploy performed.
- Temporary demo domain purchased: `signalarc.fun`.
- DNS/live deployment remains pending and not approved.

## CI / Quality Gate

Status: COMPLETE.

Done:

- Added first CI-only GitHub Actions workflow at `.github/workflows/ci.yml`.
- Workflow runs on `pull_request` and push to `main`.
- Added least-privilege workflow permission: `contents: read`.
- Frontend CI validates pnpm frozen install, lint, TypeScript check, and production build.
- Backend CI runs `go test ./...` from `backend`.
- Contracts CI installs Foundry and runs `forge test` from `contracts`.
- No deployment, production secrets, Docker publish, Vercel, or GCP automation added in this step.

Validation:

- `pnpm install --frozen-lockfile` passed.
- `pnpm --dir apps/web lint` passed.
- `pnpm --dir apps/web exec tsc --noEmit` passed.
- `pnpm --dir apps/web build` passed under Node 22.
- `cd backend && go test ./...` passed.
- `cd contracts && forge test` passed.

## Current Last Completed Step

- Multi-tenant agent onboarding session foundation added for isolated pending onboarding sessions and durable per-agent session boundaries. Circle OTP/provisioning remains not implemented.

## Live AI Agent Transaction MVP

Status: IN PROGRESS - backend Circle Agent Wallet provider trading, payout, cancel, and refund lifecycle actions are validated on Arc Testnet from the host-shell backend runtime; ChatGPT Custom GPT external client trigger is validated for health, intent preview, confirm, and real `create_market` execution through a temporary ngrok tunnel; Docker/Cloud Run runtime strategy, production readiness, and mainnet readiness remain pending.

Current checkpoint state:

- Agent Wallet Boundary documented in `project-roadmap/agent-mcp.md`.
- Previous deployer-signed Arc Testnet transactions are explicitly not accepted as agent-wallet validation.
- `contracts/.env` `PRIVATE_KEY` must not be used as an agent wallet.
- `AGENT_EXECUTOR_PRIVATE_KEY` is not the final Agent Wallet design.
- Circle Agent Wallet is the preferred provider if the official CLI path works on `ARC-TESTNET`.
- Circle CLI was installed globally as an external tool with `npm install -g @circle-fin/cli`; no SignalArc package manager files were intentionally changed for Circle CLI.
- `circle --version` returned `0.0.3`.
- Circle CLI safe help discovery was completed.
- Circle CLI auth/chain discovery is blocked by user-controlled Terms acceptance before login or OTP; `circle terms accept` must be run manually by the user before authenticated discovery continues.
- Official Circle docs show `circle wallet login you@example.com --testnet`, but the installed `circle wallet login --help` output did not list a `--testnet` option before Terms/login validation.
- User-provided Circle CLI evidence now proves a real `ARC-TESTNET` Circle Agent Wallet at `0x96d5051a005547eba149f71604ccf58ae1a7c950`, distinct from deployer/resolver `0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`.
- User-provided funding proof showed `20 USDC`.
- User-provided create-market evidence recorded tx `0x7142dbd7eebe7cbfb19199d9984efa5cef814d0e6038c17b98f2e98cc731cacf`, source address `0x96d5051a005547eba149f71604ccf58ae1a7c950`, state `COMPLETE`, `marketCount() == 8`, and `allMarkets(7) == 0xAbCf081E456C1a11106deF590666A07B76D456f8`.
- User-provided market readback proved market `0xAbCf081E456C1a11106deF590666A07B76D456f8`, collateral token `0x3600000000000000000000000000000000000000`, admin/resolver both `0x96d5051a005547eba149f71604ccf58ae1a7c950`, and `isOpen == true`.
- User-provided buy YES evidence recorded approve tx `0xeb7304b0a1be9f5dc575f62fb705dfaf384bc720da13f7e4ffe9563442c036ca`, buy tx `0xe311d999e15e6f34fa6f623a8f27bc724c665d7c3296632460339326b6094b16`, and readback `yesPositions(agent wallet) == 1000000`, `totalYes == 1000000`, `totalCollateral == 1000000`.
- User-provided buy NO evidence recorded approve tx `0x6ea6a10293a4df5d7ed50e077821115571787d8e9d6b9507a984ebf33fc52a9b`, buy tx `0xaefe8bcdcec794c811d615517f0dfa800b9e263631200a74c85d000374aa8f24`, and readback `noPositions(agent wallet) == 1000000`, `totalNo == 1000000`, `totalCollateral == 2000000`, `USDC.balanceOf(market) == 2000000`.
- Continuation from this Codex shell is blocked because Circle CLI returns `AUTH_REQUIRED` or no active agent session for status/list/balance, even when `CIRCLE_ACCEPT_TERMS=1` is set for the process.
- Added persistent `agent_wallets` migration for user-owned Circle Agent Wallet onboarding.
- Added migration `000016_create_agent_onboarding_sessions` for `agent_onboarding_sessions` and `agent_sessions`.
- `POST /agent/onboarding/register` remains a registry-only convenience endpoint for creating the final agent wallet mapping; it does not model OTP or Circle session isolation by itself.
- Added `POST /agent/onboarding/start` as the pending-session foundation for per-user, per-agent onboarding. It creates `pending_otp` state only and returns `circle_otp_verification_not_implemented`.
- Corrected `POST /agent/onboarding/start` to be agent-first and email-based: `user_wallet` is no longer required for initial agent onboarding and belongs to frontend wallet-connect or user transaction flows.
- Added disabled-by-default Circle Agent Wallet OTP start skeleton behind `CIRCLE_AGENT_ONBOARDING_OTP_START_ENABLED=false` by default. When enabled for a controlled dev runtime, it can call the Circle CLI login init runner, store only a hashed request reference plus expiry, and return `circle_otp_required` without exposing the raw request ID.
- Added disabled-by-default Circle Agent Wallet OTP verify skeleton at `POST /agent/onboarding/verify`. It uses the in-memory request ID from OTP start, consumes it on successful fake-runner verification, updates onboarding status to `verified`, and returns `agent_wallet_resolution_not_implemented`.
- Corrected OTP verify CLI completion command shape to the documented Circle form: `circle wallet login --request <request-id> --otp <code>`.
- Added sanitized server-side diagnostics for Circle OTP verify CLI failures; diagnostics redact the raw request ID and OTP and the API response remains the generic `circle_otp_verify_failed` error.
- After commit `b1280f0`, Custom GPT onboarding start returned `circle_otp_start_failed` even though the OTP email was delivered.
- Root cause identified: Circle CLI OTP init can print the request ID as text, while backend OTP start handling only accepted JSON request ID fields.
- Updated OTP start handling to accept JSON `request_id` / `requestId` and documented text-style printed request IDs.
- Raw Circle request IDs remain hidden from API responses and are stored only in memory for verify; the database stores only the hash.
- Added sanitized Circle CLI OTP start diagnostics that redact email and request ID.
- Added tests for JSON/text request ID extraction, sanitized start diagnostics, and successful `/agent/onboarding/start` with text request ID output.
- Validation: `go test ./...` passed.
- After commit `a14f17d`, Custom GPT onboarding start still returned `circle_otp_start_failed` although the OTP email was delivered.
- Observed Circle CLI output placed the request ID inside a nested message containing the documented completion command shape: `circle wallet login --request <id> --otp <code>`.
- Updated OTP start handling to extract request IDs from completion-command output and treat that as success even if the CLI runner returned an error.
- Strengthened sanitized diagnostics to redact request IDs embedded in command text.
- Validation: `go test ./...` passed.
- Runtime validation after commit `65cebdb` passed through Custom GPT/ngrok:
  - `/agent/onboarding/start` returned `pending_otp` and `next_step=circle_otp_required`.
  - Onboarding ID: `agent_onboarding_ad068b4c0538605db092819aa55df08d`.
  - Agent ID: `agent_adenhusen65_chatgpt_008`.
  - `/agent/onboarding/verify` accepted the emailed OTP and returned status `verified`.
  - Backend returned expected boundary `agent_wallet_resolution_not_implemented`.
  - No market, intent, confirmation, execution, or user wallet request was performed.
- Post-OTP wallet/session activation added.
- Verified OTP onboarding can now resolve a Circle Agent Wallet address from read-only Circle CLI list output.
- Backend registers the resolved wallet and creates an active agent session.
- Backend exposes read-only wallet/session status and balance for Custom GPT.
- No funding, transfer, market creation, or transaction execution is performed by the activation step.
- Validation: `go test ./...` passed.
- Added read-only onboarding/session status APIs: `GET /agent/onboarding/{onboarding_id}` and `GET /agent/sessions/{agent_id}`.
- Multi-tenant/session isolation state now separates user email, user wallet, source client, channel, pending onboarding, and activated agent-session boundaries without storing Circle session secrets.
- Backend now registers agent wallets through DB-backed `POST /agent/wallets` in production routing.
- Backend now returns registered wallet metadata through `GET /agent/wallets/{agent_id}`.
- Backend now disables registered wallets through `POST /agent/wallets/{agent_id}/disable`.
- Test handlers use an in-memory registry only as a test double; production registration is no longer the temporary in-memory development path.
- Backend agent intents now carry `agent_id` and optional `agent_wallet_address`.
- Backend agent intents now carry `source_client` and optional `client_request_id` so WhatsApp, Telegram, ChatGPT, Claude, and web clients can share one channel-agnostic API contract.
- Backend execution now rejects missing agent wallets, deployer/resolver wallet reuse, user-wallet reuse, disallowed actions, wrong chain, inactive wallet, mismatched wallet address, and unsupported provider execution.
- Added a guarded backend Circle CLI executor provider for `circle_agent_wallet`.
- Provider execution mode is `circle_agent_wallet_cli`.
- Provider is disabled by default with `CIRCLE_AGENT_WALLET_EXECUTION_ENABLED=false`.
- Provider config added: `CIRCLE_CLI_PATH`, `CIRCLE_AGENT_WALLET_CHAIN`, and `CIRCLE_AGENT_WALLET_TIMEOUT_SECONDS`.
- Provider supports `create_market`, `buy_yes`, `buy_no`, `close_market`, `resolve_market`, `claim_payout`, `cancel_market`, and `claim_refund` when the registered agent wallet allowlist permits the action.
- Lifecycle support uses repo-confirmed `SignalArcAgentMarket` functions only:
  - `close_market`: `closeMarket()`
  - `resolve_market`: `resolve(uint8)` from contract `resolve(Outcome winningOutcome_)`
  - `claim_payout`: `claimPayout()`
  - `cancel_market`: `cancelMarket()`
  - `claim_refund`: `claimRefund()`
- Lifecycle readbacks use repo-confirmed getters only: `status()`, `winningOutcome()`, `claimablePayout(address)`, `claimableRefund(address)`, `hasClaimed(address)`, `isOpen()`, and USDC `balanceOf(address)`.
- Provider uses Circle CLI `wallet execute` for writes and `contract query` for readbacks, with JSON-only parsing and sanitized errors.
- Provider never calls Circle login, never accepts OTP, and never stores Circle tokens, session files, private keys, or deployer keys.
- Runtime validation was performed from the host shell backend on port `4001`, not the Docker backend container, because the Docker backend container does not have Circle CLI installed.
- Host Circle Agent Wallet session was valid for the recorded backend provider run.
- DB migration was advanced to version `15`, `dirty=false`.
- Agent wallet was registered through DB-backed `POST /agent/wallets`.
- Registered agent:
  - `agent_id`: `agent_desi_001`
  - `user_email`: `desi33905@gmail.com`
  - `user_wallet`: `0x153D2Fc8334a84a37B7A7cF9deFA5Cb401a36FdC`
  - `agent_wallet_address`: `0x96d5051a005547eba149f71604ccf58ae1a7c950`
  - `wallet_provider`: `circle_agent_wallet`
  - `chain`: `ARC-TESTNET`
  - `allowed_actions`: `create_market`, `buy_yes`, `buy_no`
  - `source_client`: `manual_backend_runtime`
- Backend provider `create_market` evidence:
  - `intent_id`: `agent_intent_f76bd9653e9ce3f6269023a25e7c6b8c`
  - `execution_mode`: `circle_agent_wallet_cli`
  - Transaction: `0x7aa51a0d19b163a3a88ae16ac4a88a1cdbb3090cad5d0ccc54f828b166f74e5d`
  - Created market: `0x38aE7E0133e9594F597F913884cbDa619A950523`
  - Readback: `market_count == 9`, `is_market == true`
- Backend provider `buy_yes` evidence:
  - `intent_id`: `agent_intent_d6f9cd78896886f954a05a164c17c067`
  - `execution_mode`: `circle_agent_wallet_cli`
  - Market: `0x38aE7E0133e9594F597F913884cbDa619A950523`
  - Approve transaction: `0x09d0a418c34b0e54a31bc3a2a7bfba85218eba27e84becdcbd89e5b63b8bb387`
  - `buyYes` transaction: `0xa1fadb400aa8b4babca0c936698e686eeaac3ae408b22d1e37960901a5c48ade`
  - Readback: `yes_positions == 1000000`, `total_yes == 1000000`, `total_collateral == 1000000`, `USDC.balanceOf(market) == 1000000`
- Backend provider `buy_no` evidence:
  - `intent_id`: `agent_intent_ed7050305f62798d8472b7f48e538ff8`
  - `execution_mode`: `circle_agent_wallet_cli`
  - Market: `0x38aE7E0133e9594F597F913884cbDa619A950523`
  - Approve transaction: `0x40ba807e24e2fcfeba22f21575920d3dfe5389f7f00320ab5a74fb46f06c6dc8`
  - `buyNo` transaction: `0xb36288fc40a69765d62679982fdd3319d09b27c010e2eb9caafa9c6508d03e9c`
  - Readback: `no_positions == 1000000`, `total_no == 1000000`, `total_collateral == 2000000`, `USDC.balanceOf(market) == 2000000`
- Latest full lifecycle runtime evidence was performed from the host shell backend on `APP_PORT=4001` with execution mode `circle_agent_wallet_cli`, agent `agent_desi_001`, agent wallet `0x96d5051a005547eba149f71604ccf58ae1a7c950`, wallet provider `circle_agent_wallet`, and chain `ARC-TESTNET`.
- Docker backend was not used for the latest execution because the Docker backend container still does not include the Circle CLI/session strategy.
- Backend provider payout lifecycle evidence:
  - `create_market` intent `agent_intent_ecc88160f7e2b908fc498c3dff66fbe7`, transaction `0x7913fd51b38b147cfc6936da7eb7166a97a156351c9ccce4264d48d31fc91ae9`, market `0x38D4317fcB0C82e5EC2407a89c311b3Be8059CD0`, readback `market_count == 10`, `is_market == true`
  - `buy_yes` intent `agent_intent_e1d3889f3252ac318b7b435b2a789d9d`, approve transaction `0x9c66332ad2d798126118b961ad9005ab7f2055649b46d808cc979ccf40eee3f7`, buy transaction `0xdd4448fdf237f13bc9e90737f27b7a1e912ee8d34195bae311cc2bac15aaa95d`, readback `yes_positions == 1000000`, `total_yes == 1000000`, `total_collateral == 1000000`, `USDC.balanceOf(market) == 1000000`
  - `close_market` intent `agent_intent_6cab73b950296ac2dcc3a5414f4e1613`, transaction `0x5accb4dafee2a27be032709e427e25eede4e7eb67f54a06bdf0ba82e5f4a013e`, readback `market_status == 1`, `is_open == false`
  - `resolve_market` intent `agent_intent_4073e0f29f842a0110a1e099d9fa50b0`, transaction `0xc45d8612768591736425a743520b5e432b9a424ef9552cfbbc1bb04d785c874b`, readback `market_status == 2`, `winning_outcome == 1`, `claimable_payout == 1000000`, `has_claimed == false`, `USDC.balanceOf(market) == 1000000`
  - `claim_payout` intent `agent_intent_65284f22e52fce32d6a3efc2fa6163cd`, transaction `0xc80bb9dd7f6924c93c1722d7c5f1136c076403d052ed792f98ed2d7abd59568f`, readback `market_status == 2`, `winning_outcome == 1`, `claimable_payout == 1000000`, `has_claimed == true`, `USDC.balanceOf(market) == 0`
- Backend provider refund lifecycle evidence:
  - `create_market` intent `agent_intent_b167e466b198c7f909c388c76b683f5d`, transaction `0x1340817f922aaa7ae181789bbaf2b7bff13a426b0397524150d01ae869d2a033`, market `0xbfd93169DAFf0610EA10E1221B9a2a6552379648`, readback `market_count == 11`, `is_market == true`
  - `buy_yes` intent `agent_intent_b06dd357e62778ed1a0527800690473a`, approve transaction `0x83bdc164512979296c385f51b4b6b1df51c741f2157cc93944b7f15c1328f487`, buy transaction `0xa85c60d903a1b1e30f39da95a37ca356d61d11cc20c942aaa22f0740a06945ab`, readback `yes_positions == 1000000`, `total_yes == 1000000`, `total_collateral == 1000000`, `USDC.balanceOf(market) == 1000000`
  - First backend provider `cancel_market` attempt failed: intent `agent_intent_634dad5fc81b7a155e3854d8e14bb135`, backend response `502 agent_execution_failed`; onchain check after failure showed the market was still Open. This failed attempt is not validated backend cancel evidence.
  - Manual Circle CLI `cancelMarket` then succeeded on the refund market with transaction `0xbf8b98862ed691c0023643ab72ee71d8422e434868b2806f3390b4ffc88fe21b`; this was manual Circle CLI, not backend provider evidence.
  - Backend provider `claim_refund` succeeded after manual cancel: intent `agent_intent_c519c98114f8a5c113ab4bc6dfefbcae`, transaction `0x40f3f4e1737340dbbbff92e3020d0cfa6dbd7d7bbca8a1ecb580b1c0cdc43dfd`, readback `market_status == 3`, `claimable_refund == 1000000`, `has_claimed == true`, `USDC.balanceOf(market) == 0`
- Backend provider cancel-only validation evidence:
  - `create_market` intent `agent_intent_666c781a7a3fe2f939bd5342e331fd67`, transaction `0x62f5f43b834a09f4f4c78e3bb365403a2b76374d0030a91ccd3e3571fd3c9d12`, market `0x928F3F9Cb43811837C0e8D4FA40c24A4f083B3Ed`, readback `market_count == 12`, `is_market == true`
  - `cancel_market` intent `agent_intent_34f54290e5f59ffaff4a167986868c56`, transaction `0x3eae0d0508397e5ea515d417bdc5be5c38f40f3b0296b1cf424d99060cb92de4`, readback `market_status == 3`, `claimable_refund == 0`, `has_claimed == false`, `USDC.balanceOf(market) == 0`
- Current proven status:
  - Backend Circle provider `create_market` validated with a real Arc Testnet transaction.
  - Backend Circle provider `buy_yes` validated with a real Arc Testnet transaction.
  - Backend Circle provider `buy_no` validated with a real Arc Testnet transaction.
  - Backend Circle provider `close_market` validated with a real Arc Testnet transaction.
  - Backend Circle provider `resolve_market` validated with a real Arc Testnet transaction.
  - Backend Circle provider `claim_payout` validated with a real Arc Testnet transaction.
  - Backend Circle provider `cancel_market` validated after one failed backend provider attempt on the refund market and a later successful backend-only cancel-only test.
  - Backend Circle provider `claim_refund` validated after manual Circle CLI cancellation of the refund market.
- Backend-to-Circle-to-Arc path is proven for create, trade, close, resolve, payout claim, cancel, and refund claim flows on Arc Testnet.
- ChatGPT Custom GPT Action external client trigger evidence:
  - External client: ChatGPT Custom GPT Action.
  - Temporary tunnel: ngrok HTTPS URL to local backend.
  - Backend runtime: host shell on `APP_PORT=4001`.
  - Execution mode: `circle_agent_wallet_cli`.
  - Backend public test path: ChatGPT Custom GPT Action -> ngrok HTTPS -> `localhost:4001` SignalArc backend -> Circle Agent Wallet CLI provider -> Arc Testnet.
  - Local backend health returned 200: `http://127.0.0.1:4001/health` -> `{"status":"ok"}`.
  - Ngrok public health returned 200: `https://undamaged-commerce-juggling.ngrok-free.dev/health` -> `{"status":"ok"}`.
  - Custom GPT `getHealth` action returned `{"status":"ok"}`.
  - Custom GPT Action schema exposed `getHealth`, `createAgentIntent`, `confirmAgentIntent`, and `executeAgentIntent`.
  - Authentication was set to None only for temporary local tunnel testing. This is not production API authentication.
  - Initial intent preview and confirm were tested.
  - One confirm failed due to a typo in a manually copied intent ID. This was a user/operator copy error, not a backend failure.
  - One later execution attempt failed because `close_timestamp` was already in the past:
    - `close_timestamp` used: `1779439999`
    - current timestamp checked later: `1779452713`
    - probable root cause: `createMarket` reverted due to `closeTimestamp <= block.timestamp`
    - Circle Agent Wallet session was checked and testnet `tokenStatus` was `VALID`
    - This is recorded as stale input evidence, not a proven provider failure.
  - New valid preview used future `close_timestamp`: `1779456313`.
- Successful ChatGPT Custom GPT external execution:
  - Intent ID: `agent_intent_7d0bfd385329ba97cb7c1b88ada6f049`
  - Action: `create_market`
  - Status: `executed`
  - Execution mode: `circle_agent_wallet_cli`
  - Network: `arc_testnet`
  - Agent factory: `0x69aE770e8b2F96297101FeC4dc123B3801dA7d80`
  - Agent ID: `agent_desi_001`
  - Agent wallet: `0x96d5051a005547eba149f71604ccf58ae1a7c950`
  - Broadcast performed: `true`
  - Transaction hash: `0x1062e254f8640ffdc2d75d368754e5d698d42cb935bc7bcb113547c7a501aec2`
  - Created market: `0x6cef2f33F0F2a5e01E885176bAa17709d6A6a299`
  - Readback:
    - `market_count == 13`
    - `is_market == true`
- Current external client trigger status:
  - ChatGPT Custom GPT trigger validated for health, create intent preview, confirm intent, and execute `create_market` real Arc Testnet transaction.
  - WA, Telegram, and Claude triggers are still not tested.
- No secrets, `.env` files, frontend code, production deployment config, commits, pushes, or deploys were changed.

Current non-claims:

- WhatsApp, Telegram, and Claude live client triggers have not been tested yet.
- Docker/Cloud Run runtime does not yet include a Circle CLI/session strategy.
- No production readiness claim.
- No mainnet claim.
- No Circle policy limit claim on `ARC-TESTNET`.
- Circle session storage is not implemented.
- Circle wallet provisioning readback beyond read-only CLI wallet list output is unknown / not documented.
- Custom GPT / Claude / Telegram / WA onboarding schemas should not require `user_wallet` for start-agent onboarding.
- Backend restart before OTP verify requires onboarding restart because raw request IDs are not stored in the database.
- No real Circle CLI run was performed for the OTP start skeleton tests.
- The temporary ngrok tunnel is not a production API endpoint.
- Circle CLI command shapes in `project-roadmap/agent-mcp.md` are official-doc/help-discovery shapes only unless accompanied by exact authenticated output and onchain evidence.
- The phase is not complete.

## Circle CLI Warning-Prefix Wallet Resolver Fix

Status: COMPLETE.

Done:

- Runtime verification after commit `5697295` reached wallet resolution but failed because Circle CLI prepended Node deprecation warning text before JSON output.
- Root cause: `parseCircleAgentWallets` and `parseCircleAgentWalletBalances` in `circle_wallet_resolver.go` called `json.Unmarshal` directly on raw CLI output, which fails when Node.js prints `(node:...) [DEP0040] DeprecationWarning...` lines before the JSON payload.
- Added `extractJSONFromCLIOutput` helper that scans for the first `{` or `[` that starts valid JSON, skipping warning text and bracket characters in deprecation labels like `[DEP0040]`.
- Updated `parseCircleAgentWallets` and `parseCircleAgentWalletBalances` to use the helper before unmarshalling.
- Clean JSON output still parses identically (fast path).
- Output without any valid JSON object/array still returns a parse error.
- Sanitized diagnostics behavior unchanged.
- Validation: `go test ./...` passed.

## Cloud Run Backend Image Includes Circle CLI

Status: COMPLETE (image-level only; durable Circle CLI session persistence remains a production runtime concern).

Done:

- Cloud Run production reached the live domain `api.signalarc.fun` and the production database was migrated to version `18`.
- Production agent endpoints then reached the Circle runtime boundary because the Cloud Run runtime image did not include Node.js, npm, or the Circle CLI binary, so commands like `circle wallet login`, `circle wallet list`, and `circle wallet balance` could not be executed inside the deployed container.
- `backend/Dockerfile` runtime stage was updated to install `nodejs` and `npm` from Alpine packages and then install the Circle CLI globally with `npm install -g @circle-fin/cli`, then clean the npm cache.
- The runtime stage continues to copy the compiled `signalarc-api` binary from the Go builder stage into `/usr/local/bin/signalarc-api`.
- The image final layout exposes both `signalarc-api` and `circle` on `PATH` at `/usr/local/bin`.
- The existing `ENTRYPOINT ["signalarc-api"]` was preserved.
- No Circle credentials, OTP request IDs, session files, private keys, or API keys are baked into the image.
- This change does not solve durable Circle CLI session persistence; session strategy remains a production runtime concern outside the image.

Validation:

- `cd backend && go test ./...` passed; no Go test files were modified.
- `docker build -t signalarc-backend:circle-cli-test backend` succeeded.
- Inside the built image:
  - `command -v circle` returned `/usr/local/bin/circle`.
  - `circle --version` returned `0.0.3`.
  - `command -v signalarc-api` returned `/usr/local/bin/signalarc-api`.
  - `node --version` returned `v20.15.1` and `npm --version` returned `10.9.1`.

Not done:

- Production deploy of the updated image to Cloud Run was not performed in this step.
- Circle CLI session/login persistence inside Cloud Run is still unresolved.

## ARC-TESTNET Faucet Helper Endpoint

Status: COMPLETE.

Done:

- Added ARC-TESTNET faucet helper endpoint for active registered agent wallets at `POST /agent/wallets/{agent_id}/faucet`.
- Endpoint uses the registered `agent_wallet_address` from the SignalArc database only and does not accept any arbitrary recipient address from the request body.
- Endpoint requires the registered wallet to exist (`404 agent_wallet_not_found`), to be `active` (`409 agent_wallet_status_invalid`), and to be on `ARC-TESTNET` (`409 agent_wallet_chain_invalid`).
- Endpoint calls only the documented Circle CLI testnet faucet command shape `circle wallet fund --address <registered_agent_wallet_address> --chain ARC-TESTNET --token usdc --output json`.
- Endpoint never passes `--amount`, `--method`, `--open`, `--export`, `transfer`, `swap`, `execute`, or any mainnet funding option, and does not transfer, swap, execute contracts, create markets, or use mainnet funding.
- Faucet runner is disabled by default behind `CIRCLE_AGENT_WALLET_FAUCET_ENABLED=false`. When disabled or not wired, the endpoint returns `501 circle_agent_wallet_faucet_not_configured`.
- Provider failures return `502 circle_agent_wallet_faucet_failed` with sanitized API responses; CLI/runner output is never echoed back to clients and credential paths, tokens, OTP-like values, request IDs, session material, and emails are redacted from server-side diagnostics.
- Output parser reuses the existing warning-prefixed JSON extraction helper used by wallet list/balance, so Node deprecation warnings before JSON still parse cleanly.
- Successful JSON CLI output is returned under `result`. Successful text-only CLI output is returned under `result.message` after sanitization.
- Updated Custom GPT OpenAPI schema `project-roadmap/signalarc-custom-gpt-openapi.json` with the `requestAgentWalletFaucet` action, including `404`, `409`, `501`, and `502` error responses; server URL remains `https://api.signalarc.fun` and `/agent/onboarding/register` was not re-added.
- Validation: `go test ./...` passed.

Not done:

- Production deploy of the updated image to Cloud Run was not performed in this step.
- No funding, transfer, swap, contract execution, market creation, mainnet funding, login/logout, or session persistence was added.

## Public Documentation Refresh For Live Architecture And Faucet

Status: COMPLETE.

Done:

- Updated `docs/index.mdx`, `docs/API.md`, `docs/AGENT_API.md`, `docs/DEPLOYMENT_PLAN.md`, `docs/ARCHITECTURE.md`, and `docs/GRANT_READINESS.md` to describe the live architecture: `https://signalarc.fun` on Vercel for the frontend, `https://api.signalarc.fun` on GCP Cloud Run service `signalarc-backend-api` for the backend, GCP Cloud SQL PostgreSQL migrated to version 18, and the Cloud Run image bundling Node/npm and the Circle CLI (`@circle-fin/cli`) on `PATH`.
- Documented the agent flow `onboarding start -> OTP verify -> active session -> wallet -> balance -> faucet -> create intent -> confirm intent -> execute intent` in the agent API reference.
- Documented the new `POST /agent/wallets/{agent_id}/faucet` endpoint, response shape, failure codes (`404 agent_wallet_not_found`, `409 agent_wallet_status_invalid`, `409 agent_wallet_chain_invalid`, `501 circle_agent_wallet_faucet_not_configured`, `502 circle_agent_wallet_faucet_failed`), token-fixed-to-USDC behavior, ARC-TESTNET-only behavior, and that the endpoint never accepts arbitrary recipient addresses, transfers, swaps, contract executions, market creations, or mainnet funding.
- Added a judge/user testing guide that walks through opening the SignalArc GPT Agent, connecting an account, providing email, entering OTP, checking wallet, checking balance, requesting the faucet, creating a draft market intent, confirming, and executing only after explicit approval.
- Documented the maintainer-only OpenAPI import URL `https://raw.githubusercontent.com/wahyu241205/SignalArc/main/project-roadmap/signalarc-custom-gpt-openapi.json` and clarified that judges/users do not need to import it because the published GPT Agent is already wired to https://api.signalarc.fun.
- Refreshed the docs status table to show available capabilities (health, onboarding, OTP verify, session, wallet, balance, faucet, market intent lifecycle) and explicit out-of-scope capabilities (arbitrary transfer, withdraw/deposit, logout/session management, mainnet funding).
- ngrok references in docs are now explicitly marked as local development only.
- `docs/docs.json` navigation was reviewed; no structural change was required because API and AGENT_API are already linked.
- Validation: `grep -R "ngrok-free\\|undamaged-commerce\\|signalarc-backend-973633221696\\|localhost:4001" docs project-roadmap/signalarc-custom-gpt-openapi.json` returned no matches.

Not done:

- No backend, frontend, or contract logic was changed.
- No commits were created and nothing was pushed.

## Next Recommended Step

- Design and validate the Docker/Cloud Run Circle CLI/session strategy before treating the backend provider as deployable. Do not capture OTP or store Circle session material in SignalArc.
- Validate WA, Telegram, or Claude external client triggers through the existing backend-approved path when explicitly approved.
- Do not add Circle API keys, deployer/user private keys, DNS, live deployment, contract redeploy, frontend execution UI, or mainnet configuration yet.

Do not start unrelated coding before checking:

```bash
git status
```

## Important Rules

1. Frontend is UI only.
2. Backend owns business logic.
3. Database schema must be migration-based.
4. Arc and Circle behavior must use official documentation only.
5. Unknown or undocumented behavior must be marked as unknown / not documented.
6. Do not invent integrations.
7. Do not commit secrets.
8. Do not commit `node_modules`, `.next`, `dist`, `out`, or `.env`.
9. Keep commits small and reviewable.
10. Build MVP first, then polish.


## SignalArc Custom GPT / Agent Onboarding Stability Pass

Status: COMPLETE.

Done:

- Added backend `agent_id` validation helper in `backend/internal/api/agent_id_validation.go` that rejects empty values, generic placeholder values such as `signalarc-gpt-agent`, `agent_desi_001`, `default`, `defaultagent`, `default_agent`, `test`, `testagent`, `test_agent`, `demo`, `demoagent`, `demo_agent`, `user`, `useragent`, `user_agent`, `agent`, `chatgpt`, `chatgpt_agent`, and any case-insensitive equivalent, requires the `agent_<slug>` SignalArc shape, requires at least 10 characters, and limits the alphabet to ASCII letters, digits, underscore, and hyphen.
- Wired the helper into `validateAgentOnboardingSessionInput` and `validateAgentWalletRegistrationInput` in `backend/internal/api/agent_handlers.go`. `POST /agent/onboarding/start` now returns HTTP 400 `agent_onboarding_invalid` for generic/short/wrong-shape `agent_id` values and `POST /agent/wallets` returns HTTP 400 `agent_wallet_invalid` for the same.
- Added `normalizeCloseTimestamp` in `backend/internal/agent/intents.go` that accepts either a base-10 unix-seconds string (preserved unchanged for backward compatibility with the existing executor) or a UTC RFC3339 timestamp such as `2026-05-31T23:59:00Z`, normalizing the latter to unix seconds before the validator runs. The on-chain `SignalArcAgentMarketFactory.createMarket(uint256)` signature is unchanged because the executor still passes a base-10 unix-seconds string to `big.Int.SetString(value, 10)`.
- Tightened `validateIntent` so a `create_market` intent with a non-numeric `close_timestamp` returns a stable validation error such as `close_timestamp must be a unix-seconds integer or an RFC3339 timestamp such as 2026-05-31T23:59:00Z`, instead of a confusing `invalid_json`.
- Replaced the generic `invalid_json` message on `POST /agent/intents` with `request body must be valid JSON; for create_market send close_timestamp as RFC3339 (for example 2026-05-31T23:59:00Z) or unix-seconds integer string`, so non-developer Custom GPT users get a clear hint about how to convert natural-language dates.
- Hardened `project-roadmap/signalarc-custom-gpt-openapi.json` (now `0.2.0`):
  - Added top-level Custom GPT instruction language about generating a unique `agent_id` per user, never using generic placeholders, and converting natural-language dates to RFC3339 before calling create market intent.
  - Added explicit `pattern`, `minLength`, `description`, and `examples` for `agent_id` on `startAgentOnboarding` request body and on `agent_id` path parameters in `getAgentSession`, `getAgentWallet`, `getAgentWalletBalance`, and `requestAgentWalletFaucet`.
  - Added explicit `pattern` and `description` for `agent_id` on `createAgentIntent`, `confirmAgentIntent`, and `executeAgentIntent`.
  - Changed `createAgentIntent.close_timestamp` to a `string` with description and examples for both RFC3339 and unix-seconds shapes, and added inline `examples` for `create_market` and `buy_yes` payloads.
  - Documented the read-only market discovery surface as `listMarkets` (`GET /markets`), `getMarket` (`GET /markets/{market_id}`), and `listAgentMarkets` (`GET /agent/markets`) so the Custom GPT can show open/trending markets by calling an Action instead of refusing.
  - Added 201 Created responses for `startAgentOnboarding` and `createAgentIntent` matching the actual backend status codes.
  - Added `Market` and `AgentMarket` component schemas.
- Updated `docs/AGENT_API.md`:
  - Added a Custom GPT instruction block, an explicit `agent_id` naming rule list with the documented blocklist, and live recommended shapes `agent_sanatarau21_chatgpt_001` / `agent_adenhusen65_live_002`.
  - Documented the RFC3339 / unix-seconds rule for `create_market.close_timestamp` and the new stable error message.
  - Documented `GET /markets` and `GET /markets/{market_id}` as exposed read-only Actions, alongside `GET /agent/markets`.
  - Replaced the demo `agent_demo_001` example with the recommended `agent_sanatarau21_chatgpt_001` shape.
- Added Go tests in `backend/internal/api/agent_id_validation_test.go`:
  - `TestStartAgentOnboardingRejectsGenericAgentIDValues` covers each documented placeholder including case variants.
  - `TestStartAgentOnboardingRejectsShortAgentID` covers the minimum length rule.
  - `TestStartAgentOnboardingAcceptsRecommendedAgentIDShape` covers the recommended shape.
  - `TestRegisterAgentWalletRejectsGenericAgentID` covers the explicit `/agent/wallets` registration path.
  - `TestCreateAgentIntentRejectsInvalidJSON` covers the natural-language date case where the body itself is not valid JSON.
  - `TestCreateAgentIntentRejectsNonRFC3339CloseTimestamp` covers a body that decodes but contains a natural-language `close_timestamp`.
  - `TestCreateAgentIntentAcceptsRFC3339CloseTimestamp` covers the recommended RFC3339 input and verifies the response normalizes to unix seconds.
  - `TestCreateAgentIntentAcceptsUnixSecondsCloseTimestamp` covers the existing unix-seconds clients.
  - `TestValidateAgentIDDirect` table-tests the helper in isolation.

Validation results:

- `cd backend && go vet ./...` passed.
- `cd backend && go test ./...` passed (agent + api packages, including the new tests).
- `python3 -m json.tool project-roadmap/signalarc-custom-gpt-openapi.json > /tmp/signalarc-openapi.validated.json` passed.
- `grep -R 'ngrok-free\|undamaged-commerce' project-roadmap/signalarc-custom-gpt-openapi.json docs` returned no matches.

Out of scope for this pass:

- No frontend code change.
- No Solidity / contract change.
- No new arbitrary transfer, withdraw, deposit, swap, or mainnet feature.
- No commit and no push.

## Backend stopgap — Circle CLI session liveness

Problem:
- Production Cloud Run logs showed Circle CLI `AUTH_REQUIRED`.
- SignalArc DB could report an agent session as `active`, while Circle CLI runtime session was unavailable on the Cloud Run instance handling balance/execute.
- This caused misleading session status and opaque `agent_execution_failed` / `circle_agent_wallet_balance_failed` errors.

Completed:
- Added backend-only Circle CLI session liveness detection for active Circle Agent Wallet sessions.
- `GET /agent/sessions/{agent_id}` no longer blindly returns `active` when local Circle CLI session is unavailable.
- Balance and execute public error codes are preserved.
- Backend logs now include sanitized structured Circle provider failure details.

Important limitation:
- This does not fix Cloud Run Circle CLI session persistence.
- It only makes runtime session loss visible and non-misleading.
- Circle Agent Wallet HTTP API replacement remains unknown / not documented.

Validation:
- `go build ./...` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `gofmt` clean on changed Go files.

## Backend runtime hardening — Circle CLI Node engine

Problem:
- Cloud Build for image tag `948ca29` succeeded, but npm emitted `EBADENGINE` warnings.
- `@circle-fin/cli@0.0.3` requires Node `>=20.18.2`.
- The previous runtime image installed Node `v20.15.1` from `alpine:3.20`.
- Because Circle CLI is a backend runtime dependency for Agent Wallet balance/execute paths, this mismatch is not acceptable as a long-term baseline.

Completed:
- Updated backend runtime Docker stage to use pinned `node:20.18.2-alpine3.20`.
- Removed Alpine `nodejs`/`npm` package install from the runtime stage.
- Kept `ca-certificates` and global `@circle-fin/cli` install.
- No backend Go logic, frontend, OpenAPI, docs, contracts, migrations, or database behavior changed.

Validation required before deploy:
- Build backend image again.
- Confirm the previous `EBADENGINE` warning is gone.
- Deploy only the new image if build succeeds.

## Custom GPT Actions schema compatibility fix

Problem:
- Backend and Cloud Run were healthy.
- `curl https://api.signalarc.fun/health` returned HTTP 200 with `{"status":"ok"}`.
- A minimal Custom GPT Action schema containing only `GET /health` worked.
- The full SignalArc OpenAPI schema failed in GPT Actions with `ClientResponseError` and `Failed to Parse JSON`.
- Cloud Run logs showed that failing full-schema Action calls did not reach the backend.
- Root cause was therefore schema compatibility with GPT Actions parser/runtime, not backend, Cloud Run, custom domain, or deploy.

Completed:
- Updated `project-roadmap/signalarc-custom-gpt-openapi.json` for GPT Actions compatibility.
- Changed OpenAPI version from `3.1.0` to `3.0.3`.
- Removed JSON Schema / OpenAPI 3.1 type arrays such as `["string", "null"]`.
- Replaced nullable-style fields with plain string schemas where nullable behavior is not required for Action routing.
- Replaced empty array item schema `{}` with explicit object item schema.
- Replaced schema-level `examples` with OpenAPI 3.0-compatible `example`.
- Preserved all 14 operationIds and kept server URL as `https://api.signalarc.fun`.

Validation:
- JSON parses successfully with Python.
- All 14 operationIds are present and unique.
- No `type` arrays remain.
- No empty `items` schemas remain.
- No `oneOf`, `anyOf`, or `allOf` remain.

Important limitation:
- Exact GPT Actions parser feature support is unknown / not documented.
- This fix intentionally uses conservative OpenAPI 3.0.3-compatible schema constructs.

## Backend diagnostics — Circle CLI execute failures

Problem:
- Agent onboarding, OTP verification, session, wallet, balance, market listing, intent preview, and intent confirm work.
- `executeAgentIntent` still fails with `agent_execution_failed`.
- Production logs showed `operation=agent_execution`, `error_class=unknown`, and generic summary `Circle CLI command failed`.
- Existing logs were not enough to distinguish contract revert, ABI/function mismatch, gas/funds issue, Circle CLI syntax, process crash, or unsupported Circle CLI execute behavior.

Completed:
- Added diagnostics-only context for Circle CLI execute failures.
- Public API behavior is unchanged: execute failures still return HTTP 502 with public error code `agent_execution_failed`.
- Raw Circle CLI stdout/stderr is not returned to API callers.
- Structured logs now include sanitized execute context where available:
  - action
  - function signature
  - redacted contract address
  - redacted wallet address
  - chain
  - command category
  - exit status
  - raw output length
  - sanitized process error
  - sanitized stdout/stderr summary
- Added sanitizer behavior to preserve `0x`-prefixed 64-hex transaction hashes while still redacting bare long hex secrets.

Validation:
- `go build ./...` passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `gofmt` clean on changed Go files.

Important limitation:
- This does not fix execute failure.
- It only makes the next execute failure diagnosable from production logs.
- Exact Circle CLI stderr format for execute failure modes remains unknown / not documented.

## Backend guard — create market stale close timestamp

Problem:
- After execute diagnostics were deployed, `create_market` execution failed inside Circle CLI with `Transaction failed: ESTIMATION_ERROR`.
- Diagnostics showed the call reached `wallet_execute` for `createMarket(string,string,uint256,address,address)`.
- ABI/signature and argument order matched the deployed factory.
- Most likely revert cause was stale `close_timestamp`, because the contract requires close timestamp to be greater than block timestamp.

Completed:
- Added backend-only pre-execution guard for `create_market`.
- `ExecuteCreateMarket` now validates `close_timestamp` is still sufficiently in the future before calling Circle CLI.
- Uses a 60-second safety margin.
- If stale, backend does not call Circle CLI.
- Public response is now HTTP 400 with error code `create_market_close_timestamp_stale`.
- Existing provider/Circle CLI failures still use HTTP 502 with `agent_execution_failed`.
- Non-`create_market` actions are not affected.
- Tests use dynamic future/past timestamps to avoid time-bomb failures.
- Timestamp comparison uses `big.Int.Cmp` instead of `Int64()` to avoid overflow behavior for very large timestamp inputs.

Validation:
- `go test ./...` passed.
- `go vet ./...` passed.
- `gofmt` clean on changed Go files.

Important limitation:
- This does not fix every possible `ESTIMATION_ERROR`.
- Duplicate `marketId`, unsupported collateral, authorization, or other contract reverts remain possible and should be handled separately if observed.
