"use client"

import { Button } from "@/components/ui/button"
import { getWalletConnectionMessage } from "@/modules/wallet"

export function TradeWalletState({
  disabledReason,
  isConnected,
  isArcTestnet,
  isSwitchingChain,
  onSwitchNetwork,
}: {
  disabledReason: string | null
  isConnected: boolean
  isArcTestnet: boolean
  isSwitchingChain: boolean
  onSwitchNetwork: () => void
}) {
  if (!disabledReason) return null

  return (
    <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-4">
      <p className="text-sm font-medium text-yellow-300">{disabledReason}</p>
      {!isConnected ? (
        <p className="mt-1 text-sm text-muted-foreground">
          {getWalletConnectionMessage(isConnected)}
        </p>
      ) : null}
      {isConnected && !isArcTestnet ? (
        <Button
          className="mt-3"
          disabled={isSwitchingChain}
          onClick={onSwitchNetwork}
          size="sm"
          type="button"
          variant="outline"
        >
          {isSwitchingChain ? "Switching..." : "Switch network"}
        </Button>
      ) : null}
    </div>
  )
}
