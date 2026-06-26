import { Card, CardHeader } from "@/components/ui/card"
import type { Market } from "@/lib/api"

import { MarketDetailHeader } from "./market-detail-header"
import { MarketDetailHero } from "./market-detail-hero"

export function MarketSummaryCard({ market }: { market: Market }) {
  return (
    <Card>
      <CardHeader className="space-y-4">
        <MarketDetailHero market={market} />
        <MarketDetailHeader market={market} />
      </CardHeader>
    </Card>
  )
}
