"use client"

import { type FormEvent, useState } from "react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
  getUserPositions,
  getUserSettlements,
  type Position,
  type Settlement,
} from "@/lib/api"

type PortfolioData = {
  positions: Position[]
  settlements: Settlement[]
}

type PortfolioState =
  | { status: "idle" }
  | { status: "loading"; userId: string }
  | { status: "empty"; userId: string }
  | { status: "loaded"; userId: string; data: PortfolioData }
  | { status: "error"; userId: string; message: string; requestId: string | null }

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

function getErrorState(
  userId: string,
  error: unknown,
): Extract<PortfolioState, { status: "error" }> {
  if (error instanceof ApiError) {
    return {
      status: "error",
      userId,
      message: error.message,
      requestId: error.requestId,
    }
  }

  if (error instanceof Error) {
    return {
      status: "error",
      userId,
      message: error.message,
      requestId: null,
    }
  }

  return {
    status: "error",
    userId,
    message: "Unable to load portfolio data.",
    requestId: null,
  }
}

function PositionsTable({ positions }: { positions: Position[] }) {
  if (positions.length === 0) {
    return <p className="text-sm text-muted-foreground">No positions returned.</p>
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Market</TableHead>
          <TableHead>Outcome</TableHead>
          <TableHead>Quantity</TableHead>
          <TableHead>Average entry price</TableHead>
          <TableHead>Realized PnL</TableHead>
          <TableHead>Updated</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {positions.map((position) => (
          <TableRow key={position.id}>
            <TableCell className="font-mono text-xs">{position.market_id}</TableCell>
            <TableCell>{position.outcome}</TableCell>
            <TableCell>{position.quantity}</TableCell>
            <TableCell>{position.average_entry_price}</TableCell>
            <TableCell>{position.realized_pnl}</TableCell>
            <TableCell>{formatDate(position.updated_at)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

function SettlementsTable({ settlements }: { settlements: Settlement[] }) {
  if (settlements.length === 0) {
    return <p className="text-sm text-muted-foreground">No settlements returned.</p>
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Market</TableHead>
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
            <TableCell className="font-mono text-xs">{settlement.market_id}</TableCell>
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

function PortfolioTables({ data }: { data: PortfolioData }) {
  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Positions</CardTitle>
          <CardDescription>
            Read-only position rows returned by the backend.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <PositionsTable positions={data.positions} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Settlements</CardTitle>
          <CardDescription>
            Read-only settlement rows. This UI does not claim or settle funds.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SettlementsTable settlements={data.settlements} />
        </CardContent>
      </Card>
    </div>
  )
}

export function PortfolioView() {
  const [state, setState] = useState<PortfolioState>({ status: "idle" })

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const userId = String(formData.get("user_id") ?? "").trim()

    if (!userId) {
      setState({
        status: "error",
        userId: "",
        message: "user_id is required.",
        requestId: null,
      })
      return
    }

    setState({ status: "loading", userId })

    try {
      const [positionsResponse, settlementsResponse] = await Promise.all([
        getUserPositions(userId),
        getUserSettlements(userId),
      ])
      const data = {
        positions: positionsResponse.data.positions,
        settlements: settlementsResponse.data.settlements,
      }

      if (data.positions.length === 0 && data.settlements.length === 0) {
        setState({ status: "empty", userId })
        return
      }

      setState({ status: "loaded", userId, data })
    } catch (error) {
      setState(getErrorState(userId, error))
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Lookup by user ID</CardTitle>
          <CardDescription>
            Auth and wallet identity are not implemented yet, so enter a backend user_id manually.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="flex flex-col gap-3 sm:flex-row sm:items-end" onSubmit={handleSubmit}>
            <div className="grid flex-1 gap-2">
              <Label htmlFor="user_id">User ID</Label>
              <Input id="user_id" name="user_id" placeholder="UUID" />
            </div>
            <Button disabled={state.status === "loading"} type="submit">
              {state.status === "loading" ? "Loading..." : "Load portfolio"}
            </Button>
          </form>
        </CardContent>
      </Card>

      {state.status === "idle" ? (
        <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
          Enter a user ID to load read-only positions and settlements.
        </div>
      ) : null}

      {state.status === "loading" ? (
        <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
          Loading portfolio data for {state.userId}...
        </div>
      ) : null}

      {state.status === "empty" ? (
        <div className="rounded-lg border bg-card p-6 text-sm text-muted-foreground">
          No positions or settlements were returned for {state.userId}.
        </div>
      ) : null}

      {state.status === "error" ? (
        <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-6">
          <h2 className="text-base font-medium text-destructive">
            Unable to load portfolio data
          </h2>
          <p className="mt-2 text-sm text-muted-foreground">{state.message}</p>
          {state.requestId ? (
            <p className="mt-3 font-mono text-xs text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          ) : null}
        </div>
      ) : null}

      {state.status === "loaded" ? <PortfolioTables data={state.data} /> : null}
    </div>
  )
}
