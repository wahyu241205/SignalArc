"use client"

import Link from "next/link"
import { useEffect, useState } from "react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
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
    message: "Unable to load agent-readable markets.",
    requestId: null,
  }
}

function AgentMarketCard({ market }: { market: AgentMarket }) {
  return (
    <Card size="sm">
      <CardHeader>
        <CardTitle>
          <Link className="hover:underline" href={`/markets/${market.id}`}>
            {market.title}
          </Link>
        </CardTitle>
        <CardDescription className="flex flex-wrap items-center gap-2">
          <Badge variant="outline">{market.status}</Badge>
          {market.category ? <span>{market.category}</span> : null}
        </CardDescription>
        <CardAction>
          <Button asChild size="sm" variant="outline">
            <Link href={`/markets/${market.id}`}>View market</Link>
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <dl className="grid gap-3 text-sm text-muted-foreground sm:grid-cols-4">
          <div>
            <dt className="font-medium text-foreground">ID</dt>
            <dd className="font-mono text-xs">{market.id}</dd>
          </div>
          <div>
            <dt className="font-medium text-foreground">Collateral</dt>
            <dd>{market.collateral_asset}</dd>
          </div>
          <div>
            <dt className="font-medium text-foreground">Chain</dt>
            <dd>{market.chain}</dd>
          </div>
          <div>
            <dt className="font-medium text-foreground">Closes</dt>
            <dd>{formatDate(market.closes_at)}</dd>
          </div>
          {market.resolution_source ? (
            <div className="sm:col-span-4">
              <dt className="font-medium text-foreground">Resolution source</dt>
              <dd>{market.resolution_source}</dd>
            </div>
          ) : null}
        </dl>
      </CardContent>
    </Card>
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
    return (
      <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
        Loading agent-readable markets...
      </div>
    )
  }

  if (state.status === "empty") {
    return (
      <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
        No agent-readable markets were returned.
      </div>
    )
  }

  if (state.status === "error") {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-6">
        <h2 className="text-base font-medium text-destructive">
          Unable to load intelligence data
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">{state.message}</p>
        {state.requestId ? (
          <p className="mt-3 font-mono text-xs text-muted-foreground">
            Request ID: {state.requestId}
          </p>
        ) : null}
      </div>
    )
  }

  return (
    <div className="grid gap-4">
      {state.markets.map((market) => (
        <AgentMarketCard key={market.id} market={market} />
      ))}
    </div>
  )
}
