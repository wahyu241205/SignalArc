"use client"

import { useState } from "react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CardTitle } from "@/components/ui/card"
import { getMarketCategoryLabel } from "@/modules/categories"
import type { Market } from "@/lib/api"

import {
  marketStatusBadgeClass,
  marketStatusContext,
} from "../format"

function ShareMarketButton({ market }: { market: Market }) {
  const [status, setStatus] = useState<"idle" | "copied" | "shared">("idle")

  async function handleShare() {
    const url = window.location.href
    const text =
      market.description ||
      `Trade ${market.outcome_yes_label} or ${market.outcome_no_label} on SignalArc.`

    try {
      if (navigator.share) {
        await navigator.share({
          title: market.title,
          text,
          url,
        })
        setStatus("shared")
      } else {
        await navigator.clipboard.writeText(url)
        setStatus("copied")
      }

      window.setTimeout(() => setStatus("idle"), 2000)
    } catch (error) {
      if (error instanceof DOMException && error.name === "AbortError") {
        return
      }

      try {
        await navigator.clipboard.writeText(url)
        setStatus("copied")
        window.setTimeout(() => setStatus("idle"), 2000)
      } catch {
        setStatus("idle")
      }
    }
  }

  return (
    <Button type="button" variant="outline" size="sm" onClick={handleShare}>
      {status === "copied"
        ? "Link copied"
        : status === "shared"
          ? "Shared"
          : "Share Market"}
    </Button>
  )
}

export function MarketDetailHeader({ market }: { market: Market }) {
  const context = marketStatusContext(market.status, market.winning_outcome)

  return (
    <>
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="outline" className={marketStatusBadgeClass(market.status)}>
            {market.status}
          </Badge>
          <Badge variant="secondary">
            {getMarketCategoryLabel(market.category)}
          </Badge>
          <span className="text-xs text-muted-foreground">
            {market.collateral_asset} - {market.chain}
          </span>
        </div>
        <ShareMarketButton market={market} />
      </div>

      <CardTitle className="text-xl leading-snug sm:text-2xl">
        {market.title}
      </CardTitle>

      {market.description ? (
        <p className="max-w-3xl text-sm leading-6 text-muted-foreground">
          {market.description}
        </p>
      ) : null}

      {context ? <p className="text-sm text-muted-foreground">{context}</p> : null}
    </>
  )
}
