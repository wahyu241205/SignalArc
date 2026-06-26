import { parseUnits } from "viem"

import { USDC_ERC20_DECIMALS } from "@/lib/contracts"

import type { TradeDisabledReasonInput } from "./types"

export function parseUsdcAmount(value: string) {
  const normalized = value.trim()

  if (!normalized) {
    throw new Error("Enter a USDC amount.")
  }

  if (!/^\d+(\.\d{1,6})?$/.test(normalized)) {
    throw new Error("Enter a valid USDC amount with up to 6 decimals.")
  }

  const amount = parseUnits(normalized, USDC_ERC20_DECIMALS)
  if (amount <= BigInt(0)) {
    throw new Error("USDC amount must be greater than 0.")
  }

  return amount
}

export function getTradeDisabledReason(input: TradeDisabledReasonInput) {
  if (!input.contractAddress) {
    return "Onchain contract not deployed for this market."
  }
  if (input.isContractTradingClosed || input.hasReachedCloseTime) {
    return "Trading is closed for this market."
  }
  if (!input.isConnected) return "Connect wallet to trade."
  if (!input.isArcTestnet) return "Switch to Arc Testnet."
  if (!input.isTradingOpen) return "Trading is not open for this market."
  if (input.isContractOpen === false) {
    return "Trading is closed for this market."
  }
  if (!input.parsedAmount) return "Enter a valid USDC amount."

  return null
}
