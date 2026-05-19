"use client"

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
import {
  ApiError,
  createTradeIntent,
  type CreateTradeIntentRequest,
  type TradeIntentResponse,
} from "@/lib/api"

type SubmitState =
  | { status: "idle" }
  | { status: "submitting" }
  | { status: "success"; response: TradeIntentResponse }
  | { status: "error"; message: string; requestId: string | null }

function requiredText(formData: FormData, key: string) {
  return String(formData.get(key) ?? "").trim()
}

function selectValue<TValue extends string>(
  formData: FormData,
  key: string,
  allowedValues: readonly TValue[],
): TValue {
  const value = requiredText(formData, key)

  if (allowedValues.includes(value as TValue)) {
    return value as TValue
  }

  throw new Error(`${key} is invalid.`)
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
    message: "Unable to submit trade intent.",
    requestId: null,
  }
}

function TradeIntentResult({ response }: { response: TradeIntentResponse }) {
  const { trade, execution } = response

  return (
    <div className="rounded-lg border bg-muted/40 p-4">
      <p className="text-sm font-medium">Intent submitted</p>
      <p className="mt-1 text-sm text-muted-foreground">
        Execution status: {execution.status}. {execution.reason}
      </p>
      <dl className="mt-4 grid gap-3 text-sm sm:grid-cols-2">
        <div>
          <dt className="font-medium text-muted-foreground">Trade ID</dt>
          <dd className="font-mono text-xs text-foreground">{trade.id}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Status</dt>
          <dd>{trade.status}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Outcome</dt>
          <dd>{trade.outcome}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Side</dt>
          <dd>{trade.side}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Quantity</dt>
          <dd>{trade.quantity}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Price</dt>
          <dd>{trade.price}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Collateral amount</dt>
          <dd>{trade.collateral_amount}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground">Fee amount</dt>
          <dd>{trade.fee_amount}</dd>
        </div>
      </dl>
    </div>
  )
}

export function TradeIntentPanel({ marketId }: { marketId: string }) {
  const [state, setState] = useState<SubmitState>({ status: "idle" })

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setState({ status: "submitting" })

    const formData = new FormData(event.currentTarget)

    try {
      const input: CreateTradeIntentRequest = {
        user_id: requiredText(formData, "user_id"),
        market_id: marketId,
        outcome: selectValue(formData, "outcome", ["YES", "NO"] as const),
        side: selectValue(formData, "side", ["BUY", "SELL"] as const),
        quantity: requiredText(formData, "quantity"),
        price: requiredText(formData, "price"),
      }
      const response = await createTradeIntent(input)

      setState({ status: "success", response: response.data })
    } catch (error) {
      setState(getErrorState(error))
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Trade intent</CardTitle>
        <CardDescription>
          Submit an intent to the backend only. This does not connect a wallet, execute
          onchain, update a position, or settle a trade.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <div className="grid gap-2">
            <Label htmlFor="user_id">User ID</Label>
            <Input id="user_id" name="user_id" placeholder="UUID" required />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="outcome">Outcome</Label>
              <select
                className="h-8 w-full rounded-lg border border-input bg-transparent px-2.5 text-sm outline-none transition-colors focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
                id="outcome"
                name="outcome"
                required
              >
                <option value="YES">YES</option>
                <option value="NO">NO</option>
              </select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="side">Side</Label>
              <select
                className="h-8 w-full rounded-lg border border-input bg-transparent px-2.5 text-sm outline-none transition-colors focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
                id="side"
                name="side"
                required
              >
                <option value="BUY">BUY</option>
                <option value="SELL">SELL</option>
              </select>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="quantity">Quantity</Label>
              <Input id="quantity" inputMode="decimal" name="quantity" required />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="price">Price</Label>
              <Input id="price" inputMode="decimal" name="price" required />
            </div>
          </div>

          {state.status === "error" ? (
            <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
              <p className="text-sm font-medium text-destructive">
                Unable to submit trade intent
              </p>
              <p className="mt-1 text-sm text-muted-foreground">{state.message}</p>
              {state.requestId ? (
                <p className="mt-2 font-mono text-xs text-muted-foreground">
                  Request ID: {state.requestId}
                </p>
              ) : null}
            </div>
          ) : null}

          {state.status === "success" ? <TradeIntentResult response={state.response} /> : null}

          <Button disabled={state.status === "submitting"} type="submit">
            {state.status === "submitting" ? "Submitting intent..." : "Submit intent"}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
