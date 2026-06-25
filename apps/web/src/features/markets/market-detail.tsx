"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { MarketResolutionPanel } from "@/features/markets/market-resolution-panel"
import { OnchainMarketLifecyclePanel } from "@/features/markets/onchain-market-lifecycle-panel"
import { TradeIntentPanel } from "@/features/markets/trade-intent-panel"
import { ApiError, getMarket, type Market } from "@/lib/api"

/* ---------------------------------------------------------------------------
 * State
 * --------------------------------------------------------------------------- */

type MarketDetailState =
  | { status: "loading" }
  | { status: "error"; message: string; requestId: string | null }
  | { status: "ready"; market: Market }

/* ---------------------------------------------------------------------------
 * Helpers
 * --------------------------------------------------------------------------- */

function formatDate(value: string) {
  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date)
}

function getErrorState(error: unknown): Extract<MarketDetailState, { status: "error" }> {
  if (error instanceof ApiError) {
    return {
      status: "error",
      message: error.message,
      requestId: error.requestId,
    }
  }

  if (error instanceof Error) {
    return {
      status: "error",
      message: error.message,
      requestId: null,
    }
  }

  return {
    status: "error",
    message: "Unable to load market.",
    requestId: null,
  }
}

function statusColor(status: string) {
  switch (status.toLowerCase()) {
    case "open":
      return "border-green-500/30 bg-green-500/10 text-green-300"
    case "closed":
      return "border-yellow-500/30 bg-yellow-500/10 text-yellow-300"
    case "resolved":
      return "border-indigo-500/30 bg-indigo-500/10 text-indigo-300"
    case "cancelled":
      return "border-red-500/30 bg-red-500/10 text-red-300"
    default:
      return ""
  }
}

function statusContext(status: string, winningOutcome: string | null) {
  switch (status.toLowerCase()) {
    case "open":
      return "This market is currently accepting positions."
    case "closed":
      return "This market has closed and is pending resolution."
    case "resolved":
      return winningOutcome
        ? `This market has been resolved. Winning outcome: ${winningOutcome}.`
        : "This market has been resolved."
    case "cancelled":
      return "This market has been cancelled. Participants may be eligible for refunds."
    default:
      return null
  }
}

function arcscanContractUrl(address: string) {
  return `https://testnet.arcscan.app/address/${address}`
}

/* ---------------------------------------------------------------------------
 * Detail Item
 * --------------------------------------------------------------------------- */

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">{label}</dt>
      <dd className="mt-1 text-sm font-medium text-foreground">{value}</dd>
    </div>
  )
}

/* ---------------------------------------------------------------------------
 * Section Header
 * --------------------------------------------------------------------------- */

function SectionHeader({ children }: { children: React.ReactNode }) {
  return (
    <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
      {children}
    </h3>
  )
}

/* ---------------------------------------------------------------------------
 * Share Market Button
 * --------------------------------------------------------------------------- */

function ShareMarketButton({ market }: { market: Market }) {
  const [status, setStatus] = useState<"idle" | "copied" | "shared">("idle")

  async function handleShare() {
    const url = window.location.href
    const text =
      market.description ||
      `Trade ${market.outcome_yes_label} or ${market.outcome_no_label} on SignalArc.`

    try {
      if (navigator.share) {
        await navigator.share({
          title: market.title,
          text,
          url,
        })
        setStatus("shared")
      } else {
        await navigator.clipboard.writeText(url)
        setStatus("copied")
      }

      window.setTimeout(() => setStatus("idle"), 2000)
    } catch (error) {
      if (error instanceof DOMException && error.name === "AbortError") {
        return
      }

      try {
        await navigator.clipboard.writeText(url)
        setStatus("copied")
        window.setTimeout(() => setStatus("idle"), 2000)
      } catch {
        setStatus("idle")
      }
    }
  }

  return (
    <Button type="button" variant="outline" size="sm" onClick={handleShare}>
      {status === "copied"
        ? "Link copied"
        : status === "shared"
          ? "Shared"
          : "Share Market"}
    </Button>
  )
}

/* ---------------------------------------------------------------------------
 * Market Summary Card (main column)
 * --------------------------------------------------------------------------- */

function MarketSummaryCard({ market }: { market: Market }) {
  const context = statusContext(market.status, market.winning_outcome)

  return (
    <Card>
      <CardHeader className="space-y-4">
        {market.cover_image_url ? (
          // Plain img is intentional for v1 user-provided remote URLs.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={market.cover_image_url}
            alt={market.title}
            className="h-48 w-full rounded-xl object-cover sm:h-64"
          />
        ) : (
          <div className="h-48 w-full rounded-xl bg-muted sm:h-64" aria-hidden="true" />
        )}

        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant="outline" className={statusColor(market.status)}>
              {market.status}
            </Badge>
            <Badge variant="secondary">
              {market.category || "Uncategorized"}
            </Badge>
            <span className="text-xs text-muted-foreground">
              {market.collateral_asset} · {market.chain}
            </span>
          </div>
          <ShareMarketButton market={market} />
        </div>

        <CardTitle className="text-xl leading-snug sm:text-2xl">{market.title}</CardTitle>

        {market.description ? (
          <CardDescription className="max-w-3xl text-sm leading-6">
            {market.description}
          </CardDescription>
        ) : null}

        {context ? (
          <p className="text-sm text-muted-foreground">{context}</p>
        ) : null}
      </CardHeader>
    </Card>
  )
}

/* ---------------------------------------------------------------------------
 * Outcomes Section (main column)
 * --------------------------------------------------------------------------- */

function OutcomesSection({ market }: { market: Market }) {
  return (
    <Card>
      <CardHeader>
        <SectionHeader>Outcomes &amp; Probability</SectionHeader>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-3">
          <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4 text-center">
            <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">YES</p>
            <p className="mt-1 text-lg font-bold text-green-400">{market.outcome_yes_label}</p>
          </div>
          <div className="rounded-lg border border-red-500/20 bg-red-500/5 p-4 text-center">
            <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">NO</p>
            <p className="mt-1 text-lg font-bold text-red-400">{market.outcome_no_label}</p>
          </div>
        </div>
        <p className="text-xs text-muted-foreground">
          Probability signal will be derived from market position data when available.
        </p>
      </CardContent>
    </Card>
  )
}

/* ---------------------------------------------------------------------------
 * Market Parameters Section (main column)
 * --------------------------------------------------------------------------- */

function MarketParametersSection({ market }: { market: Market }) {
  return (
    <Card>
      <CardHeader>
        <SectionHeader>Market Parameters</SectionHeader>
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <DetailItem label="Closes" value={formatDate(market.closes_at)} />
          {market.opens_at ? <DetailItem label="Opens" value={formatDate(market.opens_at)} /> : null}
          {market.resolved_at ? (
            <DetailItem label="Resolved" value={formatDate(market.resolved_at)} />
          ) : null}
          {market.winning_outcome ? (
            <DetailItem label="Winning Outcome" value={market.winning_outcome} />
          ) : null}
          <DetailItem label="Created" value={formatDate(market.created_at)} />
          <DetailItem label="Onchain Status" value={market.onchain_deployment_status} />
          {market.market_contract_address ? (
            <div>
              <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Market Contract</dt>
              <dd className="mt-1 text-sm font-medium">
                <a
                  href={arcscanContractUrl(market.market_contract_address)}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="break-all font-mono text-indigo-400 transition-colors hover:text-indigo-300"
                >
                  {market.market_contract_address}
                </a>
              </dd>
            </div>
          ) : null}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------------------------------------------------------------------------
 * Resolution Section (main column)
 * --------------------------------------------------------------------------- */

function ResolutionSection({ market }: { market: Market }) {
  if (!market.resolution_source) {
    return null
  }

  return (
    <Card>
      <CardHeader>
        <SectionHeader>Resolution</SectionHeader>
      </CardHeader>
      <CardContent>
        <DetailItem label="Resolution Source" value={market.resolution_source} />
      </CardContent>
    </Card>
  )
}

/* ---------------------------------------------------------------------------
 * Loading Skeleton
 * --------------------------------------------------------------------------- */

function SkeletonBlock({ className }: { className?: string }) {
  return <div className={`rounded bg-muted ${className ?? ""}`} />
}

function LoadingSkeleton() {
  return (
    <div className="grid gap-6 lg:grid-cols-3">
      {/* Main column skeleton */}
      <div className="space-y-6 lg:col-span-2">
        {/* Summary card skeleton */}
        <Card className="animate-pulse">
          <CardHeader className="space-y-4">
            <SkeletonBlock className="h-48 w-full rounded-xl sm:h-64" />
            <div className="flex items-center gap-2">
              <SkeletonBlock className="h-5 w-16" />
              <SkeletonBlock className="h-5 w-24" />
            </div>
            <SkeletonBlock className="h-8 w-3/4" />
            <SkeletonBlock className="h-4 w-full max-w-md" />
          </CardHeader>
        </Card>

        {/* Outcomes skeleton */}
        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-4 w-40" />
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-3">
              <SkeletonBlock className="h-20" />
              <SkeletonBlock className="h-20" />
            </div>
          </CardContent>
        </Card>

        {/* Parameters skeleton */}
        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-4 w-36" />
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {Array.from({ length: 6 }).map((_, i) => (
                <SkeletonBlock key={i} className="h-10" />
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Side column skeleton */}
      <div className="space-y-6 lg:col-span-1">
        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-5 w-32" />
          </CardHeader>
          <CardContent className="space-y-3">
            <SkeletonBlock className="h-10" />
            <SkeletonBlock className="h-10" />
            <SkeletonBlock className="h-10" />
          </CardContent>
        </Card>

        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-5 w-44" />
          </CardHeader>
          <CardContent>
            <SkeletonBlock className="h-24" />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

/* ---------------------------------------------------------------------------
 * Error State
 * --------------------------------------------------------------------------- */

function ErrorState({
  message,
  requestId,
}: {
  message: string
  requestId: string | null
}) {
  return (
    <div className="mx-auto max-w-lg py-12">
      <Card className="border-destructive/30 bg-destructive/5">
        <CardHeader className="space-y-2 text-center">
          <CardTitle className="text-lg text-destructive">Unable to load market</CardTitle>
          <CardDescription className="text-sm text-muted-foreground">
            {message}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-center">
          {requestId ? (
            <p className="font-mono text-xs text-muted-foreground">
              Request ID: {requestId}
            </p>
          ) : null}
          <Separator />
          <div className="flex items-center justify-center gap-3">
            <Button asChild size="sm" variant="outline">
              <Link href="/markets">Back to markets</Link>
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => window.location.reload()}
            >
              Retry
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

/* ---------------------------------------------------------------------------
 * Main Component
 * --------------------------------------------------------------------------- */

export function MarketDetail({ marketId }: { marketId: string }) {
  const [state, setState] = useState<MarketDetailState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadMarket() {
      setState({ status: "loading" })

      try {
        const response = await getMarket(marketId)

        if (isActive) {
          setState({ status: "ready", market: response.data.market })
        }
      } catch (error) {
        if (isActive) {
          setState(getErrorState(error))
        }
      }
    }

    void loadMarket()

    return () => {
      isActive = false
    }
  }, [marketId])

  if (state.status === "loading") {
    return <LoadingSkeleton />
  }

  if (state.status === "error") {
    return <ErrorState message={state.message} requestId={state.requestId} />
  }

  const { market } = state

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      {/* ---- Main Column ---- */}
      <div className="space-y-6 lg:col-span-2">
        <MarketSummaryCard market={market} />
        <OutcomesSection market={market} />
        <MarketParametersSection market={market} />
        <ResolutionSection market={market} />
      </div>

      {/* ---- Side Column ---- */}
      <div className="space-y-6 lg:col-span-1">
        <TradeIntentPanel
          marketId={market.id}
          marketStatus={market.status}
          marketContractAddress={market.market_contract_address}
        />
        <OnchainMarketLifecyclePanel marketContractAddress={market.market_contract_address} />
        <MarketResolutionPanel marketId={market.id} />
      </div>
    </div>
  )
}
