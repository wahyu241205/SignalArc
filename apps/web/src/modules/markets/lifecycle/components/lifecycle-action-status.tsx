import { AlertCircle, CheckCircle2, Loader2 } from "lucide-react"

import type { LifecycleActionState } from "../types"
import { LifecycleTransactionCard } from "./lifecycle-transaction-card"

export function LifecycleActionStatus({
  actionState,
}: {
  actionState: LifecycleActionState
}) {
  if (actionState.status === "pending") {
    return (
      <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
        <div className="flex items-center gap-2">
          <Loader2 className="h-4 w-4 animate-spin text-blue-300" aria-hidden="true" />
          <p className="text-sm font-medium text-blue-200">
            {actionState.label} pending
          </p>
        </div>
        {actionState.hash ? (
          <p className="mt-3 text-sm text-muted-foreground">
            <LifecycleTransactionCard hash={actionState.hash} />
          </p>
        ) : null}
      </div>
    )
  }

  if (actionState.status === "success") {
    return (
      <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4">
        <div className="flex items-center gap-2">
          <CheckCircle2 className="h-4 w-4 text-green-400" aria-hidden="true" />
          <p className="text-sm font-medium text-green-300">
            {actionState.label} confirmed on Arc Testnet
          </p>
        </div>
        <p className="mt-2 text-sm text-muted-foreground">
          <LifecycleTransactionCard hash={actionState.hash} />
        </p>
      </div>
    )
  }

  if (actionState.status === "error") {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
        <div className="flex items-center gap-2">
          <AlertCircle className="h-4 w-4 text-destructive" aria-hidden="true" />
          <p className="text-sm font-medium text-destructive">
            Unable to run {actionState.label}
          </p>
        </div>
        <p className="mt-1 text-sm text-muted-foreground">
          {actionState.message}
        </p>
        {actionState.hash ? (
          <p className="mt-2 text-sm text-muted-foreground">
            <LifecycleTransactionCard hash={actionState.hash} />
          </p>
        ) : null}
      </div>
    )
  }

  return null
}
