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

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-sm font-medium text-muted-foreground">{label}</dt>
      <dd className="mt-1 text-sm text-foreground">{value}</dd>
    </div>
  )
}

function MarketDetailCard({ market }: { market: Market }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-2xl">{market.title}</CardTitle>
        <CardDescription className="flex flex-wrap items-center gap-2">
          <Badge variant="outline">{market.status}</Badge>
          {market.category ? <span>{market.category}</span> : null}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {market.description ? (
          <p className="max-w-3xl text-sm leading-6 text-muted-foreground">
            {market.description}
          </p>
        ) : null}

        <dl className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <DetailItem label="YES label" value={market.outcome_yes_label} />
          <DetailItem label="NO label" value={market.outcome_no_label} />
          <DetailItem label="Collateral" value={market.collateral_asset} />
          <DetailItem label="Chain" value={market.chain} />
          {market.resolution_source ? (
            <DetailItem label="Resolution source" value={market.resolution_source} />
          ) : null}
          {market.opens_at ? <DetailItem label="Opens" value={formatDate(market.opens_at)} /> : null}
          <DetailItem label="Closes" value={formatDate(market.closes_at)} />
          {market.resolved_at ? (
            <DetailItem label="Resolved" value={formatDate(market.resolved_at)} />
          ) : null}
          {market.settled_at ? (
            <DetailItem label="Settled" value={formatDate(market.settled_at)} />
          ) : null}
          {market.winning_outcome ? (
            <DetailItem label="Winning outcome" value={market.winning_outcome} />
          ) : null}
          <DetailItem label="Created" value={formatDate(market.created_at)} />
          <DetailItem label="Updated" value={formatDate(market.updated_at)} />
        </dl>
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
    return (
      <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
        Loading market...
      </div>
    )
  }

  if (state.status === "error") {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-6">
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
      </div>
    )
  }

  return (
    <div className="grid gap-6">
      <MarketDetailCard market={state.market} />
      <TradeIntentPanel marketId={state.market.id} />
      <MarketResolutionPanel marketId={state.market.id} />
    </div>
  )
}
