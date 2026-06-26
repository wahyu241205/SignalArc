"use client"

import { useRouter } from "next/navigation"
import { type FormEvent, useState } from "react"
import { type Hash } from "viem"
import { waitForTransactionReceipt, writeContract } from "wagmi/actions"
import { useAccount, useChainId, useConfig, useSwitchChain } from "wagmi"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ApiError,
  attachMarketContract,
  createMarket,
  localDemoUserId,
  type Market,
} from "@/lib/api"
import {
  ARC_TESTNET_CHAIN_ID,
  ARC_TESTNET_USDC_ADDRESS,
  SIGNAL_ARC_MARKET_FACTORY_ABI,
  SIGNAL_ARC_MARKET_FACTORY_ADDRESS,
  isArcTestnetChain,
} from "@/lib/contracts"
import { arcTestnet } from "@/lib/wagmi"
import {
  closeTimestampSeconds,
  CreateMarketAdvancedFields,
  CreateMarketCategoryField,
  CreateMarketCoverUploadField,
  CreateMarketDeployPanel,
  CreateMarketFields,
  CreateMarketOutcomeFields,
  CreateMarketScheduleFields,
  CreateMarketSubmitStatus,
  getDeployedMarketAddress,
  getDeployErrorMessage,
  optionalText,
  requiredText,
  toRfc3339,
  type DeployState,
  type SubmitState,
} from "@/modules/markets/create"

function getErrorState(
  error: unknown,
): Extract<SubmitState, { status: "error" }> {
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

export function CreateMarketForm() {
  const router = useRouter()
  const config = useConfig()
  const { address, isConnected } = useAccount()
  const chainId = useChainId()
  const { switchChain, isPending: isSwitchingChain } = useSwitchChain()
  const [state, setState] = useState<SubmitState>({ status: "idle" })
  const [deployState, setDeployState] = useState<DeployState>({
    status: "idle",
  })
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [coverImageUrl, setCoverImageUrl] = useState("")
  const [isUploadingImage, setIsUploadingImage] = useState(false)
  const isArcTestnet = isArcTestnetChain(chainId)
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
    const canDeploy = Boolean(
      isConnected &&
        isArcTestnet &&
        SIGNAL_ARC_MARKET_FACTORY_ADDRESS &&
        !isDeploying,
    )

    return (
      <CreateMarketDeployPanel
        market={state.market}
        deployState={deployState}
        canDeploy={canDeploy}
        isConnected={isConnected}
        isArcTestnet={isArcTestnet}
        isSwitchingChain={isSwitchingChain}
        isDeploying={isDeploying}
        onSwitchChain={() => switchChain({ chainId: arcTestnet.id })}
        onDeploy={() => handleDeploy(state.market)}
      />
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

          <CreateMarketFields />

          <CreateMarketCoverUploadField
            coverImageUrl={coverImageUrl}
            onCoverImageUrlChange={setCoverImageUrl}
            onUploadingChange={setIsUploadingImage}
            disabled={state.status === "submitting"}
            isUploadingImage={isUploadingImage}
          />

          <CreateMarketCategoryField
            disabled={state.status === "submitting" || isUploadingImage}
          />

          <CreateMarketOutcomeFields />

          <CreateMarketScheduleFields />

          <CreateMarketAdvancedFields
            showAdvanced={showAdvanced}
            onToggleAdvanced={() => setShowAdvanced(!showAdvanced)}
          />

          <CreateMarketSubmitStatus
            state={state}
            isUploadingImage={isUploadingImage}
          />
        </form>
      </CardContent>
    </Card>
  )
}
