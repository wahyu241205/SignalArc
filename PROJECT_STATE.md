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

## Current Last Completed Step

- Local frontend/backend MVP integration fix and usability smoke test completed.

## Live AI Agent Transaction MVP

Status: IN PROGRESS - boundary/discovery only; blocked before live onchain validation.

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
- Backend now registers agent wallets through DB-backed `POST /agent/wallets` in production routing.
- Backend now returns registered wallet metadata through `GET /agent/wallets/{agent_id}`.
- Backend now disables registered wallets through `POST /agent/wallets/{agent_id}/disable`.
- Test handlers use an in-memory registry only as a test double; production registration is no longer the temporary in-memory development path.
- Backend agent intents now carry `agent_id` and optional `agent_wallet_address`.
- Backend agent intents now carry `source_client` and optional `client_request_id` so WhatsApp, Telegram, ChatGPT, Claude, and web clients can share one channel-agnostic API contract.
- Backend execution now rejects missing agent wallets, deployer/resolver wallet reuse, user-wallet reuse, disallowed actions, wrong chain, inactive wallet, mismatched wallet address, and unsupported provider execution.
- No secrets, `.env` files, frontend code, production deployment config, commits, pushes, or deploys were changed.

Current non-claims:

- No live external ChatGPT/Claude/Telegram client has triggered the backend.
- Payout lifecycle from the Circle Agent Wallet is not complete yet.
- Cancel/refund lifecycle from the Circle Agent Wallet is not complete yet.
- Backend Circle Agent Wallet execution automation is not implemented yet; backend remains fail-closed for `circle_agent_wallet` execution until safe server-side Circle auth/session handling is designed from official docs.
- External AI client integrations are not implemented yet.
- Circle CLI command shapes in `project-roadmap/agent-mcp.md` are official-doc/help-discovery shapes only unless accompanied by exact authenticated output and onchain evidence.
- SignalArc does not claim production policy-limited execution on ARC-TESTNET; Circle CLI help says `wallet limit` is mainnet only.
- The phase is not complete.

## Next Recommended Step

- Design the safe backend Circle Agent Wallet execution provider from official Circle documentation without relying on user-local CLI sessions, OTP capture, private keys, or deployer wallet reuse.
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
