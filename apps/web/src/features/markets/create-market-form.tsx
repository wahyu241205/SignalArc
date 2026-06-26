"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { type FormEvent, useState } from "react"
import { CheckCircle2, ExternalLink, Loader2 } from "lucide-react"
import { decodeEventLog, type Address, type Hash, type TransactionReceipt } from "viem"
import { waitForTransactionReceipt, writeContract } from "wagmi/actions"
import { useAccount, useChainId, useConfig, useSwitchChain } from "wagmi"

import { MarketImageUpload } from "@/components/markets/market-image-upload"
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
import { Textarea } from "@/components/ui/textarea"
import {
  ApiError,
  attachMarketContract,
  createMarket,
  localDemoUserId,
  type Market,
} from "@/lib/api"
import { CREATE_MARKET_CATEGORIES, DEFAULT_CATEGORY_ID } from "@/modules/categories"
import {
  ARC_TESTNET_CHAIN_ID,
  ARC_TESTNET_USDC_ADDRESS,
  SIGNAL_ARC_MARKET_FACTORY_ABI,
  SIGNAL_ARC_MARKET_FACTORY_ADDRESS,
  getArcscanTxUrl,
} from "@/lib/contracts"
import { arcTestnet } from "@/lib/wagmi"

type SubmitState =
  | { status: "idle" }
  | { status: "submitting" }
  | { status: "success"; market: Market }
  | { status: "error"; message: string; requestId: string | null }

type DeployState =
  | { status: "idle" }
  | { status: "deploying"; hash?: Hash }
  | { status: "success"; hash: Hash; marketAddress: Address }
  | { status: "error"; message: string; hash?: Hash }

function optionalText(formData: FormData, key: string) {
  const value = String(formData.get(key) ?? "").trim()
  return value === "" ? undefined : value
}

function requiredText(formData: FormData, key: string) {
  return String(formData.get(key) ?? "").trim()
}

function toRfc3339(value: string) {
  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    throw new Error("Dates must be valid.")
  }

  return date.toISOString()
}

function closeTimestampSeconds(value: string) {
  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    throw new Error("Close date must be valid.")
  }

  return BigInt(Math.floor(date.getTime() / 1000))
}

function defaultCloseValue() {
  const date = new Date()
  date.setDate(date.getDate() + 7)
  date.setMinutes(date.getMinutes() - date.getTimezoneOffset())

  return date.toISOString().slice(0, 16)
}

function getErrorState(error: unknown): Extract<SubmitState, { status: "error" }> {
  if (error instanceof ApiError) {
    return {
      status: "error",
      message: error.message,
      requestId: error.requestId,
    }
  }

  if (error instanceof Error) {
    return {
      status: "error",
      message: error.message,
      requestId: null,
    }
  }

  return {
    status: "error",
    message: "Unable to create market.",
    requestId: null,
  }
}

function getDeployErrorMessage(error: unknown) {
  if (error instanceof Error) {
    const message = error.message.toLowerCase()
    if (
      message.includes("user rejected") ||
      message.includes("user denied") ||
      message.includes("rejected the request") ||
      message.includes("request rejected")
    ) {
      return "Wallet transaction was rejected. The backend market remains NOT_DEPLOYED."
    }

    return error.message
  }

  return "Unable to deploy the Arc Testnet market contract."
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

function getDeployedMarketAddress(receipt: TransactionReceipt): Address {
  for (const log of receipt.logs) {
    try {
      const decoded = decodeEventLog({
        abi: SIGNAL_ARC_MARKET_FACTORY_ABI,
        data: log.data,
        topics: log.topics,
      })

      if (decoded.eventName === "MarketDeployed") {
        return decoded.args.market
      }
    } catch {
      // Ignore logs from contracts other than the factory.
    }
  }

  throw new Error("MarketDeployed event was not found in the factory receipt.")
}

export function CreateMarketForm() {
  const router = useRouter()
  const config = useConfig()
  const { address, isConnected } = useAccount()
  const chainId = useChainId()
  const { switchChain, isPending: isSwitchingChain } = useSwitchChain()
  const [state, setState] = useState<SubmitState>({ status: "idle" })
  const [deployState, setDeployState] = useState<DeployState>({ status: "idle" })
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [coverImageUrl, setCoverImageUrl] = useState("")
  const [isUploadingImage, setIsUploadingImage] = useState(false)
  const isArcTestnet = chainId === ARC_TESTNET_CHAIN_ID
  const isDeploying = deployState.status === "deploying"

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    // Never submit while an image upload is mid-flight; the cover URL would
    // be stale or empty. The submit button is disabled in this state too.
    if (isUploadingImage) return

    setState({ status: "submitting" })
    setDeployState({ status: "idle" })

    const formData = new FormData(event.currentTarget)
    const opensAt = optionalText(formData, "opens_at")
    const closesAt = requiredText(formData, "closes_at")

    try {
      const response = await createMarket({
        creator_user_id: requiredText(formData, "creator_user_id"),
        title: requiredText(formData, "title"),
        description: optionalText(formData, "description"),
        category: optionalText(formData, "category"),
        cover_image_url: coverImageUrl.trim() || undefined,
        outcome_yes_label: optionalText(formData, "outcome_yes_label"),
        outcome_no_label: optionalText(formData, "outcome_no_label"),
        collateral_asset: optionalText(formData, "collateral_asset"),
        chain: requiredText(formData, "chain"),
        resolution_source: optionalText(formData, "resolution_source"),
        opens_at: opensAt ? toRfc3339(opensAt) : undefined,
        closes_at: toRfc3339(closesAt),
      })

      setState({ status: "success", market: response.data.market })
    } catch (error) {
      setState(getErrorState(error))
    }
  }

  async function handleDeploy(market: Market) {
    if (!address || !SIGNAL_ARC_MARKET_FACTORY_ADDRESS) return

    let hash: Hash | undefined
    try {
      setDeployState({ status: "deploying" })
      hash = await writeContract(config, {
        address: SIGNAL_ARC_MARKET_FACTORY_ADDRESS,
        abi: SIGNAL_ARC_MARKET_FACTORY_ABI,
        functionName: "createMarket",
        args: [
          market.id,
          market.title,
          closeTimestampSeconds(market.closes_at),
          address,
          ARC_TESTNET_USDC_ADDRESS,
        ],
        chainId: ARC_TESTNET_CHAIN_ID,
        account: address,
      })
      setDeployState({ status: "deploying", hash })

      const receipt = await waitForTransactionReceipt(config, {
        hash,
        chainId: ARC_TESTNET_CHAIN_ID,
      })
      const marketAddress = getDeployedMarketAddress(receipt)

      await attachMarketContract(market.id, {
        market_contract_address: marketAddress,
        market_deployment_tx_hash: hash,
        market_factory_address: SIGNAL_ARC_MARKET_FACTORY_ADDRESS,
        resolver_address: address,
      })

      setDeployState({ status: "success", hash, marketAddress })
      router.push(`/markets/${market.id}`)
    } catch (error) {
      setDeployState({
        status: "error",
        message: getDeployErrorMessage(error),
        hash,
      })
    }
  }

  if (state.status === "success") {
    const canDeploy = Boolean(isConnected && isArcTestnet && SIGNAL_ARC_MARKET_FACTORY_ADDRESS && !isDeploying)

    return (
      <Card className="border-green-500/20">
        <CardHeader>
          <div className="flex items-center gap-2">
            <CheckCircle2 className="h-5 w-5 text-green-400" aria-hidden="true" />
            <CardTitle>Market Created</CardTitle>
          </div>
          <CardDescription>{state.market.title}</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-4">
            <p className="text-sm font-medium text-yellow-300">Onchain contract not deployed.</p>
            <p className="mt-1 text-sm text-muted-foreground">
              Your connected wallet will be the resolver for this testnet market.
            </p>
            {SIGNAL_ARC_MARKET_FACTORY_ADDRESS ? (
              <p className="mt-2 font-mono text-xs text-muted-foreground">
                Factory: {SIGNAL_ARC_MARKET_FACTORY_ADDRESS}
              </p>
            ) : (
              <p className="mt-2 text-sm text-muted-foreground">Factory address not configured.</p>
            )}
          </div>

          {!isConnected ? (
            <p className="text-sm text-muted-foreground">Connect a wallet to deploy the Arc Testnet market contract.</p>
          ) : null}

          {isConnected && !isArcTestnet ? (
            <Button
              disabled={isSwitchingChain}
              onClick={() => switchChain({ chainId: arcTestnet.id })}
              type="button"
              variant="outline"
              className="w-full sm:w-fit"
            >
              {isSwitchingChain ? "Switching..." : "Switch to Arc Testnet"}
            </Button>
          ) : null}

          {deployState.status === "deploying" ? (
            <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
              <div className="flex items-center gap-2">
                <Loader2 className="h-4 w-4 animate-spin text-blue-300" aria-hidden="true" />
                <p className="text-sm font-medium text-blue-200">Deploying market contract</p>
              </div>
              {deployState.hash ? (
                <p className="mt-2 text-sm text-muted-foreground">
                  Transaction: <TxLink hash={deployState.hash} />
                </p>
              ) : null}
            </div>
          ) : null}

          {deployState.status === "error" ? (
            <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
              <p className="text-sm font-medium text-destructive">Unable to deploy onchain market</p>
              <p className="mt-1 text-sm text-muted-foreground">{deployState.message}</p>
              {deployState.hash ? (
                <p className="mt-2 text-sm text-muted-foreground">
                  Transaction: <TxLink hash={deployState.hash} />
                </p>
              ) : null}
            </div>
          ) : null}

          {deployState.status === "success" ? (
            <div className="rounded-lg border border-green-500/20 bg-green-500/5 p-4">
              <p className="text-sm font-medium text-green-300">Market contract deployed on Arc Testnet.</p>
              <p className="mt-2 font-mono text-xs text-muted-foreground">{deployState.marketAddress}</p>
              <p className="mt-2 text-sm text-muted-foreground">
                Transaction: <TxLink hash={deployState.hash} />
              </p>
            </div>
          ) : null}

          <div className="flex flex-col gap-3 sm:flex-row">
            <Button disabled={!canDeploy} onClick={() => handleDeploy(state.market)} type="button">
              {isDeploying ? "Deploying..." : "Deploy on Arc Testnet"}
            </Button>
            <Button asChild variant="outline">
              <Link href={`/markets/${state.market.id}`}>View Market</Link>
            </Button>
            <Button asChild variant="outline">
              <Link href="/markets">Back to Markets</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Create a Market</CardTitle>
        <CardDescription>
          Launch a new USDC-settled prediction market on Arc Testnet.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <input type="hidden" name="creator_user_id" value={localDemoUserId} />

          <div className="grid gap-2">
            <Label htmlFor="title">Market Question</Label>
            <Input
              id="title"
              name="title"
              placeholder="Will ETH reach $5,000 by end of 2026?"
              required
            />
            <p className="text-xs text-muted-foreground">
              Frame as a yes/no question that can be resolved objectively.
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="description">Description (optional)</Label>
            <Textarea
              id="description"
              name="description"
              rows={3}
              placeholder="Additional context, resolution criteria, or relevant links..."
            />
          </div>

          <div className="grid gap-2">
            <Label>Market Image</Label>
            <MarketImageUpload
              coverImageUrl={coverImageUrl}
              onChangeUrl={(url) => setCoverImageUrl(url ?? "")}
              onUploadingChange={setIsUploadingImage}
              disabled={state.status === "submitting"}
            />
            <p className="text-xs text-muted-foreground">
              Optional. Upload an image directly, or paste an HTTPS image URL below.
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="cover_image_url">Market Image URL (advanced)</Label>
            <Input
              id="cover_image_url"
              name="cover_image_url"
              type="url"
              inputMode="url"
              maxLength={2048}
              pattern="https://.*"
              placeholder="https://example.com/market-cover.png"
              value={coverImageUrl}
              onChange={(event) => setCoverImageUrl(event.target.value)}
              disabled={isUploadingImage || state.status === "submitting"}
            />
            <p className="text-xs text-muted-foreground">
              Paste a public HTTPS image URL as an alternative to uploading.
            </p>
          </div>

          <fieldset className="grid gap-2" disabled={state.status === "submitting" || isUploadingImage}>
            <legend className="text-sm font-medium leading-none">Category</legend>
            <div
              role="radiogroup"
              aria-label="Category"
              className="grid grid-cols-2 gap-2 sm:grid-cols-4"
            >
              {CREATE_MARKET_CATEGORIES.map((cat) => (
                <label
                  key={cat.id}
                  className="relative flex cursor-pointer items-center justify-center has-[:disabled]:cursor-not-allowed"
                >
                  <input
                    type="radio"
                    name="category"
                    value={cat.id}
                    defaultChecked={cat.id === DEFAULT_CATEGORY_ID}
                    disabled={state.status === "submitting" || isUploadingImage}
                    className="peer sr-only"
                  />
                  <span className="flex h-9 w-full items-center justify-center rounded-md border border-input bg-transparent px-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted/40 peer-checked:border-indigo-500 peer-checked:bg-indigo-500/10 peer-checked:text-indigo-300 peer-focus-visible:outline-none peer-focus-visible:ring-1 peer-focus-visible:ring-ring peer-disabled:opacity-50">
                    {cat.label}
                  </span>
                </label>
              ))}
            </div>
          </fieldset>

          <div className="grid gap-2">
            <Label htmlFor="resolution_source">Resolution Source</Label>
            <Input
              id="resolution_source"
              name="resolution_source"
              placeholder="CoinGecko, AP News, Official announcement..."
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="outcome_yes_label">YES Label</Label>
              <Input id="outcome_yes_label" name="outcome_yes_label" placeholder="Yes" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="outcome_no_label">NO Label</Label>
              <Input id="outcome_no_label" name="outcome_no_label" placeholder="No" />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="opens_at">Opens At (optional)</Label>
              <Input id="opens_at" name="opens_at" type="datetime-local" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="closes_at">Closes At</Label>
              <Input
                id="closes_at"
                name="closes_at"
                defaultValue={defaultCloseValue()}
                required
                type="datetime-local"
              />
            </div>
          </div>

          <div>
            <button
              type="button"
              onClick={() => setShowAdvanced(!showAdvanced)}
              className="text-xs font-medium text-muted-foreground hover:text-foreground transition-colors"
            >
              {showAdvanced ? "Hide" : "Show"} advanced settings
            </button>

            {showAdvanced ? (
              <div className="mt-3 grid gap-4 rounded-lg border border-border/50 p-4 sm:grid-cols-2">
                <div className="grid gap-2">
                  <Label htmlFor="collateral_asset" className="text-xs">Collateral Asset</Label>
                  <Input id="collateral_asset" name="collateral_asset" defaultValue="USDC" className="text-sm" />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="chain" className="text-xs">Chain</Label>
                  <Input id="chain" name="chain" defaultValue="Arc Testnet" className="text-sm" />
                </div>
              </div>
            ) : (
              <>
                <input type="hidden" name="collateral_asset" value="USDC" />
                <input type="hidden" name="chain" value="Arc Testnet" />
              </>
            )}
          </div>

          {state.status === "error" ? (
            <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
              <p className="text-sm font-medium text-destructive">Unable to create market</p>
              <p className="mt-1 text-sm text-muted-foreground">{state.message}</p>
              {state.requestId ? (
                <p className="mt-2 font-mono text-xs text-muted-foreground">
                  Request ID: {state.requestId}
                </p>
              ) : null}
            </div>
          ) : null}

          <div className="flex flex-col gap-3 sm:flex-row">
            <Button
              disabled={state.status === "submitting" || isUploadingImage}
              type="submit"
            >
              {state.status === "submitting"
                ? "Creating..."
                : isUploadingImage
                  ? "Uploading image..."
                  : "Create Market"}
            </Button>
            <Button asChild variant="outline">
              <Link href="/markets">Cancel</Link>
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
