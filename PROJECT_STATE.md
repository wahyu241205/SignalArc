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

Status: In progress, about 40%.

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
- `schema_migrations`

Current local database migration status:

- version: `11`
- dirty: `false`

Not done:

- `webhooks` table.
- `agent_access` table.
- Final schema review.
- Rollback/down migration test.
- Seed data for local demo.

Next step:

Push api_keys migration commit.

Then continue Phase 2 by creating webhooks table migration.

## Phase 3 — Backend API MVP

Status: Started, about 10%.

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
- Verified:

```bash
curl http://127.0.0.1:4000/health
```

Expected result:

```json
{"status":"ok"}
```

Current endpoint:

- `GET /health`

Not done:

- Backend config package.
- Database connection package.
- Readiness endpoint.
- User repository.
- Wallet repository.
- Market repository.
- Position repository.
- Trade repository.
- Market creation endpoint.
- Market listing endpoint.
- Market detail endpoint.
- Trade intent endpoint.
- Position endpoint.
- Resolver endpoint.
- Settlement status endpoint.
- Agent API endpoint.
- API error model.
- Request validation.
- Structured logging middleware.
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

- Added and applied api_keys table migration locally.

## Next Recommended Step

- Push api_keys migration commit.
- Then continue Phase 2 by creating webhooks table migration.

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
