import { PortfolioView } from "@/features/portfolio/portfolio-view"

export default function PortfolioPage() {
  return (
    <main className="min-h-screen bg-background px-4 py-8 text-foreground sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="space-y-2">
          <p className="text-sm font-medium text-muted-foreground">Portfolio</p>
          <h1 className="text-3xl font-semibold tracking-tight">User portfolio</h1>
          <p className="max-w-3xl text-sm leading-6 text-muted-foreground">
            Read-only positions and settlements from the Phase 3 backend API. This is
            not wallet balance data and does not provide a claimable settlement flow.
          </p>
        </header>
        <PortfolioView />
      </div>
    </main>
  )
}
