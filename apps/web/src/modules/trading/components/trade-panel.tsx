"use client"

import type { FormEvent } from "react"
import type { Address } from "viem"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import type { TradeOutcome, TradeSubmitState } from "../types"
import { TradeAmountInput } from "./trade-amount-input"
import { TradePositionCard } from "./trade-position-card"
import { TradePreviewCard } from "./trade-preview-card"
import { TradeSideSelector } from "./trade-side-selector"
import { TradeSubmitStatus } from "./trade-submit-status"
import { TradeWalletState } from "./trade-wallet-state"

export function TradePanel({
  marketId,
  contractAddress,
  walletAddress,
  outcome,
  onOutcomeChange,
  amount,
  onAmountChange,
  parsedAmount,
  yesPosition,
  noPosition,
  state,
  disabledReason,
  canSubmit,
  isPending,
  isConnected,
  isArcTestnet,
  isSwitchingChain,
  onSwitchNetwork,
  onSubmit,
}: {
  marketId: string
  contractAddress: Address | null
  walletAddress: Address | undefined
  outcome: TradeOutcome
  onOutcomeChange: (value: TradeOutcome) => void
  amount: string
  onAmountChange: (value: string) => void
  parsedAmount: bigint | null
  yesPosition: bigint | undefined
  noPosition: bigint | undefined
  state: TradeSubmitState
  disabledReason: string | null
  canSubmit: boolean
  isPending: boolean
  isConnected: boolean
  isArcTestnet: boolean
  isSwitchingChain: boolean
  onSwitchNetwork: () => void
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}) {
  return (
    <Card className="overflow-hidden">
      <CardHeader className="space-y-3">
        <div className="flex flex-wrap items-center gap-2">
          <CardTitle>Trade Ticket</CardTitle>
          <Badge
            variant="outline"
            className="border-green-500/30 bg-green-500/10 text-green-300 text-xs"
          >
            Arc Testnet
          </Badge>
        </div>
        <CardDescription>Choose a side and enter the USDC amount to trade.</CardDescription>
      </CardHeader>
      <CardContent className="px-4 pb-4 sm:px-6 sm:pb-6">
        <form className="grid gap-5" onSubmit={onSubmit}>
          <TradePreviewCard
            contractAddress={contractAddress}
            marketId={marketId}
            walletAddress={walletAddress}
            outcome={outcome}
            amount={amount}
            parsedAmount={parsedAmount}
          />

          <TradePositionCard
            walletAddress={walletAddress}
            yesPosition={yesPosition}
            noPosition={noPosition}
            isConnected={isConnected}
          />

          <div className="grid gap-4 sm:grid-cols-2">
            <TradeSideSelector value={outcome} onChange={onOutcomeChange} />
            <TradeAmountInput value={amount} onChange={onAmountChange} />
          </div>

          <TradeSubmitStatus state={state} />

          <TradeWalletState
            disabledReason={disabledReason}
            isConnected={isConnected}
            isArcTestnet={isArcTestnet}
            isSwitchingChain={isSwitchingChain}
            onSwitchNetwork={onSwitchNetwork}
          />

          <div className="grid gap-3 border-t border-border pt-4">
            <div className="flex items-start justify-between gap-4 text-sm">
              <span className="text-muted-foreground">Expected action</span>
              <span className="text-right font-medium text-foreground">
                Buy {outcome} with {amount || "0"} USDC
              </span>
            </div>
            <p className="text-xs text-muted-foreground">
              The wallet submits USDC approval first, then opens the market position.
            </p>
            <Button disabled={!canSubmit} type="submit" className="h-11 w-full">
              {state.status === "approving" ? "Approving USDC..." : null}
              {state.status === "opening" ? "Opening position..." : null}
              {!isPending ? `Trade ${outcome} on Arc Testnet` : null}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
