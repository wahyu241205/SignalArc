import type { Address, Hash } from "viem"

import type { Market } from "@/lib/api"

export type SubmitState =
  | { status: "idle" }
  | { status: "submitting" }
  | { status: "success"; market: Market }
  | { status: "error"; message: string; requestId: string | null }

export type DeployState =
  | { status: "idle" }
  | { status: "deploying"; hash?: Hash }
  | { status: "success"; hash: Hash; marketAddress: Address }
  | { status: "error"; message: string; hash?: Hash }
