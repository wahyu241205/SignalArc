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
import { ApiError, getMarkets, type Market } from "@/lib/api"

type MarketListState =
  | { status: "loading" }
  | { status: "error"; message: string; requestId: string | null }
  | { status: "ready"; markets: Market[] }

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

function getErrorState(error: unknown): Extract<MarketListState, { status: "error" }> {
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
    message: "Unable to load markets.",
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

function MarketCard({ market }: { market: Market }) {
  return (
    <Card className="group transition-colors hover:border-indigo-500/30">
      <CardHeader>
        <div className="flex items-start justify-between gap-4">
          <div className="space-y-1.5">
            <CardTitle className="text-base leading-snug">
              <Link className="hover:text-indigo-300 transition-colors" href={`/markets/${market.id}`}>
                {market.title}
              </Link>
            </CardTitle>
            <CardDescription className="flex flex-wrap items-center gap-2">
              <Badge variant="outline" className={statusColor(market.status)}>
                {market.status}
              </Badge>
              {market.category ? (
                <span className="text-xs text-muted-foreground">{market.category}</span>
              ) : null}
            </CardDescription>
          </div>
          <Button asChild size="sm" variant="outline" className="shrink-0 opacity-0 transition-opacity group-hover:opacity-100">
            <Link href={`/markets/${market.id}`}>View</Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <dl className="grid gap-3 text-sm text-muted-foreground sm:grid-cols-3">
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Collateral</dt>
            <dd className="mt-0.5 font-medium text-foreground">{market.collateral_asset}</dd>
          </div>
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Chain</dt>
            <dd className="mt-0.5 font-medium text-foreground">{market.chain}</dd>
          </div>
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Closes</dt>
            <dd className="mt-0.5 font-medium text-foreground">{formatDate(market.closes_at)}</dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  )
}

function LoadingSkeleton() {
  return (
    <div className="grid gap-4">
      {[1, 2, 3].map((i) => (
        <Card key={i} className="animate-pulse">
          <CardHeader>
            <div className="h-5 w-3/4 rounded bg-muted" />
            <div className="mt-2 h-4 w-1/4 rounded bg-muted" />
          </CardHeader>
          <CardContent>
            <div className="grid gap-3 sm:grid-cols-3">
              <div className="h-4 rounded bg-muted" />
              <div className="h-4 rounded bg-muted" />
              <div className="h-4 rounded bg-muted" />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

export function MarketList() {
  const [state, setState] = useState<MarketListState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadMarkets() {
      setState({ status: "loading" })

      try {
        const response = await getMarkets()

        if (isActive) {
          setState({ status: "ready", markets: response.data.markets })
        }
      } catch (error) {
        if (isActive) {
          setState(getErrorState(error))
        }
      }
    }

    void loadMarkets()

    return () => {
      isActive = false
    }
  }, [])

  if (state.status === "loading") {
    return <LoadingSkeleton />
  }

  if (state.status === "error") {
    return (
      <Card className="border-destructive/30 bg-destructive/5">
        <CardContent className="pt-6">
          <h2 className="text-base font-medium text-destructive">Unable to load markets</h2>
          <p className="mt-2 text-sm text-muted-foreground">{state.message}</p>
          {state.requestId ? (
            <p className="mt-3 font-mono text-xs text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          ) : null}
        </CardContent>
      </Card>
    )
  }

  if (state.markets.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center gap-4 py-12 text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
            <svg className="h-6 w-6 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75z" />
            </svg>
          </div>
          <div>
            <p className="font-medium">No markets yet</p>
            <p className="mt-1 text-sm text-muted-foreground">Create the first prediction market to get started.</p>
          </div>
          <Button asChild>
            <Link href="/markets/new">Create Market</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="grid gap-4">
      {state.markets.map((market) => (
        <MarketCard key={market.id} market={market} />
      ))}
    </div>
  )
}
