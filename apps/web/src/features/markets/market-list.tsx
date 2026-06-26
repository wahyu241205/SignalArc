"use client"

import Link from "next/link"
import { useEffect, useMemo, useState } from "react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { ApiError, getMarkets, type Market } from "@/lib/api"

type MarketListState =
  | { status: "loading" }
  | { status: "error"; message: string; requestId: string | null }
  | { status: "ready"; markets: Market[] }

type StatusFilter = "all" | "open" | "closed" | "resolved"

const STATUS_FILTERS: { value: StatusFilter; label: string }[] = [
  { value: "all", label: "All" },
  { value: "open", label: "Open" },
  { value: "closed", label: "Closed" },
  { value: "resolved", label: "Resolved" },
]

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
    case "cancelled":
      return "border-red-500/30 bg-red-500/10 text-red-300"
    default:
      return ""
  }
}

function deploymentColor(status: string) {
  switch (status) {
    case "DEPLOYED":
      return "border-green-500/20 bg-green-500/5 text-green-400"
    case "FAILED":
      return "border-red-500/20 bg-red-500/5 text-red-400"
    default:
      return "border-border bg-muted/30 text-muted-foreground"
  }
}

function deploymentLabel(status: string) {
  switch (status) {
    case "DEPLOYED":
      return "Onchain"
    case "FAILED":
      return "Deploy Failed"
    default:
      return "Not Deployed"
  }
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
      />
    </svg>
  )
}

function ChartIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
      />
    </svg>
  )
}

function AlertIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
      />
    </svg>
  )
}

function FilterIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z"
      />
    </svg>
  )
}

/* -------------------------------------------------------------------------- */
/*  Market card                                                               */
/* -------------------------------------------------------------------------- */

function MarketCard({ market }: { market: Market }) {
  return (
    <Card className="group transition-colors hover:border-indigo-500/30">
      <div className="px-6 pt-6">
        {market.cover_image_url ? (
          // Plain img is intentional for v1 user-provided remote URLs.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={market.cover_image_url}
            alt={market.title}
            className="h-48 w-full rounded-xl object-cover"
            loading="lazy"
          />
        ) : (
          <div className="h-48 w-full rounded-xl bg-muted" aria-hidden="true" />
        )}
      </div>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0 space-y-2.5">
            <CardTitle className="text-base font-semibold leading-snug">
              <Link
                className="transition-colors hover:text-indigo-300"
                href={`/markets/${market.id}`}
              >
                {market.title}
              </Link>
            </CardTitle>

            <CardDescription className="flex flex-wrap items-center gap-1.5">
              <Badge
                variant="outline"
                className={statusColor(market.status)}
              >
                {market.status}
              </Badge>

              <Badge
                variant="outline"
                className="border-muted-foreground/20 bg-muted/40 text-muted-foreground"
              >
                {market.category || "Uncategorized"}
              </Badge>

              <Badge
                variant="outline"
                className={`text-xs ${deploymentColor(market.onchain_deployment_status)}`}
              >
                {deploymentLabel(market.onchain_deployment_status)}
              </Badge>
            </CardDescription>
          </div>

          <Button
            asChild
            size="sm"
            variant="outline"
            className="shrink-0 opacity-0 transition-opacity group-hover:opacity-100"
          >
            <Link href={`/markets/${market.id}`}>View</Link>
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        <dl className="grid gap-4 text-sm text-muted-foreground sm:grid-cols-3">
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Collateral
            </dt>
            <dd className="font-medium text-foreground">
              {market.collateral_asset}
            </dd>
          </div>
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Chain
            </dt>
            <dd className="font-medium text-foreground">{market.chain}</dd>
          </div>
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Closes
            </dt>
            <dd className="font-medium text-foreground">
              {formatDate(market.closes_at)}
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  )
}

/* -------------------------------------------------------------------------- */
/*  Loading skeleton                                                          */
/* -------------------------------------------------------------------------- */

function LoadingSkeleton() {
  return (
    <div className="space-y-6">
      {/* Toolbar skeleton */}
      <div className="space-y-3">
        <div className="h-8 w-full max-w-sm animate-pulse rounded-lg bg-muted" />
        <div className="flex gap-2">
          {[1, 2, 3, 4].map((i) => (
            <div
              key={i}
              className="h-8 w-16 animate-pulse rounded-md bg-muted"
            />
          ))}
        </div>
        <div className="h-4 w-64 animate-pulse rounded bg-muted" />
      </div>

      {/* Card skeletons */}
      <div className="grid gap-4">
        {[1, 2, 3, 4].map((i) => (
          <Card key={i} className="animate-pulse">
            <div className="px-6 pt-6">
              <div className="h-48 w-full rounded-xl bg-muted" />
            </div>
            <CardHeader className="pb-3">
              <div className="space-y-2.5">
                <div className="h-5 w-3/4 rounded bg-muted" />
                <div className="flex gap-1.5">
                  <div className="h-5 w-14 rounded-full bg-muted" />
                  <div className="h-5 w-20 rounded-full bg-muted" />
                  <div className="h-5 w-20 rounded-full bg-muted" />
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-3">
                {[1, 2, 3].map((j) => (
                  <div key={j} className="space-y-1.5">
                    <div className="h-3 w-16 rounded bg-muted/70" />
                    <div className="h-4 w-24 rounded bg-muted" />
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}

/* -------------------------------------------------------------------------- */
/*  Filter toolbar                                                            */
/* -------------------------------------------------------------------------- */

function FilterToolbar({
  searchQuery,
  onSearchChange,
  statusFilter,
  onStatusFilterChange,
  markets,
  filteredCount,
}: {
  searchQuery: string
  onSearchChange: (value: string) => void
  statusFilter: StatusFilter
  onStatusFilterChange: (value: StatusFilter) => void
  markets: Market[]
  filteredCount: number
}) {
  const counts = useMemo(() => {
    const open = markets.filter(
      (m) => m.status.toLowerCase() === "open",
    ).length
    const closed = markets.filter(
      (m) => m.status.toLowerCase() === "closed",
    ).length
    const resolved = markets.filter(
      (m) => m.status.toLowerCase() === "resolved",
    ).length

    return { total: markets.length, open, closed, resolved }
  }, [markets])

  return (
    <div className="space-y-3">
      {/* Search input */}
      <div className="relative max-w-sm">
        <SearchIcon className="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="text"
          placeholder="Search markets by title…"
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          className="pl-8"
        />
      </div>

      {/* Status filter tabs */}
      <div className="flex items-center gap-1.5">
        <FilterIcon className="mr-1 h-3.5 w-3.5 text-muted-foreground" />
        {STATUS_FILTERS.map((filter) => (
          <Button
            key={filter.value}
            size="sm"
            variant={statusFilter === filter.value ? "default" : "outline"}
            className={
              statusFilter === filter.value
                ? "h-7 px-3 text-xs"
                : "h-7 px-3 text-xs text-muted-foreground"
            }
            onClick={() => onStatusFilterChange(filter.value)}
          >
            {filter.label}
          </Button>
        ))}
      </div>

      {/* Summary bar */}
      <p className="text-xs font-medium text-muted-foreground">
        {filteredCount === counts.total ? (
          <>
            {counts.total} {counts.total === 1 ? "market" : "markets"}
            {" · "}
            {counts.open} open · {counts.closed} closed · {counts.resolved}{" "}
            resolved
          </>
        ) : (
          <>
            Showing {filteredCount} of {counts.total}{" "}
            {counts.total === 1 ? "market" : "markets"}
          </>
        )}
      </p>
    </div>
  )
}

/* -------------------------------------------------------------------------- */
/*  Main component                                                            */
/* -------------------------------------------------------------------------- */

export function MarketList() {
  const [state, setState] = useState<MarketListState>({ status: "loading" })
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("all")

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

  const filteredMarkets = useMemo(() => {
    if (state.status !== "ready") return []

    let result = state.markets

    if (statusFilter !== "all") {
      result = result.filter(
        (m) => m.status.toLowerCase() === statusFilter,
      )
    }

    const query = searchQuery.trim().toLowerCase()
    if (query) {
      result = result.filter((m) =>
        m.title.toLowerCase().includes(query),
      )
    }

    return result
  }, [state, searchQuery, statusFilter])

  /* Loading ---------------------------------------------------------------- */

  if (state.status === "loading") {
    return <LoadingSkeleton />
  }

  /* Error ------------------------------------------------------------------ */

  if (state.status === "error") {
    return (
      <Card className="border-destructive/30 bg-destructive/5">
        <CardContent className="flex flex-col items-center gap-4 py-10 text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
            <AlertIcon className="h-6 w-6 text-destructive" />
          </div>

          <div className="space-y-1">
            <h2 className="text-base font-medium text-destructive">
              Unable to load markets
            </h2>
            <p className="text-sm text-muted-foreground">{state.message}</p>
          </div>

          {state.requestId ? (
            <p className="font-mono text-xs text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          ) : null}

          <Button
            variant="outline"
            size="sm"
            onClick={() => window.location.reload()}
          >
            Retry
          </Button>
        </CardContent>
      </Card>
    )
  }

  /* Empty (no markets at all) --------------------------------------------- */

  if (state.markets.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center gap-5 py-16 text-center">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-muted">
            <ChartIcon className="h-7 w-7 text-muted-foreground" />
          </div>

          <div className="max-w-sm space-y-1">
            <p className="font-medium">No markets available</p>
            <p className="text-sm text-muted-foreground">
              Markets will appear here once they are created through the API
              or creator dashboard. Use the platform API or web interface to
              launch a new event market.
            </p>
          </div>

          <Button asChild>
            <Link href="/markets/new">Create Market</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  /* Ready ------------------------------------------------------------------ */

  return (
    <div className="space-y-6">
      <FilterToolbar
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        statusFilter={statusFilter}
        onStatusFilterChange={setStatusFilter}
        markets={state.markets}
        filteredCount={filteredMarkets.length}
      />

      {filteredMarkets.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center gap-4 py-12 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
              <SearchIcon className="h-6 w-6 text-muted-foreground" />
            </div>

            <div className="max-w-xs space-y-1">
              <p className="font-medium">No matching markets</p>
              <p className="text-sm text-muted-foreground">
                No markets match the current search or filter criteria. Adjust
                your filters to see results.
              </p>
            </div>

            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setSearchQuery("")
                setStatusFilter("all")
              }}
            >
              Clear Filters
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {filteredMarkets.map((market) => (
            <MarketCard key={market.id} market={market} />
          ))}
        </div>
      )}
    </div>
  )
}
