import type { TradeOutcome } from "./types"

export const MARKET_STATUS_OPEN = 1

const outcomeSide: Record<TradeOutcome, 1 | 2> = {
  YES: 1,
  NO: 2,
}

export function getOutcomeSide(outcome: TradeOutcome) {
  return outcomeSide[outcome]
}
