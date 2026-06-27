import type { Address } from "viem"

export function TradePreviewCard({
  contractAddress,
  marketId,
  walletAddress,
}: {
  contractAddress: Address | null
  marketId: string
  walletAddress: Address | undefined
}) {
  return (
    <div className="grid gap-3 rounded-lg border border-border bg-muted/20 p-4 text-sm">
      <div className="grid gap-1">
        <p className="text-xs font-medium uppercase text-muted-foreground">Contract</p>
        <p className="break-all font-mono text-xs text-foreground">
          {contractAddress ?? "Not deployed"}
        </p>
      </div>
      <div className="grid gap-1">
        <p className="text-xs font-medium uppercase text-muted-foreground">Market ID</p>
        <p className="break-all font-mono text-xs text-foreground">{marketId}</p>
      </div>
      {walletAddress ? (
        <div className="grid gap-1">
          <p className="text-xs font-medium uppercase text-muted-foreground">
            Connected wallet
          </p>
          <p className="break-all font-mono text-xs text-foreground">
            {walletAddress}
          </p>
        </div>
      ) : null}
    </div>
  )
}
