"use client"

import { type FormEvent, useState } from "react"
import { useAccount } from "wagmi"

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
  localDemoUserId,
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

function truncateId(id: string) {
  if (id.length <= 12) return id
  return `${id.slice(0, 6)}…${id.slice(-4)}`
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
    return <p className="text-sm text-muted-foreground">No open positions.</p>
  }

  return (
    <div className="overflow-x-auto rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Market</TableHead>
            <TableHead>Outcome</TableHead>
            <TableHead>Quantity</TableHead>
            <TableHead>Avg Entry</TableHead>
            <TableHead>Realized PnL</TableHead>
            <TableHead>Updated</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {positions.map((position) => (
            <TableRow key={position.id}>
              <TableCell className="font-mono text-xs" title={position.market_id}>
                {truncateId(position.market_id)}
              </TableCell>
              <TableCell>{position.outcome}</TableCell>
              <TableCell>{position.quantity}</TableCell>
              <TableCell>{position.average_entry_price}</TableCell>
              <TableCell>{position.realized_pnl}</TableCell>
              <TableCell className="text-muted-foreground">{formatDate(position.updated_at)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

function SettlementsTable({ settlements }: { settlements: Settlement[] }) {
  if (settlements.length === 0) {
    return <p className="text-sm text-muted-foreground">No settlements yet.</p>
  }

  return (
    <div className="overflow-x-auto rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Market</TableHead>
            <TableHead>Outcome</TableHead>
            <TableHead>Amount</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Tx Hash</TableHead>
            <TableHead>Settled</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {settlements.map((settlement) => (
            <TableRow key={settlement.id}>
              <TableCell className="font-mono text-xs" title={settlement.market_id}>
                {truncateId(settlement.market_id)}
              </TableCell>
              <TableCell>{settlement.outcome ?? "-"}</TableCell>
              <TableCell>{settlement.amount}</TableCell>
              <TableCell>{settlement.status}</TableCell>
              <TableCell className="font-mono text-xs">{settlement.tx_hash ? truncateId(settlement.tx_hash) : "-"}</TableCell>
              <TableCell className="text-muted-foreground">{formatDate(settlement.settled_at)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

function PortfolioTables({ data }: { data: PortfolioData }) {
  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Positions</CardTitle>
          <CardDescription>
            Your current market positions.
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
            Settlement history for resolved markets.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SettlementsTable settlements={data.settlements} />
        </CardContent>
      </Card>
    </div>
  )
}

function WalletIdentityCard({ address }: { address: string }) {
  return (
    <Card className="border-indigo-500/20">
      <CardContent className="flex items-center gap-3 pt-6">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-indigo-500/10">
          <svg className="h-5 w-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 9m18 0V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v3" />
          </svg>
        </div>
        <div>
          <p className="text-sm font-medium text-foreground">Connected Wallet</p>
          <p className="font-mono text-xs text-muted-foreground">{address}</p>
        </div>
      </CardContent>
    </Card>
  )
}

export function PortfolioView() {
  const { address, isConnected } = useAccount()
  const [state, setState] = useState<PortfolioState>({ status: "idle" })
  const [showAdvanced, setShowAdvanced] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const userId = String(formData.get("user_id") ?? "").trim()

    if (!userId) {
      setState({
        status: "error",
        userId: "",
        message: "User ID is required.",
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
      {/* Wallet identity display */}
      {isConnected && address ? (
        <WalletIdentityCard address={address} />
      ) : (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
              <svg className="h-6 w-6 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 9m18 0V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v3" />
              </svg>
            </div>
            <p className="text-sm text-muted-foreground">
              Connect your wallet to view your portfolio.
            </p>
            <p className="text-xs text-muted-foreground/70">
              Positions are read-only until wallet-to-user mapping is implemented.
            </p>
          </CardContent>
        </Card>
      )}

      {/* Advanced: manual user ID lookup */}
      <div>
        <button
          type="button"
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="text-xs font-medium text-muted-foreground hover:text-foreground transition-colors"
        >
          {showAdvanced ? "▾ Hide" : "▸ Show"} demo lookup (advanced)
        </button>

        {showAdvanced ? (
          <Card className="mt-3">
            <CardHeader>
              <CardTitle className="text-sm">Demo Lookup</CardTitle>
              <CardDescription className="text-xs">
                Look up positions by backend user ID. This is a local demo fallback.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form className="flex flex-col gap-3 sm:flex-row sm:items-end" onSubmit={handleSubmit}>
                <div className="grid flex-1 gap-2">
                  <Label htmlFor="user_id" className="text-xs">User ID</Label>
                  <Input id="user_id" name="user_id" defaultValue={localDemoUserId} placeholder="UUID" className="text-sm" />
                </div>
                <Button disabled={state.status === "loading"} type="submit" size="sm">
                  {state.status === "loading" ? "Loading..." : "Load"}
                </Button>
              </form>
            </CardContent>
          </Card>
        ) : null}
      </div>

      {state.status === "idle" && !showAdvanced ? null : null}

      {state.status === "loading" ? (
        <Card className="animate-pulse">
          <CardContent className="py-8">
            <div className="h-4 w-1/2 rounded bg-muted" />
          </CardContent>
        </Card>
      ) : null}

      {state.status === "empty" ? (
        <Card>
          <CardContent className="py-8 text-center">
            <p className="text-sm text-muted-foreground">
              No positions or settlements found.
            </p>
          </CardContent>
        </Card>
      ) : null}

      {state.status === "error" ? (
        <Card className="border-destructive/30 bg-destructive/5">
          <CardContent className="pt-6">
            <h2 className="text-base font-medium text-destructive">
              Unable to load portfolio
            </h2>
            <p className="mt-2 text-sm text-muted-foreground">{state.message}</p>
            {state.requestId ? (
              <p className="mt-3 font-mono text-xs text-muted-foreground">
                Request ID: {state.requestId}
              </p>
            ) : null}
          </CardContent>
        </Card>
      ) : null}

      {state.status === "loaded" ? <PortfolioTables data={state.data} /> : null}
    </div>
  )
}
