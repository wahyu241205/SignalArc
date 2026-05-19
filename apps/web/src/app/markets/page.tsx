import Link from "next/link"

import { Button } from "@/components/ui/button"
import { MarketList } from "@/features/markets/market-list"

export default function MarketsPage() {
  return (
    <main className="min-h-screen bg-background px-4 py-8 text-foreground sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
        <header className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div className="space-y-2">
            <p className="text-sm font-medium text-muted-foreground">Markets</p>
            <h1 className="text-3xl font-semibold tracking-tight">SignalArc markets</h1>
            <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
              Read-only market data from the Phase 3 backend API.
            </p>
          </div>
          <Button asChild>
            <Link href="/markets/new">Create market</Link>
          </Button>
        </header>
        <MarketList />
      </div>
    </main>
  )
}
