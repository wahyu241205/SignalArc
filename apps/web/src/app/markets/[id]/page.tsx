import Link from "next/link"

import { Button } from "@/components/ui/button"
import { MarketDetail } from "@/features/markets/market-detail"

type MarketDetailPageProps = {
  params: Promise<{
    id: string
  }>
}

export default async function MarketDetailPage({ params }: MarketDetailPageProps) {
  const { id } = await params

  return (
    <div className="px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
        <div>
          <Button asChild size="sm" variant="ghost">
            <Link href="/markets">&larr; Back to markets</Link>
          </Button>
        </div>
        <MarketDetail marketId={id} />
      </div>
    </div>
  )
}
