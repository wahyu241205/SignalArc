"use client"

import { useEffect, useState } from "react"

import { MarketResolutionPanel } from "@/features/markets/market-resolution-panel"
import { OnchainMarketLifecyclePanel } from "@/features/markets/onchain-market-lifecycle-panel"
import { TradeIntentPanel } from "@/features/markets/trade-intent-panel"
import { ApiError, getMarket } from "@/lib/api"
import {
  MarketDetailErrorState,
  MarketDetailLoadingSkeleton,
  MarketMetadataCard,
  MarketOutcomeCard,
  MarketSummaryCard,
  type MarketDetailState,
} from "@/modules/markets/detail"

function getErrorState(
  error: unknown,
): Extract<MarketDetailState, { status: "error" }> {
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
    message: "Unable to load market.",
    requestId: null,
  }
}

export function MarketDetail({ marketId }: { marketId: string }) {
  const [state, setState] = useState<MarketDetailState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadMarket() {
      setState({ status: "loading" })

      try {
        const response = await getMarket(marketId)

        if (isActive) {
          setState({ status: "ready", market: response.data.market })
        }
      } catch (error) {
        if (isActive) {
          setState(getErrorState(error))
        }
      }
    }

    void loadMarket()

    return () => {
      isActive = false
    }
  }, [marketId])

  if (state.status === "loading") {
    return <MarketDetailLoadingSkeleton />
  }

  if (state.status === "error") {
    return (
      <MarketDetailErrorState
        message={state.message}
        requestId={state.requestId}
      />
    )
  }

  const { market } = state

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      <div className="space-y-6 lg:col-span-2">
        <MarketSummaryCard market={market} />
        <MarketOutcomeCard market={market} />
        <MarketMetadataCard market={market} />
      </div>

      <div className="space-y-6 lg:col-span-1">
        <TradeIntentPanel
          marketId={market.id}
          marketStatus={market.status}
          marketContractAddress={market.market_contract_address}
        />
        <OnchainMarketLifecyclePanel
          marketContractAddress={market.market_contract_address}
        />
        <MarketResolutionPanel marketId={market.id} />
      </div>
    </div>
  )
}
