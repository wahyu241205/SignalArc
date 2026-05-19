# SignalArc Project State

This file is the handoff source of truth for continuing SignalArc work in a new chat or coding session.

## Project

- Name: SignalArc
- Repository: https://github.com/wahyu241205/SignalArc
- Positioning: Arc-native prediction market infrastructure, not a Polymarket clone
- Core idea: Convert market opinions into structured probability signals, support USDC-settled event markets, expose market intelligence through APIs, and make the system usable by creators, institutions, developers, and AI agents.

## Working Style

- Use Indonesian for guidance.
- Use controlled debugging.
- One step at a time.
- Do not jump phases without justification.
- Do not provide multiple unrelated solutions at once.
- For Arc and Circle behavior, use official documentation only.
- If behavior is not documented, state: unknown / not documented.
- Do not invent integrations or system behavior.
- Keep commits small and reviewable.

## Current Overall Status

Foundation complete. Backend and database baseline started. Product MVP is not live yet.

Estimated progress toward live grant submission: 20-25%.

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

Start Phase 4 Frontend MVP by wiring the frontend to the Phase 3 backend API endpoints.

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

Status: Not started as product, about 5%.

Done:

- Next.js app exists.
- Frontend moved to `apps/web`.
- shadcn/ui base components exist.
- Frontend runs locally.

Not done:

- SignalArc landing page.
- Market list page.
- Market detail page.
- Create market form.
- YES/NO trade panel UI.
- User portfolio page.
- Resolver dashboard.
- Intelligence dashboard.
- Wallet connect UI integration.
- API client layer.
- Frontend env config.
- Frontend loading/error states.
- Responsive layout polish.
- Demo-ready UI.

## Phase 5 — Arc / Circle Integration

Status: Not started, 0%.

Done:

- Official documentation sources identified:
  - Arc docs: https://docs.arc.network/
  - Arc prediction market blueprint: https://www.arc.network/blog/build-institutional-grade-prediction-markets-on-arc-arc-blueprints
  - Circle developer docs: https://developers.circle.com/
  - Circle grants page: https://www.circle.com/grant
- Project rule added: Arc/Circle behavior must use official docs only.

Not done:

- Confirm current Arc testnet/mainnet setup from official docs.
- Confirm Circle wallet approach from official docs.
- Confirm USDC flow from Circle docs.
- Confirm supported Arc RPC / chain details.
- Configure wallet strategy.
- Integrate Circle Developer Platform.
- Integrate Arc transaction flow.
- Add webhook handling.
- Add testnet transaction proof.
- Document unknown / not documented behavior.

## Phase 6 — Contracts / Settlement Prototype

Status: Not started, 0%.

Done:

- Stack decision: Solidity + Foundry + OpenZeppelin.

Not done:

- Create `contracts` folder.
- Initialize Foundry.
- Define market contract scope.
- Define settlement contract scope.
- Define custody/collateral assumptions.
- Write minimal contract prototype.
- Add unit tests.
- Deploy to testnet if supported.
- Connect backend to contract flow.
- Document security assumptions.
- Avoid pretending production custody is solved before verified.

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

## Current Last Completed Step

- Phase 2 Core Database Schema completed.

## Next Recommended Step

- Start Phase 4 Frontend MVP by wiring the frontend to the Phase 3 backend API endpoints.

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
