"use client"

import Link from "next/link"
import { type FormEvent, useState } from "react"

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
import { ApiError, createMarket, localDemoUserId, type Market } from "@/lib/api"

type SubmitState =
  | { status: "idle" }
  | { status: "submitting" }
  | { status: "success"; market: Market }
  | { status: "error"; message: string; requestId: string | null }

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

export function CreateMarketForm() {
  const [state, setState] = useState<SubmitState>({ status: "idle" })
  const [showAdvanced, setShowAdvanced] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setState({ status: "submitting" })

    const formData = new FormData(event.currentTarget)
    const opensAt = optionalText(formData, "opens_at")
    const closesAt = requiredText(formData, "closes_at")

    try {
      const response = await createMarket({
        creator_user_id: requiredText(formData, "creator_user_id"),
        title: requiredText(formData, "title"),
        description: optionalText(formData, "description"),
        category: optionalText(formData, "category"),
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

  if (state.status === "success") {
    return (
      <Card className="border-green-500/20">
        <CardHeader>
          <div className="flex items-center gap-2">
            <svg className="h-5 w-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
            </svg>
            <CardTitle>Market Created</CardTitle>
          </div>
          <CardDescription>{state.market.title}</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-3 sm:flex-row">
          <Button asChild>
            <Link href={`/markets/${state.market.id}`}>View Market</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/markets">Back to Markets</Link>
          </Button>
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
          {/* Hidden creator ID — prefilled for demo */}
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

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="category">Category</Label>
              <Input id="category" name="category" placeholder="Crypto, Politics, Sports..." />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="resolution_source">Resolution Source</Label>
              <Input
                id="resolution_source"
                name="resolution_source"
                placeholder="CoinGecko, AP News, Official announcement..."
              />
            </div>
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

          {/* Advanced settings — collapsed by default */}
          <div>
            <button
              type="button"
              onClick={() => setShowAdvanced(!showAdvanced)}
              className="text-xs font-medium text-muted-foreground hover:text-foreground transition-colors"
            >
              {showAdvanced ? "▾ Hide" : "▸ Show"} advanced settings
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
            <Button disabled={state.status === "submitting"} type="submit">
              {state.status === "submitting" ? "Creating..." : "Create Market"}
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
