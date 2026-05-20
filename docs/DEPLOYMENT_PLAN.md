# Deployment Plan

This document describes planned deployment only. It does not claim that SignalArc is live.

## Current Deployment Status

| Item | Status |
| --- | --- |
| Production frontend | Not live. |
| Production backend | Not live. |
| Production database | Not configured. |
| DNS | Not configured. |
| Production approval | Not completed. |
| Arc Testnet contract | Deployed prototype reference exists. |

## Planned Domains

| Domain | Target |
| --- | --- |
| https://signalarc.fun | Frontend application. |
| https://api.signalarc.fun | Backend API. |
| https://docs.signalarc.fun | Documentation. |

## Planned Frontend

- Platform: Vercel.
- Primary domain: `signalarc.fun`.
- Optional domain: `www.signalarc.fun`.
- Build command should use the monorepo frontend build path.
- Frontend env vars must be configured in Vercel settings.

## Planned Backend

- Platform: GCP Cloud Run.
- Domain: `api.signalarc.fun`.
- Runtime: Docker container built from the Go backend.
- Backend must connect to a managed PostgreSQL database.
- Production CORS must allow the production frontend origin.

## Planned Database

- GCP Cloud SQL PostgreSQL or another hosted PostgreSQL provider.
- Apply migrations through the approved migration workflow.
- Do not manually edit migration catalog state unless no migration CLI path is available and the exception is documented.

## Planned DNS

| Record | Target |
| --- | --- |
| `signalarc.fun` | Vercel frontend target. |
| `www.signalarc.fun` | Optional Vercel frontend alias. |
| `api.signalarc.fun` | Cloud Run custom domain or load balancer target. |
| `docs.signalarc.fun` | Documentation hosting target. |

## Planned Environment Variables

Frontend public variables:

- `NEXT_PUBLIC_API_BASE_URL=https://api.signalarc.fun`
- `NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID`

Backend private variables:

- `APP_ENV`
- `APP_PORT`
- `DATABASE_URL`
- Future Circle variables only if Circle integration is implemented and approved.
- Future RPC or service credentials only through secret management.

Do not include actual production credentials in source control or documentation.

## Pre-Deploy Checklist

- Local Docker backend health passes.
- Local backend readiness passes.
- Local schema validation passes.
- Frontend lint passes.
- Frontend type check passes.
- Frontend production build passes.
- Wallet flow is tested on Arc Testnet.
- Onchain testnet transaction is tested if transaction UI remains enabled.
- No secrets are present in git.
- Production database is created and migrations are applied.
- Production CORS is configured for `https://signalarc.fun`.
- DNS targets are known.
- Rollback plan is documented.
- Monitoring/logging baseline is defined.

## Deployment Commands

Deployment commands are intentionally omitted because this plan must not deploy now.
