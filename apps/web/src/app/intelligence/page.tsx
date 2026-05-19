import { IntelligenceDashboard } from "@/features/intelligence/intelligence-dashboard"

export default function IntelligencePage() {
  return (
    <main className="min-h-screen bg-background px-4 py-8 text-foreground sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <header className="space-y-2">
          <p className="text-sm font-medium text-muted-foreground">Intelligence</p>
          <h1 className="text-3xl font-semibold tracking-tight">Agent-readable markets</h1>
          <p className="max-w-3xl text-sm leading-6 text-muted-foreground">
            Read-only market data from the Phase 3 backend agent endpoint. This is
            not live AI-agent execution, paid access, API key enforcement, or trading automation.
          </p>
        </header>
        <IntelligenceDashboard />
      </div>
    </main>
  )
}
