"use client"

import { useEffect, useState } from "react"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  ApiError,
  getMarketResolution,
  getMarketSettlements,
  type Resolution,
  type Settlement,
} from "@/lib/api"

type ResolutionState =
  | { status: "loading" }
  | { status: "empty"; settlements: Settlement[] }
  | { status: "loaded"; resolution: Resolution; settlements: Settlement[] }
  | { status: "error"; message: string; requestId: string | null }

function formatDate(value: string | null) {
  if (!value) {
    return "-"
  }

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date)
}

async function loadResolution(marketId: string) {
  try {
    const response = await getMarketResolution(marketId)
    return response.data.resolution
  } catch (error) {
    if (error instanceof ApiError && error.code === "resolution_not_found") {
      return null
    }

    throw error
  }
}

function getErrorState(error: unknown): Extract<ResolutionState, { status: "error" }> {
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
    message: "Unable to load resolution state.",
    requestId: null,
  }
}

function ResolutionDetails({ resolution }: { resolution: Resolution }) {
  return (
    <dl className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Status</dt>
        <dd className="mt-1 text-sm text-foreground">{resolution.status}</dd>
      </div>
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Winning outcome</dt>
        <dd className="mt-1 text-sm text-foreground">{resolution.winning_outcome ?? "-"}</dd>
      </div>
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Resolver type</dt>
        <dd className="mt-1 text-sm text-foreground">{resolution.resolver_type ?? "-"}</dd>
      </div>
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Evidence reference</dt>
        <dd className="mt-1 text-sm text-foreground">{resolution.evidence_reference ?? "-"}</dd>
      </div>
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Resolved</dt>
        <dd className="mt-1 text-sm text-foreground">{formatDate(resolution.resolved_at)}</dd>
      </div>
      <div>
        <dt className="text-sm font-medium text-muted-foreground">Updated</dt>
        <dd className="mt-1 text-sm text-foreground">{formatDate(resolution.updated_at)}</dd>
      </div>
    </dl>
  )
}

function SettlementsTable({ settlements }: { settlements: Settlement[] }) {
  if (settlements.length === 0) {
    return <p className="text-sm text-muted-foreground">No settlement rows returned.</p>
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>User</TableHead>
          <TableHead>Outcome</TableHead>
          <TableHead>Amount</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Tx hash</TableHead>
          <TableHead>Settled</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {settlements.map((settlement) => (
          <TableRow key={settlement.id}>
            <TableCell className="font-mono text-xs">{settlement.user_id ?? "-"}</TableCell>
            <TableCell>{settlement.outcome ?? "-"}</TableCell>
            <TableCell>{settlement.amount}</TableCell>
            <TableCell>{settlement.status}</TableCell>
            <TableCell className="font-mono text-xs">{settlement.tx_hash ?? "-"}</TableCell>
            <TableCell>{formatDate(settlement.settled_at)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function MarketResolutionPanel({ marketId }: { marketId: string }) {
  const [state, setState] = useState<ResolutionState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadState() {
      setState({ status: "loading" })

      try {
        const [resolution, settlementsResponse] = await Promise.all([
          loadResolution(marketId),
          getMarketSettlements(marketId),
        ])
        const settlements = settlementsResponse.data.settlements

        if (!isActive) {
          return
        }

        if (!resolution) {
          setState({ status: "empty", settlements })
          return
        }

        setState({ status: "loaded", resolution, settlements })
      } catch (error) {
        if (isActive) {
          setState(getErrorState(error))
        }
      }
    }

    void loadState()

    return () => {
      isActive = false
    }
  }, [marketId])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Resolution and settlements</CardTitle>
        <CardDescription>
          Read-only backend state. This does not submit resolver evidence, execute
          settlement, claim funds, or infer eligibility.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        {state.status === "loading" ? (
          <p className="text-sm text-muted-foreground">Loading resolution state...</p>
        ) : null}

        {state.status === "error" ? (
          <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
            <p className="text-sm font-medium text-destructive">
              Unable to load resolution state
            </p>
            <p className="mt-1 text-sm text-muted-foreground">{state.message}</p>
            {state.requestId ? (
              <p className="mt-2 font-mono text-xs text-muted-foreground">
                Request ID: {state.requestId}
              </p>
            ) : null}
          </div>
        ) : null}

        {state.status === "empty" ? (
          <div className="grid gap-6">
            <p className="text-sm text-muted-foreground">
              No resolution row was returned for this market.
            </p>
            <SettlementsTable settlements={state.settlements} />
          </div>
        ) : null}

        {state.status === "loaded" ? (
          <div className="grid gap-6">
            <ResolutionDetails resolution={state.resolution} />
            <SettlementsTable settlements={state.settlements} />
          </div>
        ) : null}
      </CardContent>
    </Card>
  )
}
