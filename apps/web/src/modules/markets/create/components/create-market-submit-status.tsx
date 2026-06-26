import Link from "next/link"

import { InlineErrorState } from "@/components/shared"
import { Button } from "@/components/ui/button"

import type { SubmitState } from "../types"

export function CreateMarketSubmitStatus({
  state,
  isUploadingImage,
}: {
  state: SubmitState
  isUploadingImage: boolean
}) {
  return (
    <>
      {state.status === "error" ? (
        <InlineErrorState
          title="Unable to create market"
          message={state.message}
          requestId={state.requestId}
        />
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
    </>
  )
}
