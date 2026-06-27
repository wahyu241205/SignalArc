import Link from "next/link"

import { InlineErrorState } from "@/components/shared"
import { Button } from "@/components/ui/button"
import { ChainStatusCard, TransactionLink } from "@/modules/wallet"

import type { DeployState, SubmitState } from "../types"

export function CreateMarketSubmitStatus({
  state,
  deployState,
  isUploadingImage,
  canSubmit,
  isConnected,
  isArcTestnet,
  isSwitchingChain,
  onSwitchChain,
}: {
  state: SubmitState
  deployState: DeployState
  isUploadingImage: boolean
  canSubmit: boolean
  isConnected: boolean
  isArcTestnet: boolean
  isSwitchingChain: boolean
  onSwitchChain: () => void
}) {
  return (
    <>
      {!isConnected ? (
        <p className="text-sm text-muted-foreground">
          Connect a wallet to deploy the market contract before it becomes public.
        </p>
      ) : null}

      {isConnected && !isArcTestnet ? (
        <ChainStatusCard
          message="Wallet is not on Arc Testnet."
          switchLabel="Switch to Arc Testnet"
          isSwitchingChain={isSwitchingChain}
          onSwitchNetwork={onSwitchChain}
        />
      ) : null}

      {deployState.status === "deploying" ? (
        <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-4">
          <p className="text-sm font-medium text-blue-200">
            Deploying market contract on Arc Testnet
          </p>
          {deployState.hash ? (
            <p className="mt-2 text-sm text-muted-foreground">
              Transaction: <TransactionLink hash={deployState.hash} />
            </p>
          ) : (
            <p className="mt-2 text-sm text-muted-foreground">
              Confirm the wallet transaction to continue. The backend save happens only after deployment succeeds.
            </p>
          )}
        </div>
      ) : null}

      {state.status === "error" ? (
        <InlineErrorState
          title="Unable to create market"
          message={state.message}
          requestId={state.requestId}
        />
      ) : null}

      <div className="flex flex-col gap-3 sm:flex-row">
        <Button
          disabled={!canSubmit || state.status === "submitting" || isUploadingImage}
          type="submit"
        >
          {state.status === "submitting"
            ? deployState.status === "deploying"
              ? "Deploying..."
              : "Creating..."
            : isUploadingImage
              ? "Uploading image..."
              : "Deploy and Create Market"}
        </Button>
        <Button asChild variant="outline">
          <Link href="/markets">Cancel</Link>
        </Button>
      </div>
    </>
  )
}
