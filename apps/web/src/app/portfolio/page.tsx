import { PortfolioView } from "@/features/portfolio/portfolio-view"

export default function PortfolioPage() {
  return (
    <div className="px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="space-y-1">
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Portfolio</h1>
          <p className="text-sm text-muted-foreground">
            Track your positions and settlement history across all markets.
          </p>
        </header>
        <PortfolioView />
      </div>
    </div>
  )
}
