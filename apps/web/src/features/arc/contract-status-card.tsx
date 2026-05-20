"use client"

import { useEffect, useState } from "react"

import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { ApiError, getArcContract, type ArcContractStatus } from "@/lib/api"

type ContractState =
  | { status: "loading" }
  | { status: "loaded"; contract: ArcContractStatus }
  | { status: "error"; message: string; requestId: string | null }

function getErrorState(error: unknown): Extract<ContractState, { status: "error" }> {
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
    message: "Unable to load Arc contract status.",
    requestId: null,
  }
}

function ContractField({ label, value }: { label: string; value: string }) {
  return (
    <div className="min-w-0">
      <dt className="text-sm font-medium text-muted-foreground">{label}</dt>
      <dd className="mt-1 break-all font-mono text-xs text-foreground">{value}</dd>
    </div>
  )
}

export function ContractStatusCard() {
  const [state, setState] = useState<ContractState>({ status: "loading" })

  useEffect(() => {
    let isActive = true

    async function loadContract() {
      try {
        const response = await getArcContract()
        if (isActive) {
          setState({ status: "loaded", contract: response.data })
        }
      } catch (error) {
        if (isActive) {
          setState(getErrorState(error))
        }
      }
    }

    void loadContract()

    return () => {
      isActive = false
    }
  }, [])

  if (state.status === "loading") {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Arc Testnet contract</CardTitle>
          <CardDescription>Loading prototype contract reference...</CardDescription>
        </CardHeader>
      </Card>
    )
  }

  if (state.status === "error") {
    return (
      <Card className="border-destructive/30 bg-destructive/5">
        <CardHeader>
          <CardTitle className="text-destructive">Arc contract status unavailable</CardTitle>
          <CardDescription>{state.message}</CardDescription>
        </CardHeader>
        {state.requestId ? (
          <CardContent>
            <p className="font-mono text-xs text-muted-foreground">
              Request ID: {state.requestId}
            </p>
          </CardContent>
        ) : null}
      </Card>
    )
  }

  const { contract } = state

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-wrap items-center gap-2">
          <CardTitle>Arc Testnet contract</CardTitle>
          <Badge variant="outline">Prototype</Badge>
          <Badge variant="secondary">Live use not approved</Badge>
        </div>
        <CardDescription>
          Read-only local reference for the deployed SignalArcMarket prototype.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <dl className="grid gap-4 sm:grid-cols-2">
          <ContractField label="Network" value={contract.network} />
          <ContractField label="Chain ID" value={String(contract.chain_id)} />
          <ContractField label="SignalArcMarket" value={contract.signal_arc_market} />
          <ContractField label="USDC ERC20 interface" value={contract.usdc_erc20_interface} />
          <ContractField label="Explorer" value={contract.explorer} />
          <ContractField label="Status" value={contract.status} />
        </dl>
      </CardContent>
    </Card>
  )
}
