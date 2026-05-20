"use client"

import { useEffect, useState } from "react"
import { AlertCircle, CheckCircle2, ExternalLink, Loader2 } from "lucide-react"
import { formatUnits, isAddress, zeroAddress, type Address, type Hash } from "viem"
import { waitForTransactionReceipt, writeContract } from "wagmi/actions"
import { useAccount, useChainId, useConfig, useReadContract, useSwitchChain } from "wagmi"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ARC_TESTNET_CHAIN_ID,
  SIGNAL_ARC_MARKET_ABI,
  USDC_ERC20_DECIMALS,
  getArcscanTxUrl,
} from "@/lib/contracts"
import { arcTestnet } from "@/lib/wagmi"

type ActionState =
  | { status: "idle" }
  | { status: "pending"; label: string; hash?: Hash }
  | { status: "success"; label: string; hash: Hash }
  | { status: "error"; label: string; message: string; hash?: Hash }

const marketStatusLabels: Record<number, string> = {
  0: "Draft",
  1: "Open",
  2: "Closed",
  3: "Resolved",
  4: "Cancelled",
}

const outcomeLabels: Record<number, string> = {
  0: "None",
  1: "Yes",
  2: "No",
}

const MARKET_STATUS_OPEN = 1
const MARKET_STATUS_CLOSED = 2
const MARKET_STATUS_RESOLVED = 3
const MARKET_STATUS_CANCELLED = 4

function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    const message = error.message.toLowerCase()
    if (
      message.includes("user rejected") ||
      message.includes("user denied") ||
      message.includes("rejected the request") ||
      message.includes("request rejected")
    ) {
      return "Wallet transaction was rejected."
    }

    return error.message
  }

  return "Unable to execute the Arc Testnet contract action."
}

function formatAddress(address: Address | undefined) {
  return address ?? "-"
}

function formatUsdc(value: bigint | undefined) {
  if (value === undefined) return "-"
  return `${formatUnits(value, USDC_ERC20_DECIMALS)} USDC`
}

function formatCloseTimestamp(value: bigint | undefined) {
  if (value === undefined) return "-"

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(Number(value) * 1000))
}

function isSameAddress(left: Address | undefined, right: Address | undefined) {
  return Boolean(left && right && left.toLowerCase() === right.toLowerCase())
}

function TxLink({ hash }: { hash: Hash }) {
  return (
    <a
      className="inline-flex items-center gap-1 font-mono text-xs text-green-200 underline underline-offset-4 hover:text-green-100"
      href={getArcscanTxUrl(hash)}
      rel="noreferrer"
      target="_blank"
    >
      {hash.slice(0, 10)}...{hash.slice(-8)}
      <ExternalLink className="h-3 w-3" aria-hidden="true" />
    </a>
  )
}

function OnchainField({ label, value, mono = false }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="min-w-0">
      <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">{label}</dt>
      <dd className={`mt-1 break-all text-sm text-foreground ${mono ? "font-mono text-xs" : ""}`}>{value}</dd>
    </div>
  )
}

export function OnchainMarketLifecyclePanel({
  marketContractAddress,
}: {
  marketContractAddress: string | null
}) {
  const config = useConfig()
  const chainId = useChainId()
  const { address, isConnected } = useAccount()
  const { switchChain, isPending: isSwitchingChain } = useSwitchChain()
  const [actionState, setActionState] = useState<ActionState>({ status: "idle" })
  const [nowSeconds, setNowSeconds] = useState(() => Math.floor(Date.now() / 1000))
  const isArcTestnet = chainId === ARC_TESTNET_CHAIN_ID
  const contractAddress =
    marketContractAddress && isAddress(marketContractAddress)
      ? (marketContractAddress as Address)
      : null
  const userAddress = address ?? zeroAddress
  const readsEnabled = Boolean(contractAddress)
  const walletReadsEnabled = readsEnabled && isConnected

  const resolverRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "resolver",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const statusRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "status",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const closeTimestampRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "closeTimestamp",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const winningOutcomeRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "winningOutcome",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const totalYesRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "totalYes",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const totalNoRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "totalNo",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const totalCollateralRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "totalCollateral",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: readsEnabled },
  })
  const yesPositionRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "yesPositions",
    args: [userAddress],
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: walletReadsEnabled },
  })

  const noPositionRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "noPositions",
    args: [userAddress],
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: walletReadsEnabled },
  })
  const claimableRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "claimableAmount",
    args: [userAddress],
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: walletReadsEnabled },
  })
  const hasClaimedRead = useReadContract({
    address: contractAddress ?? undefined,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "hasClaimed",
    args: [userAddress],
    chainId: ARC_TESTNET_CHAIN_ID,
    query: { enabled: walletReadsEnabled },
  })

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setNowSeconds(Math.floor(Date.now() / 1000))
    }, 30_000)

    return () => window.clearInterval(intervalId)
  }, [])

  if (!contractAddress) {
    return (
      <Card className="border-yellow-500/20">
        <CardHeader>
          <CardTitle>Onchain Status</CardTitle>
          <CardDescription>Onchain contract not deployed for this market.</CardDescription>
        </CardHeader>
      </Card>
    )
  }

  const deployedContractAddress = contractAddress
  const statusValue = statusRead.data
  const closeTimestamp = closeTimestampRead.data
  const claimableAmount = claimableRead.data ?? BigInt(0)
  const hasClaimed = hasClaimedRead.data ?? false
  const isOpen = statusValue === MARKET_STATUS_OPEN
  const isClosed = statusValue === MARKET_STATUS_CLOSED
  const isResolved = statusValue === MARKET_STATUS_RESOLVED
  const isCancelled = statusValue === MARKET_STATUS_CANCELLED
  const hasReachedCloseTime = closeTimestamp !== undefined && BigInt(nowSeconds) >= closeTimestamp
  const isResolver = isSameAddress(address, resolverRead.data)
  const isPending = actionState.status === "pending"
  const canClaim = isConnected && isArcTestnet && (isResolved || isCancelled) && claimableAmount > BigInt(0) && !hasClaimed

  const claimDisabledReason = (() => {
    if (!isConnected) return "Connect wallet to check claim eligibility."
    if (!isArcTestnet) return "Switch to Arc Testnet."
    if (!isResolved && !isCancelled) return "Claims are available only after resolution or cancellation."
    if (hasClaimed) return "This wallet has already claimed."
    if (claimableAmount === BigInt(0)) return "No claimable USDC for this wallet."
    return null
  })()

  const resolverDisabledReason = (() => {
    if (!isConnected) return "Connect the resolver wallet to manage lifecycle actions."
    if (!isArcTestnet) return "Switch to Arc Testnet."
    if (!isResolver) return "Connected wallet is not the resolver for this Arc Testnet market."
    return null
  })()

  async function refetchOnchainReads() {
    await Promise.all([
      resolverRead.refetch(),
      statusRead.refetch(),
      closeTimestampRead.refetch(),
      winningOutcomeRead.refetch(),
      totalYesRead.refetch(),
      totalNoRead.refetch(),
      totalCollateralRead.refetch(),
      yesPositionRead.refetch(),
      noPositionRead.refetch(),
      claimableRead.refetch(),
      hasClaimedRead.refetch(),
    ])
  }

  function getResolverActionDisabledReason(action: "close" | "resolve" | "cancel") {
    if (isPending) return "Wait for the current transaction to confirm."
    if (!isConnected) return "Connect the resolver wallet to manage lifecycle actions."
    if (!isArcTestnet) return "Switch to Arc Testnet."
    if (!isResolver) return "Connected wallet is not the resolver for this Arc Testnet market."
    if (isResolved || isCancelled) return "Resolver actions are disabled after resolution or cancellation."
    if (statusValue === undefined) return "Loading onchain market status."

    if (action === "close") {
      if (!isOpen) return "Close is only available while the market is Open."
      if (closeTimestamp === undefined) return "Loading market close time."
      if (!hasReachedCloseTime) return "Close is available after the market close time."
    }

    if (action === "resolve" && !isClosed) {
      return "Resolve is only available after the market is Closed."
    }

    if (action === "cancel" && !isOpen && !isClosed) {
      return "Cancel is only available while the market is Open or Closed."
    }

    return null
  }

  async function runMarketAction(label: string, functionName: "closeMarket" | "cancelMarket" | "claim") {
    if (!address) return

    let hash: Hash | undefined
    try {
      setActionState({ status: "pending", label })
      hash = await writeContract(config, {
        address: deployedContractAddress,
        abi: SIGNAL_ARC_MARKET_ABI,
        functionName,
        chainId: ARC_TESTNET_CHAIN_ID,
        account: address,
      })
      setActionState({ status: "pending", label, hash })
      await waitForTransactionReceipt(config, { hash, chainId: ARC_TESTNET_CHAIN_ID })
      setActionState({ status: "success", label, hash })
      await refetchOnchainReads()
    } catch (error) {
      setActionState({ status: "error", label, message: getErrorMessage(error), hash })
    }
  }

  async function runResolveAction(winningOutcome: 1 | 2) {
    if (!address) return

    const label = winningOutcome === 1 ? "Resolve YES" : "Resolve NO"
    let hash: Hash | undefined
    try {
      setActionState({ status: "pending", label })
      hash = await writeContract(config, {
        address: deployedContractAddress,
        abi: SIGNAL_ARC_MARKET_ABI,
        functionName: "resolve",
        args: [winningOutcome],
        chainId: ARC_TESTNET_CHAIN_ID,
        account: address,
      })
      setActionState({ status: "pending", label, hash })
      await waitForTransactionReceipt(config, { hash, chainId: ARC_TESTNET_CHAIN_ID })
      setActionState({ status: "success", label, hash })
      await refetchOnchainReads()
    } catch (error) {
      setActionState({ status: "error", label, message: getErrorMessage(error), hash })
    }
  }

  const claimButtonLabel = isCancelled ? "Claim Refund" : "Claim Payout"
  const closeDisabledReason = getResolverActionDisabledReason("close")
  const resolveDisabledReason = getResolverActionDisabledReason("resolve")
  const cancelDisabledReason = getResolverActionDisabledReason("cancel")
  const resolverActionReason = closeDisabledReason ?? resolveDisabledReason ?? cancelDisabledReason

  return (
    <Card className="border-indigo-500/20">
      <CardHeader>
        <div className="flex flex-wrap items-center gap-2">
          <CardTitle>Onchain Status</CardTitle>
          <Badge variant="outline" className="border-indigo-500/30 bg-indigo-500/10 text-indigo-300">
            Arc Testnet
          </Badge>
        </div>
        <CardDescription>
          Market-specific browser-wallet reads and lifecycle transactions. No real funds or production settlement.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        <dl className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <OnchainField label="Contract Address" value={deployedContractAddress} mono />
          <OnchainField label="Resolver Address" value={formatAddress(resolverRead.data)} mono />
          <OnchainField label="Connected Wallet" value={formatAddress(address)} mono />
          <OnchainField
            label="Status"
            value={statusValue === undefined ? "-" : marketStatusLabels[statusValue] ?? `Unknown (${statusValue})`}
          />
          <OnchainField label="Close Time" value={formatCloseTimestamp(closeTimestamp)} />
          <OnchainField
            label="Winning Outcome"
            value={isResolved ? outcomeLabels[winningOutcomeRead.data ?? 0] ?? `Unknown (${winningOutcomeRead.data})` : "-"}
          />
          <OnchainField label="User YES" value={isConnected ? formatUsdc(yesPositionRead.data) : "-"} />
          <OnchainField label="User NO" value={isConnected ? formatUsdc(noPositionRead.data) : "-"} />
          <OnchainField label="Claimable USDC" value={isConnected ? formatUsdc(claimableAmount) : "-"} />
          <OnchainField label="Claimed Status" value={isConnected ? (hasClaimed ? "Claimed" : "Not claimed") : "-"} />
          <OnchainField label="Total YES" value={formatUsdc(totalYesRead.data)} />
          <OnchainField label="Total NO" value={formatUsdc(totalNoRead.data)} />
          <OnchainField label="Total Collateral" value={formatUsdc(totalCollateralRead.data)} />
        </dl>

        {!isArcTestnet ? (
          <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-4">
            <p className="text-sm font-medium text-yellow-300">Wallet is not on Arc Testnet.</p>
            <Button
              className="mt-3"
              disabled={isSwitchingChain}
              onClick={() => switchChain({ chainId: arcTestnet.id })}
              size="sm"
              type="button"
              variant="outline"
            >
              {isSwitchingChain ? "Switching..." : "Switch network"}
            </Button>
          </div>
        ) : null}

        <div className="grid gap-3 rounded-lg border border-border bg-muted/20 p-4">
          <h3 className="text-sm font-medium text-foreground">Claim</h3>
          {claimDisabledReason ? <p className="text-sm text-muted-foreground">{claimDisabledReason}</p> : null}
          <Button
            className="w-full sm:w-fit"
            disabled={!canClaim || isPending}
            onClick={() => runMarketAction(claimButtonLabel, "claim")}
            type="button"
          >
            {isPending && actionState.label === claimButtonLabel ? "Claiming..." : claimButtonLabel}
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
            <Button disabled={Boolean(closeDisabledReason)} onClick={() => runMarketAction("Close Market", "closeMarket")} type="button" variant="outline">
              Close Market
            </Button>
            <Button disabled={Boolean(resolveDisabledReason)} onClick={() => runResolveAction(1)} type="button" variant="outline">
              Resolve YES
            </Button>
            <Button disabled={Boolean(resolveDisabledReason)} onClick={() => runResolveAction(2)} type="button" variant="outline">
              Resolve NO
            </Button>
            <Button disabled={Boolean(cancelDisabledReason)} onClick={() => runMarketAction("Cancel Market", "cancelMarket")} type="button" variant="destructive">
              Cancel Market
            </Button>
          </div>
        </div>

        {actionState.status === "pending" ? (
          <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
            <div className="flex items-center gap-2">
              <Loader2 className="h-4 w-4 animate-spin text-blue-300" aria-hidden="true" />
              <p className="text-sm font-medium text-blue-200">{actionState.label} pending</p>
            </div>
            {actionState.hash ? (
              <p className="mt-3 text-sm text-muted-foreground">
                Transaction: <TxLink hash={actionState.hash} />
              </p>
            ) : null}
          </div>
        ) : null}

        {actionState.status === "success" ? (
          <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="h-4 w-4 text-green-400" aria-hidden="true" />
              <p className="text-sm font-medium text-green-300">{actionState.label} confirmed on Arc Testnet</p>
            </div>
            <p className="mt-2 text-sm text-muted-foreground">
              Transaction: <TxLink hash={actionState.hash} />
            </p>
          </div>
        ) : null}

        {actionState.status === "error" ? (
          <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
            <div className="flex items-center gap-2">
              <AlertCircle className="h-4 w-4 text-destructive" aria-hidden="true" />
              <p className="text-sm font-medium text-destructive">Unable to run {actionState.label}</p>
            </div>
            <p className="mt-1 text-sm text-muted-foreground">{actionState.message}</p>
            {actionState.hash ? (
              <p className="mt-2 text-sm text-muted-foreground">
                Transaction: <TxLink hash={actionState.hash} />
              </p>
            ) : null}
          </div>
        ) : null}
      </CardContent>
    </Card>
  )
}
