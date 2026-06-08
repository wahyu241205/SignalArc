import Link from "next/link"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { MarketList } from "@/features/markets/market-list"

export default function MarketsPage() {
  return (
    <div className="px-4 py-6 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <h1 className="text-xl font-bold tracking-tight sm:text-2xl">Markets</h1>
              <Badge variant="outline" className="border-indigo-500/20 bg-indigo-500/5 text-indigo-300 text-[10px]">
                Arc Testnet
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              Browse USDC-settled prediction markets. Filter by status and category.
            </p>
          </div>
          <Button asChild size="sm">
            <Link href="/markets/new">Create Market</Link>
          </Button>
        </header>
        <MarketList />
      </div>
    </div>
  )
}
