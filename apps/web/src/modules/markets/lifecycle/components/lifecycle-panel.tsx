"use client"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { ChainStatusCard } from "@/modules/wallet"

import type { LifecycleActionState, LifecycleStatusData } from "../types"
import { LifecycleActionStatus } from "./lifecycle-action-status"
import { LifecycleStatusCard } from "./lifecycle-status-card"

export function LifecycleNotDeployedCard() {
  return (
    <Card className="border-yellow-500/20">
      <CardHeader>
        <CardTitle>Onchain Status</CardTitle>
        <CardDescription>Onchain contract not deployed for this market.</CardDescription>
      </CardHeader>
    </Card>
  )
}

export function LifecyclePanel({
  data,
  isArcTestnet,
  isSwitchingChain,
  onSwitchNetwork,
  claimDisabledReason,
  canClaim,
  claimButtonLabel,
  resolverDisabledReason,
  resolverActionReason,
  closeDisabledReason,
  resolveDisabledReason,
  cancelDisabledReason,
  isPending,
  actionState,
  onClaim,
  onCloseMarket,
  onResolveYes,
  onResolveNo,
  onCancelMarket,
}: {
  data: LifecycleStatusData
  isArcTestnet: boolean
  isSwitchingChain: boolean
  onSwitchNetwork: () => void
  claimDisabledReason: string | null
  canClaim: boolean
  claimButtonLabel: string
  resolverDisabledReason: string | null
  resolverActionReason: string | null
  closeDisabledReason: string | null
  resolveDisabledReason: string | null
  cancelDisabledReason: string | null
  isPending: boolean
  actionState: LifecycleActionState
  onClaim: () => void
  onCloseMarket: () => void
  onResolveYes: () => void
  onResolveNo: () => void
  onCancelMarket: () => void
}) {
  return (
    <Card className="border-indigo-500/20">
      <CardHeader>
        <div className="flex flex-wrap items-center gap-2">
          <CardTitle>Onchain Status</CardTitle>
          <Badge
            variant="outline"
            className="border-indigo-500/30 bg-indigo-500/10 text-indigo-300"
          >
            Arc Testnet
          </Badge>
        </div>
        <CardDescription>
          Market-specific browser-wallet reads and lifecycle transactions. No real funds or production settlement.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        <LifecycleStatusCard data={data} />

        {!isArcTestnet ? (
          <ChainStatusCard
            isSwitchingChain={isSwitchingChain}
            onSwitchNetwork={onSwitchNetwork}
          />
        ) : null}

        <div className="grid gap-3 rounded-lg border border-border bg-muted/20 p-4">
          <h3 className="text-sm font-medium text-foreground">Claim</h3>
          {claimDisabledReason ? (
            <p className="text-sm text-muted-foreground">{claimDisabledReason}</p>
          ) : null}
          <Button
            className="w-full sm:w-fit"
            disabled={!canClaim || isPending}
            onClick={onClaim}
            type="button"
          >
            {isPending && actionState.status === "pending" && actionState.label === claimButtonLabel
              ? "Claiming..."
              : claimButtonLabel}
          </Button>
        </div>

        <div className="grid gap-3 rounded-lg border border-border bg-muted/20 p-4">
          <h3 className="text-sm font-medium text-foreground">Resolver Actions</h3>
          {resolverDisabledReason || resolverActionReason ? (
            <p className="rounded-md border border-yellow-500/30 bg-yellow-500/5 p-3 text-sm text-yellow-300">
              {resolverDisabledReason ?? resolverActionReason}
            </p>
          ) : null}
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-4">
            <Button
              disabled={Boolean(closeDisabledReason)}
              onClick={onCloseMarket}
              type="button"
              variant="outline"
            >
              Close Market
            </Button>
            <Button
              disabled={Boolean(resolveDisabledReason)}
              onClick={onResolveYes}
              type="button"
              variant="outline"
            >
              Resolve YES
            </Button>
            <Button
              disabled={Boolean(resolveDisabledReason)}
              onClick={onResolveNo}
              type="button"
              variant="outline"
            >
              Resolve NO
            </Button>
            <Button
              disabled={Boolean(cancelDisabledReason)}
              onClick={onCancelMarket}
              type="button"
              variant="destructive"
            >
              Cancel Market
            </Button>
          </div>
        </div>

        <LifecycleActionStatus actionState={actionState} />
      </CardContent>
    </Card>
  )
}
