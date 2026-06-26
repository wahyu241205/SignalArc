import { Card, CardContent } from "@/components/ui/card"

function WalletIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={1.5}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 9m18 0V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v3"
      />
    </svg>
  )
}

export function WalletNotConnectedState() {
  return (
    <Card>
      <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
          <WalletIcon className="h-6 w-6 text-muted-foreground" />
        </div>
        <p className="text-sm text-muted-foreground">
          Connect your wallet to view your portfolio.
        </p>
        <p className="text-xs text-muted-foreground/70">
          Positions are read-only until wallet-to-user mapping is implemented.
        </p>
      </CardContent>
    </Card>
  )
}

export function WalletIdentityCard({ address }: { address: string }) {
  return (
    <Card className="border-indigo-500/20">
      <CardContent className="flex items-center gap-3 pt-6">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-indigo-500/10">
          <WalletIcon className="h-5 w-5 text-indigo-400" />
        </div>
        <div>
          <p className="text-sm font-medium text-foreground">Connected Wallet</p>
          <p className="font-mono text-xs text-muted-foreground">{address}</p>
        </div>
      </CardContent>
    </Card>
  )
}
