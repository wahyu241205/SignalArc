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
import { TradeIntentPanel } from "@/features/markets/trade-intent-panel"
import { ApiError, getMarket, type Market } from "@/lib/api"

type MarketDetailState =
  | { status: "loading" }
  | { status: "error"; message: string; requestId: string | null }
  | { status: "ready"; market: Market }

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
    default:
      return ""
  }
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">{label}</dt>
      <dd className="mt-1 text-sm font-medium text-foreground">{value}</dd>
    </div>
  )
}

function MarketDetailCard({ market }: { market: Market }) {
  return (
    <Card>
      <CardHeader className="space-y-4">
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="outline" className={statusColor(market.status)}>
            {market.status}
          </Badge>
          {market.category ? (
            <Badge variant="secondary">{market.category}</Badge>
          ) : null}
          <span className="text-xs text-muted-foreground">{market.collateral_asset} · {market.chain}</span>
        </div>
        <CardTitle className="text-xl leading-snug sm:text-2xl">{market.title}</CardTitle>
        {market.description ? (
          <CardDescription className="max-w-3xl text-sm leading-6">
            {market.description}
          </CardDescription>
        ) : null}
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Outcomes */}
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

        <Separator />

        {/* Market details */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {market.resolution_source ? (
            <DetailItem label="Resolution Source" value={market.resolution_source} />
          ) : null}
          <DetailItem label="Closes" value={formatDate(market.closes_at)} />
          {market.opens_at ? <DetailItem label="Opens" value={formatDate(market.opens_at)} /> : null}
          {market.resolved_at ? (
            <DetailItem label="Resolved" value={formatDate(market.resolved_at)} />
          ) : null}
          {market.winning_outcome ? (
            <DetailItem label="Winning Outcome" value={market.winning_outcome} />
          ) : null}
          <DetailItem label="Created" value={formatDate(market.created_at)} />
        </div>
      </CardContent>
    </Card>
  )
}

function LoadingSkeleton() {
  return (
    <Card className="animate-pulse">
      <CardHeader>
        <div className="h-6 w-1/4 rounded bg-muted" />
        <div className="mt-2 h-8 w-3/4 rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 sm:grid-cols-3">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="h-10 rounded bg-muted" />
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

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
    return (
      <Card className="border-destructive/30 bg-destructive/5">
        <CardContent className="pt-6">
          <h2 className="text-base font-medium text-destructive">Unable to load market</h2>
          <p className="mt-2 text-sm text-muted-foreground">{state.message}</p>
          {state.requestId ? (
            <p className="mt-3 font-mono text-xs text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          ) : null}
          <Button asChild className="mt-4" size="sm" variant="outline">
            <Link href="/markets">Back to markets</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="grid gap-6">
      <MarketDetailCard market={state.market} />
      <TradeIntentPanel marketId={state.market.id} marketStatus={state.market.status} />
      <MarketResolutionPanel marketId={state.market.id} />
    </div>
  )
}
