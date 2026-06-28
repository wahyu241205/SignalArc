import { AlertCircle, CheckCircle2, Loader2 } from "lucide-react"

import { TransactionLink } from "@/modules/wallet"

import type { TradeSubmitState } from "../types"

function TradeExecutionResult({
  state,
}: {
  state: Extract<TradeSubmitState, { status: "success" }>
}) {
  return (
    <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4">
      <div className="flex items-center gap-2">
        <CheckCircle2 className="h-4 w-4 text-green-400" aria-hidden="true" />
        <p className="text-sm font-medium text-green-300">
          Position opened on Arc Testnet
        </p>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">
        Your connected wallet opened this position on Arc Testnet. Track the
        transactions below before refreshing balances.
      </p>
      <dl className="mt-4 grid gap-3 text-sm sm:grid-cols-2">
        <div>
          <dt className="text-xs uppercase tracking-wider text-muted-foreground/70">
            USDC Approval
          </dt>
          <dd className="mt-0.5">
            <TransactionLink hash={state.approveHash} />
          </dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wider text-muted-foreground/70">
            Open Position
          </dt>
          <dd className="mt-0.5">
            <TransactionLink hash={state.openHash} />
          </dd>
        </div>
      </dl>
    </div>
  )
}

function PendingState({
  state,
}: {
  state: Extract<TradeSubmitState, { status: "approving" | "opening" }>
}) {
  const label = state.status === "approving" ? "Approving USDC" : "Opening position"

  return (
    <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
      <div className="flex items-center gap-2">
        <Loader2 className="h-4 w-4 animate-spin text-blue-300" aria-hidden="true" />
        <p className="text-sm font-medium text-blue-200">{label}</p>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">
        Keep this page open. The ticket will update after wallet confirmation and
        Arc Testnet receipt finality.
      </p>
      {state.approveHash ? (
        <p className="mt-3 text-sm text-muted-foreground">
          Approval: <TransactionLink hash={state.approveHash} />
        </p>
      ) : null}
      {state.status === "opening" && state.openHash ? (
        <p className="mt-2 text-sm text-muted-foreground">
          Market transaction: <TransactionLink hash={state.openHash} />
        </p>
      ) : null}
    </div>
  )
}

export function TradeSubmitStatus({ state }: { state: TradeSubmitState }) {
  return (
    <>
      {state.status === "error" ? (
        <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
          <div className="flex items-center gap-2">
            <AlertCircle className="h-4 w-4 text-destructive" aria-hidden="true" />
            <p className="text-sm font-medium text-destructive">
              Unable to execute trade
            </p>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">{state.message}</p>
          <p className="mt-2 text-xs text-muted-foreground">
            No new position is recorded unless the market transaction confirms.
          </p>
          {state.approveHash ? (
            <p className="mt-2 text-sm text-muted-foreground">
              Approval: <TransactionLink hash={state.approveHash} />
            </p>
          ) : null}
          {state.openHash ? (
            <p className="mt-2 text-sm text-muted-foreground">
              Market transaction: <TransactionLink hash={state.openHash} />
            </p>
          ) : null}
        </div>
      ) : null}

      {state.status === "approving" || state.status === "opening" ? (
        <PendingState state={state} />
      ) : null}

      {state.status === "success" ? <TradeExecutionResult state={state} /> : null}
    </>
  )
}
