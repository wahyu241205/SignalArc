"use client"

import { type FormEvent, useMemo, useState } from "react"
import { AlertCircle, CheckCircle2, ExternalLink, Loader2 } from "lucide-react"
import { parseUnits, type Hash } from "viem"
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
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  ARC_TESTNET_CHAIN_ID,
  ARC_TESTNET_USDC_ADDRESS,
  ERC20_APPROVE_ABI,
  SIGNAL_ARC_MARKET_ABI,
  SIGNAL_ARC_MARKET_ADDRESS,
  USDC_ERC20_DECIMALS,
  getArcscanTxUrl,
} from "@/lib/contracts"
import { arcTestnet } from "@/lib/wagmi"

type SubmitState =
  | { status: "idle" }
  | { status: "approving"; approveHash?: Hash }
  | { status: "opening"; approveHash: Hash; openHash?: Hash }
  | { status: "success"; approveHash: Hash; openHash: Hash }
  | { status: "error"; message: string; approveHash?: Hash; openHash?: Hash }

type Outcome = "YES" | "NO"

const outcomeSide: Record<Outcome, 1 | 2> = {
  YES: 1,
  NO: 2,
}

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

  return "Unable to execute the Arc Testnet trade."
}

function parseUsdcAmount(value: string) {
  const normalized = value.trim()

  if (!normalized) {
    throw new Error("Enter a USDC amount.")
  }

  if (!/^\d+(\.\d{1,6})?$/.test(normalized)) {
    throw new Error("Enter a valid USDC amount with up to 6 decimals.")
  }

  const amount = parseUnits(normalized, USDC_ERC20_DECIMALS)
  if (amount <= BigInt(0)) {
    throw new Error("USDC amount must be greater than 0.")
  }

  return amount
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

function TradeExecutionResult({ state }: { state: Extract<SubmitState, { status: "success" }> }) {
  return (
    <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4">
      <div className="flex items-center gap-2">
        <CheckCircle2 className="h-4 w-4 text-green-400" aria-hidden="true" />
        <p className="text-sm font-medium text-green-300">Position opened on Arc Testnet</p>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">
        This is a real connected-wallet transaction on Arc Testnet. No production settlement.
      </p>
      <dl className="mt-4 grid gap-3 text-sm sm:grid-cols-2">
        <div>
          <dt className="text-xs uppercase tracking-wider text-muted-foreground/70">USDC Approval</dt>
          <dd className="mt-0.5">
            <TxLink hash={state.approveHash} />
          </dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wider text-muted-foreground/70">Open Position</dt>
          <dd className="mt-0.5">
            <TxLink hash={state.openHash} />
          </dd>
        </div>
      </dl>
    </div>
  )
}

function PendingState({ state }: { state: Extract<SubmitState, { status: "approving" | "opening" }> }) {
  const label = state.status === "approving" ? "Approving USDC" : "Opening position"

  return (
    <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
      <div className="flex items-center gap-2">
        <Loader2 className="h-4 w-4 animate-spin text-blue-300" aria-hidden="true" />
        <p className="text-sm font-medium text-blue-200">{label}</p>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">
        Confirm the wallet prompt, then wait for the Arc Testnet transaction to confirm.
      </p>
      {state.approveHash ? (
        <p className="mt-3 text-sm text-muted-foreground">
          Approval: <TxLink hash={state.approveHash} />
        </p>
      ) : null}
      {state.status === "opening" && state.openHash ? (
        <p className="mt-2 text-sm text-muted-foreground">
          Market transaction: <TxLink hash={state.openHash} />
        </p>
      ) : null}
    </div>
  )
}

export function TradeIntentPanel({ marketId, marketStatus }: { marketId: string; marketStatus: string }) {
  const config = useConfig()
  const { address, isConnected } = useAccount()
  const chainId = useChainId()
  const { switchChain, isPending: isSwitchingChain } = useSwitchChain()
  const [outcome, setOutcome] = useState<Outcome>("YES")
  const [amount, setAmount] = useState("1")
  const [state, setState] = useState<SubmitState>({ status: "idle" })
  const isTradingOpen = marketStatus.toUpperCase() === "OPEN"
  const isArcTestnet = chainId === ARC_TESTNET_CHAIN_ID
  const isPending = state.status === "approving" || state.status === "opening"
  const parsedAmount = useMemo(() => {
    try {
      return parseUsdcAmount(amount)
    } catch {
      return null
    }
  }, [amount])
  const { data: isContractOpen } = useReadContract({
    address: SIGNAL_ARC_MARKET_ADDRESS,
    abi: SIGNAL_ARC_MARKET_ABI,
    functionName: "isOpen",
    chainId: ARC_TESTNET_CHAIN_ID,
    query: {
      enabled: isConnected && isArcTestnet,
    },
  })

  const disabledReason = (() => {
    if (!isConnected) return "Connect wallet to trade."
    if (!isArcTestnet) return "Switch to Arc Testnet."
    if (!isTradingOpen) return "Trading is not open for this market."
    if (isContractOpen === false) return "The Arc Testnet market contract is not open."
    if (!parsedAmount) return "Enter a valid USDC amount."
    return null
  })()
  const canSubmit = !disabledReason && !isPending

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!canSubmit || !address) return

    let approveHash: Hash | undefined
    let openHash: Hash | undefined
    try {
      const usdcAmount = parseUsdcAmount(amount)
      setState({ status: "approving" })

      approveHash = await writeContract(config, {
        address: ARC_TESTNET_USDC_ADDRESS,
        abi: ERC20_APPROVE_ABI,
        functionName: "approve",
        args: [SIGNAL_ARC_MARKET_ADDRESS, usdcAmount],
        chainId: ARC_TESTNET_CHAIN_ID,
        account: address,
      })
      setState({ status: "approving", approveHash })
      await waitForTransactionReceipt(config, {
        hash: approveHash,
        chainId: ARC_TESTNET_CHAIN_ID,
      })

      setState({ status: "opening", approveHash })
      openHash = await writeContract(config, {
        address: SIGNAL_ARC_MARKET_ADDRESS,
        abi: SIGNAL_ARC_MARKET_ABI,
        functionName: "openPosition",
        args: [outcomeSide[outcome], usdcAmount],
        chainId: ARC_TESTNET_CHAIN_ID,
        account: address,
      })
      setState({ status: "opening", approveHash, openHash })
      await waitForTransactionReceipt(config, {
        hash: openHash,
        chainId: ARC_TESTNET_CHAIN_ID,
      })

      setState({ status: "success", approveHash, openHash })
    } catch (error) {
      setState({
        status: "error",
        message: getErrorMessage(error),
        approveHash,
        openHash,
      })
    }
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <CardTitle>Place Trade</CardTitle>
          <Badge variant="outline" className="border-green-500/30 bg-green-500/10 text-green-300 text-xs">
            Arc Testnet
          </Badge>
        </div>
        <CardDescription>
          This executes on Arc Testnet from your connected wallet. No production settlement.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <div className="rounded-lg border border-border bg-muted/20 p-4 text-sm text-muted-foreground">
            <p>
              Contract: <span className="font-mono text-xs text-foreground">{SIGNAL_ARC_MARKET_ADDRESS}</span>
            </p>
            <p className="mt-2">
              Market ID: <span className="font-mono text-xs text-foreground">{marketId}</span>
            </p>
            {address ? (
              <p className="mt-2">
                Connected wallet: <span className="font-mono text-xs text-foreground">{address}</span>
              </p>
            ) : null}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="outcome">Outcome</Label>
              <div className="grid grid-cols-2 gap-2">
                <label className="flex cursor-pointer items-center justify-center rounded-lg border border-border bg-input px-4 py-2.5 text-sm font-medium transition-colors has-[:checked]:border-green-500/50 has-[:checked]:bg-green-500/10 has-[:checked]:text-green-300">
                  <input
                    type="radio"
                    name="outcome"
                    value="YES"
                    checked={outcome === "YES"}
                    className="sr-only"
                    onChange={() => setOutcome("YES")}
                  />
                  YES
                </label>
                <label className="flex cursor-pointer items-center justify-center rounded-lg border border-border bg-input px-4 py-2.5 text-sm font-medium transition-colors has-[:checked]:border-red-500/50 has-[:checked]:bg-red-500/10 has-[:checked]:text-red-300">
                  <input
                    type="radio"
                    name="outcome"
                    value="NO"
                    checked={outcome === "NO"}
                    className="sr-only"
                    onChange={() => setOutcome("NO")}
                  />
                  NO
                </label>
              </div>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="amount">Amount (USDC)</Label>
              <Input
                id="amount"
                inputMode="decimal"
                min="0"
                name="amount"
                onChange={(event) => setAmount(event.target.value)}
                required
                step="0.000001"
                value={amount}
              />
            </div>
          </div>

          <p className="text-xs text-muted-foreground">
            {amount || "0"} USDC will be sent as {parsedAmount ? parsedAmount.toString() : "0"} base units.
          </p>

          {state.status === "error" ? (
            <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
              <div className="flex items-center gap-2">
                <AlertCircle className="h-4 w-4 text-destructive" aria-hidden="true" />
                <p className="text-sm font-medium text-destructive">Unable to execute trade</p>
              </div>
              <p className="mt-1 text-sm text-muted-foreground">{state.message}</p>
              {state.approveHash ? (
                <p className="mt-2 text-sm text-muted-foreground">
                  Approval: <TxLink hash={state.approveHash} />
                </p>
              ) : null}
              {state.openHash ? (
                <p className="mt-2 text-sm text-muted-foreground">
                  Market transaction: <TxLink hash={state.openHash} />
                </p>
              ) : null}
            </div>
          ) : null}

          {state.status === "approving" || state.status === "opening" ? <PendingState state={state} /> : null}

          {state.status === "success" ? <TradeExecutionResult state={state} /> : null}

          {disabledReason ? (
            <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-4">
              <p className="text-sm font-medium text-yellow-300">
                {disabledReason}
              </p>
              {!isConnected ? (
                <p className="mt-1 text-sm text-muted-foreground">Use the wallet control in the header to connect.</p>
              ) : null}
              {isConnected && !isArcTestnet ? (
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
              ) : null}
            </div>
          ) : null}

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
