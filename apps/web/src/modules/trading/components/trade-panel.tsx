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
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <CardTitle>Place Trade</CardTitle>
          <Badge
            variant="outline"
            className="border-green-500/30 bg-green-500/10 text-green-300 text-xs"
          >
            Arc Testnet
          </Badge>
        </div>
        <CardDescription>
          This executes on Arc Testnet from your connected wallet. No production settlement.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="grid gap-5" onSubmit={onSubmit}>
          <TradePreviewCard
            contractAddress={contractAddress}
            marketId={marketId}
            walletAddress={walletAddress}
          />

          <div className="grid gap-4 sm:grid-cols-2">
            <TradeSideSelector value={outcome} onChange={onOutcomeChange} />
            <TradeAmountInput value={amount} onChange={onAmountChange} />
          </div>

          <p className="text-xs text-muted-foreground">
            {amount || "0"} USDC will be sent as{" "}
            {parsedAmount ? parsedAmount.toString() : "0"} base units.
          </p>

          <TradeSubmitStatus state={state} />

          <TradeWalletState
            disabledReason={disabledReason}
            isConnected={isConnected}
            isArcTestnet={isArcTestnet}
            isSwitchingChain={isSwitchingChain}
            onSwitchNetwork={onSwitchNetwork}
          />

          <Button disabled={!canSubmit} type="submit" className="w-full sm:w-auto">
            {state.status === "approving" ? "Approving USDC..." : null}
            {state.status === "opening" ? "Opening position..." : null}
            {!isPending ? `Trade ${outcome} on Arc Testnet` : null}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
