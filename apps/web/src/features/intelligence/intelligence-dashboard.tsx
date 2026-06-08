"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { ApiError, getAgentMarkets, type AgentMarket } from "@/lib/api"

type IntelligenceState =
  | { status: "loading" }
  | { status: "empty" }
  | { status: "error"; message: string; requestId: string | null }
  | { status: "loaded"; markets: AgentMarket[] }

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

function getErrorState(error: unknown): Extract<IntelligenceState, { status: "error" }> {
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
    message: "Unable to load market signals.",
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

function IntelligenceIntroPanel({ marketCount }: { marketCount: number }) {
  return (
    <Card className="border-indigo-500/10 bg-gradient-to-r from-indigo-500/5 via-transparent to-purple-500/5">
      <CardContent className="py-5">
        <div className="grid gap-6 sm:grid-cols-[1fr_auto]">
          <div className="space-y-2">
            <h2 className="text-sm font-semibold text-foreground">Market Intelligence API</h2>
            <p className="text-xs leading-relaxed text-muted-foreground">
              Structured, API-accessible probability signals derived from prediction market activity.
              Available for programmatic consumption through the{" "}
              <a
                href="https://docs.signalarc.fun/AGENT_API"
                target="_blank"
                rel="noopener noreferrer"
                className="text-indigo-400 hover:text-indigo-300 transition-colors"
              >
                Agent API
              </a>.
              Signals include market status, category, collateral, resolution source, and close timestamps.
            </p>
          </div>
          <div className="flex items-center gap-4">
            <div className="rounded-lg border border-border/50 bg-card/60 px-4 py-2 text-center">
              <p className="text-lg font-bold text-indigo-400">{marketCount}</p>
              <p className="text-[10px] font-medium uppercase tracking-widest text-muted-foreground">{marketCount === 1 ? "Signal" : "Signals"}</p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function SignalCard({ market }: { market: AgentMarket }) {
  return (
    <div className="group rounded-lg border border-border/50 bg-card/40 transition-colors hover:border-indigo-500/20">
      <div className="flex items-start justify-between gap-3 p-4">
        <div className="min-w-0 space-y-1.5">
          <Link className="text-sm font-semibold text-foreground hover:text-indigo-300 transition-colors" href={`/markets/${market.id}`}>
            {market.title}
          </Link>
          <div className="flex flex-wrap items-center gap-1.5">
            <Badge variant="outline" className={`text-[10px] ${statusColor(market.status)}`}>
              {market.status}
            </Badge>
            <Badge variant="outline" className="text-[10px] border-border bg-muted/30 text-muted-foreground">
              {market.category || "Uncategorized"}
            </Badge>
          </div>
        </div>
        <Button asChild size="sm" variant="ghost" className="shrink-0 text-xs opacity-0 transition-opacity group-hover:opacity-100">
          <Link href={`/markets/${market.id}`}>View →</Link>
        </Button>
      </div>
      <Separator className="opacity-30" />
      <dl className="grid gap-x-4 gap-y-1 px-4 py-3 text-xs sm:grid-cols-4">
        <div>
          <dt className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">Collateral</dt>
          <dd className="font-medium text-foreground">{market.collateral_asset}</dd>
        </div>
        <div>
          <dt className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">Chain</dt>
          <dd className="font-medium text-foreground">{market.chain}</dd>
        </div>
        <div>
          <dt className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">Closes</dt>
          <dd className="font-medium text-foreground">{formatDate(market.closes_at)}</dd>
        </div>
        {market.resolution_source ? (
          <div>
            <dt className="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">Resolution</dt>
            <dd className="font-medium text-foreground truncate" title={market.resolution_source}>{market.resolution_source}</dd>
          </div>
        ) : null}
      </dl>
    </div>
  )
}

function LoadingSkeleton() {
  return (
    <div className="grid gap-3">
      <div className="h-24 animate-pulse rounded-lg border border-border/50 bg-muted/20" />
      {[1, 2, 3].map((i) => (
        <div key={i} className="animate-pulse rounded-lg border border-border/30 p-4">
          <div className="h-4 w-2/3 rounded bg-muted" />
          <div className="mt-2 flex gap-2">
            <div className="h-4 w-12 rounded bg-muted" />
            <div className="h-4 w-16 rounded bg-muted" />
          </div>
          <div className="mt-3 grid grid-cols-4 gap-4">
            <div className="h-3 rounded bg-muted/60" />
            <div className="h-3 rounded bg-muted/60" />
            <div className="h-3 rounded bg-muted/60" />
            <div className="h-3 rounded bg-muted/60" />
          </div>
        </div>
      ))}
    </div>
  )
}

export function IntelligenceDashboard() {
  const [state, setState] = useState<IntelligenceState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadMarkets() {
      setState({ status: "loading" })

      try {
        const response = await getAgentMarkets()

        if (!isActive) {
          return
        }

        if (response.data.markets.length === 0) {
          setState({ status: "empty" })
          return
        }

        setState({ status: "loaded", markets: response.data.markets })
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

  if (state.status === "empty") {
    return (
      <div className="grid gap-4">
        <IntelligenceIntroPanel marketCount={0} />
        <Card className="border-border/30">
          <CardContent className="flex flex-col items-center gap-4 py-12 text-center">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
              <svg className="h-5 w-5 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" />
              </svg>
            </div>
            <div>
              <p className="text-sm font-medium">No market signals available</p>
              <p className="mt-1 text-xs text-muted-foreground">Signals will appear here as markets become active on the platform.</p>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (state.status === "error") {
    return (
      <Card className="border-destructive/30 bg-destructive/5">
        <CardContent className="pt-6">
          <h2 className="text-sm font-semibold text-destructive">Unable to load market signals</h2>
          <p className="mt-2 text-xs text-muted-foreground">{state.message}</p>
          {state.requestId ? (
            <p className="mt-3 font-mono text-[10px] text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          ) : null}
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="grid gap-4">
      <IntelligenceIntroPanel marketCount={state.markets.length} />
      {state.markets.map((market) => (
        <SignalCard key={market.id} market={market} />
      ))}
    </div>
  )
}
