import Link from "next/link"

import { Button } from "@/components/ui/button"
import { CreateMarketForm } from "@/features/markets/create-market-form"

export default function NewMarketPage() {
  return (
    <main className="min-h-screen bg-background px-4 py-8 text-foreground sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <div>
          <Button asChild size="sm" variant="ghost">
            <Link href="/markets">Back to markets</Link>
          </Button>
        </div>
        <CreateMarketForm />
      </div>
    </main>
  )
}
