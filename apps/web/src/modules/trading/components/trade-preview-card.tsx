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
    <div className="rounded-lg border border-border bg-muted/20 p-4 text-sm text-muted-foreground">
      <p>
        Contract:{" "}
        <span className="font-mono text-xs text-foreground">
          {contractAddress ?? "Not deployed"}
        </span>
      </p>
      <p className="mt-2">
        Market ID:{" "}
        <span className="font-mono text-xs text-foreground">{marketId}</span>
      </p>
      {walletAddress ? (
        <p className="mt-2">
          Connected wallet:{" "}
          <span className="font-mono text-xs text-foreground">
            {walletAddress}
          </span>
        </p>
      ) : null}
    </div>
  )
}
