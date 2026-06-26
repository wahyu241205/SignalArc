import type { ReactNode } from "react"

export function ResolutionTransactionCard({
  children,
}: {
  children: ReactNode
}) {
  return (
    <div className="rounded-lg border border-border bg-muted/20 p-4">
      {children}
    </div>
  )
}
