import Link from "next/link"

import { Button } from "@/components/ui/button"
import { MarketList } from "@/features/markets/market-list"

export default function MarketsPage() {
  return (
    <div className="px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-1">
            <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Markets</h1>
            <p className="text-sm text-muted-foreground">
              Discover and trade USDC-settled prediction markets on Arc.
            </p>
          </div>
          <Button asChild>
            <Link href="/markets/new">Create Market</Link>
          </Button>
        </header>
        <MarketList />
      </div>
    </div>
  )
}
