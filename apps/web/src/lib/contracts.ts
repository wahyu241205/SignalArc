import type { Address } from "viem"

export const ARC_TESTNET_CHAIN_ID = 5042002
export const ARC_TESTNET_EXPLORER_URL = "https://testnet.arcscan.app"
export const SIGNAL_ARC_MARKET_ADDRESS =
  "0xf4ccc11A9e24fb996679F946C23C04AFd2797F26" as Address
export const ARC_TESTNET_USDC_ADDRESS =
  "0x3600000000000000000000000000000000000000" as Address
export const USDC_ERC20_DECIMALS = 6

export const SIGNAL_ARC_MARKET_ABI = [
  {
    type: "function",
    name: "isOpen",
    stateMutability: "view",
    inputs: [],
    outputs: [{ name: "", type: "bool" }],
  },
  {
    type: "function",
    name: "openPosition",
    stateMutability: "nonpayable",
    inputs: [
      { name: "side", type: "uint8" },
      { name: "amount", type: "uint256" },
    ],
    outputs: [],
  },
] as const

export const ERC20_APPROVE_ABI = [
  {
    type: "function",
    name: "approve",
    stateMutability: "nonpayable",
    inputs: [
      { name: "spender", type: "address" },
      { name: "amount", type: "uint256" },
    ],
    outputs: [{ name: "", type: "bool" }],
  },
] as const

export function getArcscanTxUrl(hash: string) {
  return `${ARC_TESTNET_EXPLORER_URL}/tx/${hash}`
}
