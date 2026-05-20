import { IntelligenceDashboard } from "@/features/intelligence/intelligence-dashboard"

export default function IntelligencePage() {
  return (
    <div className="px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="space-y-1">
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Market Intelligence</h1>
          <p className="text-sm text-muted-foreground">
            Real-time probability signals and structured market data for agents and developers.
          </p>
        </header>
        <IntelligenceDashboard />
      </div>
    </div>
  )
}
